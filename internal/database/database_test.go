package database

import (
	"database/sql"
	"fmt"
	"reflect"
	"testing"
)

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

func NewDB() *sql.DB {
	path, ok := DBFileExist()
	if !ok {
		return nil
	}
	db, err := Connect(path)
	if err == nil {
		return db
	}
	return nil
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
