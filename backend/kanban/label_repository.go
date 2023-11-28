package kanban

import "main/db"

type LabelRepository struct{}

func (r *LabelRepository) GetOrCreate(query Label, attrs Label) (*Label, error) {
	var label Label
	if result := db.DefaultConnection.Db.Where(query).Attrs(attrs).FirstOrCreate(&label); result.Error != nil {
		return nil, result.Error
	} else {
		return &label, nil
	}
}

func (r *LabelRepository) GetLabels() ([]Label, error) {
	var labels []Label
	if result := db.DefaultConnection.Db.Order("id").Find(&labels); result.Error != nil {
		return nil, result.Error
	} else {
		return labels, nil
	}
}
