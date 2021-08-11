package models

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

// User is used by pop to map your users database table to your go code.
type User struct {
	ID              uuid.UUID `json:"id" db:"id"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
	FirstName       string    `json:"first_name" db:"first_name"`
	LastName        string    `json:"last_name" db:"last_name"`
	Username        string    `json:"username" db:"username"`
	Password        string    `json:"-" db:"-"`
	ConfirmPassword string    `json:"-" db:"-"`
	PasswordHash    string    `json:"-" db:"password_hash"`
}

// String is not required by pop and may be deleted
func (u User) String() string {
	ju, _ := json.Marshal(u)
	return string(ju)
}

// Users is not required by pop and may be deleted
type Users []User

// String is not required by pop and may be deleted
func (u Users) String() string {
	ju, _ := json.Marshal(u)
	return string(ju)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (u *User) Validate(tx *pop.Connection) (*validate.Errors, error) {
	var err error
	return validate.Validate(
		&validators.StringIsPresent{Field: u.FirstName, Name: "FirstName"},
		&validators.StringIsPresent{Field: u.LastName, Name: "LastName"},
		&validators.StringIsPresent{Field: u.Username, Name: "Username"},
		&validators.StringIsPresent{Field: u.PasswordHash, Name: "PasswordHash"},
		&validators.FuncValidator{
			Field:   u.Username,
			Name:    "Username",
			Message: "%s already taken.",
			Fn: func() bool {
				var exist bool
				chosenName := tx.Where("username = ?", u.Username)
				if u.ID != uuid.Nil {
					chosenName = chosenName.Where("id != ?", u.ID)
				}
				exist, err = chosenName.Exists(u)
				if err != nil {
					return false
				}
				return !exist
			},
		},
	), err
}

// function to create a hash password before storing in db
func (u *User) Create(tx *pop.Connection) (*validate.Errors, error) {
	u.Username = strings.ToLower(u.Username)
	passharsh, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return validate.NewErrors(), err
	}
	u.PasswordHash = string(passharsh)
	return tx.ValidateAndCreate(u)
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (u *User) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	var err error
	return validate.Validate(
		&validators.StringIsPresent{Field: u.Password, Name: "Password"},
		&validators.StringsMatch{Name: "ConfirmPassword", Field: u.Password, Field2: u.ConfirmPassword, Message: "Passwords do not match"},
	), err
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (u *User) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
