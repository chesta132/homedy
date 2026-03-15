package validatorlib

func ValidateShareName(value any) bool {
	return value != "global" && value != "printers" && value != "print$" && value != "config" && value != "backup" && value != "restore"
}
