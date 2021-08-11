// this auth package is used to authenticate the user before signing in

package actions

import (
	"database/sql"
	"myapp/models"
	"strings"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/validate"
	"github.com/markbates/errx"
	"golang.org/x/crypto/bcrypt"
)

// AuthNew loads the signin page
func AuthNew(c buffalo.Context) error {
	c.Set("user", models.User{})
	return c.Render(200, r.HTML("auth_new.plush.html"))
}

// AuthCreate attempts to log the user in with an existing account.
func AuthCreate(c buffalo.Context) error {
	u := &models.User{}
	if err := c.Bind(u); err != nil {
		return err
	}

	tx := c.Value("tx").(*pop.Connection)

	// find a user with the email
	err := tx.Where("username = ?", strings.ToLower(u.Username)).First(u)

	// helper function to handle bad attempts
	bad := func() error {
		c.Set("user", u)
		verrs := validate.NewErrors()
		verrs.Add("username", "invalid username/password")
		c.Set("errors", verrs)
		return c.Render(422, r.HTML("auth_new.plush.html"))
	}

	if err != nil {
		if errx.Unwrap(err) == sql.ErrNoRows {
			// couldn't find an user with the supplied username.
			return bad()
		}
		return err
	}

	// confirm that the given password matches the hashed password from the db
	err = bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(u.Password))
	if err != nil {
		return bad()
	}
	c.Session().Set("current_user_id", u.ID)
	c.Flash().Add("success", "Welcome Back to Todo List Online Management!")

	return c.Redirect(302, "/tasks")
}

// AuthDestroy clears the session and logs a user out
func AuthDestroy(c buffalo.Context) error {
	c.Session().Clear()
	c.Flash().Add("success", "You have been logged out!")
	return c.Redirect(302, "/")
}
