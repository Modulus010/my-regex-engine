package nfa

const EPS byte = 0

type State struct {
	transitions map[byte][]*State
}

func NewState() *State {
	return &State{
		transitions: make(map[byte][]*State),
	}
}

func (s *State) Add(c byte, to *State) {
	s.transitions[c] = append(s.transitions[c], to)
}
