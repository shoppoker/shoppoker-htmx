package models

import (
	"regexp"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const USERS_PER_PAGE = 2

type User struct {
	gorm.Model

	ID           uint
	Username     string `gorm:"unique"`
	PasswordHash string
	IsAdmin      bool
}

func NewUser(username, password string, is_admin bool) (*User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return &User{
		Username:     username,
		PasswordHash: string(hash),
		IsAdmin:      is_admin,
	}, nil
}

func (u *User) ComparePassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)) == nil
}

func (u *User) HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

// username must be from 4 to 16 characters and can only contain letters, numbers, and underscores
func IsUsernameValid(username string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	return len(username) >= 4 && len(username) <= 16 && re.MatchString(username)
}

func GetUsernameRules() string {
	return "Имя пользователя должно быть от 4 до 16 символов и содержать только буквы, цифры и нижнее подчеркивание"
}

// password must be at least 8 characters long
func IsPasswordValid(password string) bool {
	return len(password) >= 8
}

func GetPasswordRules() string {
	return "Пароль должен содержать не менее 8 символов"
}
