package tests

import (
	"fmt"
	"github.com/revel/revel"
	"github.com/richtr/baseapp/app/models"
	"github.com/richtr/baseapp/app/routes"
	"net/url"
)

type ProfileTest struct {
	revel.TestSuite
}

func (t *ProfileTest) Before() {
	// Runs before each test below is executed
	t.DoLogin(demoUser)
}

func (t *ProfileTest) After() {
	// Runs after each test below is executed
	t.DoLogout()
}

var demoUser = models.User{
	UserId:   1,
	Email:    "demo@demo.com",
	Password: "demouser",
}

var demoProfile = models.Profile{
	ProfileId:   1,
	UserId:      demoUser.UserId,
	UserName:    "demouser",
	Name:        "Demo User",
	Summary:     "Just a regular guy",
	Description: "...",
	PhotoUrl:    "http://www.gravatar.com/avatar/53444f91e698c0c7caa2dbc3bdbf93fc?s=128&d=identicon",
	User:        &demoUser,
}

func (t *ProfileTest) DoLogin(user models.User) {
	urlValues := url.Values{}
	urlValues.Add("account", user.Email)
	urlValues.Add("password", user.Password)
	urlValues.Add("remember", "0")

	t.PostForm("/account/login", urlValues)
}

func (t *ProfileTest) DoLogout() {
	t.Get(routes.Account.Logout())
}

func (t *ProfileTest) DoUpdateSettings(profileUserName string, profile models.Profile, verifyPassword string) {
	urlValues := url.Values{}
	urlValues.Add("profile.Name", profile.Name)
	urlValues.Add("profile.User.Email", profile.User.Email)
	urlValues.Add("profile.User.Password", profile.User.Password)
	urlValues.Add("verifyPassword", verifyPassword)
	urlValues.Add("profile.UserName", profile.UserName)
	urlValues.Add("profile.Summary", profile.Summary)
	urlValues.Add("profile.Description", profile.Description)
	urlValues.Add("profile.PhotoUrl", profile.PhotoUrl)

	postUrl := fmt.Sprintf("/%s/edit", profileUserName)

	t.PostForm(postUrl, urlValues)
}

func (t *ProfileTest) TestProfileDisplayPage_Own() {
	t.Get(routes.Profile.Show(demoProfile.UserName))

	t.AssertOk()
	t.AssertContains("Just a regular guy") // Part of demo user profile

	t.AssertContains("Edit Profile</a>")
	t.AssertContains("Create New Post</a>")
}

func (t *ProfileTest) TestProfileDisplayPage_Other() {
	t.Get(routes.Profile.Show("otheruser")) // access other user's profile by their profile id

	t.AssertOk()
	t.AssertContains("Just another regular guy") // Part of 'demo user 1' profile
}

func (t *ProfileTest) TestProfileSettings_Access_Success() {
	t.Get(routes.Profile.Settings(demoProfile.UserName))

	t.AssertOk()
	t.AssertContains("Edit Profile</h2>")
	t.AssertContains(" value=\"Update Profile\">")
}

func (t *ProfileTest) TestProfileSettings_Access_Error() {
	// No log in
	t.DoLogout()

	t.Get(routes.Profile.Settings(demoProfile.UserName))

	t.AssertOk()
	t.AssertContains("You must log in to access your account")
}

func (t *ProfileTest) TestOtherProfileSettings_Access_Denied() {
	t.Get(routes.Profile.Settings("otheruser")) // try to access edit page of other user's profile

	t.AssertOk()
	t.AssertContains("You must log in to access your account")
}

func (t *ProfileTest) TestProfileSettings_EmptyEmail() {
	demoProfileUpdate := demoProfile
	user := *demoProfileUpdate.User
	user.Email = ""
	demoProfileUpdate.User = &user

	t.DoUpdateSettings(demoProfile.UserName, demoProfileUpdate, demoUser.Password)

	t.AssertOk()
	t.AssertContains("Email address required")
}

func (t *ProfileTest) TestProfileSettings_EmptyPasswordAndVerifyPassword() {
	demoProfileUpdate := demoProfile
	user := *demoProfileUpdate.User
	user.Password = ""
	demoProfileUpdate.User = &user

	t.DoUpdateSettings(demoProfile.UserName, demoProfileUpdate, "")

	t.AssertOk()
	t.AssertContains("Profile has been updated") // though no changes have been made!
}

func (t *ProfileTest) TestProfileSettings_InvalidEmail() {
	demoProfileUpdate := demoProfile
	user := *demoProfileUpdate.User
	user.Email = "myemailaddress"
	demoProfileUpdate.User = &user

	t.DoUpdateSettings(demoProfile.UserName, demoProfileUpdate, demoUser.Password)

	t.AssertOk()
	t.AssertContains("You must provide a valid email address")
}

