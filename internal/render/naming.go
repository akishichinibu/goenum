package render

import (
	"fmt"
	"strings"

	"github.com/akishichinibu/goenum/internal/model"
)

type naming struct {
	Interface        string
	Builder          string
	BuilderImpl      string
	BuilderSingleton string

	EqualMethodName string
	MatchMethodName string
	ValueMethodName string
}

func newNaming(m *model.Enum) *naming {
	nativeName := m.Spec.Name.Name
	name := strings.TrimPrefix(nativeName, model.EnumDefinitionPrefix)

	return &naming{
		Interface:        fmt.Sprintf("Enum%sVariant", name),
		Builder:          fmt.Sprintf("Enum%sBuilder", name),
		BuilderImpl:      fmt.Sprintf("enum%sBuilder", name),
		BuilderSingleton: "Enum" + name,
		EqualMethodName:  "Equal",
		MatchMethodName:  "Match",
		ValueMethodName:  "Value",
	}
}

func (e *naming) VariantInterfaceName(variant *model.Variant) string {
	return fmt.Sprintf("%s%s", e.Interface, variant.Name)
}

func (e *naming) VariantImplName(variant *model.Variant) string {
	return fmt.Sprintf("_%s%s", e.Interface, variant.Name)
}

func (e *naming) VariantBuilderName(variant *model.Variant) string {
	return fmt.Sprintf("%s%sBuilder", e.Interface, variant.Name)
}

func (e *naming) VariantBuilderImplName(variant *model.Variant) string {
	return fmt.Sprintf("_%s%sBuilder", e.Interface, variant.Name)
}

func (e *naming) VariantBuilderSingletonName(variant *model.Variant) string {
	return fmt.Sprintf("%s%s_", e.BuilderSingleton, variant.Name)
}

func (e *naming) VariantBuilderSingletonImplName(variant *model.Variant) string {
	return fmt.Sprintf("%s%s_", e.BuilderImpl, variant.Name)
}

func (e *naming) ParamsPrivateFieldName(param *model.Param) string {
	return "_" + param.Name
}

func (e *naming) ParamsPublicFieldName(param *model.Param) string {
	return fmt.Sprintf("%s%s", strings.ToUpper(string(param.Name[0])), param.Name[1:])
}

func (e *naming) ParamsReturnValueName(param *model.Param) string {
	return fmt.Sprintf("%s%s", strings.ToLower(string(param.Name[0])), param.Name[1:])
}

func (e *naming) ParamsGetterName(param *model.Param) string {
	return "Get" + e.ParamsPublicFieldName(param)
}

func (e *naming) ParamsTypeName(param *model.Param) string {
	return fmt.Sprintf("%s%s", e.Interface, param.Name)
}

func (e *naming) ParamsPrivateMemberInVariant(param *model.Param) string {
	return "_" + param.Name
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
