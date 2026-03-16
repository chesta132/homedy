package validatorlib

import (
	"homedy/internal/models"
	"slices"
)

func ValidateMessageType(value any) bool {
	str, ok := validateStr(value)
	if !ok {
		return false
	}

	if !slices.Contains(models.MessageTypes, models.MessageType(str)) {
		return false
	}

	return true
}

func ValidateMemberRole(value any) bool {
	str, ok := validateStr(value)
	if !ok {
		return false
	}

	if !slices.Contains(models.MemberRoles, models.MemberRole(str)) {
		return false
	}

	return true
}
