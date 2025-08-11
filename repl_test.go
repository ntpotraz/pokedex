package main
import "testing"

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input string
		command string
		args []string
	}{
		{
			input: "  hello  world  ",
			command: "hello",
			args: []string{"world"},
		},
		{
			input: "My NAME    is NaThAN",
			command:  "my",
			args: []string{"name", "is", "nathan"},
		},
		{
			input: "ARTEMIS",
			command: "artemis",
			args: []string{},
		},
		{
			input: "   ",
			command: "",
			args: []string{},
		},
	}

	for _, c := range cases {
		command, args := cleanInput(c.input)
		if command != c.command || len(args) != len(c.args) {
			t.Errorf("actual length doesn't equal expected:")
			t.Errorf("Actual %v: %d | Expected %v: %d", command, len(command), c.args, len(c.args))
			t.Fail()
		}
		for i := range args {
			actualWord := args[i]
			expectedWord := c.args[i]
			if actualWord != expectedWord {
				t.Errorf("'%s' does not match expected: '%s'", actualWord, expectedWord)
				t.Fail()
			}
		}
	}

}
