package parse

import "github.com/Modulus010/my-regex-engine/pkg/nfa"

type literal byte

func (l literal) String() string {
	return string(l)
}

func (l literal) ToNFA() *nfa.NFA {
	start := nfa.NewState()
	accept := nfa.NewState()

	start.Add(byte(l), accept)
	return nfa.NewNFA(start, accept)
}
