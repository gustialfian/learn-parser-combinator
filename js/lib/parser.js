'use strict'

const updateParserState = (state, idx, result) => ({...state, idx, result})

const updateParserResult = (state, result) => ({...state, result})

const updateParserError = (state, error) => ({ ...state, error, isError: true })

class Parser {
  constructor(parserFn) {
    this.parserFn = parserFn
  }

  run(target) {
    return this.parserFn({ 
      target, 
      idx: 0, 
      result: '',
      isError: false, 
      error: '',
    })
  }

  map(fn) {
    return new Parser((parserState) => {
      const next = this.parserFn(parserState)
      if (next.isError) return next
      return updateParserResult(next, fn(next.result))
    })
  }

  mapError(fn) {
    return new Parser((parserState) => {
      const next = this.parserFn(parserState)
      if (!next.isError) return next
      return updateParserError(next, fn(next.error, next.idx))
    })
  }
}

const str = s => new Parser(parserState => {
  const { 
    target, 
    idx,
    isError, 
  } = parserState

  if (isError) {
    return parserState
  }

  const targetSliced = target.slice(idx)

  if (targetSliced.length === 0) {
    return updateParserError(parserState, 'str: unecpected EOF')
  }

  if (!targetSliced.startsWith(s)) {
    return updateParserError(parserState, `str: "${s}", "${targetSliced}"`)
  }

  return updateParserState(parserState, idx + s.length, s)
})

const letters = new Parser(parserState => {
  const { 
    target, 
    idx,
    isError, 
  } = parserState

  if (isError) {
    return parserState
  }

  const targetSliced = target.slice(idx)

  if (targetSliced.length === 0) {
    return updateParserError(parserState, 'letters: unecpected EOF')
  }

  const m = targetSliced.match(/^[A-Za-z]+/)
  if (m === null) {
    return updateParserError(parserState, `letters: could not match on idx ${idx}`)
  }

  return updateParserState(parserState, idx + m[0].length, m[0])
})

const digits = new Parser(parserState => {
  const { 
    target, 
    idx,
    isError, 
  } = parserState

  if (isError) {
    return parserState
  }

  const targetSliced = target.slice(idx)

  if (targetSliced.length === 0) {
    return updateParserError(parserState, 'digits: unecpected EOF')
  }

  const m = targetSliced.match(/^[0-9]+/)
  if (!m) {
    return updateParserError(parserState, `digits: could not match on idx ${idx}`)
  }

  return updateParserState(parserState, idx + m[0].length, m[0])
})

const sequenceOf = parsers => new Parser(parserState => {
  if (parserState.isError) {
    return parserState
  }

  const results = []
  let nextState = parserState

  for (let p of parsers) {
    nextState = p.parserFn(nextState)
    results.push(nextState.result)
  }

  return updateParserResult(nextState, results)
})

const choice = parsers => new Parser(parserState => {
  if (parserState.isError) {
    return parserState
  }

  for (let p of parsers) {
    const next = p.parserFn(parserState)
    if (next.isError) continue
    return updateParserResult(next, next.result)
  }

  return updateParserError(parserState, `choice: could not match on idx ${parserState.idx}`)
})

const many = parser => new Parser(parserState => {
  if (parserState.isError) {
    return parserState
  }

  let next = parserState
  const result = []

  while (true) {
    const test = parser.parserFn(next)
    if (test.isError) {
      break;
    }
    result.push(test.result)
    next = test
  }

  return updateParserResult(next, result)
})

const many1 = parser => new Parser(parserState => {
  if (parserState.isError) {
    return parserState
  }

  let next = parserState
  const result = []

  while (true) {
    const test = parser.parserFn(next)
    if (test.isError) {
      break;
    }
    result.push(test.result)
    next = test
  }

  if (result.length === 0) {
    return updateParserError(parserState, `many: could not match on idx ${parserState.idx}`)
  }

  return updateParserResult(next, result)
})

const between = (leftParser, rightParser) => (contentParser) => sequenceOf([
  leftParser,
  contentParser,
  rightParser,
]).map(result => result[1])


module.exports = {
  Parser,
  str,
  letters,
  digits,
  sequenceOf,
  choice,
  many,
  many1,
  between,
}
// main
// const betweenBrackets = between(str('('), str(')'))

// const parser = betweenBrackets(letters)
// const target = '(hello)'

// console.log(parser.run(target));