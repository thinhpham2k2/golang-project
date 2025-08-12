package utils

import "github.com/nicksnyder/go-i18n/v2/i18n"

var INTERNAL_ERROR = &i18n.Message{
	ID:    "INTERNAL_ERROR",
	Other: "Internal server error",
}

var INVALID_VALUE = &i18n.Message{
	ID:    "INVALID_VALUE",
	Other: "Invalid value",
}

var INVALID_BIRTHDAY = &i18n.Message{
	ID:    "INVALID_BIRTHDAY",
	Other: "Birthday must be in the format YYYY-MM-DD and the age must be between 5 and 100 years old",
}

var INVALID_ROLE = &i18n.Message{
	ID:    "INVALID_ROLE",
	Other: "Role must be one of the following: admin, staff, or customer",
}

var ROLE_REQUIRE = &i18n.Message{
	ID:    "ROLE_REQUIRE",
	Other: "Role is required",
}

var PASSWORD_ENCRYPTION_FAIL = &i18n.Message{
	ID:    "PASSWORD_ENCRYPTION_FAIL",
	Other: "Password encryption failed",
}

var INVALID_PASSWORD = &i18n.Message{
	ID:    "INVALID_PASSWORD",
	Other: "Password must be 8–36 characters long and contain only lowercase letters, numbers, dots, or underscores",
}

var PASSWORD_REQUIRE = &i18n.Message{
	ID:    "PASSWORD_REQUIRE",
	Other: "Password is required",
}

var DUPLICATE_USERNAME = &i18n.Message{
	ID:    "DUPLICATE_USERNAME",
	Other: "Username is already taken",
}

var INVALID_USERNAME = &i18n.Message{
	ID:    "INVALID_USERNAME",
	Other: "Username must be 3–24 characters long and contain only lowercase letters, numbers, dots, or underscores",
}

var USERNAME_REQUIRE = &i18n.Message{
	ID:    "USERNAME_REQUIRE",
	Other: "Username is required",
}

var INVALID_CLAIM = &i18n.Message{
	ID:    "INVALID_CLAIM",
	Other: "Invalid claims",
}

var PERMISSION_REQUIRE = &i18n.Message{
	ID:    "PERMISSION_REQUIRE",
	Other: "You do not have permission to access this resource",
}

var AUTHEN_REQUIRE = &i18n.Message{
	ID:    "AUTHEN_REQUIRE",
	Other: "Authentication required",
}

var INVALID_AUTHOR_HEADER = &i18n.Message{
	ID:    "INVALID_AUTHOR_HEADER",
	Other: "Missing or invalid Authorization header",
}

var INVALID_USERNAME_PASSWORD = &i18n.Message{
	ID:    "INVALID_USERNAME_PASSWORD",
	Other: "Invalid username or password",
}

var FAIL_CREATE_TOKEN = &i18n.Message{
	ID:    "FAIL_CREATE_TOKEN",
	Other: "Fail to create token",
}

var NOT_FOUND = &i18n.Message{
	ID:    "NOT_FOUND",
	Other: "Not found item",
}

var UPDATE_FAIL = &i18n.Message{
	ID:    "UPDATE_FAIL",
	Other: "Update failed",
}

var CREATE_FAIL = &i18n.Message{
	ID:    "CREATE_FAIL",
	Other: "Create failed",
}

var DELETE_FAIL = &i18n.Message{
	ID:    "DELETE_FAIL",
	Other: "Delete failed",
}
