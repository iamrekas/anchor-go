package generator

import (
	"strings"
	"testing"

	"github.com/iamrekas/anchor-go/pkg/idl/types"
)

// TestGenerateStruct_OptionFieldEmitsBorshDiscriminator pins the per-field
// option<...> marshal/unmarshal emission. The bin:"optional" struct tag is
// only consulted by ag_binary when it walks struct fields via reflection;
// our generated MarshalWithEncoder body (a BinaryMarshaler) bypasses that
// walk, so the option discriminator byte must be written/read manually.
// Without that, Rust borsh sees a wrong first byte (the start of the
// payload) and rejects the buffer with "Invalid Option representation".
func TestGenerateStruct_OptionFieldEmitsBorshDiscriminator(t *testing.T) {
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
		// Struct field declaration is unchanged: pointer to fixed array with optional tag.
		"OnceSeed *[32]uint8 `bin:\"optional\"`",

		// Marshal: discriminator + payload pattern.
		"if obj.OnceSeed == nil {",
		"err = encoder.WriteByte(0)",
		"} else {",
		"err = encoder.WriteByte(1)",
		"err = encoder.Encode(obj.OnceSeed)",

		// Unmarshal: read discriminator, conditionally decode payload into &v.
		"optTag, err := decoder.ReadByte()",
		"if optTag == 1 {",
		"var v [32]uint8",
		"if err := decoder.Decode(&v); err != nil {",
		"obj.OnceSeed = &v",
		"obj.OnceSeed = nil",
	}
	for _, want := range mustContain {
		if !strings.Contains(code, want) {
			t.Errorf("generated code missing expected snippet: %q", want)
		}
	}

	// The old broken form was an *unguarded* encoder.Encode(obj.OnceSeed) at
	// one-tab indentation (i.e. directly inside the method body, not inside
	// the discriminator-guarded `else` block). Look for that exact line.
	if strings.Contains(code, "\n\terr = encoder.Encode(obj.OnceSeed)\n") {
		t.Errorf("generated marshal still emits unguarded encoder.Encode(obj.OnceSeed)\nfull code:\n%s", code)
	}
	// Likewise the old unmarshal form was an unguarded decoder.Decode at the
	// inner-of-Remaining indentation (two tabs).
	if strings.Contains(code, "\n\t\terr = decoder.Decode(&obj.OnceSeed)\n") {
		t.Errorf("generated unmarshal still emits unguarded decoder.Decode(&obj.OnceSeed)\nfull code:\n%s", code)
	}
}

// TestGenerateStruct_OptionScalarStillCorrect ensures option<u64> and other
// non-composite option types also flow through the new discriminator-aware
// path. ag_binary's bin:"optional" tag would have been similarly bypassed
// for these, so the prior emitter was producing un-discriminated u64 bytes
// for option<u64> as well, even though it happened to round-trip with
// itself.
func TestGenerateStruct_OptionScalarStillCorrect(t *testing.T) {
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
		"Amount *uint64 `bin:\"optional\"`",
		"if obj.Amount == nil {",
		"err = encoder.WriteByte(0)",
		"err = encoder.WriteByte(1)",
		"err = encoder.Encode(obj.Amount)",
		"optTag, err := decoder.ReadByte()",
		"var v uint64",
		"obj.Amount = &v",
	} {
		if !strings.Contains(code, want) {
			t.Errorf("generated code missing snippet: %q", want)
		}
	}
}
