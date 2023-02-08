package valuesstorage

import "github.com/hatchet-dev/hatchet/internal/models"

type ValuesStorageManager interface {
	WriteValues(mvv *models.ModuleValuesVersion, values map[string]interface{}) error
	ReadValues(mvv *models.ModuleValuesVersion) (map[string]interface{}, error)
}
