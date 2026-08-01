package main

import "testing"

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    " hello   world 	",
			expected: []string{"hello", "world"},
		},
		{
			// Multiple internal spaces collapse into single separators
			input:    "charmander    bulbasaur   squirtle",
			expected: []string{"charmander", "bulbasaur", "squirtle"},
		},
		{
			// A single word with surrounding whitespace
			input:    "   pikachu   ",
			expected: []string{"pikachu"},
		},
		{
			// Tabs and newlines are treated as whitespace separators
			input:    "\tpidgey\ncaterpie\t weedle\n",
			expected: []string{"pidgey", "caterpie", "weedle"},
		},
		{
			// Already-clean input is returned unchanged
			input:    "one two three",
			expected: []string{"one", "two", "three"},
		},
		{
			// Uppercase input is lowercased
			input:    "Charmander Bulbasaur PIKACHU",
			expected: []string{"charmander", "bulbasaur", "pikachu"},
		},
		{
			// Mixed case with extra whitespace is both lowercased and split
			input:    "  HELLO   World 	",
			expected: []string{"hello", "world"},
		},
		{
			// An empty string yields no words
			input:    "",
			expected: []string{},
		},
		{
			// A whitespace-only string yields no words
			input:    "   \t  \n ",
			expected: []string{},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)
		if len(actual) != len(c.expected) {
			t.Errorf("Failed:\ninput: %v\nexpected output: %v\nactual output: %v\n", c.input, c.expected, actual)
			continue
		}
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			if word != expectedWord {
				t.Errorf("Failed:\ninput: %v\nexpected output: %v\nactual output: %v\nMISSMATCH:\nexpected word: %v\nactual word: %v\n", c.input, c.expected, actual, expectedWord, word)
			}
		}
	}

}
