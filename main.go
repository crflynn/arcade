package main

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

var DOCROOT = getenvOrPanic("GOLYGLOT_DOCROOT")
var PORT = getenvOrPanic("GOLYGLOT_PORT")
var USERNAME = getenvOrPanic("GOLYGLOT_USERNAME")
var PASSWORD = getenvOrPanic("GOLYGLOT_PASSWORD")

// panic if environment variable is not set
func getenvOrPanic(name string) string {
	value := os.Getenv(name)
	if len(value) == 0 {
		panic(name + " not set")
	}
	return value
}

// Disallow project names from clashing with other endpoints
func invalid_name(name string) bool{
	switch strings.ToLower(name) {
	case
		"login",
		"health":
		return true
	}
	return false
}

// Unpack the contents of the zip file to the folder
func put_docs(c *gin.Context) {
	name := c.Param("project")
	if invalid_name(name) {
		c.String(http.StatusBadRequest, "invalid name")
		return
	}
	path := filepath.Join(DOCROOT, name)

	// delete everything
	err := os.RemoveAll(path)

	// recreate the dir
	err = os.MkdirAll(path, os.ModePerm)

	// gunzip
	gr, err := gzip.NewReader(c.Request.Body)
	defer gr.Close()

	if err != nil {
		c.String(http.StatusBadRequest, "error reading")
		return
	}

	// untar
	tr := tar.NewReader(gr)

	// iterate through all of the files
	for {
		header, err := tr.Next()

		switch {
		// done if EOF
		case err == io.EOF:
			c.String(http.StatusOK, "PUT %s", name)
			return
		case err != nil:
			c.String(http.StatusBadRequest, "error reading")
			return
		// skip nil headers
		case header == nil:
			continue
		}

		// the name of the file or folder
		destination := filepath.Join(path, header.Name)

		// folder / file creation depending on type
		switch header.Typeflag {
		// create the directory
		case tar.TypeDir:
			if _, err := os.Stat(destination); err != nil {
				if err := os.MkdirAll(destination, 0755); err != nil {
					c.String(http.StatusBadRequest, "error creating dir")
					return
				}
			}
		// create the file
		case tar.TypeReg:
			f, err := os.OpenFile(destination, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				c.String(http.StatusBadRequest, "error creating")
				return
			}
			// copy file contents
			if _, err := io.Copy(f, tr); err != nil {
				c.String(http.StatusBadRequest, "error copying")
				return
			}
			_ = f.Close()
		}
	}

}

// Delete the folder at the project name
func delete_docs(c *gin.Context) {
	name := c.Param("project")
	if invalid_name(name) {
		c.String(http.StatusBadRequest, "invalid name")
		return
	}
	path := filepath.Join(DOCROOT, name)

	err := os.RemoveAll(path)

	if err != nil {
		c.String(http.StatusBadRequest, "Error")
	} else {
		c.String(http.StatusOK, "DELETE %s", name)
	}
}

func home(c *gin.Context) {
	c.HTML(http.StatusOK, "index.tmpl", gin.H{})
}

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	router.GET("/", home)

	// protect these resources
	private := router.Group("/docs")
	private.Use(gin.BasicAuth(gin.Accounts{
		USERNAME:    PASSWORD,
	}))
	// PUT request to upload new docs in tar gz file
	private.PUT("/:project", put_docs)
	// DELETE request to remove the project folder
	private.DELETE("/:project", delete_docs)
	// Static server for everything
	private.StaticFS("/", gin.Dir(DOCROOT, true))

	_ = router.Run(":" + PORT)
}