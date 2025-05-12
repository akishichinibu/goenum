package render

import (
	"fmt"
	"strings"

	"github.com/akishichinibu/goenum/internal/model"
)

type EnumNaming struct {
	Interface        string
	Builder          string
	BuilderImpl      string
	BuilderSingleton string

	EqualMethodName string
	MatchMethodName string
}

func NewEnumNaming(m *model.Enum) *EnumNaming {
	nativeName := m.Spec.Name.Name
	name := strings.TrimPrefix(nativeName, "_G_")
	return &EnumNaming{
		Interface:        fmt.Sprintf("Enum%sVariant", name),
		Builder:          fmt.Sprintf("Enum%sBuilder", name),
		BuilderImpl:      fmt.Sprintf("enum%sBuilder", name),
		BuilderSingleton: fmt.Sprintf("Enum%s", name),
		EqualMethodName:  "Equal",
		MatchMethodName:  "Match",
	}
}

func (e *EnumNaming) VariantInterfaceName(variant *model.Variant) string {
	return fmt.Sprintf("%s%s", e.Interface, variant.Name)
}

func (e *EnumNaming) VariantImplName(variant *model.Variant) string {
	return fmt.Sprintf("_%s%s", e.Interface, variant.Name)
}

func (e *EnumNaming) VariantBuilderName(variant *model.Variant) string {
	return fmt.Sprintf("%s%sBuilder", e.Interface, variant.Name)
}

func (e *EnumNaming) VariantBuilderImplName(variant *model.Variant) string {
	return fmt.Sprintf("_%s%sBuilder", e.Interface, variant.Name)
}

func (e *EnumNaming) VariantBuilderSingletonName(variant *model.Variant) string {
	return fmt.Sprintf("%s%s_", e.BuilderSingleton, variant.Name)
}

func (e *EnumNaming) VariantBuilderSingletonImplName(variant *model.Variant) string {
	return fmt.Sprintf("%s%s_", e.BuilderImpl, variant.Name)
}

func (e *EnumNaming) ParamsPrivateFieldName(param *model.Param) string {
	return fmt.Sprintf("_%s", param.Name)
}

func (e *EnumNaming) ParamsPublicFieldName(param *model.Param) string {
	return fmt.Sprintf("%s%s", strings.ToUpper(string(param.Name[0])), param.Name[1:])
}

func (e *EnumNaming) ParamsReturnValueName(param *model.Param) string {
	return fmt.Sprintf("%s%s", strings.ToLower(string(param.Name[0])), param.Name[1:])
}

func (e *EnumNaming) ParamsGetterName(param *model.Param) string {
	return fmt.Sprintf("Get%s", e.ParamsPublicFieldName(param))
}

func (e *EnumNaming) ParamsTypeName(param *model.Param) string {
	return fmt.Sprintf("%s%s", e.Interface, param.Name)
}

func (e *EnumNaming) ParamsPrivateMemberInVariant(param *model.Param) string {
	return fmt.Sprintf("_%s", param.Name)
}

// func (p *EnumParamMeta) PrivateName() string {
// 	return "_" + p.Name
// }

// func (p *EnumParamMeta) PublicName() string {
// 	return strings.ToUpper(string(p.Name[0])) + p.Name[1:]
// }

// func (p *EnumParamMeta) GetterName() string {
// 	return "Get" + p.PublicName()
// }

// func (p *EnumParamMeta) nameType() string {
// 	return exprToString(p.typ)
// }
