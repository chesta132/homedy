package logger

import "io"

// Color codes
const (
	colorReset = "\033[0m"
	colorBold  = "\033[1m"
	colorDim   = "\033[2m"

	colorBlack   = "\033[30m"
	colorRed     = "\033[31m"
	colorGreen   = "\033[32m"
	colorYellow  = "\033[33m"
	colorBlue    = "\033[34m"
	colorMagenta = "\033[35m"
	colorCyan    = "\033[36m"
	colorWhite   = "\033[37m"

	colorBgRed     = "\033[41m"
	colorBgGreen   = "\033[42m"
	colorBgYellow  = "\033[43m"
	colorBgBlue    = "\033[44m"
	colorBgMagenta = "\033[45m"
	colorBgCyan    = "\033[46m"
)

// Level defines log severity
type Level int

const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
	FATAL
)

func (l Level) String() string {
	switch l {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	case FATAL:
		return "FATAL"
	default:
		return "???"
	}
}

func (l Level) color() string {
	switch l {
	case DEBUG:
		return colorCyan
	case INFO:
		return colorGreen
	case WARN:
		return colorYellow
	case ERROR:
		return colorRed
	case FATAL:
		return colorBgRed + colorWhite
	default:
		return colorWhite
	}
}

func (l Level) badge() string {
	switch l {
	case DEBUG:
		return "  DBG  "
	case INFO:
		return "  INF  "
	case WARN:
		return "  WRN  "
	case ERROR:
		return "  ERR  "
	case FATAL:
		return "  FTL  "
	default:
		return "  ???  "
	}
}

// Config holds logger configuration
type Config struct {
	Output     io.Writer
	Level      Level
	TimeFormat string
	NoColor    bool
	Prefix     string
}

// Logger is the main logger struct
type Logger struct {
	config Config
}
