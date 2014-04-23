package controllers

import (
	r "github.com/revel/revel"
)

type Application struct {
	*r.Controller
	Account
}

func (c Application) Index() r.Result {
	return c.Render()
}

func (c Application) About() r.Result {
	return c.Render()
}

func (c Application) Contact() r.Result {
	return c.Render()
}
