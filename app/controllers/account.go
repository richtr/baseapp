package controllers

import (
	"code.google.com/p/go.crypto/bcrypt"
	r "github.com/revel/revel"
	m "github.com/revel/revel/mail"
	gr "github.com/ftrvxmtrx/gravatar"
	"github.com/richtr/baseapp/app/models"
	"github.com/richtr/baseapp/app/routes"
	"fmt"
	"crypto/rand"
	"time"
	"strings"
)

type Account struct {
	GorpController
}

const alphanum = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789abcdefghijklmnopqrstuvwxyz"

func (c Account) AddUser() r.Result {
	if profile := c.connected(); profile != nil {
		c.RenderArgs["user"] = profile
	}
	return nil
}

// Add .AppName to RenderArgs
func (c Account) AddAppName() r.Result {
	c.RenderArgs["AppName"] = r.Config.StringDefault("app.name", "BaseApp")
	return nil
}

func (c Account) connected() *models.Profile {
	if c.RenderArgs["user"] != nil {
		return c.RenderArgs["user"].(*models.Profile)
	}
	if email, ok := c.Session["userEmail"]; ok {
		profile := c.getProfileByEmailAddress(email)
		return profile
	}
	return nil
}

func (c Account) getProfileByEmailAddress(email string) *models.Profile {
	user := &models.User{}
	err := c.Txn.SelectOne(user, `select * from User where Email = ?`, email)
	if err != nil || user == nil {
		return nil
	}

	profile := c.getProfileByUserId(user.UserId)

	return profile
}

func (c Account) getProfileByUserName(username string) *models.Profile {
	profile := &models.Profile{}
	err := c.Txn.SelectOne(profile, `select * from Profile where UserName = ?`, username)
	if err != nil || profile == nil {
		return nil
	}
	return profile
}

func (c Account) getProfileByUserId(userId int) *models.Profile {
	// Bind associated user profile for admin purposes
	profile := &models.Profile{}
	err := c.Txn.SelectOne(profile, `select * from Profile where UserId = ?`, userId)
	if err != nil || profile == nil {
		return nil
	}

	return profile
}

func (c Account) Index() r.Result {
	profile := c.connected();
	if profile == nil {
		c.Flash.Error("You must log in to access your account");
		return c.Redirect(routes.Account.Logout())
	}

	return c.Redirect(routes.Profile.Show(profile.UserName))
}

func (c Account) Register() r.Result {
	return c.Render()
}

func (c Account) SaveUser(user models.User, username, name string) r.Result {

	// Lower case username
	username = strings.ToLower(username)
	lcPass := strings.ToLower(user.Password)

	// Validate User components
	models.ValidateUserEmail(c.Validation, user.Email).Key("user.Email")
	models.ValidateUserPassword(c.Validation, user.Password).Key("user.Password")

	// Additional user components verification
	c.Validation.Required(user.Password != user.Email).Message("Password cannot be the same as your email address").Key("user.Password")
	c.Validation.Required(lcPass != username).Message("Password cannot be the same as your user name").Key("user.Password")

	// Validate Profile components
	models.ValidateProfileUserName(c.Validation, username).Key("username")
	models.ValidateProfileName(c.Validation, name).Key("name")

	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		c.Flash.Error("Registration failed.")
		return c.Redirect(routes.Account.Register())
	}

	userExists := c.getProfileByEmailAddress(user.Email)

	if userExists != nil {
		c.Validation.Keep()
		c.FlashParams()
		c.Flash.Error("Email '" + user.Email + "' is already registered.")
		return c.Redirect(routes.Account.Register())
	}

	userExists = c.getProfileByUserName(username)

	if userExists != nil {
		c.Validation.Keep()
		c.FlashParams()
		c.Flash.Error("User name '" + username + "' is already taken.")
		return c.Redirect(routes.Account.Register())
	}

	user.HashedPassword, _ = bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	user.Created = time.Now()

	user.Confirmed = false

	err := c.Txn.Insert(&user)
	if err != nil {
		panic(err)
	}

	// Create profile (and assign correct UserId)
	profile := &models.Profile{0, user.UserId, username, name, "", "", "", 0, 0, &user}

	// Get Gravatar Icon
	emailHash := gr.EmailHash(user.Email)
	gravatarUrl := gr.GetAvatarURL("https", emailHash, gr.DefaultIdentIcon, 128)

	if gravatarUrl != nil {
		profile.PhotoUrl = gravatarUrl.String()
	}

	err = c.Txn.Insert(profile)
	if err != nil {
		panic(err)
	}

	// Send out confirmation email
	err = c.sendAccountConfirmEmail(&user)

	if err != nil {
		c.Flash.Error("Could not send confirmation email")
		fmt.Println(err.Error())
	}

	c.Session["userEmail"] = string(user.Email)
	c.Flash.Success("Welcome, " + profile.Name)
	return c.Redirect(routes.Profile.Show(profile.UserName))
}

