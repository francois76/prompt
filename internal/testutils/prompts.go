package testutils

import "fmt"

// Prompts Utility function to simulate prompts
func Prompts(data ...string) func(a ...any) (n int, err error) {
	counter := 0
	return func(a ...any) (n int, err error) {
		if len(a) == 0 {
			return 0, nil
		}
		if counter+len(a) > len(data) {
			return 0, fmt.Errorf("not enough data passed as input, either some input is missing in the unit test or something is wrong with the code")
		}
		for _, d := range a {
			p := d.(*string)

			*p = data[counter]
			counter++
		}

		return 1, nil
	}
}
