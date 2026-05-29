package factories

import (
	"main/internal/domain/models"
	"main/internal/domain/repo"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/goioc/di"
)

type LabelFactory struct {
	repo repo.LabelRepo
}

func NewLabelFactory() *LabelFactory {
	return &LabelFactory{
		repo: di.GetInstance("LabelRepository").(repo.LabelRepo),
	}
}

func (f *LabelFactory) Build(data map[string]any) *models.Label {
	label := &models.Label{}

	if val, ok := data["id"].(string); ok && val != "" {
		label.Id = models.LabelId(val)
	} else {
		label.Id = models.LabelId(gofakeit.UUID())
	}

	if val, ok := data["name"].(string); ok && val != "" {
		label.Name = val
	} else {
		label.Name = gofakeit.Word()
	}

	if val, ok := data["color"].(string); ok && val != "" {
		label.Color = models.Color(val)
	} else {
		label.Color = models.Color(gofakeit.HexColor())
	}

	if val, ok := data["text_color"].(string); ok && val != "" {
		label.TextColor = models.Color(val)
	} else {
		label.TextColor = models.Color(gofakeit.HexColor())
	}

	if val, ok := data["alt_name"].(*string); ok {
		label.AltName = val
	}

	if val, ok := data["binding_status"].(*models.IssueBindingStatus); ok {
		label.BindingStatus = val
	}

	if val, ok := data["priority"].(*models.IssueBindingPriority); ok {
		label.Priority = val
	}

	return label
}

func (f *LabelFactory) Create(data map[string]any) (*models.Label, error) {
	label := f.Build(data)
	return f.repo.Create(label)
}