func (c Account) Login() r.Result {
	if c.connected() != nil {
		return c.Logout()
	}

	hasEmailCapability := hasEmailCapability()

	return c.Render(hasEmailCapability)
}

func (c Account) DoLogin(user *models.User, remember bool) {
	c.Session["userEmail"] = user.Email
	if remember {
		c.Session.SetDefaultExpiration()
	} else {
		c.Session.SetNoExpiration()
	}
}

func (c Account) LoginAccount(account, password string, remember bool) r.Result {

	var profile *models.Profile

	// If account is a valid email address, retrieve account by email
	// otherwise, retrieve account by username
	models.ValidateUserEmail(c.Validation, account).Key("account")

	if c.Validation.HasErrors() {
		c.Validation.Clear()
		profile = c.getProfileByUserName(account)
	} else {
		profile = c.getProfileByEmailAddress(account)
	}

	if profile != nil {
		err := bcrypt.CompareHashAndPassword(profile.User.HashedPassword, []byte(password))
		if err == nil {
			c.DoLogin(profile.User, remember)
			c.Flash.Success("Welcome back, " + profile.Name)
			return c.Redirect(routes.Profile.Show(profile.UserName))
		}
	}

	c.Flash.Error("Sign In failed.")
	return c.Redirect(routes.Account.Login())
}

func (c Account) Recover() r.Result {
	hasEmailCapability := hasEmailCapability()

	if hasEmailCapability == false {
		return c.Redirect(routes.Account.Login())
	}

	return c.Render()
}

func (c Account) RetrieveAccount(email string) r.Result {

	models.ValidateUserEmail(c.Validation, email).Key("email")

	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		c.Flash.Error("Account recovery failed.")
		return c.Redirect(routes.Account.Recover())
	}

	profile := c.getProfileByEmailAddress(email)

	if profile == nil {
		// Return a false positive response to requestors at this point
		c.Flash.Success("A password reset request has been sent to " + email + ".")
		return c.Redirect(routes.Application.Index())
	}

	err := c.sendAccountRecoverEmail(profile.User)

	if err != nil {
		panic(err)
	}

	c.Flash.Success("A password reset request has been sent to " + email + ".")
	return c.Redirect(routes.Account.Login())
}


func (c Account) ConfirmEmail(token string) r.Result {

	existingToken := c.getVerifyHashRecord("confirm", token)

	if existingToken == nil {
		c.Flash.Error("Token invalid or used");
		return c.Redirect(routes.Application.Index())
	}

	existingProfile := c.getProfileByEmailAddress(existingToken.Email)

	if existingProfile == nil {
		c.Flash.Error("Token invalid or used");
		return c.Redirect(routes.Application.Index())
	}

	// Update user record to indicate registered email address has been confirmed
	existingProfile.User.Confirmed = true

	_, err := c.Txn.Update(existingProfile.User)
	if err != nil {
		panic(err)
	}

	// Delete used token
	c.deleteVerifyHashRecord(existingToken)

	return c.Render(existingProfile)
}

func (c Account) PasswordReset(token string) r.Result {

	existingToken := c.getVerifyHashRecord("reset", token)

	if existingToken == nil {
		c.Flash.Error("Token invalid or used");
		return c.Redirect(routes.Application.Index())
	}

	existingProfile := c.getProfileByEmailAddress(existingToken.Email)

	if existingProfile == nil {
		c.Flash.Error("Token invalid or used");
		return c.Redirect(routes.Application.Index())
	}

	// Delete used token
	c.deleteVerifyHashRecord(existingToken)

	// Log user in, flash message and redirect to change password page
	c.DoLogin(existingProfile.User, false)
	c.Flash.Success("Please now enter a new password")
	return c.Redirect(routes.Profile.Password(existingProfile.UserName))
}

