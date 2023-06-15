package main

import (
	"fmt"
	"regexp"
)

func main() {
	fmt.Println("learn parser combinator")

	p := Many1(Choice([]*Parser{
		Letters(),
		Digits(),
	}))

	res := p.Run("2d8")
	fmt.Printf("%+v", res)
}

type State struct {
	target  string
	idx     int
	result  interface{}
	isError bool
	err     string
}

func UpdateParserState(state State, result interface{}, idx int) State {
	return State{
		target:  state.target,
		idx:     idx,
		result:  result,
		isError: state.isError,
		err:     state.err,
	}
}

func UpdateParserResult(state State, result interface{}) State {
	return State{
		target:  state.target,
		idx:     state.idx,
		result:  result,
		isError: state.isError,
		err:     state.err,
	}
}

func UpdateParserError(state State, err string) State {
	return State{
		target:  state.target,
		idx:     state.idx,
		result:  state.result,
		isError: true,
		err:     err,
	}
}

type Parser struct {
	fn func(State) State
}

func (p *Parser) Run(target string) State {
	return p.fn(State{
		target:  target,
		idx:     0,
		result:  nil,
		isError: false,
		err:     "",
	})
}

func (p *Parser) Map(fn func(v interface{}) interface{}) *Parser {
	return &Parser{
		fn: func(state State) State {
			next := p.fn(state)
			if next.isError {
				return next
			}
			return UpdateParserResult(next, fn(next.result))
		},
	}
}

func (p *Parser) MapErr(fn func(v string) string) *Parser {
	return &Parser{
		fn: func(state State) State {
			next := p.fn(state)
			if !next.isError {
				return next
			}
			return UpdateParserError(next, fn(next.err))
		},
	}
}

func Str(s string) *Parser {
	return &Parser{
		fn: func(state State) State {
			if state.isError {
				return state
			}

			target := state.target[state.idx:]

			if len(target) == 0 {
				return UpdateParserError(state, "str: unexpected end of input")
			}

			if s != target[:len(s)] {
				return UpdateParserError(state, fmt.Sprintf("str: could not match on idx %d", state.idx))
			}

			return UpdateParserState(state, s, len(s)+state.idx)
		},
	}
}

func Letters() *Parser {
	return &Parser{
		fn: func(state State) State {
			if state.isError {
				return state
			}

			target := state.target[state.idx:]

			if len(target) == 0 {
				return UpdateParserError(state, "letters: unexpected end of input")
			}

			re := regexp.MustCompile(`^[A-Za-z]+`)
			m := re.FindString(target)
			if len(m) == 0 {
				return UpdateParserError(state, fmt.Sprintf("letters: could not match on idx %d", state.idx))
			}

			return UpdateParserState(state, m, len(m)+state.idx)
		},
	}
}

func Digits() *Parser {
	return &Parser{
		fn: func(state State) State {
			if state.isError {
				return state
			}

			target := state.target[state.idx:]

			if len(target) == 0 {
				return UpdateParserError(state, "digits: unexpected end of input")
			}

			re := regexp.MustCompile(`^[0-9]+`)
			m := re.FindString(target)
			if len(m) == 0 {
				return UpdateParserError(state, fmt.Sprintf("digits: could not match on idx %d", state.idx))
			}

			return UpdateParserState(state, m, len(m)+state.idx)
		},
	}
}

func SequenceOf(ps []*Parser) *Parser {
	return &Parser{
		fn: func(state State) State {
			if state.isError {
				return state
			}

			var res []interface{}
			next := state

			for _, p := range ps {
				next = p.fn(next)
				res = append(res, next.result)
			}

			return UpdateParserResult(next, res)
		},
	}
}

func Choice(ps []*Parser) *Parser {
	return &Parser{
		fn: func(state State) State {
			if state.isError {
				return state
			}

			for _, p := range ps {
				next := p.fn(state)
				if next.isError {
					continue
				}
				return UpdateParserResult(next, next.result)
			}

			return UpdateParserError(state, `choice: could not match on idx ${parserState.idx}`)
		},
	}
}

func Many(p *Parser) *Parser {
	return &Parser{
		fn: func(state State) State {
			if state.isError {
				return state
			}

			var res []interface{}
			next := state

			for {
				test := p.fn(next)
				if test.isError {
					break
				}
				res = append(res, test.result)
				next = test
			}

			return UpdateParserResult(next, res)
		},
	}
}

func Many1(p *Parser) *Parser {
	return &Parser{
		fn: func(state State) State {
			if state.isError {
				return state
			}

			var res []interface{}
			next := state

			for {
				test := p.fn(next)
				if test.isError {
					break
				}
				res = append(res, test.result)
				next = test
			}

			if len(res) == 0 {
				return UpdateParserError(state, fmt.Sprintf(`many: could not match on idx %d`, state.idx))
			}

			return UpdateParserResult(next, res)
		},
	}
}
