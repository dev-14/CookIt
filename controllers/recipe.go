package controllers

import (
	"CookIt/models"
	"fmt"
	"net/http"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

type tempRecipe struct {
	Name        string
	Description string
	Ingredients map[string]string
}

// CreateBook godoc
// @Summary CreateBook endpoint is used by the supervisor role user to create a new book.
// @Description CreateBook endpoint is used by the supervisor role user to create a new book
// @Router /api/v1/auth/books/create [post]
// @Tags book
// @Accept json
// @Produce json
// @Param name formData string true "name of the book"
// @Param category_id formData string true "category_id of the book"
func CreateRecipe(c *gin.Context) {

	var existingRecipe models.Recipe
	var tempRecipe models.Recipe
	claims := jwt.ExtractClaims(c)
	user_email, _ := claims["email"]
	//var ingredients []string
	//var User models.User
	//var category models.Category
	// user_email, _ := Rdb.HGet("user", "email").Result()
	fmt.Println(user_email)
	// // Check if the current user had admin role.
	// if err := models.DB.Where("email = ? AND user_role_id=2", user_email).First(&User).Error; err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Recipe can only be added by user"})
	// 	return
	// }

	// id, _ := models.Rdb.HGet("user", "ID").Result()

	// ID, _ := strconv.Atoi(id)
	// roleId, _ := models.Rdb.HGet("user", "RoleID").Result()

	// if roleId != "2" {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Books can only be added by supervisor"})
	// 	return
	// }

	// c.Request.ParseForm()

	// if c.PostForm("name") == "" {
	// 	ReturnParameterMissingError(c, "name")
	// 	return
	// }
	// // if c.PostForm("category_id") == "" {
	// // 	ReturnParameterMissingError(c, "category_id")
	// // 	return
	// // }

	// name := template.HTMLEscapeString(c.PostForm("name"))
	// description := template.HTMLEscapeString(c.PostForm("description"))
	// ingredients = template.HTMLEscapeString(c.PostForm("ingredients"))

	//unmarshaled := json.Unmarshal([]byte(ingred[]), &ingredients)

	//ingredients = template.HTMLEscapeString(c.PostForm("ingredients"))

	if err := c.ShouldBindJSON(&tempRecipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var b []byte
	// json.Unmarshal(b, &tempRecipe)
	fmt.Println(tempRecipe)
	//Check if the product already exists.
	err := models.DB.Where("name = ?", tempRecipe.Name).First(&existingRecipe).Error
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "recipe already exists."})
		return
	}

	// Check if the category exists
	// err = models.DB.First(&category, category_id).Error
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "category does not exists."})
	// 	return
	// }

	// book := models.Recipe{
	// 	Id:          0,
	// 	Name:        name,
	// 	Description: description,
	// 	Ingredients: ingredients,
	// 	CreatedAt:   time.Time{},
	// 	UpdatedAt:   time.Time{},
	// }

	fmt.Println(b)

	newRecipe := models.Recipe{
		Name:        tempRecipe.Name,
		Description: tempRecipe.Description,
		Ingredients: tempRecipe.Ingredients,
	}
	_ = newRecipe

	// err = models.DB.Create(&newRecipe).Error
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 	return
	// }
	c.JSON(http.StatusCreated, gin.H{
		"message": "recipe uploaded successfully",
	})

}

// func (h Hstore) Value() (driver.Value, error) {
// 	hstore := hstore.Hstore{Map: map[string]sql.NullString{}}
// 	if len(h) == 0 {
// 		return nil, nil
// 	}

// 	for key, value := range h {
// 		var s sql.NullString
// 		if value != nil {
// 			s.String = *value
// 			s.Valid = true
// 		}
// 		hstore.Map[key] = s
// 	}
// 	return hstore.Value()
// }
