package users

import (
	"net/mail"
	"fmt"
	"strings"
	"crypto/md5"
	"golang.org/x/crypto/bcrypt"
)

//bcryptCost is the default bcrypt cost to use when hashing passwords
var bcryptCost = 13

//User represents a user account in the database
type User struct {
	UserID    int64  `json:"userID"`
	Email     string `json:"-"` //never JSON encoded/decoded
	PassHash  []byte `json:"-"` //never JSON encoded/decoded
	UserName  string `json:"username"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

//Credentials represents user sign-in credentials
type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

//NewUser represents a new user signing up for an account
type NewUser struct {
	Email        string `json:"email"`
	Password     string `json:"password"`
	PasswordConf string `json:"passwordConf"`
	UserName     string `json:"userName"`
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
}

//Updates represents allowed updates to a user profile
type Updates struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

//Validate validates the new user and returns an error if
//any of the validation rules fail, or nil if its valid
func (nu *NewUser) Validate() error {
	//TODO: validate the new user according to these rules:
	//- Email field must be a valid email address (hint: see mail.ParseAddress)
	//- Password must be at least 6 characters
	//- Password and PasswordConf must match
	//- UserName must be non-zero length and may not contain spaces
	//use fmt.Errorf() to generate appropriate error messages if
	//the new user doesn't pass one of the validation rules
	_ , err := mail.ParseAddress(nu.Email)
	if err != nil {
		return fmt.Errorf("error in parsing email address: %v\n", err)
	}
	passLength := len(nu.Password)
	if passLength < 6 {
		return fmt.Errorf("password must be at least 6 characters long, but got %d", passLength)
	}
	if nu.Password != nu.PasswordConf {
		return fmt.Errorf("password and password confirmation must match. Expected %s, but got %s", nu.Password, nu.PasswordConf)
	}
	if len(nu.UserName) == 0{
		return fmt.Errorf("username must be non-zero length")
	}
	if strings.Contains(nu.UserName, " ") {
		return fmt.Errorf("username cannot contain spaces")
	}
	return nil
}

//ToUser converts the NewUser to a User, setting the
//PassHash fields appropriately
func (nu *NewUser) ToUser() (*User, error) {
	//TODO: call Validate() to validate the NewUser and
	//return any validation errors that may occur.
	//if valid, create a new *User and set the fields
	//based on the field values in `nu`.
	//Leave the ID field as the zero-value; your Store
	//implementation will set that field to the DBMS-assigned
	//primary key value.
	err := nu.Validate()
	if err != nil {
		return nil, err
	}
	processedEmail := strings.ToLower(strings.TrimSpace(nu.Email))
	h := md5.New()
	h.Write([]byte(processedEmail))
	user := &User{
		UserID: 0,
		Email: nu.Email,
		UserName: nu.UserName,
		FirstName: nu.FirstName,
		LastName: nu.LastName,
	}
	//TODO: also call .SetPassword() to set the PassHash
	//field of the User to a hash of the NewUser.Password
	err = user.SetPassword(nu.Password)
	if err != nil {
		return nil, err
	}
	return user, nil
}

//FullName returns the user's full name, in the form:
// "<FirstName> <LastName>"
//If either first or last name is an empty string, no
//space is put between the names. If both are missing,
//this returns an empty string
func (u *User) FullName() string {
	//TODO: implement according to comment above
	firstName := u.FirstName
	lastName := u.LastName
	if len(firstName) > 0 && len(lastName) > 0 {
		return firstName + " " + lastName
	} else {
		return firstName + lastName
	}
}

//SetPassword hashes the password and stores it in the PassHash field
func (u *User) SetPassword(password string) error {
	//TODO: use the bcrypt package to generate a new hash of the password
	//https://godoc.org/golang.org/x/crypto/bcrypt
	bCryptPass, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return err
	}
	u.PassHash = bCryptPass
	return nil
}

//Authenticate compares the plaintext password against the stored hash
//and returns an error if they don't match, or nil if they do
func (u *User) Authenticate(password string) error {
	//TODO: use the bcrypt package to compare the supplied
	//password with the stored PassHash
	//https://godoc.org/golang.org/x/crypto/bcrypt
	err := bcrypt.CompareHashAndPassword(u.PassHash, []byte(password))
	if err != nil {
		return err
	}
	return nil
}

//ApplyUpdates applies the updates to the user. An error
//is returned if the updates are invalid
func (u *User) ApplyUpdates(updates *Updates) error {
	//TODO: set the fields of `u` to the values of the related
	//field in the `updates` struct
	if len(updates.FirstName) == 0 && len(updates.LastName) == 0 {
		return fmt.Errorf("first and last name needs to be non-zero value to be updated")
	}
	u.FirstName = updates.FirstName
	u.LastName = updates.LastName
	return nil
}
