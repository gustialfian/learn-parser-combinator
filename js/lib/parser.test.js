'use strict'

const tap = require('tap')
const P = require('./parser')

tap.test('str', async (t) => {
  const parser = P.str('abc123')

  const got = parser.run('abc123')

  t.has(got, {
    target: 'abc123',
    idx: 6,
    result: 'abc123',
    isError: false,
    error: ''
  })
})