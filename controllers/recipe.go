package controllers

import (
	"CookIt/models"
	"bufio"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

type temp struct {
	key   string
	value string
}

type JSON json.RawMessage

// Scan scan value into Jsonb, implements sql.Scanner interface
func (j *JSON) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}

	result := json.RawMessage{}
	err := json.Unmarshal(bytes, &result)
	*j = JSON(result)
	return err
}

// Value return json value, implement driver.Valuer interface
func (j JSON) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}
	return json.RawMessage(j).MarshalJSON()
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
	//var ingred JSON
	//var tempRecipe models.Recipe
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

	name := template.HTMLEscapeString(c.PostForm("name"))
	description := template.HTMLEscapeString(c.PostForm("description"))
	ingredient := template.HTMLEscapeString(c.PostForm("ingredients"))

	var jsonMap json.RawMessage
	json.Unmarshal([]byte(ingredient), &jsonMap)
	//unmarshaled := json.Unmarshal([]byte(ingred[]), &ingredients)

	//ingredients = template.HTMLEscapeString(c.PostForm("ingredients"))
	//c.Value()
	// func (c models.Recipe) Value() (driver.Value, error) {
	// 	return json.Marshal(c)
	// }

	// func (c *models.Recipe) Scan(value interface{}) error {
	// 	b, ok := value.([]byte)
	// 	if !ok {
	// 	  return errors.New("type assertion to []byte failed")
	// 	}
	// 	return json.Unmarshal(b, &c)
	// }
	//fmt.Println(name, description)

	// if err := c.BindJSON(&tempRecipe); err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 	return
	// }

	//var b []byte
	// json.Unmarshal(b, &tempRecipe)
	// fmt.Println(tempRecipe)
	// //Check if the product already exists.
	err := models.DB.Where("name = ?", name).First(&existingRecipe).Error
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "recipe already exists."})
		return
	}
	// r := csv.NewReader(strings.NewReader(ingredient))
	// list, e := r.ReadAll()
	// if e != nil {
	// 	panic(e)
	// }
	// var ing []string
	// for _, List := range list {
	// 	ing = append(ing, List...)
	// }
	var ing []string
	// reader := strings.NewReader(ingredient)
	// scanner := bufio.NewScanner(reader)

	// // optionally, resize scanner's capacity for lines over 64K, see next example
	// for scanner.Scan() {
	// 	ing = append(ing, scanner.Text())
	// 	fmt.Println(scanner.Text())
	// }

	// if err := scanner.Err(); err != nil {
	// 	log.Fatal(err)
	// }

	reader := strings.NewReader(ingredient)
	in := bufio.NewScanner(reader)
	for in.Scan() {
		line := in.Text()
		if len(line) == 1 {
			// Group Separator (GS ^]): ctrl-]
			if line[0] == '\x1D' {
				break
			}
		}
		ing = append(ing, line)
	}

	//ing = strings.Split(ingredient, "\n")
	//result, err := in.ReadString('\n')
	//fmt.Println(ing)

	// Check if the category exists
	// err = models.DB.First(&category, category_id).Error
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "category does not exists."})
	// 	return
	// }

	//fmt.Println(b)
	fmt.Println(name, description, ing)

	newRecipe := models.Recipe{
		Name:        name,
		Description: description,
		Ingredients: ing,
	}

	err = models.DB.Create(&newRecipe).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
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
