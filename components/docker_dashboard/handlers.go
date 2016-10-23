package docker_dashboard

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/fsouza/go-dockerclient"
	"github.com/gin-gonic/gin"
)

func getImages(c *gin.Context) {
	client, ok := getDockerClient(c)
	if !ok {
		return
	}

	images, err := client.ListImages(docker.ListImagesOptions{All: false})
	if err != nil {
		response(c, err.Error(), nil)
		return
	}

	response(c, "", images)
}

func inspectImage(c *gin.Context) {
	id := strings.TrimSpace(c.PostForm("id"))
	client, ok := getDockerClient(c)
	if !ok {
		return
	}

	image, err := client.InspectImage(id)
	if err != nil {
		response(c, err.Error(), nil)
		return
	}
	imageJSON, err := json.MarshalIndent(image, "", "    ")
	if err != nil {
		response(c, err.Error(), nil)
		return
	}

	response(c, "", string(imageJSON))
}

func deleteImage(c *gin.Context) {
	id := strings.TrimSpace(c.PostForm("id"))
	client, ok := getDockerClient(c)
	if !ok {
		return
	}

	// delete image
	if err := client.RemoveImage(id); err != nil {
		response(c, err.Error(), nil)
		return
	}

	response(c, "", nil)
}

func response(c *gin.Context, errorMessage string, data interface{}) {
	c.JSON(http.StatusOK, map[string]interface{}{
		"Error": errorMessage,
		"Data":  data,
	})
}

func getDockerClient(c *gin.Context) (client *docker.Client, ok bool) {
	addr := strings.TrimSpace(c.PostForm("docker_api_server"))
	if addr == "" {
		response(c, "Docker API Server must be not empty", nil)
		return
	}

	a := strings.ToLower(addr)
	if !strings.HasPrefix(a, "http://") && !strings.HasPrefix(a, "https://") {
		addr = "http://" + addr
	}
	if !strings.HasSuffix(addr, "/") {
		addr += "/"
	}

	dc, err := docker.NewClient(addr)
	if err != nil {
		response(c, err.Error(), nil)
		return
	}

	return dc, true
}
