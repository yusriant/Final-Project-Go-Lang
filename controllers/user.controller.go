package controllers

import (
	"net/http"

	"final-project-golang/models"
	"final-project-golang/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserController struct {
	DB *gorm.DB
}

func NewUserController(DB *gorm.DB) UserController {
	return UserController{DB}
}

// func (uc *UserController) GetMe(ctx *gin.Context) {
// 	currentUser := ctx.MustGet("currentUser").(models.User)

// 	userResponse := &models.UserResponse{
// 		ID:        currentUser.ID,
// 		Username:  currentUser.Username,
// 		Email:     currentUser.Email,
// 		CreatedAt: currentUser.CreatedAt,
// 		UpdatedAt: currentUser.UpdatedAt,
// 	}

// 	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"user": userResponse}})
// }

func (uc *UserController) UpdateMe(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)
	var payload models.UpdateCurrentUserRequest

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	// Validasi email
	if !utils.IsValidEmail(payload.Email) {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid email format"})
		return
	}

	// Validasi username
	if payload.Username == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Username is required"})
		return
	}

	// Validasi age
	if payload.Age < 8 {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Age must be at least 8"})
		return
	}

	// Validasi URL profil
	if payload.ProfileImageURL != "" && !utils.IsValidURL(payload.ProfileImageURL) {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid profile image URL format"})
		return
	}

	// Update current user information
	currentUser.Username = payload.Username
	currentUser.Email = payload.Email
	currentUser.Age = payload.Age
	currentUser.ProfileImageURL = payload.ProfileImageURL

	// Validasi Unik Username
	existingUser := models.User{}
	if err := uc.DB.Where("username = ?", payload.Username).First(&existingUser).Error; err == nil && existingUser.ID != currentUser.ID {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Username is already taken"})
		return
	}

	// Validasi Unik Email
	if err := uc.DB.Where("email = ?", payload.Email).First(&existingUser).Error; err == nil && existingUser.ID != currentUser.ID {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Email is already taken"})
		return
	}

	// Simpan perubahan ke dalam database
	if err := uc.DB.Save(&currentUser).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update user information"})
		return
	}

	// Siapkan respons sesuai spesifikasi OpenAPI
	responseData := gin.H{
		"id":                currentUser.ID,
		"email":             currentUser.Email,
		"username":          currentUser.Username,
		"age":               currentUser.Age,
		"profile_image_url": currentUser.ProfileImageURL,
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": responseData})
}

func (uc *UserController) DeleteMe(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)

	// Hapus pengguna saat ini dari database
	if err := uc.DB.Delete(&currentUser).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to delete user"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}
