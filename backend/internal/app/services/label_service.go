package app_services

import (
	"main/internal/app/queries"
	domain "main/internal/domain/models"
	"main/internal/domain/repo"
	"main/pkg/utils"
	"slices"
)

type LabelService struct {
	labelRepo  repo.LabelRepo     `di.inject:"LabelRepository"`
	labelQuery queries.LabelQuery `di.inject:"LabelQuery"`
}

func (s *LabelService) ExistNames(names []string) (bool, []string, error) {
	labelNames := utils.Unique(names, func(l string) string {
		return l
	})

	labels, err := s.labelRepo.List(s.labelQuery.GetSpec(queries.LabelFilter{
		Names: labelNames,
	}))
	if err != nil {
		return false, []string{}, err
	}
	existNames := utils.Unique(
		utils.Map(*labels, func(l domain.Label) string {
			return l.Name
		}), func(name string) string {
			return name
		},
	)
	if len(existNames) != len(labelNames) {
		diff := utils.Filter(labelNames, func(name string) bool {
			return !slices.Contains(existNames, name)
		})
		return false, diff, nil
	}

	return true, []string{}, nil
}
