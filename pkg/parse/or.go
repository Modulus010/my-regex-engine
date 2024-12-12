package parse

import "github.com/Modulus010/my-regex-engine/pkg/nfa"

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

func (o *or) ToNFA() *nfa.NFA {
	start := nfa.NewState()
	accept := nfa.NewState()
	for _, c := range o.child {
		chNFA := c.ToNFA()
		start.Add(nfa.EPS, chNFA.Start)
		chNFA.Accept.Add(nfa.EPS, accept)
	}
	return nfa.NewNFA(start, accept)
}
