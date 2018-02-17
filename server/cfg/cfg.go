package cfg

import (
	"fmt"
	"os"
)

func GetRoot() string {
	env := os.Getenv("GO_ENV")
	if env == "" {
		env = "dev"
	}

	return fmt.Sprintf("./config/%s/", env)
}
