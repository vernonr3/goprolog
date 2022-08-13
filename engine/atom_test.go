package engine

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAtom_WriteTerm(t *testing.T) {
	tests := []struct {
		atom   Atom
		opts   WriteOptions
		output string
	}{
		{atom: `a`, opts: WriteOptions{Quoted: false}, output: `a`},
		{atom: `a`, opts: WriteOptions{Quoted: true}, output: `a`},
		{atom: "\a\b\f\n\r\t\v\x00\\'\"`", opts: WriteOptions{Quoted: false}, output: "\a\b\f\n\r\t\v\x00\\'\"`"},
		{atom: "\a\b\f\n\r\t\v\x00\\'\"`", opts: WriteOptions{Quoted: true}, output: "'\\a\\b\\f\\n\\r\\t\\v\\x0\\\\\\\\'\"`'"},
		{atom: `,`, opts: WriteOptions{Quoted: false}, output: `,`},
		{atom: `,`, opts: WriteOptions{Quoted: true}, output: `','`},
		{atom: `[]`, opts: WriteOptions{Quoted: false}, output: `[]`},
		{atom: `[]`, opts: WriteOptions{Quoted: true}, output: `[]`},
		{atom: `{}`, opts: WriteOptions{Quoted: false}, output: `{}`},
		{atom: `{}`, opts: WriteOptions{Quoted: true}, output: `{}`},

		{atom: `-`, output: `-`},
		{atom: `-`, opts: WriteOptions{ops: operators{"+": {}, "-": {}}, left: operator{specifier: operatorSpecifierFY, name: "+"}}, output: ` (-)`},
		{atom: `-`, opts: WriteOptions{ops: operators{"+": {}, "-": {}}, right: operator{name: "+"}}, output: `(-)`},

		{atom: `X`, opts: WriteOptions{Quoted: true, left: operator{name: `F`}}, output: ` 'X'`},  // So that it won't be 'F''X'.
		{atom: `X`, opts: WriteOptions{Quoted: true, right: operator{name: `F`}}, output: `'X' `}, // So that it won't be 'X''F'.

		{atom: `foo`, opts: WriteOptions{left: operator{name: `bar`}}, output: ` foo`},  // So that it won't be barfoo.
		{atom: `foo`, opts: WriteOptions{right: operator{name: `bar`}}, output: `foo `}, // So that it won't be foobar.
	}

	var buf bytes.Buffer
	for _, tt := range tests {
		t.Run(string(tt.atom), func(t *testing.T) {
			buf.Reset()
			assert.NoError(t, tt.atom.WriteTerm(&buf, &tt.opts, nil))
			assert.Equal(t, tt.output, buf.String())
		})
	}
}

func TestAtom_Compare(t *testing.T) {
	var m mockTerm
	defer m.AssertExpectations(t)

	assert.Equal(t, int64(-1), Atom("a").Compare(&m, nil))
	assert.Equal(t, int64(-1), Atom("a").Compare(Atom("b"), nil))
	assert.Equal(t, int64(0), Atom("a").Compare(Atom("a"), nil))
	assert.Equal(t, int64(1), Atom("b").Compare(Atom("a"), nil))
	assert.Equal(t, int64(1), Atom("a").Compare(Variable("X"), nil))
	assert.Equal(t, int64(1), Atom("a").Compare(Float(0), nil))
	assert.Equal(t, int64(1), Atom("a").Compare(Integer(0), nil))
}
