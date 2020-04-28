package prompt

import (
	"errors"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/tcnksm/go-input"
)

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

// Select is responsible for providing a list of options to remove
func Select(ui *input.UI, query, def string, list []string) (string, int, error) {
	selected, err := ui.Select(query, list, &input.Options{
		Required: true,
		Default:  def,
	})
	if err != nil {
		return "", 0, err
	}

	for i, s := range list {
		if s == selected {
			return s, i, nil
		}
	}

	return "", 0, errors.New("unable to find the selected option")
}

func SelectWithDefault(label, def string, options []string) (int, string) {
	p := promptui.Select{
		Label: label + " [" + def + "]",
		Items: options,
	}

	i, selected, _ := p.Run()

	return i, selected
}

func Verify(ui *input.UI, query, def string) (bool, error) {
	a, err := ui.Ask(query, &input.Options{
		Default:  def,
		Required: true,
		Loop:     true,
	})
	if err != nil {
		return false, err
	}

	if strings.ContainsAny(a, "y") {
		return true, nil
	}

	return false, nil
}
