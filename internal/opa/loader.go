package opa

import (
	"context"
	"fmt"
	"io/ioutil"

	"github.com/open-policy-agent/opa/rego"
)

const PACKAGE_HATCHET_MODULE = "hatchet.module"
const PACKAGE_HATCHET_ORGANIZATION = "hatchet.organization"
const PACKAGE_HATCHET_TEAM = "hatchet.team"

func LoadQueryFromFile(name, filepath string) (*OPAQuery, error) {
	fileBytes, err := ioutil.ReadFile(filepath)

	if err != nil {
		return nil, err
	}

	return LoadQueryFromBytes(name, fileBytes)
}

func LoadQueryFromBytes(name string, fileBytes []byte) (*OPAQuery, error) {
	preparedQuery, err := rego.New(
		rego.Query(fmt.Sprintf("data.%s", name)),
		rego.Module(name, string(fileBytes)),
	).PrepareForEval(context.Background())

	if err != nil {
		// Handle error.
		return nil, err
	}

	return &OPAQuery{&preparedQuery}, nil
}
