package kubernetes_dashboard

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	k "github.com/simpleelegant/devtools/plugins/kubernetes"
)

func getJobs(c *gin.Context) {
	namespace := strings.TrimSpace(c.PostForm("namespace"))
	labelSelector := strings.TrimSpace(c.PostForm("labelSelector"))
	client, ok := getKubernetesClient(c)
	if !ok {
		return
	}

	jobs, err := client.GetJobs(namespace, labelSelector)
	if err != nil {
		response(c, err.Error(), nil)
		return
	}

	response(c, "", jobs)
}

func describeJob(c *gin.Context) {
	name := strings.TrimSpace(c.PostForm("name"))
	namespace := strings.TrimSpace(c.PostForm("namespace"))
	client, ok := getKubernetesClient(c)
	if !ok {
		return
	}

	// get job
	job, err := client.GetJob(namespace, name)
	if err != nil {
		if err == k.ErrNotFound {
			err = errors.New("Job not found")
		}
		response(c, err.Error(), nil)
		return
	}
	jobJSON, err := json.MarshalIndent(job, "", "    ")
	if err != nil {
		response(c, err.Error(), nil)
		return
	}

	// get pods
	podList, err := client.ListObjectsOfKindPod(
		namespace,
		&k.ListObjectsOfKindPodOptions{LabelSelector: "job-name=" + name})
	if err != nil {
		response(c, err.Error(), nil)
		return
	}

	response(c, "", map[string]interface{}{
		"Job":     string(jobJSON),
		"PodList": podList,
	})
}

func describePod(c *gin.Context) {
	name := strings.TrimSpace(c.PostForm("name"))
	namespace := strings.TrimSpace(c.PostForm("namespace"))
	client, ok := getKubernetesClient(c)
	if !ok {
		return
	}

	// get pod
	pod, err := client.GetPod(namespace, name)
	if err != nil {
		if err == k.ErrNotFound {
			err = errors.New("Pod not found")
		}
		response(c, err.Error(), nil)
		return
	}
	podJSON, err := json.MarshalIndent(pod, "", "    ")
	if err != nil {
		response(c, "Fail to get logs: "+err.Error(), nil)
		return
	}

	// get logs
	logs, err := client.GetPodLog(namespace, name, &k.GetPodLogOptions{})
	if err != nil {
		response(c, err.Error(), nil)
		return
	}

	response(c, "", map[string]interface{}{
		"Pod":  string(podJSON),
		"Logs": logs,
	})
}

func deleteJob(c *gin.Context) {
	name := strings.TrimSpace(c.PostForm("name"))
	namespace := strings.TrimSpace(c.PostForm("namespace"))
	client, ok := getKubernetesClient(c)
	if !ok {
		return
	}

	// delete job
	err := client.DeleteJob(namespace, name)
	switch err {
	case nil:
		// delete relatied pods
		err := client.DeletePods(namespace, &k.DeletePodsOptions{LabelSelector: "job-name=" + name})
		if err != nil {
			response(c, err.Error(), nil)
			return
		}

		response(c, "", nil)
	case k.ErrNotFound:
		response(c, "Job not found", nil)
	default:
		response(c, err.Error(), nil)
	}
}

func deletePod(c *gin.Context) {
	name := strings.TrimSpace(c.PostForm("name"))
	namespace := strings.TrimSpace(c.PostForm("namespace"))
	client, ok := getKubernetesClient(c)
	if !ok {
		return
	}

	err := client.DeletePod(namespace, name, &k.DeleteOptions{GracePeriodSeconds: 0})
	switch err {
	case nil:
		response(c, "", nil)
	case k.ErrNotFound:
		response(c, "Pod not found", nil)
	default:
		response(c, err.Error(), nil)
	}
}

func response(c *gin.Context, errorMessage string, data interface{}) {
	c.JSON(http.StatusOK, map[string]interface{}{
		"Error": errorMessage,
		"Data":  data,
	})
}

func getKubernetesClient(c *gin.Context) (client *k.Client, ok bool) {
	addr := strings.TrimSpace(c.PostForm("kubernetes_api_server"))
	if addr == "" {
		response(c, "Kubernetes API Server must be not empty", nil)
		return
	}

	a := strings.ToLower(addr)
	if !strings.HasPrefix(a, "http://") && !strings.HasPrefix(a, "https://") {
		addr = "http://" + addr
	}
	if !strings.HasSuffix(addr, "/") {
		addr += "/"
	}

	return &k.Client{BaseURL: addr}, true
}
