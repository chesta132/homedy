package config

import (
	"bufio"
	"fmt"
	"homedy/flags"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

// ---- DO NOT CHANGE ----

func IsEnvProd() bool {
	return GO_ENV == "production"
}

func IsEnvDev() bool {
	return GO_ENV == "development"
}

func SplitByEnv[T any](prodValue, devValue T) T {
	if IsEnvProd() {
		return prodValue
	} else {
		return devValue
	}
}

func init() {
	envPath := *flags.EnvPath
	if envPath == "" {
		// cmdlib.Input
		fmt.Print("env path: ")
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		envPath = strings.TrimSpace(input)
	}

	err := godotenv.Load(envPath)
	if err != nil {
		panic(err)
	}
	ReloadEnv()
	if IsEnvDev() || IsEnvProd() {
		return
	}
	panic("[ENVIRONMENT] invalid environment, must be 'development' or 'production'")
}
