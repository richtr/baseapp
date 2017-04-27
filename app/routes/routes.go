// GENERATED CODE - DO NOT EDIT
package routes

import "github.com/revel/revel"


type tApplication struct {}
var Application tApplication


func (_ tApplication) Index(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("Application.Index", args).URL
}

func (_ tApplication) About(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("Application.About", args).URL
}

func (_ tApplication) Contact(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("Application.Contact", args).URL
}

func (_ tApplication) Search(
		query string,
		page int,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "query", query)
	revel.Unbind(args, "page", page)
	return revel.MainRouter.Reverse("Application.Search", args).URL
}

func (_ tApplication) SwitchToDesktop(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("Application.SwitchToDesktop", args).URL
}

func (_ tApplication) SwitchToMobile(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("Application.SwitchToMobile", args).URL
}


type tGorpController struct {}
var GorpController tGorpController


func (_ tGorpController) Begin(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("GorpController.Begin", args).URL
}

func (_ tGorpController) Commit(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("GorpController.Commit", args).URL
}

func (_ tGorpController) Rollback(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("GorpController.Rollback", args).URL
}


type tJobs struct {}
var Jobs tJobs


func (_ tJobs) Status(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("Jobs.Status", args).URL
}


type tStatic struct {}
var Static tStatic


func (_ tStatic) Serve(
		prefix string,
		filepath string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "prefix", prefix)
	revel.Unbind(args, "filepath", filepath)
	return revel.MainRouter.Reverse("Static.Serve", args).URL
}

func (_ tStatic) ServeModule(
		moduleName string,
		prefix string,
		filepath string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "moduleName", moduleName)
	revel.Unbind(args, "prefix", prefix)
	revel.Unbind(args, "filepath", filepath)
	return revel.MainRouter.Reverse("Static.ServeModule", args).URL
}


type tTestRunner struct {}
var TestRunner tTestRunner


func (_ tTestRunner) Index(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("TestRunner.Index", args).URL
}

func (_ tTestRunner) Suite(
		suite string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "suite", suite)
	return revel.MainRouter.Reverse("TestRunner.Suite", args).URL
}

func (_ tTestRunner) Run(
		suite string,
		test string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "suite", suite)
	revel.Unbind(args, "test", test)
	return revel.MainRouter.Reverse("TestRunner.Run", args).URL
}

func (_ tTestRunner) List(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("TestRunner.List", args).URL
}


type tAccount struct {}
var Account tAccount


func (_ tAccount) AddUser(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("Account.AddUser", args).URL
}

func (_ tAccount) AddAppName(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("Account.AddAppName", args).URL
}

func (_ tAccount) AddRenderMode(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("Account.AddRenderMode", args).URL
}

func (_ tAccount) Index(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("Account.Index", args).URL
}

func (_ tAccount) Register(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("Account.Register", args).URL
}

func (_ tAccount) SaveUser(
		user interface{},
		username string,
		name string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "user", user)
	revel.Unbind(args, "username", username)
	revel.Unbind(args, "name", name)
	return revel.MainRouter.Reverse("Account.SaveUser", args).URL
}

func (_ tAccount) Login(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("Account.Login", args).URL
}

func (_ tAccount) LoginAccount(
		account string,
		password string,
		remember bool,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "account", account)
	revel.Unbind(args, "password", password)
	revel.Unbind(args, "remember", remember)
	return revel.MainRouter.Reverse("Account.LoginAccount", args).URL
}

func (_ tAccount) Recover(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("Account.Recover", args).URL
}

func (_ tAccount) RetrieveAccount(
		email string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "email", email)
	return revel.MainRouter.Reverse("Account.RetrieveAccount", args).URL
}

func (_ tAccount) ConfirmEmail(
		token string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "token", token)
	return revel.MainRouter.Reverse("Account.ConfirmEmail", args).URL
}

func (_ tAccount) PasswordReset(
		token string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "token", token)
	return revel.MainRouter.Reverse("Account.PasswordReset", args).URL
}

func (_ tAccount) Logout(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("Account.Logout", args).URL
}

func (_ tAccount) CheckUserName(
		username string,
		currentUsername string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "username", username)
	revel.Unbind(args, "currentUsername", currentUsername)
	return revel.MainRouter.Reverse("Account.CheckUserName", args).URL
}

func (_ tAccount) CheckEmail(
		email string,
		currentEmail string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "email", email)
	revel.Unbind(args, "currentEmail", currentEmail)
	return revel.MainRouter.Reverse("Account.CheckEmail", args).URL
}


type tProfile struct {}
var Profile tProfile


func (_ tProfile) Show(
		username string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "username", username)
	return revel.MainRouter.Reverse("Profile.Show", args).URL
}

func (_ tProfile) Settings(
		username string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "username", username)
	return revel.MainRouter.Reverse("Profile.Settings", args).URL
}

func (_ tProfile) UpdateSettings(
		username string,
		profile interface{},
		verifyPassword string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "username", username)
	revel.Unbind(args, "profile", profile)
	revel.Unbind(args, "verifyPassword", verifyPassword)
	return revel.MainRouter.Reverse("Profile.UpdateSettings", args).URL
}

func (_ tProfile) Password(
		username string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "username", username)
	return revel.MainRouter.Reverse("Profile.Password", args).URL
}

func (_ tProfile) UpdatePassword(
		username string,
		password string,
		verifyPassword string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "username", username)
	revel.Unbind(args, "password", password)
	revel.Unbind(args, "verifyPassword", verifyPassword)
	return revel.MainRouter.Reverse("Profile.UpdatePassword", args).URL
}

func (_ tProfile) FollowUser(
		username string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "username", username)
	return revel.MainRouter.Reverse("Profile.FollowUser", args).URL
}

func (_ tProfile) Followers(
		username string,
		page int,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "username", username)
	revel.Unbind(args, "page", page)
	return revel.MainRouter.Reverse("Profile.Followers", args).URL
}

func (_ tProfile) Following(
		username string,
		page int,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "username", username)
	revel.Unbind(args, "page", page)
	return revel.MainRouter.Reverse("Profile.Following", args).URL
}


type tPost struct {}
var Post tPost


func (_ tPost) Show(
		username string,
		id int,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "username", username)
	revel.Unbind(args, "id", id)
	return revel.MainRouter.Reverse("Post.Show", args).URL
}

func (_ tPost) Create(
		username string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "username", username)
	return revel.MainRouter.Reverse("Post.Create", args).URL
}

func (_ tPost) Save(
		username string,
		post interface{},
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "username", username)
	revel.Unbind(args, "post", post)
	return revel.MainRouter.Reverse("Post.Save", args).URL
}

func (_ tPost) Edit(
		username string,
		id int,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "username", username)
	revel.Unbind(args, "id", id)
	return revel.MainRouter.Reverse("Post.Edit", args).URL
}

func (_ tPost) Update(
		username string,
		id int,
		post interface{},
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "username", username)
	revel.Unbind(args, "id", id)
	revel.Unbind(args, "post", post)
	return revel.MainRouter.Reverse("Post.Update", args).URL
}

func (_ tPost) Remove(
		username string,
		id int,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "username", username)
	revel.Unbind(args, "id", id)
	return revel.MainRouter.Reverse("Post.Remove", args).URL
}

func (_ tPost) Delete(
		username string,
		id int,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "username", username)
	revel.Unbind(args, "id", id)
	return revel.MainRouter.Reverse("Post.Delete", args).URL
}

func (_ tPost) Like(
		username string,
		id int,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "username", username)
	revel.Unbind(args, "id", id)
	return revel.MainRouter.Reverse("Post.Like", args).URL
}


