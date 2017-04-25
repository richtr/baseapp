package controllers

import (
	"golang.org/x/crypto/bcrypt"
	r "github.com/revel/revel"
	"github.com/richtr/baseapp/app/models"
	"github.com/richtr/baseapp/app/routes"
	"strings"
)

type Profile struct {
	Account
}

func (c Profile) loadProfileByUserName(username string) *models.Profile {
	var profile models.Profile
	err := c.Txn.SelectOne(&profile, "select * from Profile where UserName=?", username)
	if err != nil {
		return nil
	}

	return &profile
}

func (c Profile) getProfileShowParams(username string) (profiles *models.Profile, owner, following bool) {

	profile := c.loadProfileByUserName(username)

	if profile == nil {
		return nil, false, false
	}

	user := c.connected()

	isOwner := false
	isFollowing := false
	if user != nil {
		if user.UserId == profile.User.UserId { // Check if logged in user owns the current profile
			isOwner = true
		} else { // Check if logged in user is following the current profile
			fErr := c.Txn.SelectOne(&models.Follower{}, `select * from Follower where UserId = ? and FollowUserId = ?`, user.UserId, profile.User.UserId)
			if fErr == nil {
				isFollowing = true
			}
		}
	}

	return profile, isOwner, isFollowing

}

func (c Profile) Show(username string) r.Result {
	profile, isOwner, isFollowing := c.getProfileShowParams(username)

	if profile == nil {
		return c.NotFound("Profile does not exist")
	}

	// Retrieve all posts for profile
	var posts []*models.Post
	results, err := c.Txn.Select(models.Post{}, `select * from Post where ProfileId = ?`, profile.ProfileId)
	if err == nil {
		for _, r := range results {
			posts = append(posts, r.(*models.Post))
		}
	}

	appName := r.Config.StringDefault("app.name", "BaseApp")

	title := profile.Name + " on " + appName

	return c.Render(title, profile, posts, isOwner, isFollowing)
}

func (c Profile) Settings(username string) r.Result {
	profile := c.connected()
	if profile == nil || profile.UserName != username {
		c.Flash.Error("You must log in to access your account")
		return c.Redirect(routes.Account.Logout())
	}

	return c.Render(profile)
}

func (c Profile) UpdateSettings(username string, profile *models.Profile, verifyPassword string) r.Result {
	existingProfile := c.connected()
	if existingProfile == nil || existingProfile.UserName != username {
		c.Flash.Error("You must log in to access your account")
		return c.Redirect(routes.Account.Logout())
	}

	email := profile.User.Email

	// Step 1: Validate data

	if email != existingProfile.User.Email {
		// Validate email
		models.ValidateUserEmail(c.Validation, profile.User.Email).Key("profile.User.Email")
	}

	if profile.User.Password != "" || verifyPassword != "" {
		models.ValidateUserPassword(c.Validation, profile.User.Password).Key("profile.User.Password")

		// Additional password verification
		c.Validation.Required(profile.User.Password != profile.User.Email).Message("Password cannot be the same as your email address").Key("profile.User.Password")
		c.Validation.Required(verifyPassword).Message("Password verification required").Key("verifyPassword")
		c.Validation.Required(verifyPassword == profile.User.Password).Message("Provided passwords do not match").Key("verifyPassword")
	}

	// Validate profile components
	models.ValidateProfileUserName(c.Validation, profile.UserName).Key("profile.UserName")
	models.ValidateProfileName(c.Validation, profile.Name).Key("profile.Name")
	models.ValidateProfileSummary(c.Validation, profile.Summary).Key("profile.Summary")
	models.ValidateProfileDescription(c.Validation, profile.Description).Key("profile.Description")

	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		c.Flash.Error("Profile could not be updated")
		return c.Redirect(routes.Profile.Settings(username))
	}

	// Step 2: Commit data

	if email != existingProfile.User.Email {
		userExists := c.getProfileByEmailAddress(email)

		if userExists != nil {
			c.Flash.Error("Email address is already registered to another account")
			return c.Redirect(routes.Profile.Settings(username))
		}

		// Re-send email confirmation
		existingProfile.User.Email = email
		existingProfile.User.Confirmed = false

		// FIXME
		/*// Send out confirmation email
		eErr := c.sendAccountConfirmEmail(existingProfile.User)

		if eErr != nil {
			c.Flash.Error("Could not send confirmation email")
		} else {*/

		// Update email address in database
		_, err := c.Txn.Exec("update User set Email = ?, Confirmed = ? where UserId = ?",
			email, 0, existingProfile.User.UserId)
		if err != nil {
			panic(err)
		}

		// Update session value
		c.Session["userEmail"] = email

		//}

	}

	profile.UserName = strings.ToLower(profile.UserName)

	// Update username?
	if profile.UserName != existingProfile.UserName {
		userExists := c.getProfileByUserName(profile.UserName)

		if userExists != nil {
			c.Flash.Error("User name is already registered to another account")
			return c.Redirect(routes.Profile.Settings(username))
		}

		// Change username
		existingProfile.UserName = profile.UserName
	}

	// Update password?
	if profile.User.Password != "" || verifyPassword != "" {
		c.CommitPassword(existingProfile.User, profile.User.Password)
	}

	// Update other profile components
	existingProfile.Name = profile.Name
	existingProfile.Summary = profile.Summary
	existingProfile.Description = profile.Description

	_, err := c.Txn.Update(existingProfile)
	if err != nil {
		c.Flash.Error("Profile could not be updated")
		return c.Redirect(routes.Profile.Settings(username))
	}

	c.Flash.Success("Profile has been updated")
	return c.Redirect(routes.Profile.Show(existingProfile.UserName))
}

