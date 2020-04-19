package main

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

var DOCROOT = getenvOrPanic("ARCADE_DOCROOT")
var PORT = getenvOrPanic("ARCADE_PORT")
var USERNAME = os.Getenv("ARCADE_USERNAME")
var PASSWORD = os.Getenv("ARCADE_PASSWORD")

// panic if environment variable is not set
func getenvOrPanic(name string) string {
	value := os.Getenv(name)
	if len(value) == 0 {
		panic(name + " not set")
	}
	return value
}

// Unpack the contents of the zip file to the folder
func putDocs(c *gin.Context) {
	name := c.Param("project")
	version := c.Param("version")

	path := filepath.Join(DOCROOT, name)
	path = filepath.Join(path, version)

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
			c.String(http.StatusOK, "PUT %s", name+"/"+version)
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
func deleteDocs(c *gin.Context) {
	name := c.Param("project")
	version := c.Param("version")

	path := filepath.Join(DOCROOT, name)
	if version != "" {
		path = filepath.Join(path, version)
	}

	err := os.RemoveAll(path)

	if err != nil {
		c.String(http.StatusBadRequest, "Error")
	} else {
		c.String(http.StatusOK, "DELETE %s", name+"/"+version)
	}
}

func main() {
	router := gin.Default()
	router.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})
	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/docs")
	})

	// optionally protected resources
	private := router.Group("/docs")
	if USERNAME != "" {
		private.Use(gin.BasicAuth(gin.Accounts{
			USERNAME: PASSWORD,
		}))
	}

	// PUT request to upload new docs in tar gz file
	private.PUT("/:project/:version", putDocs)
	// DELETE request to remove the project+version docs
	private.DELETE("/:project/:version", deleteDocs)
	// DELETE request to remove the entire project
	private.DELETE("/:project", deleteDocs)

	// Static server for docs
	private.StaticFS("/", gin.Dir(DOCROOT, true))

	_ = router.Run(":" + PORT)
}
