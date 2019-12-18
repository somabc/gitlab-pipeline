package gitlab

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"
)

// Pipeline models a gitlab pipeline entity
type Pipeline struct {
	ID         int        `json:"id"`
	WebURL     string     `json:"web_url"`
	Tag        bool       `json:"tag"`
	Status     string     `json:"status"`
	BeforeSHA  string     `json:"before_sha"`
	FinishedAt *time.Time `json:"finished_at"`
	SHA        string     `json:"sha"`
	User       User       `json:"user"`
	CreatedAt  string     `json:"created_at"`
	StartedAt  string     `json:"started_at"`
	Coverage   *string    `json:"coverage"`
	Duration   *int       `json:"duration"`
}

// TriggerPipeline calls gitlab APIs to run a pipeline in a project
// a Trigger Token is required to trigger the pipeline.
// Extra parameters can be passed to the pipeline, to be part of it's
// environment variables
func (c *Client) TriggerPipeline(projectID int, triggerToken string,
	targetBranch string,
	extraParams map[string]string) (*Pipeline, error){

	data := url.Values{
		"token": {triggerToken},
		"ref":   {targetBranch},
	}
	// append extra params
	for k, v := range extraParams {
		data.Add(k, v)
	}

	triggerURI := buildTriggerURL(c.Host, projectID)

	resp, err := c.client.PostForm(triggerURI, data)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		err := fmt.Errorf("failed to trigger pipeline, got status: %d", resp.StatusCode)
		return nil, err
	}

	pipeline := Pipeline{}
	err = json.NewDecoder(resp.Body).Decode(&pipeline)
	if err != nil {
		return nil, err
	}

	return &pipeline, nil
}

// CheckPipelineStatus calls gitlab APIs to check the status of a pipeline inside a project
func (c *Client) CheckPipelineStatus(pipelineID int, projectID int) (*Pipeline, error) {
	statusAPI := buildPipelineStatusURL(c.Host, projectID, pipelineID)

	req, err := http.NewRequest("GET", statusAPI, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("PRIVATE-TOKEN", c.APIToken)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var pipeline Pipeline

	err = json.NewDecoder(resp.Body).Decode(&pipeline)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if resp.StatusCode >= 300 {
		err = fmt.Errorf("failed to get pipeline status, got error code: %d", resp.StatusCode)
		return nil, err
	}

	return &pipeline, nil
}

func buildTriggerURL(host string, projectID int) string {
	return fmt.Sprintf("%s/api/v4/projects/%d/trigger/pipeline", host, projectID)
}

func buildPipelineStatusURL(host string, projectID int, pipelineID int) string {
	return fmt.Sprintf("%s/api/v4/projects/%d/pipelines/%d", host, projectID, pipelineID)
}