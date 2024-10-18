package models

type User struct {
	Id        uint     `gorm:"primarykey" json:"id"`
	Name      string   `gorm:"name" json:"name"`
	Username  string   `gorm:"username" json:"username"`
	AvatarUrl string   `gorm:"avatar_url" json:"avatar_url"`
	IsVisible bool     `gorm:"is_visible;default:true;not null" json:"is_visible"`
	Groups    []*Group `gorm:"many2many:user_groups" json:"groups"`
}
