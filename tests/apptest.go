package tests

import (
	"github.com/revel/revel/testing"
)

type ApplicationTest struct {
	testing.TestSuite
}

func (t *ApplicationTest) Before() {
	// Runs before each test below is executed
}

func (t *ApplicationTest) After() {
	// Runs after each test below is executed
}

func (t *ApplicationTest) TestIndex() {
	t.Get("/")
	t.AssertOk()
	t.AssertContentType("text/html; charset=utf-8")
}

func (t *ApplicationTest) TestAbout() {
	t.Get("/about")
	t.AssertOk()
	t.AssertContentType("text/html; charset=utf-8")
}

func (t *ApplicationTest) TestContact() {
	t.Get("/contact")
	t.AssertOk()
	t.AssertContentType("text/html; charset=utf-8")
}
