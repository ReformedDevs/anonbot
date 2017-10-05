package db

import (
	"encoding/base64"

	"golang.org/x/crypto/bcrypt"
)

const (
	UserTypeRegular = "regular"
	UserTypeStaff   = "staff"
	UserTypeAdmin   = "admin"
)

// User represents an individual user that can login to the website. Regular
// users are able to suggest tweets and staff are able to edit and queue them.
// Administrators are able to add and remove users.
//
// Passwords are salted and hashed with bcrypt.
type User struct {
	ID       int64
	Username string `gorm:"not null;unique_index"`
	Password string `gorm:"not null"`
	Type     string `gorm:"not null"`
}

// Authenticate hashes the password and compares it to the value stored in the
// database. An error is returned if the values do not match.
func (u *User) Authenticate(password string) error {
	h, err := base64.StdEncoding.DecodeString(u.Password)
	if err != nil {
		return err
	}
	return bcrypt.CompareHashAndPassword(h, []byte(password))
}

// SetPassword salts and hashes the user's password. It does not store the new
// value in the database.
func (u *User) SetPassword(password string) error {
	h, err := bcrypt.GenerateFromPassword([]byte(password), 0)
	if err != nil {
		return err
	}
	u.Password = base64.StdEncoding.EncodeToString(h)
	return nil
}
