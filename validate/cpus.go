package validate

import (
	"errors"
	"strconv"
)

// ValidCPUCount is used to contain the host CPU count and
// provide a receiver to validate the CPU count.
type ValidCPUCount struct {
	Actual int
}

// Validate performs the work to validate the number of
// requested CPUs against the hosts number of CPUs.
func (c ValidCPUCount) Validate(requested string) error {
	n, err := strconv.Atoi(requested)
	if err != nil {
		return err
	}

	// if the cpu count is set
	if c.Actual == 0 {
		return errors.New("unable to determine the host machines CPU count")
	}

	// if its a match or higher than
	if n >= c.Actual {
		return errors.New("the number of CPUs cannot match or exceed the host CPU count")
	}

	return nil
}

// NewCPUValidator is used to setup a validator struct
// that reduces the code setup.
func NewCPUValidator(actual int) ValidCPUCount {
	return ValidCPUCount{Actual: actual}
}
