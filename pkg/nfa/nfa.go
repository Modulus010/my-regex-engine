package nfa

type NFA struct {
	Start  *State
	Accept *State
}

func NewNFA(start, accept *State) *NFA {
	return &NFA{start, accept}
}
