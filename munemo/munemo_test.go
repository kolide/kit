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
		{-110000, "xadibaba"},
		{111000, "didaba"},
		{112674, "dihisho"},
		{0, "ba"},
		{725973, "shuposhe"},
		{-1, "xabi"},
	}

	for _, tt := range tests {
		t.Run(tt.s, func(t *testing.T) {
			str := Munemo(tt.i)
			assert.Equal(t, tt.s, str)
		})
	}
}
