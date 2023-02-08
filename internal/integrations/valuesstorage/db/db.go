package db

import (
	"encoding/json"

	"github.com/hatchet-dev/hatchet/internal/models"
	"github.com/hatchet-dev/hatchet/internal/repository"
)

// DatabaseValuesStore represents a values storage manager which uses the repository for
// retrieving and writing values
type DatabaseValuesStore struct {
	repo repository.Repository
}

func NewDatabaseValuesStore(repo repository.Repository) *DatabaseValuesStore {
	return &DatabaseValuesStore{repo}
}

func (d *DatabaseValuesStore) WriteValues(mvv *models.ModuleValuesVersion, values map[string]interface{}) error {
	valuesBytes, err := json.Marshal(&values)

	if err != nil {
		return err
	}

	// create new module values object in the database
	mv := &models.ModuleValues{
		ModuleValuesVersionID: mvv.ID,
		Values:                valuesBytes,
	}

	mv, err = d.repo.ModuleValues().CreateModuleValues(mv)

	return err
}

func (d *DatabaseValuesStore) ReadValues(mvv *models.ModuleValuesVersion) (map[string]interface{}, error) {
	mv, err := d.repo.ModuleValues().ReadModuleValuesByVersionID(mvv.ModuleID, mvv.ID)

	if err != nil {
		return nil, err
	}

	res := make(map[string]interface{})

	err = json.Unmarshal(mv.Values, &res)

	if err != nil {
		return nil, err
	}

	return res, nil
}
