package DnsCompare

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func testIsIPv4Func(input string, expected bool) func(*testing.T) {
	return func(t *testing.T) {
		actual := IsIPv4(input)
		assert.Equal(t, expected, actual, input)
	}
}
func TestIsIPv4(t *testing.T) {
	t.Run("Empty", testIsIPv4Func("", false))
	t.Run("Valid", testIsIPv4Func("1.2.3.4", true))
	t.Run("Short", testIsIPv4Func("1.2.3", false))
	t.Run("Long", testIsIPv4Func("1.2.3.4.5", false))
	t.Run("WithPort", testIsIPv4Func("1.2.3.4:80", false))
	t.Run("String", testIsIPv4Func("foobar", false))
}

func TestIsIPv6(t *testing.T) {
	testCases := []struct {
		name, input string
		expected    bool
	}{
		{"Empty", "", false},
		{"Valid", "1:2:3:4:5:6:7:8", true},
		{"Short", "1:2:3:4:5:6:7", false},
		{"Long", "1:2:3:4:5:6:7:8:9", false},
		{"WithPort", "[1:2:3:4:5:6:7:8]:80", false},
		{"String", "foobar", false},
		{"v6 Abbreviation", "::1", true},
		{"v6 Full", "0001:0002:0003:0004:0005:0006:0007:0008", true},
		{"v6 EmbeddedIPv4", "ff::1.2.3.4", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(tt *testing.T) {
			actual := IsIPv6(tc.input)
			assert.Equal(t, tc.expected, actual, tc.input)
		})
	}
}

func TestChooseServer(t *testing.T) {
	t.Run("Methods", func(tt *testing.T) {
		testList := []string{"1.1.1.1", "1:2:3:4:5:6:7:8", "2.2.2.2", "8:7:6:5:4:3:2:1", "foobar"}

		tt.Run("First", func(ttt *testing.T) {
			items := chooseServerAddress("first", testList)
			assert.Equal(ttt, 1, len(items), "returned list size")
			assert.Equal(ttt, "1.1.1.1", items[0])
		})
		tt.Run("First4", func(ttt *testing.T) {
			items := chooseServerAddress("first4", testList)
			assert.Equal(ttt, 1, len(items), "returned list size")
			assert.Equal(ttt, "1.1.1.1", items[0])
		})
		tt.Run("First4_No4", func(ttt *testing.T) {
			items := chooseServerAddress("first4", []string{"abc", "def", "hij"})
			assert.Equal(ttt, 0, len(items), "returned list size")
		})
		tt.Run("First6", func(ttt *testing.T) {
			items := chooseServerAddress("first6", testList)
			assert.Equal(ttt, 1, len(items), "returned list size")
			assert.Equal(ttt, "1:2:3:4:5:6:7:8", items[0])
		})
		tt.Run("All", func(ttt *testing.T) {
			items := chooseServerAddress("all", testList)
			assert.Equal(ttt, 5, len(items), "returned list size")
		})
		tt.Run("All4", func(ttt *testing.T) {
			items := chooseServerAddress("all4", testList)
			assert.Equal(ttt, 2, len(items), "returned list size")
			for _, item := range items {
				assert.True(ttt, IsIPv4(item))
			}
		})
		tt.Run("All6", func(ttt *testing.T) {
			items := chooseServerAddress("all6", testList)
			assert.Equal(ttt, 2, len(items), "returned list size")
			for _, item := range items {
				assert.True(ttt, IsIPv6(item))
			}
		})
		tt.Run("Random", func(ttt *testing.T) {
			items := chooseServerAddress("random", testList)
			assert.Equal(ttt, 1, len(items), "returned list size")
		})
		tt.Run("Random4", func(ttt *testing.T) {
			items := chooseServerAddress("random4", testList)
			assert.Equal(ttt, 1, len(items), "returned list size")
			assert.True(ttt, IsIPv4(items[0]))
		})
		tt.Run("Random4_no4", func(ttt *testing.T) {
			items := chooseServerAddress("random4", []string{"one", "two", "three"})
			assert.Equal(ttt, 0, len(items), "returned list size")
		})
		tt.Run("Random6", func(ttt *testing.T) {
			items := chooseServerAddress("random6", testList)
			assert.Equal(ttt, 1, len(items), "returned list size")
			assert.True(ttt, IsIPv6(items[0]))
		})
	})
}
