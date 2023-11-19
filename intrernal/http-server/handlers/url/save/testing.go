package save

import (
	"fmt"
	"os"
)

func ConfigureConfigPath() {
	if _, ok := os.LookupEnv("CONFIG_PATH"); !ok {
		dir, err := os.Getwd()
		if err != nil {
			fmt.Println(err)
		}
		os.Setenv("CONFIG_PATH", dir+"/../../../../../config/local.yaml")
	}
}
