package meander

type Facade interface {
	Pubclic() interface{}
}

func Public(o interface{}) interface{} {
	if p, ok := o.(Facade); ok {
		return p.Public()
	}
	return o
}
