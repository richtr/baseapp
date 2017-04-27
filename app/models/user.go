package models

import (
	"fmt"
	"github.com/go-gorp/gorp"
	"github.com/revel/revel"
	"regexp"
	"time"
)

const (
	DATE_FORMAT     = "Jan _2, 2006"
	SQL_DATE_FORMAT = "2006-01-02"
)

type User struct {
	UserId         int
	Email          string
	HashedPassword []byte
	CreatedStr     string
	Confirmed      bool

	// Transient
	Password string
	Created  time.Time
	Profile  *Profile
}

func (u *User) String() string {
	return fmt.Sprintf("User(%s)", u.Email)
}

var EmailRegex = regexp.MustCompile("^[a-z0-9!#$\\%&'*+\\/=?^_`{|}~.-]+@[a-z0-9-]+(\\.[a-z0-9-]+)*$")

func (user *User) Validate(v *revel.Validation) {
	ValidateUserEmail(v, user.Email)
	ValidateUserPassword(v, user.Password)
}

func ValidateUserEmail(v *revel.Validation, email string) *revel.ValidationResult {
	result := v.Required(email).Message("Email address required")
	if !result.Ok {
		return result
	}

	result = v.MinSize(email, 6).Message("Email address can not be less than 6 characters")
	if !result.Ok {
		return result
	}

	result = v.MaxSize(email, 200).Message("Email address can not exceed 200 characters")
	if !result.Ok {
		return result
	}

	result = v.Match(email, EmailRegex).Message("You must provide a valid email address")

	return result
}

func ValidateUserPassword(v *revel.Validation, password string) *revel.ValidationResult {
	result := v.Required(password).Message("Password required")
	if !result.Ok {
		return result
	}

	result = v.MinSize(password, 6).Message("Password must be at least 6 characters")
	if !result.Ok {
		return result
	}

	result = v.MaxSize(password, 200).Message("Password must be at most 200 characters")

	return result
}

// These hooks work around two things:
// - Gorp's lack of support for loading relations automatically.
// - Sqlite's lack of support for datetimes.

func (user *User) PreInsert(_ gorp.SqlExecutor) error {
	user.CreatedStr = user.Created.Format(SQL_DATE_FORMAT)
	return nil
}

func (user *User) PostGet(exe gorp.SqlExecutor) error {
	var (
		err error
	)

	if user.Created, err = time.Parse(SQL_DATE_FORMAT, user.CreatedStr); err != nil {
		return fmt.Errorf("Error parsing created date '%s':", user.CreatedStr, err)
	}

	return nil
}
