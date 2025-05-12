package render

import j "github.com/dave/jennifer/jen"

type Emitter func(s ...*j.Statement)

type Renderable func(Emitter) error

func ChainRender(es ...Renderable) (ss []*j.Statement, err error) {
	var emitter = func(s ...*j.Statement) {
		ss = append(ss, j.Line())
		ss = append(ss, s...)
		ss = append(ss, j.Line())
	}
	for _, e := range es {
		if err := e(emitter); err != nil {
			return nil, err
		}
	}
	return ss, nil
}

type Renderer interface {
	Render() ([]*j.Statement, error)
}

func ToCode(ss []*j.Statement) []j.Code {
	codes := make([]j.Code, len(ss))
	for i, s := range ss {
		codes[i] = s
	}
	return codes
}
