package prompt

import (
	"errors"
	"strings"

	"github.com/pixelandtonic/go-input"
)

func Ask(ui *input.UI, query, def string, req bool) (string, error) {
	a, err := ui.Ask(query, &input.Options{
		Default:  def,
		Required: req,
		Loop:     true,
	})
	if err != nil {
		return "", err
	}

	return a, nil
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
