package models

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type UserRole string

const (
	OwnerUserRole UserRole = "OWNER"
	UserUserRole  UserRole = "USER"
)

type UserStatus string

const (
	AvailableUserStatus UserStatus = "AVAILABLE"
	ReservedUserStatus  UserStatus = "RESERVED"
	CheckedInUserStatus UserStatus = "CHECKED-IN"
)

type User struct {
	ID                  string     `json:"_id,omitempty"`
	Username            string     `json:"username"`
	Name                string     `json:"name"`
	Surname             string     `json:"surname"`
	Tel                 string     `json:"tel"`
	Status              UserStatus `json:"status" default:"AVAILABLE"`
	Role                string     `json:"role"`
	Password            string     `json:"password,omitempty"`
	ResetPasswordToken  string     `json:"resetPasswordToken,omitempty"`
	ResetPasswordExpire time.Time  `json:"resetPasswordExpire,omitempty"`
	CreatedAt           time.Time  `json:"createdAt,omitempty"`
}

// EncryptPassword encrypts the user's password using bcrypt
func (u *User) EncryptPassword() error {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	return nil
}

// GenerateJWT generates a JSON Web Token for the user
func (u *User) GenerateJWT() (string, error) {
	claims := jwt.MapClaims{
		"id":  u.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(), // 30 days
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte("your-secret-key")) // Replace with your secret key
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

// ComparePassword compares the user's entered password with the hashed password
func (u *User) ComparePassword(enteredPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(enteredPassword))
	return err == nil
}
