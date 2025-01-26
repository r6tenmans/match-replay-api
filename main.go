package main

import (
	"archive/zip"
	"bytes"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redraskal/r6-dissect/dissect"
)

type RoundFile struct {
	FileName   string                `json:"fileName"`
	Error      string                `json:"error"`
	Round      dissect.Header        `json:"round"`
	Activities []dissect.MatchUpdate `json:"activities"`
}

func main() {
	router := gin.Default()

	router.Use(CORSMiddleware())
	router.GET("/test", getTest)
	router.POST("/round", postRound)
	router.POST("/replay", postReplayZip)

	router.Run(":8080")
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func getFile(c *gin.Context) (multipart.File, error) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "context": "Error fetching form file"})
		return nil, err
	}

	f, err := file.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "context": "Error opening file stream"})
		return nil, err
	}

	return f, nil
}

func getFileBuffer(c *gin.Context) (*bytes.Buffer, error) {
	f, err := getFile(c)
	if f == nil {
		return nil, err
	}

	defer f.Close()

	var buf bytes.Buffer
	_, err = io.Copy(&buf, f)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "context": "Error copying file to buffer"})
		return nil, err
	}

	return &buf, nil
}

func getFileZip(c *gin.Context) (*zip.Reader, error) {
	f, err := getFileBuffer(c)
	if f == nil {
		return nil, err
	}

	readerAt := bytes.NewReader(f.Bytes())
	archive, err := zip.NewReader(readerAt, int64(f.Len()))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "context": "Error creating zip reader"})
		return nil, err
	}

	return archive, nil
}

func postRound(c *gin.Context) {
	f, _ := getFile(c)
	if f == nil {
		return
	}

	defer f.Close()
	r, err := dissect.NewReader(f)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := r.Read(); !dissect.Ok(err) {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, r)
}

func postReplayZip(c *gin.Context) {
	archive, _ := getFileZip(c)
	if archive == nil {
		return
	}

	var rounds []RoundFile
	for _, file := range archive.File {

		if file.FileInfo().IsDir() {
			continue
		}

		b, err := file.Open()
		if err != nil {
			rounds = append(rounds, RoundFile{FileName: file.Name, Error: err.Error()})
			continue
		}

		defer b.Close()

		r, err := dissect.NewReader(b)
		if err != nil {
			rounds = append(rounds, RoundFile{FileName: file.Name, Error: err.Error()})
			continue
		}
		if err := r.Read(); !dissect.Ok(err) {
			rounds = append(rounds, RoundFile{FileName: file.Name, Error: err.Error()})
			continue
		}

		rounds = append(rounds, RoundFile{FileName: file.Name, Round: r.Header, Activities: r.MatchFeedback})
	}

	c.JSON(http.StatusOK, rounds)
}

func getTest(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Hello, World!"})
}
