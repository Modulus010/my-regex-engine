package nfa

const EPS byte = 0
const WILD byte = '.'

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

func (s *State) match(ctx *Context) bool {
	for _, nxt := range s.transitions[EPS] {
		if nxt.match(ctx) {
			return true
		}
	}

	if ctx.isEnd() {
		return s == ctx.accept
	}

	ch := ctx.cur()
	ctx.pos++
	for _, nxt := range s.transitions[ch] {
		if nxt.match(ctx) {
			return true
		}
	}

	for _, nxt := range s.transitions[WILD] {
		if nxt.match(ctx) {
			return true
		}
	}
	ctx.pos--

	return false
}
