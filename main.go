package main

import (
	"fmt"
	"strconv"
	"strings"
)

type context struct {
	regex string
	pos   int
}

func (ctx context) cur() byte {
	return ctx.regex[ctx.pos]
}

func (ctx *context) adv() {
	ctx.pos++
}

func (ctx context) isEnd() bool {
	return ctx.pos >= len(ctx.regex)
}

func (ctx *context) consume(ch byte) {
	if ctx.isEnd() || ctx.cur() != ch {
		panic("syntax error")
	}
	ctx.adv()
}

type Node interface {
}

type exp struct {
	children Node
}

func newExp(ctx *context) exp {
	return exp{newOr(ctx)}
}

type or struct {
	children []Node
}

func newOr(ctx *context) (nd or) {
	nd.children = []Node{newConcat(ctx)}
	for ctx.cur() == '|' {
		ctx.adv()
		nd.children = append(nd.children, newConcat(ctx))
	}
	return nd
}

type concat struct {
	children []Node
}

func newConcat(ctx *context) (nd concat) {
	nd.children = []Node{newRepeat(ctx)}
	for !ctx.isEnd() {
		nd.children = append(nd.children, newRepeat(ctx))
	}
	return nd
}

type repeat struct {
	children   Node
	quantifier quantifier
}

func newRepeat(ctx *context) (nd repeat) {
	nd.children = newAtom(ctx)
	nd.quantifier = quantifier{1, 1}
	for ctx.cur() == '*' || ctx.cur() == '+' || ctx.cur() == '?' || ctx.cur() == '{' {
		nd = repeat{nd.children, newQuantifier(ctx)}
	}
	return nd
}

type quantifier struct {
	min, max int
}

const timeInf = -1

func newQuantifier(ctx *context) quantifier {
	switch ctx.cur() {
	case '*':
		ctx.adv()
		return quantifier{0, timeInf}
	case '+':
		ctx.adv()
		return quantifier{1, timeInf}
	case '?':
		ctx.adv()
		return quantifier{0, 1}
	case '{':
		return parseQuantifierSpecified(ctx)
	}
	return quantifier{1, 1}
}

func parseQuantifierSpecified(ctx *context) quantifier {
	ctx.consume('{')
	defer ctx.consume('}')

	substring := ""
	for ctx.cur() != '}' {
		substring += string(ctx.cur())
		ctx.adv()
	}

	pieces := strings.Split(substring, ",")
	if len(pieces) == 1 {
		if value, err := strconv.Atoi(pieces[0]); err != nil {
			panic(err.Error())
		} else {
			return quantifier{value, value}
		}
	} else if len(pieces) == 2 {
		min, err := strconv.Atoi(pieces[0])
		if err != nil {
			panic(err.Error())
		}

		if pieces[1] == "" {
			return quantifier{min, timeInf}
		} else if value, err := strconv.Atoi(pieces[1]); err != nil {
			panic(err.Error())
		} else {
			return quantifier{min, value}
		}
	} else {
		panic("syntax error")
	}
}

type atom struct {
	children Node
}

func newAtom(ctx *context) (nd atom) {
	switch ctx.cur() {
	case '(':
		nd.children = newGroup(ctx)
	case '[':
		nd.children = newSet(ctx)
	default:
		nd.children = newLiteral(ctx)
	}
	return nd
}

type literal byte

func newLiteral(ctx *context) (nd literal) {
	nd = literal(ctx.cur())
	ctx.adv()
	return nd
}

type set struct {
	val map[byte]bool
}

func newSet(ctx *context) (nd set) {
	ctx.consume('[')
	defer ctx.consume(']')

	var literals []string
	for ctx.cur() != ']' {
		ch := ctx.cur()
		if ch == '-' {
			ctx.adv()
			if ctx.isEnd() || ctx.cur() == ']' || len(literals) == 0 {
				panic("syntax error")
			}
			next := ctx.cur()
			prev := literals[len(literals)-1][0]
			literals[len(literals)-1] = fmt.Sprintf("%c%c", prev, next)
			ctx.adv()
		} else {
			literals = append(literals, string(ch))
		}
		ctx.adv()
	}

	nd.val = make(map[byte]bool)
	for _, l := range literals {
		for i := l[0]; i <= l[len(l)-1]; i++ {
			nd.val[i] = true
		}
	}
	return nd
}

type group struct {
	children Node
}

func newGroup(ctx *context) (nd group) {
	ctx.consume('(')
	defer ctx.consume(')')
	nd.children = newExp(ctx)
	return nd
}
