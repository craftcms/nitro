package prompt

import (
	"github.com/manifoldco/promptui"
)

func Ask(label, def string, validator promptui.ValidateFunc) (string, error) {
	p := promptui.Prompt{
		Label:    label,
		Default:  def,
		Validate: validator,
	}

	v, err := p.Run()
	if err != nil {
		return "", err
	}

	return v, nil
}

func AskWithDefault(label, def string, validator promptui.ValidateFunc) (string, error) {
	p := promptui.Prompt{
		Label:    label + " [" + def + "]",
		Validate: validator,
	}

	v, err := p.Run()
	if err != nil {
		return "", err
	}

	switch v {
	case "":
		v = def
	}

	return v, nil
}

func SelectWithDefault(label, def string, options []string) (int, string) {
	p := promptui.Select{
		Label: label + " [" + def + "]",
		Items: options,
	}

	i, selected, _ := p.Run()

	return i, selected
}

func Select(label string, options []string) (int, string) {
	p := promptui.Select{
		Label: label,
		Items: options,
	}

	i, selected, _ := p.Run()

	return i, selected
}

func Verify(label string) bool {
	verify := promptui.Prompt{
		Label: label,
	}

	answer, err := verify.Run()
	if err != nil {
		return false
	}

	if answer == "" {
		return true
	}

	return false
}
