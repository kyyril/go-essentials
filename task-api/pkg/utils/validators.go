package utils

import (
	"regexp"
	"strings"
)

// IsValidEmail checks if the email format is valid
func IsValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// IsValidPassword checks if the password meets minimum requirements
func IsValidPassword(password string) bool {
	// Minimum 6 characters
	if len(password) < 6 {
		return false
	}
	return true
}

// IsValidUserRole checks if the role is valid
func IsValidUserRole(role string) bool {
	validRoles := map[string]bool{
		"admin":  true,
		"user":   true,
		"viewer": true,
	}
	return validRoles[strings.ToLower(role)]
}

// IsValidProjectStatus checks if the project status is valid
func IsValidProjectStatus(status string) bool {
	validStatuses := map[string]bool{
		"active":   true,
		"archived": true,
	}
	return validStatuses[strings.ToLower(status)]
}

// IsValidTaskStatus checks if the task status is valid
func IsValidTaskStatus(status string) bool {
	validStatuses := map[string]bool{
		"todo":        true,
		"in_progress": true,
		"done":        true,
	}
	return validStatuses[strings.ToLower(status)]
}

// IsValidTaskPriority checks if the task priority is valid
func IsValidTaskPriority(priority string) bool {
	validPriorities := map[string]bool{
		"low":    true,
		"medium": true,
		"high":   true,
	}
	return validPriorities[strings.ToLower(priority)]
}

// IsValidProjectMemberRole checks if the project member role is valid
func IsValidProjectMemberRole(role string) bool {
	validRoles := map[string]bool{
		"owner":  true,
		"editor": true,
		"viewer": true,
	}
	return validRoles[strings.ToLower(role)]
}

// IsValidUUID checks if the string is a valid UUID format
func IsValidUUID(id string) bool {
	uuidRegex := regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)
	return uuidRegex.MatchString(strings.ToLower(id))
}

// ValidateUserName checks if user name is valid
func ValidateUserName(name string) bool {
	return len(strings.TrimSpace(name)) >= 2 && len(name) <= 255
}

// ValidateProjectName checks if project name is valid
func ValidateProjectName(name string) bool {
	return len(strings.TrimSpace(name)) >= 1 && len(name) <= 255
}

// ValidateTaskTitle checks if task title is valid
func ValidateTaskTitle(title string) bool {
	return len(strings.TrimSpace(title)) >= 1 && len(title) <= 255
}
