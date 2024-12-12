package nfa

type NFA struct {
	Start  *State
	Accept *State
}

func NewNFA(start, accept *State) *NFA {
	return &NFA{start, accept}
}

type Context struct {
	text   string
	pos    int
	accept *State
}

func (c *Context) isEnd() bool {
	return c.pos >= len(c.text)
}

func (c *Context) cur() byte {
	return c.text[c.pos]
}

// func (n *NFA) Match(input string) bool {
// 	return n.Start.match(
// 		&Context{
// 			text:   input,
// 			pos:    0,
// 			accept: n.Accept,
// 		})
// }

type stateSet map[*State]bool

func (n *NFA) Match(input string) bool {
	curSet := stateSet(make(map[*State]bool))
	curSet.add(n.Start)
	for i := range input {
		curSet = curSet.step(input[i])
	}
	for state := range curSet {
		if state == n.Accept {
			return true
		}
	}
	return false
}

func (s stateSet) add(state *State) {
	if s[state] {
		return
	}
	s[state] = true
	for _, nextState := range state.transitions[EPS] {
		s.add(nextState)
	}
}

func (s stateSet) step(c byte) stateSet {
	res := stateSet(make(map[*State]bool))
	for state := range s {
		for _, nextState := range state.transitions[c] {
			res.add(nextState)
		}
		for _, nextState := range state.transitions[WILD] {
			res.add(nextState)
		}
	}
	return res
}
