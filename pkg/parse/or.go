package parse

type or struct {
	child []Node
}

func (o *or) String() string {
	if len(o.child) == 0 {
		return ""
	}
	res := "(" + o.child[0].String()
	for _, c := range o.child[1:] {
		res += "|" + c.String()
	}
	return res + ")"
}
