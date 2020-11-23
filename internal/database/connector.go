package database

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"io/ioutil"
	"log"
	"os"
)

type Admin struct {
	UID  int64
	Name string
}

func DBFileExist() (string, bool) {
	_, err := ioutil.ReadFile("./bot.db")
	if err != nil {
		if env := os.Getenv("SPT_BOT_PATH"); env != "" {
			_, err = ioutil.ReadFile(env + "/bot.db")
			if err != nil {
				// If bot.db can't be found in running path or environment path return false.
				log.Println("[DatabaseError]Can't find bot.db")
				return "", false
			}
			return env + "bot.db", true
		}
	}
	return "./bot.db", true
}

func Connect(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", path)
	log.Println("Load database at", path)
	if err != nil {
		log.Printf("[DatabaseError]%v", err)
		os.Exit(-1)
	}
	err = db.Ping()
	if err != nil {
		log.Println("[DatabaseError]Unable to connect to database")
		os.Exit(-1)
	}
	return db, nil
}

func GetTable(db *sql.DB) (string, error) {
	row, err := db.Query("SELECT name FROM sqlite_master where type='table'")
	if err != nil {
		log.Println("[DatabaseError]Unable to get table name,", err.Error())
		return "", err
	}
	var table string
	for row.Next() {
		err = row.Scan(&table)

		if err != nil {
			log.Println("[DatabaseError]Unable Scan value")
			return "", err
		}
	}
	return table, nil
}

func GetAdmin(db *sql.DB) ([]Admin, error) {
	row, err := db.Query("SELECT UID, Name FROM manager")
	if err != nil {
		log.Println("[DatabaseError]Unable to get manager info,", err)
		return nil, err
	}
	var name string
	var uid int64
	var admins []Admin
	for row.Next() {
		err := row.Scan(&uid, &name)
		if err != nil {
			log.Println("[DatabaseError]Unable to scan value, ", err)
			continue
		}
		admins = append(admins, Admin{Name: name, UID: uid})
	}
	err = row.Err()
	if err != nil {
		log.Println("[DatabaseError]Unable to get next row.")
		return nil, err
	}
	return admins, nil
}
