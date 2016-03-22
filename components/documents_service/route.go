package documentsservice

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strings"

	"yujian/devtools/plugins/conf"

	"github.com/gin-gonic/gin"
	"github.com/russross/blackfriday"
)

// Route register routes
func Route(r *gin.Engine) {
	r.GET("/documents_service/", func(c *gin.Context) {
		c.File("./components/documents_service/index.html")
	})

	r.GET("/documents_service/document.html", func(c *gin.Context) {
		c.File("./components/documents_service/document.html")
	})

	r.GET("/documents_service/project-list", func(c *gin.Context) {
		var projects []string

		for _, v := range conf.Options.DocumentsService {
			projects = append(projects, v.Project)
		}

		c.JSON(http.StatusOK, projects)
	})

	r.GET("/documents_service/project-doc-list", func(c *gin.Context) {
		project := c.Query("project")
		path, err := getPath(project)
		if err != nil {
			c.JSON(http.StatusOK, map[string]interface{}{"Error": err.Error()})
			return
		}

		fileInfos, _ := ioutil.ReadDir(path)
		var files []string
		for _, f := range fileInfos {
			if !f.IsDir() {
				files = append(files, f.Name())
			}
		}

		c.JSON(http.StatusOK, map[string]interface{}{
			"Project": project,
			"Files":   files,
		})
	})

	r.GET("/documents_service/project-doc", func(c *gin.Context) {
		project := c.Query("project")
		file := c.Query("file")

		if file == "" {
			c.JSON(http.StatusOK, map[string]interface{}{
				"Error":   "请指定一个文档",
				"Project": project,
			})
			return
		}

		path, err := getPath(project)
		if err != nil {
			c.JSON(http.StatusOK, map[string]interface{}{"Error": err.Error()})
			return
		}

		content, err := getMarkdownFileToHTML(path, file)
		if err != nil {
			c.JSON(http.StatusOK, map[string]interface{}{"Error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, map[string]interface{}{
			"Error":   "",
			"Project": project,
			"File":    file,
			"Content": string(blackfriday.MarkdownCommon([]byte(content))),
		})
	})
}

func getPath(project string) (string, error) {
	for _, v := range conf.Options.DocumentsService {
		if v.Project == project {
			return v.Path, nil
		}
	}

	return "", errors.New("指定的项目不存在")
}

func getMarkdownFileToHTML(dir, file string) (string, error) {
	if !strings.HasSuffix(dir, "/") {
		dir += "/"
	}

	data, err := ioutil.ReadFile(dir + file)
	if err != nil {
		return "", errors.New(file + " 不存在")
	}

	return string(data), nil
}
