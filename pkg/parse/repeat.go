package parse

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Modulus010/my-regex-engine/pkg/nfa"
)

type repeat struct {
	child      Node
	quantifier *quantifier
}

const quantifierInf = 1000

type quantifier struct {
	min int
	max int
}

func (r *repeat) parseQuantifier(ctx *Parser) error {
	if ctx.cur() == '{' {
		if err := r.parseQuantifierSpecified(ctx); err != nil {
			return err
		}
		return nil
	}
	switch ctx.cur() {
	case '*':
		r.quantifier = &quantifier{0, quantifierInf}
	case '+':
		r.quantifier = &quantifier{1, quantifierInf}
	case '?':
		r.quantifier = &quantifier{0, 1}
	default:
		r.quantifier = nil
		return nil
	}
	ctx.adv()
	return nil
}

func (r *repeat) parseQuantifierSpecified(ctx *Parser) error {
	if err := ctx.consume('{'); err != nil {
		return err
	}

	substring := ""
	for !ctx.isEnd() && ctx.cur() != '}' {
		substring += string(ctx.cur())
		ctx.adv()
	}
	var minv, maxv int
	pieces := strings.Split(substring, ",")
	if len(pieces) == 1 {
		if value, err := strconv.Atoi(pieces[0]); err != nil {
			return err
		} else {
			minv = value
			maxv = value
		}
	} else if len(pieces) == 2 {
		if min, err := strconv.Atoi(pieces[0]); err != nil {
			return err
		} else {
			minv = min
		}
		if pieces[1] == "" {
			maxv = quantifierInf
		} else if max, err := strconv.Atoi(pieces[1]); err != nil {
			return err
		} else {
			maxv = max
		}
	} else {
		return fmt.Errorf("invalid quantifier: %s", substring)
	}
	if err := ctx.consume('}'); err != nil {
		return err
	}
	if minv > maxv || minv < 0 || maxv < 0 || maxv > quantifierInf {
		return fmt.Errorf("invalid quantifier: %d, %d", minv, maxv)
	}
	r.quantifier = &quantifier{minv, maxv}
	return nil
}

func (r *repeat) String() string {
	res := ""
	res += r.child.String()
	if r.quantifier != nil {
		if r.quantifier.min == r.quantifier.max {
			res += fmt.Sprintf("{%d}", r.quantifier.min)
		} else if r.quantifier.max == quantifierInf {
			res += fmt.Sprintf("{%d,}", r.quantifier.min)
		} else {
			res += fmt.Sprintf("{%d,%d}", r.quantifier.min, r.quantifier.max)
		}
	}
	return res
}

func (r *repeat) ToNFA() *nfa.NFA {
	start := nfa.NewState()
	accept := start
	for i := 0; i < r.quantifier.min; i++ {
		chNFA := r.child.ToNFA()
		accept.Add(nfa.EPS, chNFA.Start)
		accept = chNFA.Accept
	}

	if r.quantifier.max == quantifierInf {
		chNFA := r.child.ToNFA()
		repeatStart := chNFA.Accept
		repeatAccept := chNFA.Start
		repeatStart.Add(nfa.EPS, repeatAccept)
		accept.Add(nfa.EPS, repeatStart)
		accept = repeatAccept
	} else {
		var midStates []*nfa.State
		for i := r.quantifier.min; i < r.quantifier.max; i++ {
			midStates = append(midStates, accept)
			chNFA := r.child.ToNFA()
			accept.Add(nfa.EPS, chNFA.Start)
			accept = chNFA.Accept
		}
		for _, midState := range midStates {
			midState.Add(nfa.EPS, accept)
		}
	}
	return nfa.NewNFA(start, accept)
}
