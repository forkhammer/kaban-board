package models

type LabelId string
type Color string

type Label struct {
	Id            LabelId
	Name          string
	Color         Color
	TextColor     Color
	AltName       *string
	BindingStatus *IssueBindingStatus
	Priority      *IssueBindingPriority
}

func (l Label) Validate() error {
	return nil
}
