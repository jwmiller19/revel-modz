package controllers

import (
	"github.com/revel/revel"

	"github.com/iassic/revel-modz/modules/maillist"

	"github.com/iassic/revel-modz/sample/app/routes"
)

type Admin struct {
	User
}

// moving towards RBAC here...
func (c Admin) CheckLoggedIn() revel.Result {
	u := c.connected()
	if u == nil {
		c.Flash.Error("Please log in first")
		return c.Redirect(routes.App.Login())
	}

	// look up role in RBAC module
	isAdmin := u.UserName == "admin@domain.com"

	if !isAdmin {
		return c.Redirect(routes.App.Index())
	}

	// set up things for an admin role
	c.Session["admin"] = "true"

	return nil
}

func (c Admin) Index() revel.Result {
	return c.Render()
}

func (c Admin) MaillistView() revel.Result {
	println("TOG HERE")
	maillist_users, err := maillist.GetAllUsers(c.Txn)
	if err != nil {
		revel.ERROR.Println(err)
		return c.Render()
	}
	return c.Render(maillist_users)
}

func (c Admin) MaillistCompose() revel.Result {
	return c.Render()
}
