package munemo

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// testCase contains a single munemo test case. If the string and int
// are defined, they are expected to convert between them. If the
// string and error are defined, it represents an error state. (All
// ints are expected to be able to translate)
type testCase struct {
	i       int
	s       string
	e       string
	skipInt bool // special handle a zero conversion
}

var originalTests = []testCase{
	{s: "dibaba", i: 110000},
	{s: "xadibaba", i: -110000},
	{s: "didaba", i: 111000},
	{s: "dihisho", i: 112674},
	{s: "ba", i: 0},
	{s: "xaba", i: 0, skipInt: true},
	{s: "shuposhe", i: 725973},
	{s: "xabi", i: -1},
	{s: "babi", i: 1, skipInt: true}, // leading zero

	{s: "hello", e: "decode failed: unknown syllable llo"},
}

var munemo2Tests = []testCase{
	{i: -1, s: "aabe"},
	{i: -100, s: "aabeba"},
	{i: -32, s: "aabaji", skipInt: true}, // leading zero
	{i: -99, s: "aazu"},
	{i: 0, s: "ba"},
	{i: 1, s: "be"},
	{i: 100, s: "beba"},
	{i: 101, s: "bebe"},
	{i: 25437225, s: "halotiha"},
	{i: 33, s: "bajo", skipInt: true}, // leading zero
	{i: 392406, s: "kuguce"},
	{i: 73543569, s: "tonukasu"},
	{i: 936710, s: "yosida"},
	{i: 99, s: "zu"},

	{s: "hello", e: "decode failed: unknown syllable llo"},
	{s: "qabixabi", e: "decode failed: unknown syllable qabixabi"},
}

func TestMunemoMunemo2(t *testing.T) {
	t.Parallel()
	mg := New()
	testMunemo(t, mg, munemo2Tests)
}

func TestMunemoOriginal(t *testing.T) {
	t.Parallel()
	mg := New(WithDialect(Original))
	testMunemo(t, mg, originalTests)
}

func testMunemo(t *testing.T, mg *munemoGenerator, tests []testCase) {
	for _, tt := range tests {
		tt := tt
		if tt.e == "" {
			// If we lack an error, this is a legit conversion. Try both ways
			if !tt.skipInt {
				t.Run(fmt.Sprintf("string/%d", tt.i), func(t *testing.T) {
					t.Parallel()

					ret := mg.String(tt.i)
					require.Equal(t, tt.s, ret)
				})
			}

			t.Run(fmt.Sprintf("int/%s", tt.s), func(t *testing.T) {
				t.Parallel()

				ret, err := mg.Int(tt.s)
				assert.Equal(t, tt.i, ret)
				assert.NoError(t, err)
			})
		} else {
			// Having an error, means we're looking for an error
			t.Run(fmt.Sprintf("interr/%s", tt.s), func(t *testing.T) {
				t.Parallel()

				ret, err := mg.Int(tt.s)
				require.Equal(t, tt.i, ret)
				require.EqualError(t, err, tt.e)
			})
		}

	}

}

func TestLegacyInterfaces(t *testing.T) {
	t.Parallel()

	for _, tt := range originalTests {
		tt := tt
		if tt.e == "" {
			// If we lack an error, this is a legit conversion. Try both ways
			if !tt.skipInt {
				t.Run(fmt.Sprintf("Munemo/%d", tt.i), func(t *testing.T) {
					t.Parallel()

					ret := Munemo(tt.i)
					require.Equal(t, tt.s, ret)
				})
			}

			t.Run(fmt.Sprintf("UnMunemo/%s", tt.s), func(t *testing.T) {
				t.Parallel()

				ret, err := UnMunemo(tt.s)
				assert.Equal(t, tt.i, ret)
				assert.NoError(t, err)
			})
		} else {
			// Having an error, means we're looking for an error
			t.Run(fmt.Sprintf("UnMunemo/%s", tt.s), func(t *testing.T) {
				t.Parallel()

				ret, err := UnMunemo(tt.s)
				require.Equal(t, tt.i, ret)
				require.EqualError(t, err, tt.e)
			})
		}
	}
}
