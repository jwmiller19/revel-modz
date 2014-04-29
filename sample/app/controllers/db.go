package controllers

import (
	gorm "github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"github.com/revel/revel"

	"github.com/iassic/revel-modz/modules/analytics"
	"github.com/iassic/revel-modz/modules/auth"
	"github.com/iassic/revel-modz/modules/forum"
	"github.com/iassic/revel-modz/modules/maillist"
	"github.com/iassic/revel-modz/modules/stripe"
	"github.com/iassic/revel-modz/modules/user"
	"github.com/iassic/revel-modz/modules/user-files"
)

var (
	TestDB *gorm.DB

	fill = false
)

type DbController struct {
	*revel.Controller
	Txn *gorm.DB
}

func (c *DbController) Begin() revel.Result {
	txn := TestDB.Begin()
	err := txn.Error
	checkPANIC(err)
	c.Txn = txn
	return nil
}

func (c *DbController) Commit() revel.Result {
	if c.Txn == nil {
		return nil
	}
	err := c.Txn.Commit().Error
	checkPANIC(err)

	c.Txn = nil
	return nil
}

func (c *DbController) Rollback() revel.Result {
	if c.Txn == nil {
		return nil
	}
	err := c.Txn.Rollback().Error
	checkPANIC(err)

	c.Txn = nil
	return nil
}

func InitDB() {

	var driver, spec string
	var found bool
	if driver, found = revel.Config.String("db.driver"); !found {
		revel.ERROR.Fatal("No db.driver found.")
	}
	if spec, found = revel.Config.String("db.spec"); !found {
		revel.ERROR.Fatal("No db.spec found.")
	}

	// Open a connection.
	ndb, err := gorm.Open(driver, spec)
	checkPANIC(err)

	ndb.SetLogger(gorm.Logger{revel.INFO})
	ndb.LogMode(true)

	TestDB = &ndb

	revel.INFO.Println("Connection made to DB")
}

func SetupTables() {
	revel.INFO.Println("Setting up Prod DB")
	addTables()
}

func SetupDevDB() {
	revel.INFO.Println("Setting up Dev DB")
	dropTables()
	addTables()

	user.FillTables(TestDB)
	forum.FillTables(TestDB)
	stripe.TestTables(TestDB)
	// fillMailTables()
	// testUserDB()

}

func dropTables() {
	revel.INFO.Println("Dropping tables")
	// analytics.DropTables(TestDB)
	auth.DropTables(TestDB)
	user.DropTables(TestDB)
	maillist.DropTables(TestDB)
	forum.DropTables(TestDB)
	stripe.DropTables(TestDB)
	// userfiles.DropTables(TestDB)
}

func addTables() {
	revel.INFO.Println("AutoMigrate tables")
	analytics.AddTables(TestDB)
	user.AddTables(TestDB)
	auth.AddTables(TestDB)
	forum.AddTables(TestDB)
	userfiles.AddTables(TestDB)
	stripe.AddTables(TestDB)
	maillist.AddTables(TestDB)
}
