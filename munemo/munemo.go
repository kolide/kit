// Munemo is based off of the ruby library https://github.com/jmettraux/munemo.
// It provides a deterministic way to generate a string from a number.
// it uses 100 syllables, and chunks numbers appropriately
package munemo

import (
	"bytes"
	"fmt"

	"github.com/pkg/errors"
)

// Munemo is based off of the ruby library https://github.com/jmettraux/munemo.
// It provides a deterministic way to generate a string from a number.
func Munemo(id int) string {
	m := newMunemo()
	m.calculate(id)
	return m.string()
}

func UnMunemo(s string) (int, error) {
	m := newMunemo()
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

func newMunemo() *munemo {
	m := &munemo{
		symbols: []string{
			"ba", "bi", "bu", "be", "bo",
			"cha", "chi", "chu", "che", "cho",
			"da", "di", "du", "de", "do",
			"fa", "fi", "fu", "fe", "fo",
			"ga", "gi", "gu", "ge", "go",
			"ha", "hi", "hu", "he", "ho",
			"ja", "ji", "ju", "je", "jo",
			"ka", "ki", "ku", "ke", "ko",
			"la", "li", "lu", "le", "lo",
			"ma", "mi", "mu", "me", "mo",
			"na", "ni", "nu", "ne", "no",
			"pa", "pi", "pu", "pe", "po",
			"ra", "ri", "ru", "re", "ro",
			"sa", "si", "su", "se", "so",
			"sha", "shi", "shu", "she", "sho",
			"ta", "ti", "tu", "te", "to",
			"tsa", "tsi", "tsu", "tse", "tso",
			"wa", "wi", "wu", "we", "wo",
			"ya", "yi", "yu", "ye", "yo",
			"za", "zi", "zu", "ze", "zo",
		},
		sign:           1,
		symbolValues:   make(map[string]int),
		negativeSymbol: "xa",
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
		s = s[2:len(s)]
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
			s = s[2:len(s)]
		} else if val, ok := m.symbolValues[s[0:3]]; ok {
			m.number = len(m.symbols)*m.number + val
			s = s[3:len(s)]
		} else {
			m.number = 0
			return errors.New(fmt.Sprintf("decode failed: unknown syllable %s", s))
		}
	}
	// No errors!
	return nil
}

func (m *munemo) calculate(number int) {
	if number < 0 {
		m.buffer.Write([]byte(m.negativeSymbol))
		number = number * -1
	}

	modulo := number % len(m.symbols)
	result := number / len(m.symbols)

	if result > 0 {
		m.calculate(result)
	}

	m.buffer.Write([]byte(m.symbols[modulo]))
}
