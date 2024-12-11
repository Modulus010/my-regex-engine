package parse

import (
	"fmt"
	"strconv"
	"strings"
)

type repeat struct {
	child      Node
	quantifier *quantifier
}

const quantifierInf = -1

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

	pieces := strings.Split(substring, ",")
	if len(pieces) == 1 {
		if value, err := strconv.Atoi(pieces[0]); err != nil {
			return err
		} else {
			r.quantifier = &quantifier{value, value}
		}
	} else if len(pieces) == 2 {
		min, err := strconv.Atoi(pieces[0])
		if err != nil {
			return err
		}
		if pieces[1] == "" {
			r.quantifier = &quantifier{min, quantifierInf}
		} else if max, err := strconv.Atoi(pieces[1]); err != nil {
			return err
		} else {
			r.quantifier = &quantifier{min, max}
		}
	} else {
		return fmt.Errorf("invalid quantifier: %s", substring)
	}
	if err := ctx.consume('}'); err != nil {
		return err
	}
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
