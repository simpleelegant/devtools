package kubernetes_dashboard

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	k "github.com/simpleelegant/devtools/plugins/kubernetes"
)

func listJobs(c *gin.Context) (int, interface{}) {
	namespace := formValue(c, "namespace")
	labelSelector := formValue(c, "labelSelector")
	client, err := getClient(c)
	if err != nil {
		return http.StatusBadRequest, err
	}

	jobs, err := client.GetJobs(namespace, labelSelector)
	if err != nil {
		return http.StatusBadRequest, err
	}

	var res []map[string]interface{}
	for _, j := range jobs.Items {
		r := map[string]interface{}{
			"namespace":         j.Metadata.Namespace,
			"name":              j.Metadata.Name,
			"creationTimestamp": j.Metadata.CreationTimestamp,
		}

		if j.Status != nil {
			r["succeeded"] = j.Status.Succeeded
		} else {
			r["succeeded"] = 0
		}

		res = append(res, r)
	}

	return http.StatusOK, res
}

func listPods(c *gin.Context) (int, interface{}) {
	namespace := formValue(c, "namespace")
	if namespace == "" {
		namespace = "default"
	}

	labelSelector := formValue(c, "labelSelector")
	client, err := getClient(c)
	if err != nil {
		return http.StatusBadRequest, err
	}

	pods, err := client.ListObjectsOfKindPod(namespace,
		&k.ListObjectsOfKindPodOptions{LabelSelector: labelSelector})
	if err != nil {
		return http.StatusBadRequest, err
	}

	var res []map[string]interface{}
	for _, p := range pods.Items {
		res = append(res, map[string]interface{}{
			"namespace":         p.Metadata.Namespace,
			"name":              p.Metadata.Name,
			"creationTimestamp": p.Metadata.CreationTimestamp,
			"nodeName":          p.Spec.NodeName,
		})
	}

	return http.StatusOK, res
}

func describeJob(c *gin.Context) (int, interface{}) {
	namespace := formValue(c, "namespace")
	name := formValue(c, "name")
	client, err := getClient(c)
	if err != nil {
		return http.StatusBadRequest, err
	}

	// get job
	job, err := client.GetJob(namespace, name)
	if err != nil {
		return http.StatusBadRequest, err
	}

	jobJSON, err := json.MarshalIndent(job, "", "    ")
	if err != nil {
		return http.StatusBadRequest, err
	}

	// get pods
	pods, err := client.ListObjectsOfKindPod(namespace,
		&k.ListObjectsOfKindPodOptions{LabelSelector: "job-name=" + name})
	if err != nil {
		return http.StatusBadRequest, err
	}

	return http.StatusOK, map[string]interface{}{
		"job":  string(jobJSON),
		"pods": pods,
	}
}

func describePod(c *gin.Context) (int, interface{}) {
	namespace := formValue(c, "namespace")
	name := formValue(c, "name")
	client, err := getClient(c)
	if err != nil {
		return http.StatusBadRequest, err
	}

	// get pod
	pod, err := client.GetPod(namespace, name)
	if err != nil {
		return http.StatusBadRequest, err
	}

	podJSON, err := json.MarshalIndent(pod, "", "    ")
	if err != nil {
		return http.StatusBadRequest, err
	}

	// get logs
	logs, err := client.GetPodLog(namespace, name, &k.GetPodLogOptions{})
	if err != nil {
		return http.StatusBadRequest, err
	}

	return http.StatusOK, map[string]interface{}{
		"pod":  string(podJSON),
		"logs": logs,
	}
}

func deleteJob(c *gin.Context) (int, interface{}) {
	namespace := formValue(c, "namespace")
	name := formValue(c, "name")
	client, err := getClient(c)
	if err != nil {
		return http.StatusBadRequest, err
	}

	// delete job
	if err := client.DeleteJob(namespace, name); err != nil {
		return http.StatusBadRequest, err
	}

	// delete relatied pods
	err = client.DeletePods(namespace, &k.DeletePodsOptions{LabelSelector: "job-name=" + name})
	if err != nil {
		return http.StatusBadRequest, err
	}

	return http.StatusOK, "ok"
}

func deletePod(c *gin.Context) (int, interface{}) {
	namespace := formValue(c, "namespace")
	name := formValue(c, "name")
	client, err := getClient(c)
	if err != nil {
		return http.StatusBadRequest, err
	}

	err = client.DeletePod(namespace, name, &k.DeleteOptions{GracePeriodSeconds: 0})
	if err != nil {
		return http.StatusBadRequest, err
	}

	return http.StatusOK, "ok"
}

func formValue(c *gin.Context, key string) string {
	return strings.TrimSpace(c.Request.FormValue(key))
}

func getClient(c *gin.Context) (*k.Client, error) {
	host := formValue(c, "server")
	if host == "" {
		return nil, errors.New("Kubernetes API Server must be not empty")
	}

	host = strings.ToLower(host)
	switch {
	case strings.HasPrefix(host, "http://"):
	case strings.HasPrefix(host, "https://"):
	default:
		host = "http://" + host
	}

	if !strings.HasSuffix(host, "/") {
		host += "/"
	}

	return &k.Client{BaseURL: host}, nil
}
