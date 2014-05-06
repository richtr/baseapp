BaseApp
=======

#### A bootstrap web application built in [Go](http://golang.org) on top of the [Revel Web Framework](https://revel.github.io) ####

BaseApp is a bootstrap Revel web application that provides a baseline web application starter kit with the following features:

* Basic pages (Home, About Us, Contact Us, etc)
* Account registration (including safe storage of passwords with [bcrypt](https://en.wikipedia.org/wiki/Bcrypt))
* Account confirmation via email (via configurable e-mail settings)
* Account recovery via email (via configurable e-mail settings)
* Account log in (by username or email address)
* Public user profiles (with Gravatar integration for profile photos and profile edit functionality)
* Public user posts (with full Markdown support and edit + delete functionality)
* Social Features (Profile Follower and Following Counts and Post Likes)
* People finder search (by username or full name)
* Auto-provisioned backend data store (SQLite3, MySQL, PostgreSQL)
* Full live client-side and static server-side form validation
* Full non-interactive and interactive testing framework (Test-driven development process)
* Responsive front-end design (looks good on mobile, tablet and desktop screen sizes)

### Contents ###

* [Screenshots](#screenshots)
* [Quick Start](#quick-start)
* [BaseApp Run Modes](#baseapp-run-modes)
* [Docker Usage](#docker-usage)
* [Feedback](#feedback)
* [License](#license)

### Screenshots ###

Example Home Page:

<img src="https://github.com/richtr/baseapp/raw/master/screenshots/01.HomePage.png" style="max-width: 100%"/>

User Registration Form:

<img src="https://github.com/richtr/baseapp/raw/master/screenshots/02.Register.png" style="max-width: 100%"/>

User Registration Form - Form Validation Failed:

<img src="https://github.com/richtr/baseapp/raw/master/screenshots/09.Register_Fail.png" style="max-width: 100%"/>

User Login Form:

<img src="https://github.com/richtr/baseapp/raw/master/screenshots/03.Login.png" style="max-width: 100%"/>

Display User Profile:

<img src="https://github.com/richtr/baseapp/raw/master/screenshots/04.Profile.png" style="max-width: 100%"/>

Followers Display:

<img src="https://github.com/richtr/baseapp/raw/master/screenshots/11.Followers.png" style="max-width: 100%"/>

Following Display:

<img src="https://github.com/richtr/baseapp/raw/master/screenshots/12.Following.png" style="max-width: 100%"/>

Edit User Profile Form:

<img src="https://github.com/richtr/baseapp/raw/master/screenshots/05.EditProfile.png" style="max-width: 100%"/>

New Post Form:

<img src="https://github.com/richtr/baseapp/raw/master/screenshots/06.NewPost.png" style="max-width: 100%"/>

Post Display:

<img src="https://github.com/richtr/baseapp/raw/master/screenshots/10.Post.png" style="max-width: 100%"/>

Reset Password Form (if email settings are provided in `conf/app.conf`):

<img src="https://github.com/richtr/baseapp/raw/master/screenshots/07.ResetPwd.png" style="max-width: 100%"/>

Interactive Test Suite (at `http://localhost:9000/@tests`)

<img src="https://github.com/richtr/baseapp/raw/master/screenshots/08.TestRunner.png" style="max-width: 100%"/>

### Quick Start ####

#### Prerequisites ####

You will need a [functioning Revel installation](https://revel.github.io/tutorial/gettingstarted.html) for BaseApp to work.

#### Installing BaseApp ####

To get [BaseApp](https://github.com/richtr/baseapp), run

    go get github.com/richtr/baseapp

This command does a couple things:

* Go uses git to clone the repository into $GOPATH/src/github.com/richtr/baseapp/
* Go transitively finds all of the dependencies and runs go get on them as well.

#### Running BaseApp ####

Before you can start using BaseApp you will need to add your own `app.conf` file in your `conf/` directory. You can copy and use the <a href="https://github.com/richtr/baseapp/blob/master/conf/app.conf.default">default configuration file</a> with

    cp ./github.com/richtr/baseapp/conf/app.conf.default ./github.com/richtr/baseapp/conf/app.conf

Note: It is highly recommended that you review the configuration options available in `app.conf.default` before you run your project for the first time!

Once you have setup an `app.conf` file you can start BaseApp [_in test mode_](#baseapp-run-modes), with

    revel run github.com/richtr/baseapp test

Point your browser to your BaseApp installation at `http://localhost:9000` (or the path you specified in your `app.conf` file) and away you go.

### BaseApp Run Modes ###

BaseApp can be run in three different modes that are each useful for different stages of application development:

1. `test` mode: Uses an *in-memory* sqlite3 datastore that is created and populated with basic data which is always wiped when your application ends.

    You can run BaseApp in `test` mode as follows:

        $> revel run github.com/richtr/baseapp/ test

    Once BaseApp is up and running you can point your browser to `http://localhost:9000` to use the application or `http://localhost:9000/@tests` to run the interactive test suite.

    To run the BaseApp test suite in non-interactive mode you can use:

        $> revel test github.com/richtr/baseapp/ test

    If testing is successful then a `test-results/result.passed` is written. Otherwise a `test-results/result.failed` is written. You can use check for these files when testing before deployment within your own continuous integration system.

2. `dev` mode [default mode]: Uses a blank sqlite3/mysql/postgres datastore that is created if it does not already exist at your chosen location and persists whenever your application is restarted. Outputs detailed error messages if something goes wrong in your application.

    You can run BaseApp in `dev` mode as follows:

        $> revel run github.com/richtr/baseapp/ dev

    Or, simply:

        $> revel run github.com/richtr/baseapp/

3. `prod` mode: Uses a blank sqlite3/mysql/postgres datastore that is created if it does not already exist at your chosen location and persists whenever your application is restarted. Outputs user-friendly error messages if something goes wrong in your application.

    You can run BaseApp in `prod` mode as follows:

        $> revel run github.com/richtr/baseapp/ prod

Note: Both `dev` and `prod` modes require a configured backend DB. See [app.conf.default](https://github.com/richtr/baseapp/blob/master/conf/app.conf.default). The `test` mode creates an in-memory database representation that dies when the app dies.

### Docker Usage ###

This project comes with its own `Dockerfile` for building your own LXC container via [Docker](http://www.docker.io) to host BaseApp.

You will need a [functioning Docker installation](http://docs.docker.io/installation/) to build and use BaseApp as an LXC container.

To build your own BaseApp LXC container via Docker is easy:

```
$ git clone https://github.com/richtr/baseapp.git
$ cd baseapp
$ sudo docker build -t="richtr/baseapp" .
```

To spawn a new instance of BaseApp:

```bash
$ sudo docker run -p 80 -d richtr/baseapp
```

This will return an ID output on the command-line like:

```
d404cc2fa27b
```

This is your newly created container ID. You can use this to find the external port you can use to access BaseApp from your host machine:

```bash
$ sudo docker port d404cc2fa27b 80 # Make sure to change the ID to yours!
```

Then you can the visit the following URL in a browser on your host machine to get started:

```
http://127.0.0.1:<port>
```

### Feedback ###

If you find any bugs or issues please report them on the [BaseApp Issue Tracker](https://github.com/richtr/baseapp/issues).

If you would like to contribute to this project please consider [forking this repo](https://github.com/richtr/baseapp/fork), making your changes and then creating a new [Pull Request](https://github.com/richtr/baseapp/pulls) back to the main code repository.

### License ###

MIT. Copyright (c) Rich Tibbett.

See the [LICENSE](https://github.com/richtr/baseapp/blob/master/LICENSE) file.
