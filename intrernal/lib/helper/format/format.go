package format

import "fmt"

func Err(operator string, err error) error {
	return fmt.Errorf("%s: %w", operator, err)
}