func (t *ProfileTest) TestProfileSettings_TooShortEmail() {
	demoProfileUpdate := demoProfile
	user := *demoProfileUpdate.User
	user.Email = "a@b.c"
	demoProfileUpdate.User = &user

	t.DoUpdateSettings(demoProfile.UserName, demoProfileUpdate, demoUser.Password)

	t.AssertOk()
	t.AssertContains("Email address can not be less than 6 characters")
}

func (t *ProfileTest) TestProfileSettings_TooLongEmail() {
	demoProfileUpdate := demoProfile
	user := *demoProfileUpdate.User
	user.Email = "myemailaddress_myemailaddress_myemailaddress_myemailaddress_myemailaddress_myemailaddress_myemailaddress_myemailaddress_myemailaddress_myemailaddress_myemailaddress_myemailaddress_myemailaddress@gmail.com"
	demoProfileUpdate.User = &user

	t.DoUpdateSettings(demoProfile.UserName, demoProfileUpdate, demoUser.Password)

	t.AssertOk()
	t.AssertContains("Email address can not exceed 200 characters")
}

func (t *ProfileTest) TestProfileSettings_EmptyName() {
	demoProfileUpdate := demoProfile
	demoProfileUpdate.Name = ""

	t.DoUpdateSettings(demoProfile.UserName, demoProfileUpdate, demoUser.Password)

	t.AssertOk()
	t.AssertContains("Name required")
}

func (t *ProfileTest) TestProfileSettings_TooShortName() {
	demoProfileUpdate := demoProfile
	demoProfileUpdate.Name = "User1"

	t.DoUpdateSettings(demoProfile.UserName, demoProfileUpdate, demoUser.Password)

	t.AssertOk()
	t.AssertContains("Name must be at least 6 characters")
}

func (t *ProfileTest) TestProfileSettings_TooLongName() {
	demoProfileUpdate := demoProfile
	demoProfileUpdate.Name = "Test User 1 Test User 1 Test User 1 Test User 1 Test User 1 Test User 1 Test User 1 Test User 1 Test User 1"

	t.DoUpdateSettings(demoProfile.UserName, demoProfileUpdate, demoUser.Password)

	t.AssertOk()
	t.AssertContains("Name must be at most 100 characters")
}

func (t *ProfileTest) TestProfileSettings_EmptyPassword() {
	demoProfileUpdate := demoProfile
	user := *demoProfileUpdate.User
	verifyPassword := user.Password
	user.Password = ""
	demoProfileUpdate.User = &user

	t.DoUpdateSettings(demoProfile.UserName, demoProfileUpdate, verifyPassword)

	t.AssertOk()
	t.AssertContains("Password required")
}

func (t *ProfileTest) TestProfileSettings_EmptyVerifyPassword() {
	demoProfileUpdate := demoProfile
	user := *demoProfileUpdate.User
	demoProfileUpdate.User = &user

	t.DoUpdateSettings(demoProfile.UserName, demoProfileUpdate, "")

	t.AssertOk()
	t.AssertContains("Password verification required")
}

func (t *ProfileTest) TestProfileSettings_TooShortPassword() {
	demoProfileUpdate := demoProfile
	user := *demoProfileUpdate.User
	user.Password = "pw1"
	demoProfileUpdate.User = &user

	t.DoUpdateSettings(demoProfile.UserName, demoProfileUpdate, user.Password)

	t.AssertOk()
	t.AssertContains("Password must be at least 6 characters")
}

func (t *ProfileTest) TestProfileSettings_TooLongPassword() {
	demoProfileUpdate := demoProfile
	user := *demoProfileUpdate.User
	user.Password = "testuser1testuser1testuser1testuser1testuser1testuser1testuser1testuser1testuser1testuser1testuser1testuser1testuser1testuser1testuser1testuser1testuser1testuser1testuser1testuser1testuser1testuser1testuser1"
	demoProfileUpdate.User = &user

	t.DoUpdateSettings(demoProfile.UserName, demoProfileUpdate, user.Password)

	t.AssertOk()
	t.AssertContains("Password must be at most 200 characters")
}

func (t *ProfileTest) TestProfileSettings_PasswordMismatch() {
	demoProfileUpdate := demoProfile
	user := *demoProfileUpdate.User
	user.Password = "testuser1"
	demoProfileUpdate.User = &user

	t.DoUpdateSettings(demoProfile.UserName, demoProfileUpdate, "testuser2")

	t.AssertOk()
	t.AssertContains("Provided passwords do not match")
}

func (t *ProfileTest) TestProfileSettings_PasswordSameAsEmail() {
	demoProfileUpdate := demoProfile
	user := *demoProfileUpdate.User
	user.Password = user.Email
	demoProfileUpdate.User = &user

	t.DoUpdateSettings(demoProfile.UserName, demoProfileUpdate, user.Password)

	t.AssertOk()
	t.AssertContains("Password cannot be the same as your email address")
}
