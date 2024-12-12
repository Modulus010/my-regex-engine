package parse

import "github.com/Modulus010/my-regex-engine/pkg/nfa"

type sequence struct {
	child []Node
}

func (s *sequence) String() string {
	res := ""
	for _, c := range s.child {
		res += c.String()
	}
	return res
}

func (s *sequence) ToNFA() *nfa.NFA {
	start := nfa.NewState()
	accept := start

	for _, c := range s.child {
		chNFA := c.ToNFA()
		accept.Add(nfa.EPS, chNFA.Start)
		accept = chNFA.Accept
	}
	return nfa.NewNFA(start, accept)
}
