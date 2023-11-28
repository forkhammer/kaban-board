package kanban

import "main/db"

func MigrateModels() {
	if err := db.DefaultConnection.Db.AutoMigrate(&Project{}, &User{}, &Team{}, &Label{}, &Column{}); err != nil {
		panic(err)
	}
}