func (c Profile) Password(username string) r.Result {
	profile := c.connected()
	if profile == nil || profile.UserName != username {
		c.Flash.Error("You must log in to access your account")
		return c.Redirect(routes.Account.Logout())
	}

	return c.Render()
}

func (c Profile) CommitPassword(user *models.User, password string) error {
	bcryptPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	_, err := c.Txn.Exec("update User set HashedPassword = ? where UserId = ?",
		bcryptPassword, user.UserId)
	if err != nil {
		return err
	}

	return nil
}

func (c Profile) UpdatePassword(username string, password, verifyPassword string) r.Result {
	profile := c.connected()
	if profile == nil || profile.UserName != username {
		c.Flash.Error("You must log in to access your account")
		return c.Redirect(routes.Account.Logout())
	}

	// Validate password
	models.ValidateUserPassword(c.Validation, password).Key("password")

	// Additional password verification
	c.Validation.Required(password != profile.User.Email).Message("Password cannot be the same as your email address").Key("password")
	c.Validation.Required(verifyPassword).Message("Password verification required").Key("verifyPassword")
	c.Validation.Required(verifyPassword == password).Message("Provided passwords do not match").Key("verifyPassword")

	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		c.Flash.Error("Password could not be updated")
		return c.Redirect(routes.Profile.Settings(username))
	}

	err := c.CommitPassword(profile.User, password)

	if err != nil {
		c.Flash.Error("Password validation failed")
		return c.Redirect(routes.Profile.Settings(username))
	}

	c.Flash.Success("Account settings updated")
	return c.Redirect(routes.Profile.Show(username))
}

func (c Profile) FollowUser(username string) r.Result {
	followResponse := models.SimpleJSONResponse{"fail", ""}

	profile := c.connected()
	if profile == nil {
		followResponse.Message = "You must log in to follow another user"
		return c.RenderJSON(followResponse)
	}

	if profile.UserName == username {
		followResponse.Message = "You cannot follow yourself"
		return c.Render(followResponse)
	}

	// Get followed user profile
	followProfile := c.loadProfileByUserName(username)
	if followProfile == nil {
		followResponse.Message = "User with that username not found"
		return c.RenderJSON(followResponse)
	}

	var followerObj models.Follower
	err := c.Txn.SelectOne(&followerObj, `select * from Follower where UserId = ? and FollowUserId = ?`, profile.User.UserId, followProfile.User.UserId)

	if err != nil {

		// Add new follower

		followerObj = models.Follower{
			UserId:       profile.User.UserId,
			FollowUserId: followProfile.User.UserId,
		}

		lErr := c.Txn.Insert(&followerObj)
		if lErr != nil {
			panic(lErr)
		}

		// Update aggregate follower count on Followed Profile
		followProfile.AggregateFollowers += 1

		_, pErr := c.Txn.Update(followProfile)
		if pErr != nil {
			panic(pErr)
		}

		// Update aggregate following count on Current User Profile
		profile.AggregateFollowing += 1

		_, p2Err := c.Txn.Update(profile)
		if p2Err != nil {
			panic(p2Err)
		}

		followResponse.Message = "You are now following this user"
		followResponse.Status = "success"

	} else {

		// Remove existing follower

		_, dErr := c.Txn.Delete(&followerObj)
		if dErr != nil {
			panic(dErr)
		}

		// Update aggregate follower count on Followed Profile
		followProfile.AggregateFollowers -= 1

		_, pErr := c.Txn.Update(followProfile)
		if pErr != nil {
			panic(pErr)
		}

		// Update aggregate following count on Current User Profile
		profile.AggregateFollowing -= 1

		_, p2Err := c.Txn.Update(profile)
		if p2Err != nil {
			panic(p2Err)
		}

		followResponse.Message = "You are no longer following this user"
		followResponse.Status = "success"

	}

	return c.RenderJSON(followResponse)

}

func (c Profile) Followers(username string, page int) r.Result {

	profile, isOwner, isFollowing := c.getProfileShowParams(username)

	if profile == nil {
		return c.NotFound("Profile does not exist")
	}

	if page == 0 {
		page = 1
	}
	nextPage := page + 1
	size := 50 // results per page

	// Retrieve all profiles of followers
	var followerProfiles []*models.Profile
	results, err := c.Txn.Select(models.Profile{}, `select * from Profile where UserId in (select UserId from Follower where FollowUserId = ?) limit ?, ?`, profile.User.UserId, (page-1)*size, size)
	if err == nil {
		for _, r := range results {
			followerProfiles = append(followerProfiles, r.(*models.Profile))
		}
	}

	if len(followerProfiles) == 0 && page != 1 {
		return c.Redirect(routes.Profile.Followers(username, 1))
	}

	return c.Render(profile, isOwner, isFollowing, followerProfiles, page, nextPage)
}

func (c Profile) Following(username string, page int) r.Result {

	profile, isOwner, isFollowing := c.getProfileShowParams(username)

	if profile == nil {
		return c.NotFound("Profile does not exist")
	}

	if page == 0 {
		page = 1
	}
	nextPage := page + 1
	size := 50 // results per page

	// Retrieve all profiles of followers
	var followingProfiles []*models.Profile
	results, err := c.Txn.Select(models.Profile{}, `select * from Profile where UserId in (select FollowUserId from Follower where UserId = ?) limit ?, ?`, profile.User.UserId, (page-1)*size, size)
	if err == nil {
		for _, r := range results {
			followingProfiles = append(followingProfiles, r.(*models.Profile))
		}
	}

	if len(followingProfiles) == 0 && page != 1 {
		return c.Redirect(routes.Profile.Following(username, 1))
	}

	return c.Render(profile, isOwner, isFollowing, followingProfiles, page, nextPage)
}
