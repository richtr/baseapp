package models

import (
	"fmt"
	"github.com/go-gorp/gorp"
	"github.com/revel/revel"
	"time"
)

type Post struct {
	PostId         int
	ProfileId      int
	Title          string
	Content        []byte // mediumblob
	Status         string
	Date           string
	AggregateLikes int // to keep things quick in our app

	// Transient
	ContentStr string
	DateObj    time.Time
}

func (post *Post) Validate(v *revel.Validation) {
	ValidatePostTitle(v, post.Title)
}

func ValidatePostTitle(v *revel.Validation, title string) *revel.ValidationResult {
	result := v.Required(title).Message("A post must have a title")
	if !result.Ok {
		return result
	}

	result = v.MinSize(title, 3).Message("Post title must exceed 2 characters")
	if !result.Ok {
		return result
	}

	result = v.MaxSize(title, 200).Message("Post title cannot exceed 200 characters")

	return result
}

// These hooks work around two things:
// - Gorp's lack of support for loading relations automatically.
// - Sqlite's lack of support for datetimes.

func (post *Post) PreInsert(_ gorp.SqlExecutor) error {
	post.Date = post.DateObj.Format(SQL_DATE_FORMAT)
	post.Content = []byte(post.ContentStr)
	return nil
}

func (post *Post) PreUpdate(_ gorp.SqlExecutor) error {
	post.Content = []byte(post.ContentStr)
	return nil
}

func (post *Post) PostGet(_ gorp.SqlExecutor) error {
	var (
		err error
	)

	if post.DateObj, err = time.Parse(SQL_DATE_FORMAT, post.Date); err != nil {
		return fmt.Errorf("Error parsing post created date '%s':", post.DateObj, err)
	}

	post.ContentStr = string(post.Content)

	return nil
}
