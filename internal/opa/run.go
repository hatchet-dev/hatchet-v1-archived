package opa

import (
	"context"
	"fmt"

	"github.com/hatchet-dev/hatchet/api/v1/types"
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

// RunMonitorQuery runs a query with output that can be decoded to MonitorActionOutput
func RunMonitorQuery(query *OPAQuery, input map[string]interface{}) (*types.CreateMonitorResultRequest, error) {
	results, err := query.Eval(
		context.Background(),
		rego.EvalInput(input),
	)

	if err != nil {
		return nil, err
	}

	if len(results) == 1 {
		queryRes := &types.CreateMonitorResultRequest{}

		err = mapstructure.Decode(results[0].Expressions[0].Value, queryRes)

		if err != nil {
			return nil, err
		}

		return queryRes, nil
	}

	return nil, fmt.Errorf("no results returned")
}
