package config

import (
	"bufio"
	"fmt"
	"homedy/flags"
	"net"
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

func getMachineRunOnUrl(env string) string {
	// set url to machine ip
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		panic(fmt.Sprintf("[ENVIRONMENT] invalid %s and failed to get %s fallback (self IP)", env, env))
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ip := ipnet.IP.String()
				if SERVER_PORT != "" {
					ip += ":" + SERVER_PORT
				}
				return "http://" + ip
			}
		}
	}
	panic(fmt.Sprintf("[ENVIRONMENT] invalid %s and failed to get %s fallback (self IP)", env, env))
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
	if !IsEnvDev() && !IsEnvProd() {
		panic("[ENVIRONMENT] invalid environment, must be 'development' or 'production'")
	}

	if APP_URL == "" {
		APP_URL = getMachineRunOnUrl("APP_URL") + "/api"
	}

	if FRONTEND_URL == "" {
		FRONTEND_URL = getMachineRunOnUrl("FRONTEND_URL")
	}
}
