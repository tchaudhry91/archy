package store

import "golang.org/x/crypto/bcrypt"

// User is a struct to hold basic login information
type User struct {
	User     string
	Password string
}

// NewUser returns a user with hashed password
func NewUser(user, password string) (*User, error) {
	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return &User{
		user,
		string(passHash),
	}, nil
}

// CheckPassword compares the hash of the supplied password against the expected hash
func (u *User) CheckPassword(password string) bool {
	if bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)) != nil {
		return false
	}
	return true
}
