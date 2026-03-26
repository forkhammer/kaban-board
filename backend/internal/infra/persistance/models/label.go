package models

type Label struct {
	Id            string  `gorm:"id;primaryKey"`
	Name          string  `gorm:"name;not null"`
	Color         string  `gorm:"color"`
	TextColor     string  `gorm:"text_color"`
	AltName       *string `gorm:"alt_name"`
	BindingStatus *string `gorm:"binding_status"`
	Priority      *string `gorm:"priority"`
}
