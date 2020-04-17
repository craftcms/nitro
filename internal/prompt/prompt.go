package prompt

import "github.com/manifoldco/promptui"

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
