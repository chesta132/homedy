package logger

var (
	// Default logger instance
	Default = New()
	Cmd     = New().With("CMD")
)
