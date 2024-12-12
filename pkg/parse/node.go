package parse

import (
	"fmt"

	"github.com/Modulus010/my-regex-engine/pkg/nfa"
)

type Node interface {
	fmt.Stringer

	ToNFA() *nfa.NFA
}
