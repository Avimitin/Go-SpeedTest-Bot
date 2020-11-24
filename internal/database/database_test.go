package database

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func TestSetup(t *testing.T) {
	err := Setup(os.Getenv("SPT_BOT_PATH") + "/bot.db")
	if err != nil {
		t.Error(err)
		return
	}
	_, err = ioutil.ReadFile(os.Getenv("SPT_BOT_PATH") + "/bot.db")
	if err != nil {
		t.Errorf("Expect bot.db exist but get nothing")
		return
	}
	TestConnect(t)
}

func TestConnect(t *testing.T) {
	path, ok := DBFileExist()
	if !ok {
		t.Errorf("Want bot path got nothing")
	}
	_, err := Connect(path)
	if err != nil {
		t.Fail()
	}
	table, err := GetTable(NewDB())
	if err != nil {
		t.Fail()
	}
	if table == "" {
		t.Errorf("Want manager got nothing")
		t.Fail()
	}
	if table != "manager" {
		t.Errorf("want manager got %s", table)
	}
}

func TestGetAdmin(t *testing.T) {
	admin, err := GetAdmin(NewDB())
	if err != nil {
		t.Errorf("Error occur, %s", err.Error())
		return
	}
	if len(admin) == 0 {
		t.Errorf("Want list of admin but got nothing.")
		return
	}
	if !reflect.DeepEqual(admin, []Admin{{649191333, "SaitoAsuka_kksk"}}) {
		t.Errorf("Expected 649191333 but got %d", admin[0].UID)
		t.Fail()
	}
}

func TestDBFileExist(t *testing.T) {
	path, ok := DBFileExist()
	if !ok {
		t.Errorf("Want bot.db got %s", path)
	}
	fmt.Println(path)
}

func TestAddAdmin(t *testing.T) {
	NewUser := Admin{
		UID:  112233,
		Name: "test user",
	}
	db := NewDB()
	if db == nil {
		t.Error("Database is nil")
		return
	}
	err := AddAdmin(db, NewUser)
	if err != nil {
		t.Errorf("%v", err)
		return
	}
	admins, err := GetAdmin(db)
	if err != nil {
		t.Errorf("%v", err)
		return
	}
	for _, admin := range admins {
		if admin.UID == NewUser.UID {
			return
		}
	}
	t.Errorf("Can't find the user just add")
}
