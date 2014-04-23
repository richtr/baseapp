package tests

import (
  "github.com/revel/revel"
)

type ApplicationTest struct {
	revel.TestSuite
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
