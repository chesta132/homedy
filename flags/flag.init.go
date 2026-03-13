package flags

import "flag"

var (
	EnvPath = flag.String("env", "", "env path")
)

func init() {
	flag.Parse()
}
