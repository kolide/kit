// Munemo is a reversible numeric encoding library. It's designed to
// take id numbers, and present them in more human friendly forms. It
// does this by using a set of symbols, and doing a base conversion to
// them. There are a couple of known dialects.
//
// Original: This is compatible with the original symbol set. This has
// the disadvantage of being variable length, and non-sortable.
//
// Munemo2: This symbol set was developed as a replacement. All
// symbols are 2 characters, and it creates sortable strings.
//
// It is inspired by the ruby library
// https://github.com/jmettraux/munemo.
package munemo

import (
	"bytes"
	"fmt"
)

type munemoGenerator struct {
	dialect dialect
}

// Option is the functional option type for munemoGenerator
type Option func(*munemoGenerator)

// WithDialect defines the dialect to be used
func WithDialect(d dialect) Option {
	return func(mg *munemoGenerator) {
		mg.dialect = d
	}
}

// New returns a munemo generator for the given dialect.
func New(opts ...Option) *munemoGenerator {
	mg := &munemoGenerator{
		dialect: Munemo2,
	}
	for _, opt := range opts {
		opt(mg)
	}

	return mg
}

// String takes an integer and returns the mumemo encoded string
func (mg *munemoGenerator) String(id int) string {
	m := newMunemo(mg.dialect)
	m.calculate(id)
	return m.string()
}

// Int takes a string, and returns an integer. In the case of error,
// an error is returned.
func (mg *munemoGenerator) Int(s string) (int, error) {
	m := newMunemo(mg.dialect)
	err := m.decode(s)
	return m.int(), err
}

// Munemo is a legacy interface to munemo encoding. It defaults to the
// original dialect
func Munemo(id int) string {
	m := newMunemo(Original)
	m.calculate(id)
	return m.string()
}

// UnMunemo is a legacy interface to reverse munemo encoding. It
// defaults to the original dialect.
func UnMunemo(s string) (int, error) {
	m := newMunemo(Original)
	err := m.decode(s)
	return m.int(), err
}

type munemo struct {
	negativeSymbol string
	symbols        []string
	buffer         *bytes.Buffer
	number         int
	symbolValues   map[string]int
	sign           int
}

func newMunemo(d dialect) *munemo {
	m := &munemo{
		symbols:        d.symbols,
		negativeSymbol: d.negativeSymbol,
		sign:           1,
		symbolValues:   make(map[string]int),
		buffer:         new(bytes.Buffer),
	}
	for k, v := range m.symbols {
		m.symbolValues[v] = k
	}

	return m
}

func (m *munemo) string() string {
	return m.buffer.String()
}

func (m *munemo) int() int {
	return m.number * m.sign
}

func (m *munemo) decode(s string) error {
	// negative if the first two bytes match the negative symbol
	if s[0:2] == m.negativeSymbol {
		m.sign = -1
		s = s[2:]
	}

	// As long as there are characters, parse them
	// Read the first syllable, interpret, remove.
	for {
		if s == "" {
			break
		}

		// Syllables are 2 or 3 letters. Check to see if the first 2 or 3
		// characters are in our array of syllables.
		if val, ok := m.symbolValues[s[0:2]]; ok {
			m.number = len(m.symbols)*m.number + val
			s = s[2:]
		} else if val, ok := m.symbolValues[s[0:3]]; ok {
			m.number = len(m.symbols)*m.number + val
			s = s[3:]
		} else {
			m.number = 0
			return fmt.Errorf("decode failed: unknown syllable %s", s)
		}
	}
	// No errors!
	return nil
}

func (m *munemo) calculate(number int) {
	if number < 0 {
		m.buffer.Write([]byte(m.negativeSymbol))
		number = -number
	}

	modulo := number % len(m.symbols)
	result := number / len(m.symbols)

	if result > 0 {
		m.calculate(result)
	}

	m.buffer.Write([]byte(m.symbols[modulo]))
}
