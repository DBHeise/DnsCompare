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

func TestChooseServer_Method_First(t *testing.T) {
	testList := []string{"1.1.1.1", "1:2:3:4:5:6:7:8", "2.2.2.2", "8:7:6:5:4:3:2:1", "foobar"}
	items := chooseServerAddress("first", testList)
	assert.Equal(t, 1, len(items), "returned list size")
	assert.Equal(t, "1.1.1.1", items[0])
}
func TestChooseServer_Method_First4(t *testing.T) {
	testList := []string{"1:2:3:4:5:6:7:8", "1.1.1.1", "2.2.2.2", "8:7:6:5:4:3:2:1", "foobar"}
	items := chooseServerAddress("first4", testList)
	assert.Equal(t, 1, len(items), "returned list size")
	assert.Equal(t, "1.1.1.1", items[0])
}
func TestChooseServer_Method_First6(t *testing.T) {
	testList := []string{"1.1.1.1", "1:2:3:4:5:6:7:8", "2.2.2.2", "8:7:6:5:4:3:2:1", "foobar"}
	items := chooseServerAddress("first6", testList)
	assert.Equal(t, 1, len(items), "returned list size")
	assert.Equal(t, "1:2:3:4:5:6:7:8", items[0])
}
func TestChooseServer_Method_All(t *testing.T) {
	testList := []string{"1.1.1.1", "1:2:3:4:5:6:7:8", "2.2.2.2", "8:7:6:5:4:3:2:1", "foobar"}
	items := chooseServerAddress("all", testList)
	assert.Equal(t, 5, len(items), "returned list size")
}
func TestChooseServer_Method_All4(t *testing.T) {
	testList := []string{"1.1.1.1", "1:2:3:4:5:6:7:8", "2.2.2.2", "8:7:6:5:4:3:2:1", "foobar"}
	items := chooseServerAddress("all4", testList)
	assert.Equal(t, 2, len(items), "returned list size")
	for _, item := range items {
		assert.True(t, IsIPv4(item))
	}
}
func TestChooseServer_Method_All6(t *testing.T) {
	testList := []string{"1.1.1.1", "1:2:3:4:5:6:7:8", "2.2.2.2", "8:7:6:5:4:3:2:1", "foobar"}
	items := chooseServerAddress("all6", testList)
	assert.Equal(t, 2, len(items), "returned list size")
	for _, item := range items {
		assert.True(t, IsIPv6(item))
	}
}
func TestChooseServer_Method_Random(t *testing.T) {
	testList := []string{"1.1.1.1", "1:2:3:4:5:6:7:8", "2.2.2.2", "8:7:6:5:4:3:2:1", "foobar"}
	items := chooseServerAddress("random", testList)
	assert.Equal(t, 1, len(items), "returned list size")
}
func TestChooseServer_Method_Random4(t *testing.T) {
	testList := []string{"1:2:3:4:5:6:7:8", "1.1.1.1", "2.2.2.2", "8:7:6:5:4:3:2:1", "foobar"}
	items := chooseServerAddress("random4", testList)
	assert.Equal(t, 1, len(items), "returned list size")
	assert.True(t, IsIPv4(items[0]))
}
func TestChooseServer_Method_Random6(t *testing.T) {
	testList := []string{"1:2:3:4:5:6:7:8", "1.1.1.1", "2.2.2.2", "8:7:6:5:4:3:2:1", "foobar"}
	items := chooseServerAddress("random6", testList)
	assert.Equal(t, 1, len(items), "returned list size")
	assert.True(t, IsIPv6(items[0]))
}
