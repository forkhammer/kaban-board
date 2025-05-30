package models

type User struct {
	Id        uint     `gorm:"primarykey"`
	Name      string   `gorm:"name"`
	Username  string   `gorm:"username"`
	AvatarUrl string   `gorm:"avatar_url"`
	IsVisible bool     `gorm:"is_visible;default:true;not null"`
	Groups    []*Group `gorm:"many2many:user_groups"`
}
