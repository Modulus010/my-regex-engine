package parse

type sequence struct {
	child []Node
}

func (s *sequence) String () string {
    res := ""
    for _, c := range s.child {
        res += c.String()
    }
	return res
}