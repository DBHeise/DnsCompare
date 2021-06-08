package DnsCompare

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandInt(t *testing.T) {
	testCases := []struct {
		name            string
		count, min, max int
	}{
		{"Simple", 100, 0, 10},
		{"Normal", 100, 1, 5},
		{"Small", 100, 42, 43},
		{"Long", 100, 100, 1000000},
		{"Negative", 100, -20, -10},
		//{"MinMaxReversed", 100, 20, 10},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(tt *testing.T) {
			for i := 0; i < tc.count; i++ {
				tt.Run("Iteration "+strconv.Itoa(i), func(ttt *testing.T) {
					actual := RandomInt(tc.min, tc.max)
					assert.GreaterOrEqual(ttt, actual, tc.min)
					assert.LessOrEqual(ttt, actual, tc.max)
				})
			}
		})
	}
}

func TestRandomFromStringSlice(t *testing.T) {
	testCases := []struct {
		name       string
		inputSlice []string
	}{
		{"Simple", []string{"one", "two", "three"}},
		{"Singular", []string{"one"}},
		//{"Empty", []string{}},
		{"Five", []string{"one", "two", "three", "four", "five"}},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(tt *testing.T) {
			actual := RandomFromStringSlice(tc.inputSlice)
			assert.NotEmpty(tt, actual)
			assert.Contains(tt, tc.inputSlice, actual)
		})
	}
}
