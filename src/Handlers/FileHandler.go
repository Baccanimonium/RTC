package Handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"path/filepath"
)

func (h Handler) handleUploadFile(c *gin.Context) {
	file, err := c.FormFile("file")

	if err != nil {
		logrus.Print("upload file operation has failed: ", err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "No file is received",
		})
		return
	}

	file.Filename = uuid.New().String() + filepath.Ext(file.Filename)

	pwd, err := os.Getwd()
	if err != nil {
		logrus.Print("failed to get current filepath: ", err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}

	err = c.SaveUploadedFile(file, pwd+"\\public\\"+file.Filename)

	if err != nil {
		logrus.Print("upload file operation has failed: ", err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"fileName": file.Filename,
	})
}
