package prompt

import (
	"github.com/fgognet/prompt/internal/prompt"
	"github.com/go-playground/validator/v10"
)

func Prompt[T any](validate *validator.Validate, data T) (T, error) {
	return prompt.InternalPrompt(validate, data)
}
