package vcs

import (
	"encoding/json"
	"fmt"

	"github.com/hatchet-dev/hatchet/internal/integrations/vcs"
	"github.com/hatchet-dev/hatchet/internal/models"
)

// VCSValuesStore represents a values storage manager which uses a path in the VCS
// repository to reference values
type VCSValuesStore struct {
	vcsRepo vcs.VCSRepository
	ref     string
}

func NewGithubValuesStore(vcsRepo vcs.VCSRepository, ref string) *VCSValuesStore {
	return &VCSValuesStore{vcsRepo, ref}
}

func (d *VCSValuesStore) WriteValues(mvv *models.ModuleValuesVersion, values map[string]interface{}) error {
	return fmt.Errorf("vcs-based values storage does not support writing values")
}

func (d *VCSValuesStore) ReadValues(mvv *models.ModuleValuesVersion) (map[string]interface{}, error) {
	file, err := d.vcsRepo.ReadFile(d.ref, mvv.GitValuesPath)

	if err != nil {
		return nil, err
	}

	res := make(map[string]interface{})

	err = json.NewDecoder(file).Decode(&res)

	if err != nil {
		return nil, err
	}

	return res, nil
}
