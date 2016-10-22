package kubernetes

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

// ErrNotFound represent resource not found
var ErrNotFound = errors.New("Not Found")

// CreateJob create a job
func (c *Client) CreateJob(namespace string, job *Job) error {
	path := fmt.Sprintf("apis/batch/v1/namespaces/%s/jobs", namespace)

	r, err := ToJSONReader(job)
	if err != nil {
		return err
	}

	resp, err := http.Post(c.BaseURL+path, JSONContentType, r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// decode body
	if resp.StatusCode != http.StatusCreated {
		s := new(Status)
		if err := json.Unmarshal(b, s); err != nil {
			return errors.New(string(b))
		}

		return errors.New(s.Message)
	}

	return json.Unmarshal(b, job)
}

// DeleteJob delete a job
func (c *Client) DeleteJob(namespace string, name string) error {
	path := fmt.Sprintf("apis/batch/v1/namespaces/%s/jobs/%s", namespace, name)

	client := &http.Client{}

	req, err := http.NewRequest(http.MethodDelete, c.BaseURL+path, nil)
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", JSONContentType)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	switch resp.StatusCode {
	case http.StatusOK:
		return nil
	case http.StatusNotFound:
		return ErrNotFound
	default:
		defer resp.Body.Close()
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		// decode body
		result := new(Status)
		if err := json.Unmarshal(b, result); err != nil {
			return errors.New(string(b))
		}

		return errors.New(result.Message)
	}
}

// GetJob ...
func (c *Client) GetJob(namespace, name string) (*Job, error) {
	path := fmt.Sprintf("apis/batch/v1/namespaces/%s/jobs/%s", namespace, name)
	resp, err := http.Get(c.BaseURL + path)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	switch resp.StatusCode {
	case http.StatusOK:
		job := new(Job)
		if err := json.Unmarshal(b, job); err != nil {
			return nil, fmt.Errorf("Unmarshal json to Job fails, %s, json:\n%s", err, b)
		}

		return job, nil
	case http.StatusNotFound:
		return nil, ErrNotFound
	default:
		return nil, fmt.Errorf("%s %s", resp.Status, b)
	}
}

// GetJobs ...
// namespace is optional
func (c *Client) GetJobs(namespace, labelSelector string) (*JobList, error) {
	path := "apis/batch/v1/jobs"

	if namespace != "" {
		path = fmt.Sprintf("apis/batch/v1/namespaces/%s/jobs", namespace)
	}

	if labelSelector != "" {
		v := url.Values{"labelSelector": {labelSelector}}
		path += "?" + v.Encode()
	}

	resp, err := http.Get(c.BaseURL + path)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(string(b))
	}

	// decode body
	jl := new(JobList)
	if err := json.Unmarshal(b, jl); err != nil {
		return nil, fmt.Errorf("Unmarshal json to JobList fails, %s, json: \n%s",
			err.Error(), string(b))
	}

	return jl, nil
}
