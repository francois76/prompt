package prompt_test

import (
	"fmt"
	"testing"

	public "github.com/fgognet/prompt"
	"github.com/fgognet/prompt/internal/prompt"
	"github.com/fgognet/prompt/internal/testutils"
	"github.com/go-playground/validator/v10"
	"github.com/maxatome/go-testdeep/td"
)

func TestPrompt(t *testing.T) {
	oldScanLn := prompt.ScanLn
	oldPrintF := prompt.Printf
	defer func() {
		prompt.ScanLn = oldScanLn
		prompt.Printf = oldPrintF
	}()

	{
		type hobby struct {
			Name string `validate:"oneof=fencing football"`
		}

		type example struct {
			MainHobby hobby
			Name      string  `validate:"required"`
			Age       int     `validate:"required,gte=18" prompt:"What is your current age?"`
			FirstName string  `validate:"required"`
			Hobbies   []hobby `validate:"required,dive"`
		}
		t.Run("success_full", testInternalPromptSuccess(testInternalPromptSuccessParams[example]{
			// in this test, we try a wront type and we expect to be prompted again
			input: example{
				MainHobby: hobby{},
				Name:      "hello",
				Hobbies: []hobby{{
					Name: "invalid",
				}, {
					Name: "invalid",
				}}},
			expectedStructure: example{
				MainHobby: hobby{Name: "fencing"},
				Name:      "hello",
				Age:       18,
				FirstName: "francois",
				Hobbies: []hobby{{
					Name: "fencing",
				}, {
					Name: "football",
				}}},
			prompt: []string{"fencing", "5", "francois", "fencing", "football", "18"},
			expectedOutput: []string{
				"Validation error :\n        => Key: 'example.MainHobby.Name' Error:Field validation for 'Name' failed on the 'oneof' tag\n",
				"please input a value for the field Name : ",
				"Validation error :\n        => Key: 'example.Age' Error:Field validation for 'Age' failed on the 'required' tag\n",
				"What is your current age?",
				"Validation error :\n        => Key: 'example.FirstName' Error:Field validation for 'FirstName' failed on the 'required' tag\n",
				"please input a value for the field FirstName : ",
				"Validation error :\n        => Key: 'example.Hobbies[0].Name' Error:Field validation for 'Name' failed on the 'oneof' tag\n",
				"please input a value for the field Name : ",
				"Validation error :\n        => Key: 'example.Hobbies[1].Name' Error:Field validation for 'Name' failed on the 'oneof' tag\n",
				"please input a value for the field Name : ",
				"Validation error :\n        => Key: 'example.Age' Error:Field validation for 'Age' failed on the 'gte' tag\n",
				"What is your current age?",
			},
		}))
	}

	{
		type typ struct {
			ValueInt  int  `validate:"required"`
			ValueBool bool `validate:"required"`
		}
		t.Run("success_wrong_type", testInternalPromptSuccess(testInternalPromptSuccessParams[typ]{
			// in this test, we try a wront type and we expect to be prompted again
			input: typ{},
			expectedStructure: typ{
				ValueInt:  18,
				ValueBool: true,
			},
			prompt: []string{"not_a_number", "not_a_boolean", "18", "true"},
			expectedOutput: []string{
				"Validation error :\n        => Key: 'typ.ValueInt' Error:Field validation for 'ValueInt' failed on the 'required' tag\n",
				"please input a value for the field ValueInt : ",
				"Validation error :\n        => Key: 'typ.ValueBool' Error:Field validation for 'ValueBool' failed on the 'required' tag\n",
				"please input a value for the field ValueBool : ",
				"Validation error :\n        => Key: 'typ.ValueInt' Error:Field validation for 'ValueInt' failed on the 'required' tag\n",
				"please input a value for the field ValueInt : ",
				"Validation error :\n        => Key: 'typ.ValueBool' Error:Field validation for 'ValueBool' failed on the 'required' tag\n",
				"please input a value for the field ValueBool : ",
			},
		}))
	}

	{
		type typ struct {
			Value string `validate:"required" prompt:"this question shouldn't be asked"`
			Inner struct {
				Value string `validate:"required" prompt:"this question SHOULD be asked"`
			}
		}
		t.Run("success_duplicate_field", testInternalPromptSuccess(testInternalPromptSuccessParams[typ]{
			// this test check that the correct prompt is displayed
			input: typ{
				Value: "something",
			},
			expectedStructure: typ{
				Value: "something",
				Inner: struct {
					Value string `validate:"required" prompt:"this question SHOULD be asked"`
				}{
					Value: "result",
				},
			},
			prompt: []string{"result"},
			expectedOutput: []string{
				"Validation error :\n        => Key: 'typ.Inner.Value' Error:Field validation for 'Value' failed on the 'required' tag\n",
				"this question SHOULD be asked",
			},
		}))
	}

	{
		type typ struct {
		}
		t.Run("success_data_already_valid", testInternalPromptSuccess(testInternalPromptSuccessParams[typ]{
			input:             typ{},
			expectedStructure: typ{},
			prompt:            []string{},
			expectedOutput:    []string{},
		}))
	}

}

type testInternalPromptSuccessParams[T any] struct {
	input             T
	expectedStructure T
	prompt            []string
	expectedOutput    []string
}

func testInternalPromptSuccess[T any](params testInternalPromptSuccessParams[T]) func(t *testing.T) {
	return func(t *testing.T) {
		prompt.ScanLn = testutils.Prompts(params.prompt...)
		output := []string{}
		prompt.Printf = func(format string, a ...any) (int, error) {
			output = append(output, fmt.Sprintf(format, a...))
			return 0, nil
		}
		validate := validator.New(validator.WithRequiredStructEnabled())
		got, err := public.Prompt(validate, params.input)
		td.CmpNil(t, err)
		td.Cmp(t, got, params.expectedStructure)
		td.Cmp(t, output, params.expectedOutput)
	}
}
