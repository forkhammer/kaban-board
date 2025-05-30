package usecases

import (
	"main/internal/app/queries"
	domain "main/internal/domain/models"
	"main/internal/domain/repo"
	"main/pkg/utils"
)

type UpdateTeamRequest struct {
	Id     uint
	Title  string
	Groups []int
}

type CreateTeamRequest struct {
	Title  string
	Groups []int
}

type TeamUseCases struct {
	teamRepo   repo.TeamRepo      `di.inject:"TeamRepository"`
	groupRepo  repo.GroupRepo     `di.inject:"GroupRepository"`
	groupQuery queries.GroupQuery `di.inject:"GroupQuery"`
}

func (uc *TeamUseCases) GetTeams() (*[]domain.Team, error) {
	return uc.teamRepo.List(nil)
}

func (uc *TeamUseCases) GetTeam(id uint) (*domain.Team, error) {
	return uc.teamRepo.Get(domain.TeamId(id))
}

func (uc *TeamUseCases) Create(request *CreateTeamRequest) (*domain.Team, error) {
	groupIds := utils.Map(request.Groups, func(g int) domain.GroupId {
		return domain.GroupId(g)
	})

	var groups *[]domain.Group
	if len(groupIds) > 0 {
		var err error
		if groups, err = uc.groupRepo.List(uc.groupQuery.GetSpec(queries.GroupFilter{Ids: groupIds})); err != nil {

		}
	} else {
		groups = &[]domain.Group{}
	}

	team, err := domain.NewTeam(0, request.Title, *groups)
	if err != nil {
		return nil, err
	}
	return uc.teamRepo.Create(team)
}

func (uc *TeamUseCases) Update(request *UpdateTeamRequest) (*domain.Team, error) {
	team, err := uc.teamRepo.Get(domain.TeamId(request.Id))
	if err != nil {
		return nil, err
	}
	team.Title = request.Title

	groupIds := utils.Map(request.Groups, func(g int) domain.GroupId {
		return domain.GroupId(g)
	})

	var groups *[]domain.Group
	if len(groupIds) > 0 {
		var err error
		if groups, err = uc.groupRepo.List(uc.groupQuery.GetSpec(queries.GroupFilter{Ids: groupIds})); err != nil {

		}
	} else {
		groups = &[]domain.Group{}
	}

	team.Groups = *groups

	if err := team.Validate(); err != nil {
		return nil, err
	}

	return uc.teamRepo.Update(team)
}

func (uc *TeamUseCases) Delete(id uint) error {
	team, err := uc.teamRepo.Get(domain.TeamId(id))
	if err != nil {
		return err
	}
	return uc.teamRepo.Delete(team.Id)
}
