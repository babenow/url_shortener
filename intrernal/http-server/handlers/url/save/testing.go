package save

import (
	"fmt"
	"os"
)

func ConfigureConfigPath() error {
	if _, ok := os.LookupEnv("CONFIG_PATH"); !ok {
		dir, err := os.Getwd()
		if err != nil {
			fmt.Println(err)
		}
		if err := os.Setenv("CONFIG_PATH", dir+"/../../../../../config/local.yaml"); err != nil {

			return err
		}
	}

	return nil
}