func (c Account) Logout() r.Result {
	for k := range c.Session {
		delete(c.Session, k)
	}
	if c.Flash.Data["error"] == "" {
		c.Flash.Success("You have been successfully logged out")
	} else {
		c.Flash.Error(c.Flash.Data["error"]) // propagate error
	}
	return c.Redirect(routes.Account.Login())
}

func (c Account) CheckUserName(username, currentUsername string) r.Result {

	if(username == currentUsername) {
		return c.RenderText("true");
	}

	// Validate Profile components
	models.ValidateProfileUserName(c.Validation, username).Key("username")

	if c.Validation.HasErrors() {
		c.Validation.Clear()
		return c.RenderText("false");
	}

	userExists := c.getProfileByUserName(username)

	if userExists != nil {
		return c.RenderText("false");
	}

	// User name is available
	return c.RenderText("true");
}

func (c Account) CheckEmail(email, currentEmail string) r.Result {

	if(email == currentEmail) {
		return c.RenderText("true");
	}

	userExists := c.getProfileByEmailAddress(email)

	if userExists != nil {
		return c.RenderText("false");
	}

	// Email address is new (and therefore valid)
	return c.RenderText("true");
}

// VERIFICATION HASH STUFF

func (e Account) generateVerifyHash(n int) []byte {
	bytes := make([]byte, n)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alphanum[b % byte(len(alphanum))]
	}
	return bytes
}

func (e Account) storeVerifyHashRecord(email string, tokenType string, hash []byte) error {
	err := e.Txn.Insert( &models.Token{0, email, tokenType, string(hash)} )
	return err
}

func (e Account) getVerifyHashRecord(hashType string, hashToken string) *models.Token {
	token := &models.Token{}
	err := e.Txn.SelectOne(token, `select * from Token where Type = ? and Hash = ?`, hashType, hashToken)
	if err != nil {
		return nil
	}
	if token == nil {
		return nil
	}
	return token
}

func (e Account) deleteVerifyHashRecord(token *models.Token) {
	_, err := e.Txn.Delete(token)
	if err != nil {
		panic(err)
	}
}

// EMAIL STUFF

// Check if we have capability to send emails
func hasEmailCapability() bool {
	mailerServer := r.Config.StringDefault("mailer.server", "smtp.example.org")
	if mailerServer != "smtp.example.org" {
		return true
	}
	return false
}

func (e Account) sendAccountRecoverEmail(user *models.User) error {
	host := r.Config.StringDefault("http.addr", "localhost")
	return e.sendEmail(user, "reset", "Reset your password at " + host)
}

func (e Account) sendAccountConfirmEmail(user *models.User) error {
	host := r.Config.StringDefault("http.addr", "localhost")
	return e.sendEmail(user, "confirm", "Welcome to " + host)
}

func (e Account) sendEmail(user *models.User, verifyType, subject string) error {

	var (
		mailerServer    = r.Config.StringDefault("mailer.server", "smtp.example.org")
		mailerPort      = r.Config.IntDefault("mailer.port", 25)
		mailerUsername  = r.Config.StringDefault("mailer.username", "<username>")
		mailerPassword  = r.Config.StringDefault("mailer.password", "<password>")
		mailerFromAddr  = r.Config.StringDefault("mailer.fromaddress", "no-reply@example.org")
		mailerReplyAddr = r.Config.StringDefault("mailer.replyaddress", "support@example.org")
		callbackHost = r.Config.StringDefault("http.host", "http://localhost:9000")
	)

	// If mail has not been configured, don't try to send a confirmation email
	if mailerServer == "smtp.example.org" {
		return nil
	}

	mailer := m.Mailer{Server: mailerServer, Port: mailerPort, UserName: mailerUsername, Password: mailerPassword}
	mailer.Sender = &m.Sender{From: mailerFromAddr, ReplyTo: mailerReplyAddr}

	// Generate a new token and store against the user's id
	verifyEmailToken := e.generateVerifyHash(16);

	e.storeVerifyHashRecord(user.Email, verifyType, verifyEmailToken);

	// arguments used for template rendering
	var args = make(map[string]interface{})

	args["user"] = user

	args["Link"] = callbackHost + "/account/" + verifyType + "/" + string(verifyEmailToken)

	message := &m.Message{To: []string{user.Email}, Subject: subject}

	rErr := message.RenderTemplate("email/" + verifyType, args)
	if rErr != nil {
		return rErr
	}

	sErr := mailer.SendMessage(message)
	if sErr != nil {
		return sErr
	}

	fmt.Println("Mail sent to " + string(user.Email))

	return nil
}
