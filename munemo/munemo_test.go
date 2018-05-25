package munemo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMunemo(t *testing.T) {

	var tests = []struct {
		i int
		s string
	}{
		{110000, "dibaba"},
		{111000, "didaba"},
		{112674, "dihisho"},
	}

	for _, tt := range tests {
		t.Run(tt.s, func(t *testing.T) {
			str := Munemo(tt.i)
			assert.Equal(t, tt.s, str)
		})
	}
}
