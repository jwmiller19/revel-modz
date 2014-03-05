package controllers

import (
	"github.com/iassic/revel-modz/modules/auth"
	"github.com/iassic/revel-modz/modules/maillist"
	"github.com/iassic/revel-modz/modules/user"
	"github.com/revel/revel"

	"github.com/iassic/revel-modz/sample/app/models"
	"github.com/iassic/revel-modz/sample/app/routes"
)

type App struct {
	DbController
}

func (c App) RenderArgsFill() revel.Result {
	u := c.connected()
	if u != nil {
		c.RenderArgs["user_basic"] = u

		// look up role in RBAC module
		isAdmin := u.UserName == "admin@domain.com"
		if isAdmin {
			// set up things for an admin role
			c.Session["admin"] = "true"
		}
	}

	return nil
}

func (c App) connected() *user.UserBasic {
	if c.RenderArgs["user_basic"] != nil {
		return c.RenderArgs["user_basic"].(*user.UserBasic)
	}
	if username, ok := c.Session["user"]; ok {
		u := user.GetUserBasicByName(c.Txn, username)
		if u == nil {
			revel.ERROR.Println("user field in Session[] not found in DB")
			return nil
		}
		// revel.WARN.Printf("connected :: %+v", *u)
		return u
	}
	return nil
}

func (c App) Index() revel.Result {
	return c.Render()
}

func (c App) Result() revel.Result {
	return c.Render()
}

func (c App) Signup() revel.Result {
	return c.Render()
}

func (c App) SignupPost(usersignup *models.UserSignup) revel.Result {
	usersignup.Validate(c.Validation)

	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(routes.App.Signup())
	}

	// check that this email is not in the DB already
	UB := user.GetUserBasicByName(c.Txn, usersignup.Email)
	if UB != nil {
		c.Validation.Error("Email already taken").Key("usersignup.Email")
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(routes.App.Signup())
	}

	UB, err := c.addNewUser(usersignup.Email, usersignup.Password)
	checkERROR(err)

	c.Flash.Out["heading"] = "Thanks for Joining!"
	c.Flash.Out["message"] = "Signup successful for " + usersignup.Email

	c.Session["user"] = UB.UserName
	c.RenderArgs["user_basic"] = UB
	return c.Redirect(routes.User.Result())

}

func (c App) Maillist() revel.Result {
	return c.Render()
}

func (c App) MaillistPost(usermaillist *models.UserMaillist) revel.Result {
	usermaillist.Validate(c.Validation)

	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(routes.App.Maillist())
	}

	// check that this email is not in the DB already
	UB := user.GetUserBasicByName(c.Txn, usermaillist.Email)
	if UB != nil {
		c.Validation.Error("Email already taken").Key("usermaillist.Email")
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(routes.App.Signup())
	}

	_, err := c.addNewMaillistUser(usermaillist.Email, "MaillistPost()")
	checkERROR(err)

	c.Flash.Out["heading"] = "Thanks for Joining!"
	c.Flash.Out["message"] = usermaillist.Email + " is now subscribed to the mailing list."

	return c.Redirect(routes.App.Result())

}

func (c App) Register() revel.Result {
	return c.Render()
}

func (c App) RegisterPost(userregister *models.UserRegister) revel.Result {
	userregister.Validate(c.Validation)

	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(routes.App.Maillist())
	}

	// check that this email is not in the DB already
	UB := user.GetUserBasicByName(c.Txn, userregister.Email)
	if UB != nil {
		c.Validation.Error("Email already taken").Key("userregister.Email")
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(routes.App.Signup())
	}

	var err error
	UB, err = c.addNewUser(userregister.Email, userregister.Password)
	checkERROR(err)

	// TODO  which mailing lists did they check off?
	// ALSO  user Basic will be added twice if this current call is made
	// _, err = c.addNewMaillistUser(userregister.Email)
	// checkERROR(err)

	// TODO add profile DB insert

	c.Flash.Out["heading"] = "Thanks for Joining!"
	c.Flash.Out["message"] = userregister.Email + " is now subscribed to the mailing list."

	return c.Redirect(routes.App.Result())

}

func (c App) Login() revel.Result {
	return c.Render()
}

func (c App) LoginPost(userlogin *models.UserLogin) revel.Result {
	userlogin.Validate(c.Validation)
	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(routes.App.Login())
	}

	var found, valid bool

	// check for user in basic table
	UB := user.GetUserBasicByName(c.Txn, userlogin.Email)
	if UB != nil {
		found = true
	} else {
		println("FLASH ERROR USER")
		c.Flash.Error("unknown user")
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(routes.App.Login())
	}

	// check for user in auth table
	P := user.UserPass{UB.UserId, userlogin.Email, userlogin.Password}
	U, err := auth.Authenticate(c.Txn, &P)
	if err != nil || U == nil {
		println("FLASH ERROR PASSWORD")
		c.Flash.Error("bad password")
	} else {
		valid = true
	}

	if found && valid {
		c.Session["user"] = UB.UserName
		c.RenderArgs["user_basic"] = UB
		return c.Redirect(routes.User.Result())

	} else {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(routes.App.Login())
	}
}

func (c App) Logout() revel.Result {
	for k := range c.Session {
		delete(c.Session, k)
	}
	return c.Redirect(routes.App.Index())
}

func (c App) addNewUser(email, password string) (*user.UserBasic, error) {

	// uuid := get random number (that isn't used already)
	uuid := user.GenerateNewUserId(c.Txn)
	UB := &user.UserBasic{
		UserId:   uuid,
		UserName: email,
	}
	UP := &user.UserPass{UB.UserId, email, password}

	// add user to tables
	// TODO do something more with the errosr
	err := user.AddUserBasic(TestDB, UB)
	checkERROR(err)

	_, err = auth.AddUserAuth(TestDB, UP)
	checkERROR(err)

	return UB, nil
}

func (c App) addNewMaillistUser(email, list string) (*maillist.MaillistUser, error) {

	// uuid := get random number (that isn't used already)
	uuid := user.GenerateNewUserId(c.Txn)
	UB := &user.UserBasic{
		UserId:   uuid,
		UserName: email,
	}

	err := user.AddUserBasic(TestDB, UB)
	checkERROR(err)

	MA, err := maillist.AddUser(TestDB, uuid, email, list)
	checkERROR(err)

	return MA, nil
}
