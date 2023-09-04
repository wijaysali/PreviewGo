package controllers

import (
	"net/http"
	"photo-app/database"
	"photo-app/helpers"
	"photo-app/models"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func CreatePhoto(c *gin.Context) {
	db := database.GetDB()
	userData := c.MustGet("userData").(jwt.MapClaims)
	contentType := helpers.GetContentType(c)
	_, _ = db, contentType
	Photo := models.Photo{}
	userID := uint(userData["id"].(float64))
	if contentType == appJSON {
		c.ShouldBindJSON(&Photo)
	} else {
		c.ShouldBind(&Photo)
	}
	Photo.UserID = userID

	Result := map[string]interface{}{}

	err := db.Raw(
		"INSERT into photos (title, caption, photo_url, user_id, created_at) VALUES(?,?,?,?,?) Returning id,title,photo_url, user_id, created_at, caption",
		Photo.Title, Photo.Caption, Photo.PhotoUrl, Photo.UserID, time.Now(),
	).Scan(&Result).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err":     "Bad Request",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, Result)

}

func UpdatePhoto(c *gin.Context) {
	db := database.GetDB()
	userData := c.MustGet("userData").(jwt.MapClaims)
	contentType := helpers.GetContentType(c)
	_, _ = db, contentType
	Photo := models.Photo{}
	photoId, _ := strconv.Atoi(c.Param("photoId"))
	userID := uint(userData["id"].(float64))
	if contentType == appJSON {
		c.ShouldBindJSON(&Photo)
	} else {
		c.ShouldBind(&Photo)
	}
	Photo.UserID = userID

	Result := map[string]interface{}{}
	SqlStatement := "Update photos SET title = ?, caption = ?, photo_url = ?, updated_at = ? WHERE id = ? RETURNING id, title, caption, photo_url, updated_at"
	err := db.Raw(
		SqlStatement,
		Photo.Title, Photo.Caption, Photo.PhotoUrl, time.Now(), uint(photoId),
	).Scan(&Result).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err":     "Bad Request",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, Result)
}

func GetAllPhotos(c *gin.Context) {
	db := database.GetDB()
	Photos := []models.Photo{}
	err := db.Debug().Preload("User").Preload("Comments").Find(&Photos).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err":     "Bad Request",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Photos)
}
func GetOnePhoto(c *gin.Context) {
	db := database.GetDB()
	photoId, _ := strconv.Atoi(c.Param("photoId"))
	Photos := models.Photo{}
	err := db.Preload("User").Find(&Photos, uint(photoId)).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err":     "Bad Request",
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, Photos)
}

func DeletePhoto(c *gin.Context) {
	db := database.GetDB()
	photoId, _ := strconv.Atoi(c.Param("photoId"))
	Photo := models.Photo{}
	err := db.Delete(Photo, uint(photoId)).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err":     "Bad Request",
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Your photo has been successfully deleted",
	})
}

// `db.Raw("SELECT a.serviceId as service_id, a.serviceName as service_name, b.systemId as system_id, b.systemName as system_name FROM go_service_info a LEFT JOIN go_system_info b ON a.systemId = b.systemId").Scan(&results)`
