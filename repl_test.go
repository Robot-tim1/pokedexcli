package main

import (
	"reflect"
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected []string
	}{
		{name: "whitespace", input: "  hello  world  ", expected: []string{"hello", "world"}},
		{name: "words together", input: "fromboot.dev we learn ", expected: []string{"fromboot.dev", "we", "learn"}},
		{name: "noise", input: "/..ie vowwo vn .", expected: []string{"/..ie", "vowwo", "vn", "."}},
	}

	for i, c := range cases {
		actual := cleanInput(c.input)
		if !reflect.DeepEqual(c.expected, actual) {
			t.Errorf("Test %d, %s: Expected: %v Got: %v", i+1, c.name, c.expected, actual)
		}
	}
}
