package generator

import (
	"strings"
	"testing"

	"github.com/iamrekas/anchor-go/pkg/idl/types"
)

// TestGenerateStruct_OptionFieldUsesOptionalWrapper pins the new emission shape
// for option<T> fields. The wrapper Optional[T] (emitted into the same package
// by OptionalGenerator) owns the Borsh discriminator, so the generated struct
// stores it by value and the (un)marshal bodies do nothing more than call the
// reflective Encode / Decode — ag_binary dispatches to the wrapper's own
// MarshalWithEncoder / UnmarshalWithDecoder via BinaryMarshaler.
func TestGenerateStruct_OptionFieldUsesOptionalWrapper(t *testing.T) {
	g := NewTypesGenerator()
	g.SetConfig(&Config{})

	typeDef := types.TypeDef{
		Name: "SwapParams",
		Kind: types.TypeDefKindStruct,
		Fields: []types.Field{
			{
				Name: "once_seed",
				Type: &types.OptionType{
					ElementType: &types.ArrayType{
						ElemType:  &types.BasicType{TypeName: "u8"},
						ArraySize: 32,
					},
				},
			},
			{
				Name: "tax_bps",
				Type: &types.BasicType{TypeName: "u16"},
			},
		},
	}

	code := g.generateStruct(typeDef)

	mustContain := []string{
		// Field is stored by value via the generic wrapper. No `*`, no bin tag.
		"OnceSeed Optional[[32]uint8]",
		// Marshal collapses to a plain Encode of the wrapper.
		"err = encoder.Encode(obj.OnceSeed)",
		// Unmarshal collapses to a plain Decode of the wrapper.
		"err = decoder.Decode(&obj.OnceSeed)",
	}
	for _, want := range mustContain {
		if !strings.Contains(code, want) {
			t.Errorf("generated code missing expected snippet: %q\nfull code:\n%s", want, code)
		}
	}

	// Negative: per-call-site option discriminator code must be gone. Any of
	// these would indicate the old logic survived in this code path.
	mustNotContain := []string{
		"if obj.OnceSeed == nil {",
		"WriteByte(0)",
		"WriteByte(1)",
		"optTag",
		"obj.OnceSeed = &v",
		"obj.OnceSeed = nil",
		"`bin:\"optional\"`",
	}
	for _, bad := range mustNotContain {
		if strings.Contains(code, bad) {
			t.Errorf("generated code still emits old option logic: %q\nfull code:\n%s", bad, code)
		}
	}
}

// TestGenerateStruct_OptionScalarUsesOptionalWrapper covers option<u64> — the
// case that previously produced **uint64 in instruction structs because the
// pointer layer was added twice. With the wrapper, the field is just a value.
func TestGenerateStruct_OptionScalarUsesOptionalWrapper(t *testing.T) {
	g := NewTypesGenerator()
	g.SetConfig(&Config{})

	typeDef := types.TypeDef{
		Name: "MaybeAmount",
		Kind: types.TypeDefKindStruct,
		Fields: []types.Field{
			{
				Name: "amount",
				Type: &types.OptionType{ElementType: &types.BasicType{TypeName: "u64"}},
			},
		},
	}

	code := g.generateStruct(typeDef)

	for _, want := range []string{
		"Amount Optional[uint64]",
		"err = encoder.Encode(obj.Amount)",
		"err = decoder.Decode(&obj.Amount)",
	} {
		if !strings.Contains(code, want) {
			t.Errorf("generated code missing snippet: %q\nfull code:\n%s", want, code)
		}
	}
	for _, bad := range []string{"*uint64", "**uint64", "WriteByte(0)", "WriteByte(1)", "optTag"} {
		if strings.Contains(code, bad) {
			t.Errorf("generated code still emits old option logic: %q\nfull code:\n%s", bad, code)
		}
	}
}
