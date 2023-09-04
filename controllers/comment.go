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

func CreateComment(c *gin.Context) {
	db := database.GetDB()
	userData := c.MustGet("userData").(jwt.MapClaims)
	contentType := helpers.GetContentType(c)
	_, _ = db, contentType
	Comment := models.Comment{}
	userID := uint(userData["id"].(float64))
	if contentType == appJSON {
		c.ShouldBindJSON(&Comment)
	} else {
		c.ShouldBind(&Comment)
	}
	Comment.UserID = userID

	Result := map[string]interface{}{}

	SqlStatement := `
		INSERT into comments 
		(user_id, photo_id, message, created_at) VALUES(?,?,?,?) 
		Returning id,photo_id, user_id, message, created_at
	`

	err := db.Raw(
		SqlStatement,
		Comment.UserID, Comment.PhotoID, Comment.Message, time.Now(),
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

func UpdateComment(c *gin.Context) {
	db := database.GetDB()
	contentType := helpers.GetContentType(c)
	commentId, _ := strconv.Atoi(c.Param("commentId"))
	Comment := models.Comment{}
	if contentType == appJSON {
		c.ShouldBindJSON(&Comment)
	} else {
		c.ShouldBind(&Comment)
	}
	Result := map[string]interface{}{}

	SqlStatement := `
		Update comments 
		SET message = ?, updated_at = ? 
		WHERE id = ? 
		RETURNING id, message, updated_at`
	err := db.Raw(
		SqlStatement,
		Comment.Message, time.Now(), commentId,
	).Scan(&Result).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err":     "Bad Request",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Result)
}

func GetAllComments(c *gin.Context) {
	db := database.GetDB()
	Comments := []models.Comment{}
	err := db.Debug().Preload("User").Find(&Comments).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err":     "Bad Request",
			"message": err.Error(),
		})
		return
	}
	result := []interface{}{}
	_ = result

	for _, v := range Comments {
		temp := map[string]interface{}{
			"id":         v.ID,
			"created_at": v.CreatedAt,
			"message":    v.Message,
			"User": map[string]interface{}{
				"username": v.User.Username,
				"email":    v.User.Email,
			},
		}
		result = append(result, temp)
	}

	c.JSON(http.StatusOK, result)
}

func GetOneComment(c *gin.Context) {
	db := database.GetDB()
	Comments := []models.Comment{}
	commentId, _ := strconv.Atoi(c.Param("commentId"))
	err := db.Debug().Preload("User").Find(&Comments, commentId).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err":     "Bad Request",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Comments)
}

func DeleteComment(c *gin.Context) {
	db := database.GetDB()
	commentId, _ := strconv.Atoi(c.Param("commentId"))
	Comment := models.Comment{}
	err := db.Delete(Comment, uint(commentId)).Error

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
