package validate

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Validator interface {
	Validate(input string) error
}

// DatabaseName is used to validate a database names
type DatabaseName struct{}

func (v *DatabaseName) Validate(input string) error {
	// check length
	if len(input) < 3 {
		return fmt.Errorf("database must be more than 3 characters")
	}

	// check for spaces
	if strings.Contains(input, " ") {
		return fmt.Errorf("database must not include spaces")
	}

	// check for special characters
	if strings.ContainsAny(input, "!@#$%^&*()-.") {
		return fmt.Errorf("database must not include any special characters except underscores")
	}

	return nil
}

// HostnameValidator is used to validate a provided hostname
type HostnameValidator struct{}

func (v *HostnameValidator) Validate(input string) error {
	// check length
	if len(input) < 3 {
		return fmt.Errorf("hostname must be more than 3 characters")
	}

	// check for spaces
	if strings.Contains(input, " ") {
		return fmt.Errorf("hostname must not include spaces")
	}

	// check for special characters
	if strings.ContainsAny(input, "!@#$%^&*(),") {
		return fmt.Errorf("hostname must not include any special characters")
	}

	return nil
}

// IntegerValidator validates if the input is a valid integer
type IntegerValidator struct{}

func (v *IntegerValidator) Validate(input string) error {
	if _, err := strconv.Atoi(input); err != nil {
		return err
	}

	return nil
}

// MultipleHostnameValidator validates a comma separated list of hostnames
type MultipleHostnameValidator struct{}

func (v *MultipleHostnameValidator) Validate(input string) error {
	_, err := v.Parse(input)
	return err
}

func (v *MultipleHostnameValidator) Parse(input string) ([]string, error) {
	rawHosts := strings.Split(input, ",")
	hostV := &HostnameValidator{}
	var hosts []string

	for _, h := range rawHosts {
		h := strings.TrimSpace(h)
		if err := hostV.Validate(h); err != nil {
			return nil, err
		}
		hosts = append(hosts, h)
	}
	return hosts, nil
}

type PHPVersionValidator struct{}

func (v *PHPVersionValidator) Validate(input string) error {
	switch input {
	case "8.1", "8.0", "7.4", "7.3", "7.2", "7.1", "7.0":
		return nil
	}

	return fmt.Errorf("the PHP version %q is not valid", input)
}

type IsBoolean struct{}

func (v *IsBoolean) Validate(input string) error {
	if _, err := strconv.ParseBool(input); err != nil {
		return err
	}

	return nil
}

type IsMegabyte struct{}

func (v *IsMegabyte) Validate(input string) error {
	return isMegabytes(input)
}

type MaxExecutionTime struct{}

func (v *MaxExecutionTime) Validate(input string) error {
	return maxExecutionTime(input)
}

func maxExecutionTime(v string) error {
	_, err := strconv.Atoi(v)
	if err != nil {
		return errors.New("max_execution_time must be a valid integer")
	}

	return nil
}

func MaxInputVars(v string) error {
	num, err := strconv.Atoi(v)
	if err != nil {
		return errors.New("max_input_vars must be a valid integer")
	}

	if num >= 10000 {
		return errors.New("max_input_vars must be less than 10000")
	}

	return nil
}

func isMegabytes(v string) error {
	if len(v) == 1 {
		return errors.New("memory must be larger than 1 character (e.g. 256M)")
	}

	if !strings.HasSuffix(v, "M") {
		return errors.New("memory must end with a M")
	}

	return nil
}

func PhpMaxFileUploads(v string) error {
	num, err := strconv.Atoi(v)
	if err != nil {
		return errors.New("max_input_vars must be a valid integer")
	}

	if num >= 500 {
		return errors.New("max_file_uploads must be less than 500")
	}

	return nil
}
