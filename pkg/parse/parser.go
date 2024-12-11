package parse

import (
	"fmt"
	"strings"
)

type Parser struct {
	regex string
	pos   int
}

func NewParser(regex string) *Parser {
	return &Parser{regex: regex}
}

func (p *Parser) Parse() (Node, error) {
	nd, err := p.parseOr()
	if err != nil {
		return nil, err
	}
	if !p.isEnd() {
		return nil, parseErr{"EOF", p.cur(), p.pos}
	}
	return nd, nil
}

func (p *Parser) parseOr() (Node, error) {
	nd := &or{[]Node{}}
	for {
		if ch, err := p.parseSequence(); err != nil {
			return nil, err
		} else {
			nd.child = append(nd.child, ch)
		}
		if err := p.consume('|'); err != nil {
			break
		}
	}
	return nd, nil
}

func (p *Parser) parseSequence() (Node, error) {
	nd := &sequence{[]Node{}}
	for {
		if ch, err := p.parseRepeat(); err != nil {
			return nil, err
		} else if ch == nil {
			break
		} else {
			nd.child = append(nd.child, ch)
		}
	}
	return nd, nil
}

func (p *Parser) parseRepeat() (Node, error) {
	nd := &repeat{}
	if ch, err := p.parseAtom(); err != nil {
		return nil, err
	} else if ch == nil {
		return nil, nil
	} else {
		nd.child = ch
	}

	for {
		if err := nd.parseQuantifier(p); err != nil {
			return nil, err
		} else if nd.quantifier == nil {
			break
		} else {
			nd = &repeat{child: nd}
		}
	}
	return nd.child, nil
}

func (p *Parser) parseAtom() (Node, error) {
	switch p.cur() {
	case '(':
		return p.parseGroup()
	case '[':
		return p.parseSet()
	default:
		return p.parseLiteral()
	}
}

func (p *Parser) parseSet() (Node, error) {
	if err := p.consume('['); err != nil {
		return nil, err
	}

	var literals []string
	state := false
	for !p.isEnd() && p.cur() != ']' {
		ch := p.cur()
		if ch == '-' {
			state = true
		} else if state {
			literals[len(literals)-1] += string(ch)
			state = false
		} else {
			literals = append(literals, string(ch))
		}
		p.adv()
	}
	if err := p.consume(']'); err != nil {
		return nil, err
	}

	nd := &set{chr: make(map[byte]bool)}
	for _, l := range literals {
		for i := l[0]; i <= l[len(l)-1]; i++ {
			nd.chr[i] = true
		}
	}

	return nd, nil
}

func (p *Parser) parseLiteral() (Node, error) {
	if !isLiteral(p.cur()) {
		return nil, nil
	}
	nd := literal(p.cur())
	p.adv()
	return nd, nil
}

func (p *Parser) parseGroup() (Node, error) {
	if err := p.consume('('); err != nil {
		return nil, err
	}
	nd, err := p.parseOr()
	if err != nil {
		return nil, err
	}
	if err := p.consume(')'); err != nil {
		return nil, err
	}
	return nd, nil
}

func (p *Parser) cur() byte {
	if p.isEnd() {
		return 0
	}
	return p.regex[p.pos]
}

func (p *Parser) adv() {
	if !p.isEnd() {
		p.pos++
	}
}

func (p *Parser) isEnd() bool {
	return p.pos >= len(p.regex)
}

func (p *Parser) consume(b byte) error {
	if p.cur() != b {
		return parseErr{expected: string(b), got: p.cur(), pos: p.pos}
	}
	p.adv()
	return nil
}

func isLiteral(ch byte) bool {
	return !strings.Contains("|*+?{}[]()", string(ch)) && ch != 0
}

type parseErr struct {
	expected string
	got      byte
	pos      int
}

func (e parseErr) Error() string {
	got := "EOF"
	if e.got != 0 {
		got = string(e.got)
	}
	return fmt.Sprintf("parse error: expected %s, got %s at %d", e.expected, got, e.pos)
}
