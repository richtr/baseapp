package tests

import (
	"github.com/revel/revel"
	"github.com/richtr/baseapp/app/routes"
	"net/url"
)

type AccountsTest struct {
	revel.TestSuite
}

func (t *AccountsTest) Before() {
	// Runs before each test below is executed
}

func (t *AccountsTest) After() {
	// Runs after each test below is executed
	t.Get(routes.Account.Logout())
}

func (t *AccountsTest) TestLoginSuccess_ByEmail() {
	loginData := url.Values{}
	loginData.Add("account", "demo@demo.com")
	loginData.Add("password", "demouser")
	loginData.Add("remember", "0")

	t.PostForm("/account/login", loginData)

	// Part of login information (after redirect to profile page)
	t.AssertOk()
	t.AssertContains("Welcome back, Demo User")

	// Do Logout
	t.Get(routes.Account.Logout())
}

func (t *AccountsTest) TestLoginSuccess_ByUsername() {
	loginData := url.Values{}
	loginData.Add("account", "demouser")
	loginData.Add("password", "demouser")
	loginData.Add("remember", "0")

	t.PostForm("/account/login", loginData)

	// Part of login information (after redirect to profile page)
	t.AssertOk()
	t.AssertContains("Welcome back, Demo User")

	// Do Logout
	t.Get(routes.Account.Logout())
}

func (t *AccountsTest) TestLogoutSuccess() {
	// Log in
	loginData := url.Values{}
	loginData.Add("account", "demo@demo.com")
	loginData.Add("password", "demouser")
	loginData.Add("remember", "0")

	t.PostForm("/account/login", loginData)

	// Log out
	t.Get(routes.Account.Logout())

	t.AssertOk()
	t.AssertContains("You have been successfully logged out")
}

func (t *AccountsTest) TestLoginFail_EmptyEmail() {
	loginData := url.Values{}
	loginData.Add("account", "")
	loginData.Add("password", "demouser")
	loginData.Add("remember", "0")

	t.PostForm("/account/login", loginData)

	t.AssertOk()
	t.AssertContains("Sign In failed")
}

func (t *AccountsTest) TestLoginFail_EmptyPassword() {
	loginData := url.Values{}
	loginData.Add("account", "demo@demo.com")
	loginData.Add("password", "")
	loginData.Add("remember", "0")

	t.PostForm("/account/login", loginData)

	t.AssertOk()
	t.AssertContains("Sign In failed")
}

func (t *AccountsTest) TestSignupFail_EmptyEmail() {
	registerData := url.Values{}
	registerData.Add("user.Email", "")
	registerData.Add("username", "testuser1")
	registerData.Add("name", "Test User 1")
	registerData.Add("user.Password", "testuser1")

	t.PostForm("/account/register", registerData)

	t.AssertOk()
	t.AssertContains("Email address required")
}

func (t *AccountsTest) TestSignupFail_InvalidEmail() {
	registerData := url.Values{}
	registerData.Add("user.Email", "myemailaddress")
	registerData.Add("username", "testuser1")
	registerData.Add("name", "Test User 1")
	registerData.Add("user.Password", "testuser1")

	t.PostForm("/account/register", registerData)

	t.AssertOk()
	t.AssertContains("You must provide a valid email address")
}

func (t *AccountsTest) TestSignupFail_TooShortEmail() {
	registerData := url.Values{}
	registerData.Add("user.Email", "a@b.c")
	registerData.Add("username", "testuser1")
	registerData.Add("name", "Test User 1")
	registerData.Add("user.Password", "testuser1")

	t.PostForm("/account/register", registerData)

	t.AssertOk()
	t.AssertContains("Email address can not be less than 6 characters")
}

func (t *AccountsTest) TestSignupFail_TooLongEmail() {
	registerData := url.Values{}
	registerData.Add("user.Email", "myemailaddress_myemailaddress_myemailaddress_myemailaddress_myemailaddress_myemailaddress_myemailaddress_myemailaddress_myemailaddress_myemailaddress_myemailaddress_myemailaddress_myemailaddress@gmail.com")
	registerData.Add("username", "testuser1")
	registerData.Add("name", "Test User 1")
	registerData.Add("user.Password", "testuser1")

	t.PostForm("/account/register", registerData)

	t.AssertOk()
	t.AssertContains("Email address can not exceed 200 characters")
}

func (t *AccountsTest) TestSignupFail_EmptyUsername() {
	registerData := url.Values{}
	registerData.Add("user.Email", "testuser@example.com")
	registerData.Add("username", "")
	registerData.Add("name", "Test User 1")
	registerData.Add("user.Password", "testuser1")

	t.PostForm("/account/register", registerData)

	t.AssertOk()
	t.AssertContains("User name required")
}

func (t *AccountsTest) TestSignupFail_InvalidUsername1() {
	registerData := url.Values{}
	registerData.Add("user.Email", "testuser@example.com")
	registerData.Add("username", "@")
	registerData.Add("name", "Test User 1")
	registerData.Add("user.Password", "testuser1")

	t.PostForm("/account/register", registerData)

	t.AssertOk()
	t.AssertContains("Invalid User name. Alphanumerics allowed only")
}

