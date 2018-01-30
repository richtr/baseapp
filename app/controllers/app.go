package controllers

import (
	r "github.com/revel/revel"
	"github.com/richtr/baseapp/app/models"
	"github.com/richtr/baseapp/app/routes"
	"net/url"
	"strings"
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

func (c Application) Search(query string, page int) r.Result {

	var matchedProfiles []*models.Profile

	if page == 0 {
		page = 1
	}
	nextPage := page + 1
	size := 50 // results per page

	if query != "" {

		// Format query value
		sql_query_value, _ := url.QueryUnescape(query)
		sql_query_value = strings.Trim(sql_query_value, " @#")
		sql_query_value = strings.ToLower(sql_query_value)

		sql_query_value_full := "%" + sql_query_value + "%"
		sql_query_value_front := sql_query_value + "%"
		sql_query_value_back := "%" + sql_query_value

		sql_query_string := "SELECT * FROM Profile WHERE username LIKE ? OR name LIKE ? ORDER BY CASE WHEN username LIKE ? THEN 1 WHEN name LIKE ? THEN 2 WHEN username LIKE ? THEN 4 WHEN name LIKE ? THEN 5 ELSE 3 END LIMIT ?, ?"

		// Retrieve all profiles loosely matching search term
		results, err := c.Txn.Select(models.Profile{}, sql_query_string, sql_query_value_full, sql_query_value_full, sql_query_value_front, sql_query_value_front, sql_query_value_back, sql_query_value_back, (page-1)*size, size)

		if err == nil {
			for _, r := range results {
				matchedProfiles = append(matchedProfiles, r.(*models.Profile))
			}
		}

		if len(matchedProfiles) == 0 && page != 1 {
			return c.Redirect(routes.Application.Search(query, 1))
		}

	}

	return c.Render(query, matchedProfiles, page, nextPage)
}

func (c Application) SwitchToDesktop() r.Result {
	// Add desktop mode cookie
	c.Session["desktopmode"] = "1"

	referer, err := url.Parse(c.Request.Header.Get("Referer"))
	if err != nil || referer.String() == "" {
		return c.Redirect(routes.Application.Index())
	}

	return c.Redirect(referer.String())
}

func (c Application) SwitchToMobile() r.Result {
	// Remove desktop mode cookie
	delete(c.Session, "desktopmode")

	referer, err := url.Parse(c.Request.Header.Get("Referer"))
	if err != nil || referer.String() == "" {
		return c.Redirect(routes.Application.Index())
	}

	return c.Redirect(referer.String())
}
