package munemo

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMunemo(t *testing.T) {

	var tests = []struct {
		i int
		s string
	}{
		{110000, "dibaba"},
		{-110000, "xadibaba"},
		{111000, "didaba"},
		{112674, "dihisho"},
		{0, "ba"},
		{725973, "shuposhe"},
		{-1, "xabi"},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("Munemo/%d", tt.i), func(t *testing.T) {
			ret := Munemo(tt.i)
			assert.Equal(t, tt.s, ret)
		})

		t.Run(fmt.Sprintf("UnMunemo/%s", tt.s), func(t *testing.T) {
			ret := UnMunemo(tt.s)
			assert.Equal(t, tt.i, ret)
		})
	}
}