func (t *AccountsTest) TestSignupFail_InvalidUsername2() {
	registerData := url.Values{}
	registerData.Add("user.Email", "testuser@example.com")
	registerData.Add("username", "test_user-1")
	registerData.Add("name", "Test User 1")
	registerData.Add("user.Password", "testuser1")

	t.PostForm("/account/register", registerData)

	t.AssertOk()
	t.AssertContains("Invalid User name. Alphanumerics allowed only")
}

func (t *AccountsTest) TestSignupFail_TooLongUsername() {
	registerData := url.Values{}
	registerData.Add("user.Email", "testuser@example.com")
	registerData.Add("username", "testuser1testuser1testuser1testuser1testuser1testuser1testuser1testuser1")
	registerData.Add("name", "Test User 1")
	registerData.Add("user.Password", "testuser1")

	t.PostForm("/account/register", registerData)

	t.AssertOk()
	t.AssertContains("User name can not exceed 64 characters")
}

func (t *AccountsTest) TestSignupFail_ReservedUsername() {
	registerData := url.Values{}
	registerData.Add("user.Email", "testuser@example.com")
	registerData.Add("username", "account")
	registerData.Add("name", "Test User 1")
	registerData.Add("user.Password", "testuser1")

	t.PostForm("/account/register", registerData)

	t.AssertOk()
	t.AssertContains("Invalid User name. Reserved keywords not allowed")
}

func (t *AccountsTest) TestSignupFail_EmptyName() {
	registerData := url.Values{}
	registerData.Add("user.Email", "testuser@example.com")
	registerData.Add("username", "testuser1")
	registerData.Add("name", "")
	registerData.Add("user.Password", "testuser1")

	t.PostForm("/account/register", registerData)

	t.AssertOk()
	t.AssertContains("Name required")
}

func (t *AccountsTest) TestSignupFail_TooShortName() {
	registerData := url.Values{}
	registerData.Add("user.Email", "testuser@example.com")
	registerData.Add("username", "testuser1")
	registerData.Add("name", "User1")
	registerData.Add("user.Password", "testuser1")

	t.PostForm("/account/register", registerData)

	t.AssertOk()
	t.AssertContains("Name must be at least 6 characters")
}


func (t *AccountsTest) TestSignupFail_TooLongName() {
	registerData := url.Values{}
	registerData.Add("user.Email", "testuser@example.com")
	registerData.Add("username", "testuser1")
	registerData.Add("name", "Test User 1 Test User 1 Test User 1 Test User 1 Test User 1 Test User 1 Test User 1 Test User 1 Test User 1")
	registerData.Add("user.Password", "testuser1")

	t.PostForm("/account/register", registerData)

	t.AssertOk()
	t.AssertContains("Name must be at most 100 characters")
}

func (t *AccountsTest) TestSignupFail_EmptyPassword() {
	registerData := url.Values{}
	registerData.Add("user.Email", "testuser@example.com")
	registerData.Add("username", "testuser1")
	registerData.Add("name", "Test User 1")
	registerData.Add("user.Password", "")

	t.PostForm("/account/register", registerData)

	t.AssertOk()
	t.AssertContains("Password required")
}

func (t *AccountsTest) TestSignupFail_TooShortPassword() {
	registerData := url.Values{}
	registerData.Add("user.Email", "testuser@example.com")
	registerData.Add("username", "testuser1")
	registerData.Add("name", "Test User 1")
	registerData.Add("user.Password", "pw1")

	t.PostForm("/account/register", registerData)

	t.AssertOk()
	t.AssertContains("Password must be at least 6 characters")
}

func (t *AccountsTest) TestSignupFail_TooLongPassword() {
	registerData := url.Values{}
	registerData.Add("user.Email", "testuser@example.com")
	registerData.Add("username", "testuser1")
	registerData.Add("name", "Test User 1")
	registerData.Add("user.Password", "testuser1testuser1testuser1testuser1testuser1testuser1testuser1testuser1testuser1testuser1testuser1testuser1testuser1testuser1testuser1testuser1testuser1testuser1testuser1testuser1testuser1testuser1testuser1")

	t.PostForm("/account/register", registerData)

	t.AssertOk()
	t.AssertContains("Password must be at most 200 characters")
}

func (t *AccountsTest) TestSignupFail_PasswordSameAsEmail() {
	registerData := url.Values{}
	registerData.Add("user.Email", "test@foo.com")
	registerData.Add("username", "testuser1")
	registerData.Add("name", "Test User 1")
	registerData.Add("user.Password", "test@foo.com")

	t.PostForm("/account/register", registerData)

	t.AssertOk()
	t.AssertContains("Password cannot be the same as your email address")
}

func (t *AccountsTest) TestSignupFail_PasswordSameAsUsername() {
	registerData := url.Values{}
	registerData.Add("user.Email", "test@foo.com")
	registerData.Add("username", "testuser1")
	registerData.Add("name", "Test User 1")
	registerData.Add("user.Password", "testuser1")

	t.PostForm("/account/register", registerData)

	t.AssertOk()
	t.AssertContains("Password cannot be the same as your user name")
}

func (t *AccountsTest) TestSignupSuccess() {
	registerData := url.Values{}
	registerData.Add("user.Email", "testuser@example.com")
	registerData.Add("username", "testuser1")
	registerData.Add("name", "Test User 1")
	registerData.Add("user.Password", "pass@testuser1")

	t.PostForm("/account/register", registerData)

	t.AssertOk()
	t.AssertContains("Welcome, Test User 1")

	// Now log out
	t.Get(routes.Account.Logout())
}
