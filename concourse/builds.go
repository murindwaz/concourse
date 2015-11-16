package concourse

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/concourse/atc"
)

func (client *client) CreateBuild(plan atc.Plan) (atc.Build, error) {
	var build atc.Build

	buffer := &bytes.Buffer{}
	err := json.NewEncoder(buffer).Encode(plan)
	if err != nil {
		return build, fmt.Errorf("Unable to marshal plan: %s", err)
	}

	err = client.connection.Send(Request{
		RequestName: atc.CreateBuild,
		Body:        buffer,
		Headers: map[string][]string{
			"Content-Type": {"application/json"},
		},
	}, &Response{
		Result: &build,
	})

	if ure, ok := err.(UnexpectedResponseError); ok {
		if ure.StatusCode == http.StatusNotFound {
			return build, errors.New("build not found")
		}
	}

	return build, err
}

func (client *client) CreateJobBuild(pipelineName string, jobName string) (atc.Build, error) {
	params := map[string]string{"job_name": jobName, "pipeline_name": pipelineName}

	var build atc.Build
	err := client.connection.Send(Request{
		RequestName: atc.CreateJobBuild,
		Params:      params,
	}, &Response{
		Result: &build,
	})

	return build, err
}

func (client *client) JobBuild(pipelineName, jobName, buildName string) (atc.Build, bool, error) {
	if pipelineName == "" {
		return atc.Build{}, false, NameRequiredError("pipeline")
	}

	params := map[string]string{"job_name": jobName, "build_name": buildName, "pipeline_name": pipelineName}
	var build atc.Build
	err := client.connection.Send(Request{
		RequestName: atc.GetJobBuild,
		Params:      params,
	}, &Response{
		Result: &build,
	})

	switch err.(type) {
	case nil:
		return build, true, nil
	case ResourceNotFoundError:
		return build, false, nil
	default:
		return build, false, err
	}
}

func (client *client) Build(buildID string) (atc.Build, bool, error) {
	params := map[string]string{"build_id": buildID}
	var build atc.Build
	err := client.connection.Send(Request{
		RequestName: atc.GetBuild,
		Params:      params,
	}, &Response{
		Result: &build,
	})

	switch err.(type) {
	case nil:
		return build, true, nil
	case ResourceNotFoundError:
		return build, false, nil
	default:
		return build, false, err
	}
}

func (client *client) AllBuilds() ([]atc.Build, error) {
	var builds []atc.Build
	err := client.connection.Send(Request{
		RequestName: atc.ListBuilds,
	}, &Response{
		Result: &builds,
	})
	return builds, err
}

func (client *client) AbortBuild(buildID string) error {
	params := map[string]string{"build_id": buildID}
	return client.connection.Send(Request{
		RequestName: atc.AbortBuild,
		Params:      params,
	}, nil)
}
