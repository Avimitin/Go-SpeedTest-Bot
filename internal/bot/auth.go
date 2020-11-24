package bot

import (
	"go-speedtest-bot/internal/database"
	"log"
)

var admins []database.Admin

func Auth(id int64) bool {
	if len(admins) == 0 {
		return false
	}
	for _, admin := range admins {
		if admin.UID == id {
			return true
		}
	}
	return false
}

func LoadAdmin() error {
	db := database.NewDB()
	if db == nil {
		return &database.DatabaseNotFound{
			Text: "Database not found",
		}
	}
	result, err := database.GetAdmin(db)
	if err != nil {
		log.Println("[AuthError]Can't fetch userdata")
		return err
	}
	if len(result) == 0 {
		log.Println("[AuthError]Can't fetch userdata")
		return &database.DatabaseNotFound{
			Text: "Can't fetch userdata",
		}
	}
	for _, r := range result {
		admins = append(admins, r)
	}
	return nil
}
