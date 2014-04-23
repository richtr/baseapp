package controllers

import (
	"database/sql"
	"github.com/coopernurse/gorp"
	r "github.com/revel/revel"
	"github.com/revel/revel/modules/db/app"
  "code.google.com/p/go.crypto/bcrypt"
	"github.com/richtr/baseapp/app/models"
  "time"
)

var (
	Dbm *gorp.DbMap
)

func InitDB() {
	db.Init()

	dbDriver := r.Config.StringDefault("db.driver", "sqlite3")

	if dbDriver == "mysql" {
		Dbm = &gorp.DbMap{Db: db.Db, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}
	} else if dbDriver == "postgres" {
		Dbm = &gorp.DbMap{Db: db.Db, Dialect: gorp.PostgresDialect{}}
	} else {
		Dbm = &gorp.DbMap{Db: db.Db, Dialect: gorp.SqliteDialect{}}
	}

	setColumnSizes := func(t *gorp.TableMap, colSizes map[string]int) {
		for col, size := range colSizes {
			t.ColMap(col).MaxSize = size
		}
	}

	t := Dbm.AddTable(models.User{}).SetKeys(true, "UserId")
	t.ColMap("Password").Transient = true
	t.ColMap("Created").Transient = true
	t.ColMap("Profile").Transient = true
	setColumnSizes(t, map[string]int{
		"Email":    200,
	})

	t = Dbm.AddTable(models.Token{}).SetKeys(true, "TokenId")
	setColumnSizes(t, map[string]int{
		"Email":       200,
		"Type":        20,
		"Hash":        16,
	})

	t = Dbm.AddTable(models.Profile{}).SetKeys(true, "ProfileId")
	t.ColMap("User").Transient = true
	setColumnSizes(t, map[string]int{
		"Name":        100,
		"Summary":     140,
		"Description": 400,
		"PhotoUrl":    200,
	})

	t = Dbm.AddTable(models.Post{}).SetKeys(true, "PostId")
	t.ColMap("DateObj").Transient = true
	t.ColMap("ContentStr").Transient = true
	setColumnSizes(t, map[string]int{
		"Title":       400,
		"Content":     16777212, // mediumblob storage capacity
	})

	Dbm.TraceOn("[gorp]", r.INFO)

	// Create tables in datastore if they don't already exist
	Dbm.CreateTablesIfNotExists()

	// Set up database test data in 'test' run mode
	// e.g. `$> revel run baseapp test`
	if r.RunMode == "test" {

		bcryptPassword, _ := bcrypt.GenerateFromPassword(
			[]byte("demouser"), bcrypt.DefaultCost)
		demoUser := &models.User{
			Email: "demo@demo.com",
			HashedPassword: bcryptPassword,
			Confirmed: false,
		}
		if err := Dbm.Insert(demoUser); err != nil {
			panic(err)
		}

		demoProfile := &models.Profile{
			UserId: demoUser.UserId,
			Name: "Demo User",
			Summary: "Just a regular guy",
			Description: "...",
			PhotoUrl: "http://www.gravatar.com/avatar/53444f91e698c0c7caa2dbc3bdbf93fc?s=128&d=identicon",
			User: demoUser,
		}
		if err := Dbm.Insert(demoProfile); err != nil {
			panic(err)
		}

		demoPost := &models.Post{
			ProfileId: demoProfile.ProfileId,
			Title: "Hello World",
			ContentStr: "Full markdown support with things like [links](http://example.org)!",
			Status: "public",
			DateObj: time.Now(),
		}
		if err := Dbm.Insert(demoPost); err != nil {
			panic(err)
		}

  }

}

type GorpController struct {
	*r.Controller
	Txn *gorp.Transaction
}

func (c *GorpController) Begin() r.Result {
	txn, err := Dbm.Begin()
	if err != nil {
		panic(err)
	}
	c.Txn = txn
	return nil
}

func (c *GorpController) Commit() r.Result {
	if c.Txn == nil {
		return nil
	}
	if err := c.Txn.Commit(); err != nil && err != sql.ErrTxDone {
		panic(err)
	}
	c.Txn = nil
	return nil
}

func (c *GorpController) Rollback() r.Result {
	if c.Txn == nil {
		return nil
	}
	if err := c.Txn.Rollback(); err != nil && err != sql.ErrTxDone {
		panic(err)
	}
	c.Txn = nil
	return nil
}
