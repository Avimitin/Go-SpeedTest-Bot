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

// Setup will generate a new database; require db file absolute path.
func Setup(dbPath string) error {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Println("[DatabaseError]Error occur when open a database.", err)
		return err
	}
	defer db.Close()
	table := `
CREATE TABLE IF NOT EXISTS manager (UID BIGINT NOT NULL PRIMARY KEY, name VARCHAR(255) NOT NULL);
`
	_, err = db.Exec(table)
	if err != nil {
		log.Println("[DatabaseError]Error occur when creating database,", err)
		return err
	}
	return nil
}

func DBFileExist() (string, bool) {
	_, err := ioutil.ReadFile("./bot.db")
	// If bot not running in project root path, try to get project environment path.
	if err != nil {
		// If found environment value
		if env := os.Getenv("SPT_BOT_PATH"); env != "" {
			// Try to get bot.db in project path, else create it.
			_, err = ioutil.ReadFile(env + "/bot.db")
			if err != nil {
				log.Println("[Database]Can't find bot.db, creating now.")
				err = Setup(env + "/bot.db")
				if err != nil {
					log.Println("[DatabaseError]Error occur when setup database.")
					os.Exit(-1)
				}
				return env + "/bot.db", true
			}
			return env + "bot.db", true
		}
		log.Println("[DatabaseError]Can't found environment path.")
		os.Exit(-1)
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
	defer row.Close()
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

func AddAdmin(db *sql.DB, user Admin) error {
	tx, err := db.Begin()
	if err != nil {
		log.Println("[DatabaseError]Fail to start a transaction, ", err)
		return err
	}
	stmt, err := tx.Prepare("INSERT INTO manager(UID, name) VALUES (?, ?)")
	if err != nil {
		log.Println("[DatabaseError]Fail to prepare a insertion.", err)
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(user.UID, user.Name)
	if err != nil {
		log.Println("[DatabaseError]Fail to execute value into manage table.", err)
		return err
	}
	err = tx.Commit()
	if err != nil {
		log.Println("[DatabaseError]Fail to commit a transaction.", err)
		return err
	}
	return nil
}

type DatabaseNotFound struct {
	Text string
}

func (d *DatabaseNotFound) Error() string {
	return d.Text
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
