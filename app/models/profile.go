package models

import (
	"fmt"
	"github.com/go-gorp/gorp"
	"github.com/revel/revel"
	"regexp"
	"strings"
)

type Profile struct {
	ProfileId          int
	UserId             int
	UserName           string
	Name               string
	Summary            string
	Description        string
	PhotoUrl           string
	AggregateFollowers int
	AggregateFollowing int

	// Transient
	User *User
}

var NameRegex = regexp.MustCompile("^[^#@]+$")

var UserNameRegex = regexp.MustCompile("^[a-zA-Z0-9]+$")

var UserNameBlacklistRegex = regexp.MustCompile("^(account|contact|about|find|search|public|to(desktop|mobile)|log(in|out)|sign(in|up|out)|register|home|index|default|post(s)?|user(name)?(s)?|i)$")

func (p *Profile) String() string {
	return fmt.Sprintf("Profile(%s)", p.Summary)
}

func (profile *Profile) Validate(v *revel.Validation) {
	ValidateProfileName(v, profile.Name)
	ValidateProfileSummary(v, profile.Summary)
	ValidateProfileDescription(v, profile.Description)
	ValidateProfilePhotoUrl(v, profile.PhotoUrl)
}

func ValidateProfileUserName(v *revel.Validation, username string) *revel.ValidationResult {
	result := v.Required(username).Message("User name required")
	if !result.Ok {
		return result
	}

	result = v.MaxSize(username, 64).Message("User name can not exceed 64 characters")
	if !result.Ok {
		return result
	}

	result = v.Match(username, UserNameRegex).Message("Invalid User name. Alphanumerics allowed only")
	if !result.Ok {
		return result
	}

	// Inverse regexp matcher (username cannot be the same as any blacklisted usernames)
	if blacklistMatcher := UserNameBlacklistRegex.FindString(username); blacklistMatcher != "" {
		result = v.Error("Invalid User name. Reserved keywords not allowed")
	}

	return result
}

func ValidateProfileName(v *revel.Validation, name string) *revel.ValidationResult {
	result := v.Required(name).Message("Name required")
	if !result.Ok {
		return result
	}

	result = v.MinSize(name, 2).Message("Name must be at least 2 characters")
	if !result.Ok {
		return result
	}

	result = v.MaxSize(name, 100).Message("Name must be at most 100 characters")
	if !result.Ok {
		return result
	}

	// Inverse regexp matcher (name cannot contain reserved # or @ symbols)
	if invalidNameMatcher := NameRegex.FindString(name); invalidNameMatcher == "" {
		result = v.Error("Invalid Name. Reserved characters ('#' and '@') are not allowed")
	}

	return result
}

func ValidateProfileSummary(v *revel.Validation, summary string) *revel.ValidationResult {
	result := v.MaxSize(summary, 140).Message("Profile summary cannot exceed 140 characters")

	return result
}

func ValidateProfileDescription(v *revel.Validation, description string) *revel.ValidationResult {
	result := v.MaxSize(description, 400).Message("Profile description cannot exceed 400 characters")

	return result
}

func ValidateProfilePhotoUrl(v *revel.Validation, photoUrl string) *revel.ValidationResult {
	result := v.MaxSize(photoUrl, 200).Message("Photo URL cannot exceed 200 characters")

	return result
}

func (p *Profile) PreInsert(_ gorp.SqlExecutor) error {
	p.UserName = strings.ToLower(p.UserName)
	return nil
}

func (p *Profile) PreUpdate(_ gorp.SqlExecutor) error {
	p.UserName = strings.ToLower(p.UserName)
	return nil
}

func (p *Profile) PostGet(exe gorp.SqlExecutor) error {
	var (
		obj interface{}
		err error
	)

	obj, err = exe.Get(User{}, p.UserId)
	if err != nil {
		return fmt.Errorf("Error loading a profile's user (%d): %s", p.UserId, err)
	}
	p.User = obj.(*User)

	/*obj, err = exe.Get(Post{}, p.ProfileId)
	if err != nil {
		return fmt.Errorf("Error loading a profile's posts (%d): %s", p.ProfileId, err)
	}
	var posts []*Post
	for _, post := range obj {
		posts = append(posts, post.(*Post))
	}
	p.Posts = posts*/

	return nil
}
