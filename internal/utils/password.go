package utils

const (
	passwordPrefix = "prefix"
	passwordSuffix = "suffix"
)

func GeneratePassword(text string) string {
	password := MD5(passwordPrefix + text + passwordSuffix)
	return password
}
