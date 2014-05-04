package models

import (
	"fmt"
	"github.com/coopernurse/gorp"
	"github.com/revel/revel"
)

type Profile struct {
	ProfileId          int
	UserId             int
	Name               string
	Summary            string
	Description        string
	PhotoUrl           string
	AggregateFollowers int
	AggregateFollowing int

	// Transient
	User               *User
}

func (p *Profile) String() string {
	return fmt.Sprintf("Profile(%s)", p.Summary)
}

func (profile *Profile) Validate(v *revel.Validation) {
	ValidateProfileName(v, profile.Name)
	ValidateProfileSummary(v, profile.Summary)
	ValidateProfileDescription(v, profile.Description)
	ValidateProfilePhotoUrl(v, profile.PhotoUrl)
}

func ValidateProfileName(v *revel.Validation, name string) *revel.ValidationResult {
	result := v.Required(name).Message("Name required")
	if !result.Ok {
		return result
	}

	result = v.MinSize(name, 6).Message("Name must be at least 6 characters")
	if !result.Ok {
		return result
	}

	result = v.MaxSize(name, 100).Message("Name must be at most 100 characters")

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
