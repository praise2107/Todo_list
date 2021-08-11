package actions

import (
	"github.com/gobuffalo/buffalo"
)

// NewsHandler is set up to serve up the news page
func NewsHandler(c buffalo.Context) error {
	return c.Render(200, r.HTML("news.html"))
}
