package prompt

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/fgognet/prompt/constants"
	"github.com/fgognet/prompt/internal/reflection"
	"github.com/go-playground/validator/v10"
)

var ScanLn = fmt.Scanln
var Printf = fmt.Printf

func InternalPrompt[T any](validate *validator.Validate, data T) (T, error) {
	vErr := validate.Struct(data)
	for vErr != nil && errors.As(vErr, &validator.ValidationErrors{}) {
		err := handlingValidationError(&data, vErr)
		if err != nil {
			return data, err
		}
		vErr = validate.Struct(data)
	}
	return data, vErr
}

func handlingValidationError[T any](data *T, vErr error) error {
	var vErrs validator.ValidationErrors
	if errors.As(vErr, &vErrs) {
		ref := reflect.ValueOf(&data).Elem()
		for _, v := range vErrs {
			// printing the error
			_, err := Printf(constants.DefaultErrorMessage, v.Error())
			if err != nil {
				return err
			}

			// assigning the prompt value
			err = reflection.Assign(ref, v.StructNamespace(), func(fieldStruct reflect.StructField) (string, error) {
				// asking the user what to do
				_, err := Printf(getPromptSentence(fieldStruct))
				if err != nil {
					return "", err
				}
				// reading input
				var input string
				if _, err := ScanLn(&input); err != nil {
					return "", err
				}
				return input, nil
			})

			if err != nil {
				return err
			}
		}
	}
	return nil
}

func getPromptSentence(fieldStruct reflect.StructField) string {
	promptSentence := fieldStruct.Tag.Get(constants.PromptTag)
	if promptSentence == "" {
		promptSentence = fmt.Sprintf(constants.DefaultPromptSentence, fieldStruct.Name)
	}
	return promptSentence
}
