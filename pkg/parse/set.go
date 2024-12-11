package parse

import "github.com/Modulus010/my-regex-engine/pkg/nfa"

type set struct {
	chr map[byte]bool
}

func (s *set) String() string {
	res := "["
	for k := range s.chr {
		res += string(k)
	}
	res += "]"
	return res
}

func (s *set) toNFA() *nfa.NFA {
	start := nfa.NewState()
	accept := nfa.NewState()
	for k := range s.chr {
		start.Add(k, accept)
	}
	return nfa.NewNFA(start, accept)
}
