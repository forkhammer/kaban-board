package kanban

import "main/db"

type UserRepository struct{}

func (r *UserRepository) GetUsers() ([]User, error) {
	var users = make([]User, 0)
	result := db.DefaultConnection.Db.Order("name").Find(&users)
	return users, result.Error
}

func (r *UserRepository) GetVisibleUsers() ([]User, error) {
	var users = make([]User, 0)
	result := db.DefaultConnection.Db.Where(&User{IsVisible: true}).Order("name").Find(&users)
	return users, result.Error
}

func (r *UserRepository) GetOrCreate(query, attrs User) (*User, error) {
	var user User

	if result := db.DefaultConnection.Db.Where(query).Attrs(attrs).FirstOrCreate(&user); result.Error != nil {
		return nil, result.Error
	} else {
		return &user, nil
	}
}

func (r *UserRepository) GetUserBydId(id int) (*User, error) {
	var user User
	if result := db.DefaultConnection.Db.Where(&User{Id: uint(id)}).First(&user); result.Error != nil {
		return nil, result.Error
	} else {
		return &user, nil
	}
}

func (r *UserRepository) SaveUser(user *User) error {
	result := db.DefaultConnection.Db.Save(user)
	return result.Error
}
