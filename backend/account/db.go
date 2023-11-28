package account

import "main/db"

func MigrateModels() {
	db.DefaultConnection.Db.AutoMigrate(&Account{})
}
