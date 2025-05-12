package model

type Param struct {
	Variant *Variant

	Name string
	Type Type
}

func newParam(variant *Variant, name string, typ Type) *Param {
	return &Param{
		Variant: variant,
		Name:    name,
		Type:    typ,
	}
}
