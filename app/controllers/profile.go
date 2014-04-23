package controllers

import (
	"code.google.com/p/go.crypto/bcrypt"
	r "github.com/revel/revel"
	"github.com/richtr/baseapp/app/routes"
	"github.com/richtr/baseapp/app/models"
)

type Profile struct {
	Account
}

func (c Profile) Index() r.Result {
	return c.NotFound("Profile does not exist", 404)
}

func (c Profile) loadProfileById(id int) *models.Profile {
	p, err := c.Txn.Get(models.Profile{}, id)

	if err != nil || p == nil {
		return nil
	}

	return p.(*models.Profile)
}

func (c Profile) Show(id int) r.Result {
	profile := c.loadProfileById(id)

	if profile == nil {
		return c.NotFound("Profile does not exist")
	}

	appName := r.Config.StringDefault("app.name", "BaseApp")

	title := profile.Name + " on " + appName

	isOwner := false
	if user := c.connected(); user != nil && user.UserId == profile.User.UserId {
		isOwner = true
	}

	// Retrieve all posts for profile
	var posts []*models.Post
	results, err := c.Txn.Select(models.Post{}, `select * from Post where ProfileId = ?`, profile.ProfileId)
	if err == nil {
		for _, r := range results {
			posts = append(posts, r.(*models.Post))
		}
	}

	return c.Render(title, profile, posts, isOwner)
}

func (c Profile) Settings(id int) r.Result {
	profile := c.connected();
	if profile == nil || profile.UserId != id {
		c.Flash.Error("You must log in to access your account");
		return c.Redirect(routes.Account.Logout())
	}

	return c.Render()
}

func (c Profile) UpdateSettings(id int, user *models.Profile, verifyPassword string) r.Result {
	profile := c.connected();
	if profile == nil || profile.UserId != id {
		c.Flash.Error("You must log in to access your account");
		return c.Redirect(routes.Account.Logout())
	}

	email := user.User.Email

	// Step 1: Validate data

	if email != profile.User.Email {
		// Validate email
		models.ValidateUserEmail(c.Validation, user.User.Email).Key("user.User.Email")
	}

	if user.User.Password != "" || verifyPassword != "" {
		models.ValidateUserPassword(c.Validation, user.User.Password).Key("user.User.Password")

		// Additional password verification
		c.Validation.Required(user.User.Password != user.User.Email).Message("Password cannot be the same as your email address").Key("user.User.Password")
		c.Validation.Required(verifyPassword).Message("Password verification required").Key("verifyPassword")
		c.Validation.Required(verifyPassword == user.User.Password).Message("Provided passwords do not match").Key("verifyPassword")
	}

	// Validate profile components
	models.ValidateProfileName(c.Validation, user.Name).Key("user.Name")
	models.ValidateProfileSummary(c.Validation, user.Summary).Key("user.Summary")
	models.ValidateProfileDescription(c.Validation, user.Description).Key("user.Description")

	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		c.Flash.Error("Profile could not be updated");
		return c.Redirect(routes.Profile.Settings(id))
	}

	// Step 2: Commit data

	if email != profile.User.Email {
		userExists := c.getProfile(email)

		if userExists != nil {
			c.Flash.Error("Email address is already registered to another account");
			return c.Redirect(routes.Profile.Settings(id))
		}

		// Re-send email confirmation
		profile.User.Email = email
		profile.User.Confirmed = false

		// FIXME
		/*// Send out confirmation email
		eErr := c.sendAccountConfirmEmail(profile.User)

		if eErr != nil {
			c.Flash.Error("Could not send confirmation email")
		} else {*/

			// Update email address in database
			_, err := c.Txn.Exec("update User set Email = ?, Confirmed = ? where UserId = ?",
			  email, 0, profile.User.UserId)
			if err != nil {
			  panic(err)
			}

			// Update session value
			c.Session["userEmail"] = email

		//}

	}

	// Update password?
	if user.User.Password != "" || verifyPassword != "" {
		c.CommitPassword(profile.User, user.User.Password)
	}

	// Update profile components
	profile.UserId = id
	profile.Name = user.Name
	profile.Summary = user.Summary
	profile.Description = user.Description

	_, err := c.Txn.Update(profile)
	if err != nil {
		c.Flash.Error("Profile could not be updated");
		return c.Redirect(routes.Profile.Settings(id));
	}

	c.Flash.Success("Profile has been updated");
	return c.Redirect(routes.Profile.Show(id))
}

func (c Profile) Password(id int) r.Result {
	profile := c.connected();
	if profile == nil || profile.UserId != id {
		c.Flash.Error("You must log in to access your account");
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


func (c Profile) UpdatePassword(id int, password, verifyPassword string) r.Result {
	profile := c.connected();
	if profile == nil || profile.UserId != id {
		c.Flash.Error("You must log in to access your account");
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
		return c.Redirect(routes.Profile.Settings(id))
	}

	err := c.CommitPassword(profile.User, password)

	if err != nil {
		c.Flash.Error("Password validation failed")
		return c.Redirect(routes.Profile.Settings(id))
	}

	c.Flash.Success("Account settings updated")
	return c.Redirect(routes.Profile.Show(id))
}
