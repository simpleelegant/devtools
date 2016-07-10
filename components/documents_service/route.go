package documentsservice

import (
	"errors"
	"io/ioutil"
	"net/http"
	"path"
	"strings"

	"github.com/simpleelegant/devtools/plugins/conf"

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
		path, err := getProjectPath(project)
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

		projectPath, err := getProjectPath(project)
		if err != nil {
			c.JSON(http.StatusOK, map[string]interface{}{"Error": err.Error()})
			return
		}

		filePath, err := getFilePath(projectPath, file)
		if err != nil {
			c.JSON(http.StatusOK, map[string]interface{}{
				"Error":   err.Error(),
				"Project": project,
			})
			return
		}

		responseFile(c, project, file, filePath)
	})
}

func getProjectPath(project string) (string, error) {
	for _, v := range conf.Options.DocumentsService {
		if v.Project == project {
			return v.Path, nil
		}
	}

	return "", errors.New("指定的项目不存在")
}

func getFilePath(projectPath, file string) (string, error) {
	if !strings.HasSuffix(projectPath, "/") {
		projectPath += "/"
	}

	filePath := projectPath + file
	if _, err := ioutil.ReadFile(filePath); err != nil {
		return "", errors.New(file + " 不存在")
	}

	return filePath, nil
}

func responseFile(c *gin.Context, project, file, filePath string) {
	var err, content string

	switch strings.ToLower(path.Ext(file)) {
	case ".md", ".mdown", ".markdown":
		data, _ := ioutil.ReadFile(filePath)
		content = string(blackfriday.MarkdownCommon(data))
		goto JSON
	case ".png", ".jpg", ".jpeg", ".gif":
		if c.Query("binary") == "" {
			content = `<img src="` + c.Request.URL.String() + `&binary=1">`
			goto JSON
		}

		c.File(filePath)
		return
	case ".html", ".htm":
		if c.Query("binary") == "" {
			content = `<iframe src="` + c.Request.URL.String() + `&binary=1">点这里在新窗口打开此文件</iframe>`
			goto JSON
		}

		c.File(filePath)
		return
	default:
		err = "此文件不支持在线查看"
		goto JSON
	}

JSON:
	c.JSON(http.StatusOK, map[string]interface{}{
		"Error":   err,
		"Project": project,
		"File":    file,
		"Content": content,
	})
}
