package controllers

import (
	"CookIt/models"
	"fmt"
	"html/template"
	"net/http"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

const SecretKey = "secret"

type tempUser struct {
	FirstName       string `json:"first_name" binding:"required"`
	LastName        string `json:"last_name" binding:"required"`
	Email           string `json:"email" binding:"required"`
	Password        string `json:"password" binding:"required"`
	ConfirmPassword string `json:"confirmpassword" binding:"required"`
}

func ReturnParameterMissingError(c *gin.Context, parameter string) {
	var err = fmt.Sprintf("Required parameter %s missing.", parameter)
	c.JSON(http.StatusBadRequest, gin.H{"error": err})
}

// @Summary register endpoint is used for customer registration. ( Supervisors/admin can be added only by admin. )
// @Description API Endpoint to register the user with the role of customer.
// @Router /api/v1/register [post]
// @Tags auth
// @Accept json
// @Produce json
func Register(c *gin.Context) {
	var tempUser tempUser
	var Role models.UserRole

	c.Request.ParseForm()
	paramList := []string{"email", "first_name", "last_name", "password", "confirmpassword"}

	for _, param := range paramList {
		if c.PostForm(param) == "" {
			ReturnParameterMissingError(c, param)
		}
	}

	// if err := c.ShouldBindJSON(&tempUser); err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 	return
	// }

	tempUser.Email = template.HTMLEscapeString(c.PostForm("email"))
	tempUser.FirstName = template.HTMLEscapeString(c.PostForm("first_name"))
	tempUser.LastName = template.HTMLEscapeString(c.PostForm("last_name"))
	tempUser.Password = template.HTMLEscapeString(c.PostForm("password"))
	tempUser.ConfirmPassword = template.HTMLEscapeString(c.PostForm("confirmpassword"))

	if tempUser.Password != tempUser.ConfirmPassword {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Both passwords do not match."})
		return
	}

	ispasswordstrong, _ := IsPasswordStrong(tempUser.Password)
	if ispasswordstrong == false {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password is not strong."})
		return
	}

	// Check if the user already exists.
	if DoesUserExist(tempUser.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User already exists."})
		return
	}

	encryptedPassword, error := HashPassword(tempUser.Password)
	if error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Some error occoured."})
		return
	}

	err := models.DB.Where("role= ?", "customer").First(&Role).Error
	if err != nil {
		fmt.Println("err ", err.Error())
		return
	}

	SanitizedUser := models.User{
		FirstName:  tempUser.FirstName,
		LastName:   tempUser.LastName,
		Email:      tempUser.Email,
		Password:   encryptedPassword,
		UserRoleID: Role.Id, //This endpoint will be used only for customer registration.
		CreatedAt:  time.Now(),
		IsActive:   true,
	}

	errs := models.DB.Create(&SanitizedUser).Error
	if errs != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Some error occoured in creating entry to DB."})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"msg": "User created successfully"})

}

type login struct {
	Email    string `form:"email" json:"email" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

// redisClient := Redis.createclient()

// Login godoc
// @Summary Login endpoint is used by the user to login.
// @Description API Endpoint to register the user with the role of customer.
// @Router /api/v1/login [post]
// @Tags auth
// @Accept json
// @Produce json
func Login(c *gin.Context) (interface{}, error) {
	var loginVals login
	// var User User
	var user models.User
	var count int64
	// var user models.User
	if err := c.ShouldBind(&loginVals); err != nil {
		return "", jwt.ErrMissingLoginValues
	}
	email := loginVals.Email
	// First check if the user exist or not...
	models.DB.Where("email = ?", email).First(&user).Count(&count)
	if count == 0 {
		return nil, jwt.ErrFailedAuthentication
	}
	if CheckCredentials(loginVals.Email, loginVals.Password, models.DB) == true {
		return &models.User{
			Email: email,
		}, nil
	}
	// fmt.Println("set value ", loginVals.Email)
	// err := rdb.Set("email", loginVals.Email, 0).Err()
	// if err != nil {
	// 	c.JSON(http.StatusNotFound, gin.H{
	// 		"error": "error in redis",
	// 	})
	// }

	return nil, jwt.ErrFailedAuthentication
}
