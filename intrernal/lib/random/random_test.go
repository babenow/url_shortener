package random_test

import (
	"testing"

	"github.com/babenow/url_shortener/intrernal/lib/random"
	"github.com/stretchr/testify/assert"
)

func TestRandom_NewRandomString(t *testing.T) {
	testCases := []struct {
		desc string
		size int
	}{
		{desc: "size=1", size: 1},
		{desc: "size=2", size: 2},
		{desc: "size=3", size: 3},
		{desc: "size=4", size: 4},
		{desc: "size=5", size: 5},
		{desc: "size=6", size: 6},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			str := random.NewRandomString(tC.size)
			assert.Equal(t, len(str), tC.size)
		})
	}
}
