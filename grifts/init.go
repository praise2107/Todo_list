package grifts

import (
  "github.com/gobuffalo/buffalo"
	"myapp/actions"
)

func init() {
  buffalo.Grifts(actions.App())
}
