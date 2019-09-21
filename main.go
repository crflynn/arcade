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
	"github.com/zalando/gin-oauth2/github"
)

// const DOCROOT = "/tmp/docs"
const DOCROOT = "_docs"
const PORT = "6060"
var PROTO = os.Getenv("GOLYGLOT_PROTO")
var HOST = os.Getenv("GOLYGLOT_HOST")
var SECRET = os.Getenv("GOLYGLOT_SECRET")
var ALLOWED_TEAMS = os.Getenv("GITHUB_ALLOWED_TEAMS")

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

func main() {
	redirectURL := PROTO + "://" + HOST + ":" + PORT + "/docs"
	credFile := "./github.json"
	scopes :=[]string{"user"}
	secret := []byte(SECRET)
	sessionName := "golyglotsession"
	github.Setup(redirectURL, credFile, scopes, secret)

	router := gin.Default()
	router.Use(github.Session(sessionName))
	router.GET("/", github.LoginHandler)

	private := router.Group("/docs")

	// protect these resources
	// TODO wrap this handler to check github teams membership
	private.Use(github.Auth())
	// PUT request to upload new docs in tar gz file
	private.PUT("/:project", put_docs)
	// DELETE request to kill the project folder
	private.DELETE("/:project", delete_docs)
	// Static server for everything
	private.StaticFS("/", gin.Dir(DOCROOT, true))

	_ = router.Run(":" + PORT)
}
