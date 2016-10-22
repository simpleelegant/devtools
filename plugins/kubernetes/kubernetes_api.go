package kubernetes

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/google/go-querystring/query"
)

// GetPodLogOptions ...
type GetPodLogOptions struct {
	Pretty       string `url:"pretty,omitempty"`
	Container    string `url:"container,omitempty"`
	Follow       bool   `url:"follow,omitempty"`
	Previous     bool   `url:"previous,omitempty"`
	SinceSeconds int    `url:"sinceSeconds,omitempty"`
	SinceTime    string `url:"sinceTime,omitempty"`
	Timestamps   bool   `url:"timestamps,omitempty"`
	TailLines    int    `url:"tailLines,omitempty"`
	LimitBytes   int    `url:"limitBytes,omitempty"`
	Namespace    string `url:"namespace,omitempty"`
	Name         string `url:"name,omitempty"`
}

// ListObjectsOfKindPodOptions ...
type ListObjectsOfKindPodOptions struct {
	Pretty          string `url:"pretty,omitempty"`
	LabelSelector   string `url:"labelSelector,omitempty"`
	FieldSelector   string `url:"fieldSelector,omitempty"`
	Watch           bool   `url:"watch,omitempty"`
	ResourceVersion string `url:"resourceVersion,omitempty"`
	TimeoutSeconds  int    `url:"timeoutSeconds,omitempty"`
}

// DeletePodsOptions ...
type DeletePodsOptions struct {
	Pretty          string `url:"pretty,omitempty"`
	LabelSelector   string `url:"labelSelector,omitempty"`
	FieldSelector   string `url:"fieldSelector,omitempty"`
	Watch           bool   `url:"watch,omitempty"`
	ResourceVersion string `url:"resourceVersion,omitempty"`
	TimeoutSeconds  int    `url:"timeoutSeconds,omitempty"`
}

// GetPod ...
func (c *Client) GetPod(namespace, name string) (*Pod, error) {
	path := fmt.Sprintf("api/v1/namespaces/%s/pods/%s", namespace, name)
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
		pod := new(Pod)
		if err := json.Unmarshal(b, pod); err != nil {
			return nil, fmt.Errorf("Unmarshal json to Pod fails, %v, json:\n%s", err, b)
		}

		return pod, nil
	case http.StatusNotFound:
		return nil, ErrNotFound
	default:
		s := new(Status)
		if err := json.Unmarshal(b, s); err != nil {
			return nil, fmt.Errorf("%s %s", resp.Status, b)
		}

		return nil, errors.New(s.Message)
	}
}

// GetPodLog read log of the specified Pod
func (c *Client) GetPodLog(namespace, podName string, opt *GetPodLogOptions) (string, error) {
	path := fmt.Sprintf("api/v1/namespaces/%s/pods/%s/log", namespace, podName)

	v, err := query.Values(opt)
	if err != nil {
		return "", err
	}

	resp, err := http.Get(c.BaseURL + path + "?" + v.Encode())
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		s := new(Status)
		if err := json.Unmarshal(b, s); err != nil {
			return "", fmt.Errorf("%s %s", resp.Status, b)
		}

		return "", errors.New(s.Message)
	}

	return string(b), nil
}

// ListObjectsOfKindPod list or watch object of kind Pod
func (c *Client) ListObjectsOfKindPod(namespace string, opt *ListObjectsOfKindPodOptions) (*PodList, error) {
	path := fmt.Sprintf("api/v1/namespaces/%s/pods", namespace)

	v, err := query.Values(opt)
	if err != nil {
		return nil, err
	}

	resp, err := http.Get(c.BaseURL + path + "?" + v.Encode())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// parse response
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s %s", resp.Status, b)
	}
	podList := &PodList{}
	if err := json.Unmarshal(b, podList); err != nil {
		return nil, err
	}

	return podList, nil
}

// DeletePods delete collection of Pod
func (c *Client) DeletePods(namespace string, opt *DeletePodsOptions) error {
	v, err := query.Values(opt)
	if err != nil {
		return err
	}

	path := fmt.Sprintf("api/v1/namespaces/%s/pods?%s", namespace, v.Encode())
	req, err := http.NewRequest(http.MethodDelete, c.BaseURL+path, nil)
	if err != nil {
		return err
	}

	client := &http.Client{}
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

		return fmt.Errorf("%s %s", resp.Status, b)
	}

	return nil
}

// DeletePod delete a Pod
func (c *Client) DeletePod(namespace, name string, opt *DeleteOptions) error {
	b, err := json.Marshal(opt)
	if err != nil {
		return err
	}

	path := fmt.Sprintf("api/v1/namespaces/%s/pods/%s", namespace, name)
	req, err := http.NewRequest(http.MethodDelete, c.BaseURL+path, bytes.NewReader(b))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
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

		return fmt.Errorf("%s %s", resp.Status, b)
	}

	return nil
}
