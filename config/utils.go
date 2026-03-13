package config

import (
	"bufio"
	"flag"
	"fmt"
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
	envPathPtr := flag.String("env", "", "env path")
	flag.Parse()

	envPath := *envPathPtr
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
