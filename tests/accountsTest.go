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

func (t *AccountsTest) TestLoginSuccess() {
  loginData := url.Values{}
  loginData.Add("email", "demo@demo.com")
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
  loginData.Add("email", "demo@demo.com")
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
  loginData.Add("email", "")
  loginData.Add("password", "demouser")
  loginData.Add("remember", "0")

  t.PostForm("/account/login", loginData)

  t.AssertOk()
  t.AssertContains("Sign In failed")
}

func (t *AccountsTest) TestLoginFail_EmptyPassword() {
  loginData := url.Values{}
  loginData.Add("email", "demo@demo.com")
  loginData.Add("password", "")
  loginData.Add("remember", "0")

  t.PostForm("/account/login", loginData)

  t.AssertOk()
  t.AssertContains("Sign In failed")
}

func (t *AccountsTest) TestSignupFail_EmptyEmail() {
  registerData := url.Values{}
  registerData.Add("user.Email", "")
  registerData.Add("name", "Test User 1")
  registerData.Add("user.Password", "testuser1")
  registerData.Add("verifyPassword", "testuser1")

  t.PostForm("/account/register", registerData)

  t.AssertOk()
  t.AssertContains("Email address required")
}

func (t *AccountsTest) TestSignupFail_InvalidEmail() {
  registerData := url.Values{}
  registerData.Add("user.Email", "myemailaddress")
  registerData.Add("name", "Test User 1")
  registerData.Add("user.Password", "testuser1")
  registerData.Add("verifyPassword", "testuser1")

  t.PostForm("/account/register", registerData)

  t.AssertOk()
  t.AssertContains("You must provide a valid email address")
}

func (t *AccountsTest) TestSignupFail_TooShortEmail() {
  registerData := url.Values{}
  registerData.Add("user.Email", "a@b.c")
  registerData.Add("name", "Test User 1")
  registerData.Add("user.Password", "testuser1")
  registerData.Add("verifyPassword", "testuser1")

  t.PostForm("/account/register", registerData)

  t.AssertOk()
  t.AssertContains("Email address can not be less than 6 characters")
}

func (t *AccountsTest) TestSignupFail_TooLongEmail() {
  registerData := url.Values{}
  registerData.Add("user.Email", "myemailaddress_myemailaddress_myemailaddress_myemailaddress_myemailaddress_myemailaddress_myemailaddress_myemailaddress_myemailaddress_myemailaddress_myemailaddress_myemailaddress_myemailaddress@gmail.com")
  registerData.Add("name", "Test User 1")
  registerData.Add("user.Password", "testuser1")
  registerData.Add("verifyPassword", "testuser1")

  t.PostForm("/account/register", registerData)

  t.AssertOk()
  t.AssertContains("Email address can not exceed 200 characters")
}

func (t *AccountsTest) TestSignupFail_EmptyName() {
  registerData := url.Values{}
  registerData.Add("user.Email", "testuser@example.com")
  registerData.Add("name", "")
  registerData.Add("user.Password", "testuser1")
  registerData.Add("verifyPassword", "testuser1")

  t.PostForm("/account/register", registerData)

  t.AssertOk()
  t.AssertContains("Name required")
}

func (t *AccountsTest) TestSignupFail_TooShortName() {
  registerData := url.Values{}
  registerData.Add("user.Email", "testuser@example.com")
  registerData.Add("name", "User1")
  registerData.Add("user.Password", "testuser1")
  registerData.Add("verifyPassword", "testuser1")

  t.PostForm("/account/register", registerData)

  t.AssertOk()
  t.AssertContains("Name must be at least 6 characters")
}


func (t *AccountsTest) TestSignupFail_TooLongName() {
  registerData := url.Values{}
  registerData.Add("user.Email", "testuser@example.com")
  registerData.Add("name", "Test User 1 Test User 1 Test User 1 Test User 1 Test User 1 Test User 1 Test User 1 Test User 1 Test User 1")
  registerData.Add("user.Password", "testuser1")
  registerData.Add("verifyPassword", "testuser1")

  t.PostForm("/account/register", registerData)

  t.AssertOk()
  t.AssertContains("Name must be at most 100 characters")
}

func (t *AccountsTest) TestSignupFail_EmptyPassword() {
  registerData := url.Values{}
  registerData.Add("user.Email", "testuser@example.com")
  registerData.Add("name", "Test User 1")
  registerData.Add("user.Password", "")
  registerData.Add("verifyPassword", "testuser1")

  t.PostForm("/account/register", registerData)

  t.AssertOk()
  t.AssertContains("Password required")
}

func (t *AccountsTest) TestSignupFail_EmptyVerifyPassword() {
  registerData := url.Values{}
  registerData.Add("user.Email", "testuser@example.com")
  registerData.Add("name", "Test User 1")
  registerData.Add("user.Password", "testuser1")
  registerData.Add("verifyPassword", "")

  t.PostForm("/account/register", registerData)

  t.AssertOk()
  t.AssertContains("Password verification required")
}

func (t *AccountsTest) TestSignupFail_TooShortPassword() {
  registerData := url.Values{}
  registerData.Add("user.Email", "testuser@example.com")
  registerData.Add("name", "Test User 1")
  registerData.Add("user.Password", "pw1")
  registerData.Add("verifyPassword", "pw1")

  t.PostForm("/account/register", registerData)

  t.AssertOk()
  t.AssertContains("Password must be at least 6 characters")
}

func (t *AccountsTest) TestSignupFail_TooLongPassword() {
  registerData := url.Values{}
  registerData.Add("user.Email", "testuser@example.com")
  registerData.Add("name", "Test User 1")
  registerData.Add("user.Password", "testuser1testuser1testuser1testuser1")
  registerData.Add("verifyPassword", "testuser1testuser1testuser1testuser1")

  t.PostForm("/account/register", registerData)

  t.AssertOk()
  t.AssertContains("Password must be at most 15 characters")
}

func (t *AccountsTest) TestSignupFail_PasswordMismatch() {
  registerData := url.Values{}
  registerData.Add("user.Email", "testuser@example.com")
  registerData.Add("name", "Test User 1")
  registerData.Add("user.Password", "testuser1")
  registerData.Add("verifyPassword", "testuser2")

  t.PostForm("/account/register", registerData)

  t.AssertOk()
  t.AssertContains("Provided passwords do not match")
}

func (t *AccountsTest) TestSignupFail_PasswordSameAsEmail() {
  registerData := url.Values{}
  registerData.Add("user.Email", "test@foo.com")
  registerData.Add("name", "Test User 1")
  registerData.Add("user.Password", "test@foo.com")
  registerData.Add("verifyPassword", "test@foo.com")

  t.PostForm("/account/register", registerData)

  t.AssertOk()
  t.AssertContains("Password cannot be the same as your email address")
}

func (t *AccountsTest) TestSignupSuccess() {
  registerData := url.Values{}
  registerData.Add("user.Email", "testuser@example.com")
  registerData.Add("name", "Test User 1")
  registerData.Add("user.Password", "testuser1")
  registerData.Add("verifyPassword", "testuser1")

  t.PostForm("/account/register", registerData)

  t.AssertOk()
  t.AssertContains("Welcome, Test User 1")

  // Now log out
  t.Get(routes.Account.Logout())
}
