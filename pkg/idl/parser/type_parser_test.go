package parser

import (
	"encoding/json"
	"testing"

	"github.com/iamrekas/anchor-go/pkg/idl/types"
)

// TestParseFieldType_CompositeOptions is a regression test for the bug where
// IDL fields typed as option<array<T, N>> (and other composites such as
// option<vec<T>>, vec<array<T, N>>, etc.) were not recognized by the parser
// and fell back to types.BasicType{TypeName: "string"}.
func TestParseFieldType_CompositeOptions(t *testing.T) {
	tests := []struct {
		name     string
		json     string
		validate func(t *testing.T, got types.Type)
	}{
		{
			name: "option<array<u8, 32>>",
			json: `{"option":{"array":["u8",32]}}`,
			validate: func(t *testing.T, got types.Type) {
				opt, ok := got.(*types.OptionType)
				if !ok {
					t.Fatalf("want *OptionType, got %T (%s)", got, got.String())
				}
				arr, ok := opt.ElementType.(*types.ArrayType)
				if !ok {
					t.Fatalf("inner: want *ArrayType, got %T", opt.ElementType)
				}
				if arr.ArraySize != 32 {
					t.Errorf("size: want 32, got %d", arr.ArraySize)
				}
				elem, ok := arr.ElemType.(*types.BasicType)
				if !ok || elem.TypeName != "u8" {
					t.Errorf("elem: want BasicType{u8}, got %#v", arr.ElemType)
				}
			},
		},
		{
			name: "option<vec<u8>>",
			json: `{"option":{"vec":"u8"}}`,
			validate: func(t *testing.T, got types.Type) {
				opt, ok := got.(*types.OptionType)
				if !ok {
					t.Fatalf("want *OptionType, got %T", got)
				}
				vec, ok := opt.ElementType.(*types.VectorType)
				if !ok {
					t.Fatalf("inner: want *VectorType, got %T", opt.ElementType)
				}
				elem, ok := vec.ElementType.(*types.BasicType)
				if !ok || elem.TypeName != "u8" {
					t.Errorf("elem: want BasicType{u8}, got %#v", vec.ElementType)
				}
			},
		},
		{
			name: "vec<array<u8, 4>>",
			json: `{"vec":{"array":["u8",4]}}`,
			validate: func(t *testing.T, got types.Type) {
				vec, ok := got.(*types.VectorType)
				if !ok {
					t.Fatalf("want *VectorType, got %T", got)
				}
				arr, ok := vec.ElementType.(*types.ArrayType)
				if !ok {
					t.Fatalf("inner: want *ArrayType, got %T", vec.ElementType)
				}
				if arr.ArraySize != 4 {
					t.Errorf("size: want 4, got %d", arr.ArraySize)
				}
			},
		},
		{
			name: "option<defined> (regression: existing case must still work)",
			json: `{"option":{"defined":{"name":"Foo"}}}`,
			validate: func(t *testing.T, got types.Type) {
				opt, ok := got.(*types.OptionType)
				if !ok {
					t.Fatalf("want *OptionType, got %T", got)
				}
				def, ok := opt.ElementType.(*types.DefinedType)
				if !ok || def.Name != "Foo" {
					t.Errorf("inner: want DefinedType{Foo}, got %#v", opt.ElementType)
				}
			},
		},
		{
			name: "option<u128> (regression: scalar inner must still work)",
			json: `{"option":"u128"}`,
			validate: func(t *testing.T, got types.Type) {
				opt, ok := got.(*types.OptionType)
				if !ok {
					t.Fatalf("want *OptionType, got %T", got)
				}
				b, ok := opt.ElementType.(*types.BasicType)
				if !ok || b.TypeName != "u128" {
					t.Errorf("inner: want BasicType{u128}, got %#v", opt.ElementType)
				}
			},
		},
		{
			name: "vec<defined> (regression: existing case must still work)",
			json: `{"vec":{"defined":{"name":"Bar"}}}`,
			validate: func(t *testing.T, got types.Type) {
				vec, ok := got.(*types.VectorType)
				if !ok {
					t.Fatalf("want *VectorType, got %T", got)
				}
				def, ok := vec.ElementType.(*types.DefinedType)
				if !ok || def.Name != "Bar" {
					t.Errorf("inner: want DefinedType{Bar}, got %#v", vec.ElementType)
				}
			},
		},
		{
			name: "array<defined, N> (regression: existing case must still work)",
			json: `{"array":[{"defined":{"name":"Baz"}},8]}`,
			validate: func(t *testing.T, got types.Type) {
				arr, ok := got.(*types.ArrayType)
				if !ok {
					t.Fatalf("want *ArrayType, got %T", got)
				}
				if arr.ArraySize != 8 {
					t.Errorf("size: want 8, got %d", arr.ArraySize)
				}
				def, ok := arr.ElemType.(*types.DefinedType)
				if !ok || def.Name != "Baz" {
					t.Errorf("elem: want DefinedType{Baz}, got %#v", arr.ElemType)
				}
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := ParseFieldType(tc.name, json.RawMessage(tc.json), false)
			if !ok {
				t.Fatalf("ParseFieldType returned ok=false for %s", tc.json)
			}
			tc.validate(t, got)
		})
	}
}
