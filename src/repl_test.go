package main

import "testing"

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "  hello  world  ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "  he llo  worl d  ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "  Charmander Bulbasaur PIKACHU  ",
			expected: []string{"charmander", "bulbasaur", "pikachu"},
		},
		// add more cases here
	}

	for _, c := range cases {
		actual := cleanInput(c.input)
		if len(actual) != len(c.expected) {
			t.Errorf("Test Failed - expected %d words, got %d", len(c.expected), len(actual))
		}

		// Проверьте длину фактического фрагмента
		// если они не совпадают, используйте t.Errorf для вывода сообщения об ошибке
		// и не пройдете тест
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			if word != expectedWord {
				t.Errorf("Test Failed - expected %q, got %q", expectedWord, word)
			}
			// Проверьте каждое слово в фрагменте
			// если они не совпадают, используйте t.Errorf для вывода сообщения об ошибке
			// и не пройдете тест
		}
	}
}
