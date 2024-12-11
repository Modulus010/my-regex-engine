package parse

type literal byte

func (l literal) String() string {
    return string(l)
}
    