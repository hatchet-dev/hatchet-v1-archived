package opa

import (
	"context"

	"github.com/mitchellh/mapstructure"
	"github.com/open-policy-agent/opa/rego"
)

type allowQueryRes struct {
	Allow bool `mapstructure:"allow"`
}

type OPAQuery struct {
	*rego.PreparedEvalQuery
}

// RunAllowQuery runs a simple query with a simple "allow" results that evaluates to true or
// false
func RunAllowQuery(query *OPAQuery, input map[string]interface{}) (bool, error) {
	results, err := query.Eval(
		context.Background(),
		rego.EvalInput(input),
	)

	if err != nil {
		return false, err
	}

	if len(results) == 1 {
		queryRes := &allowQueryRes{}

		err = mapstructure.Decode(results[0].Expressions[0].Value, queryRes)

		if err != nil {
			return false, err
		}

		return queryRes.Allow, nil
	}

	return false, nil
}
