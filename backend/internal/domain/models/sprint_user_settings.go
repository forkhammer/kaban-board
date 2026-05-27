package models

import "github.com/go-playground/validator/v10"

type SprintUserSettingsId uint

type SprintUserSettings struct {
	Id           SprintUserSettingsId
	SprintId     uint
	UserId       uint
	User         User
	HoursPerUser uint `validate:"required,gt=0"`
}

func (s *SprintUserSettings) Validate() error {
	validator := validator.New()
	return validator.Struct(s)
}
