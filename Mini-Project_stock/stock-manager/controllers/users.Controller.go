package controllers

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"example.com/stock-manager/initializers"
	"example.com/stock-manager/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

func Signup(c *gin.Context) {
	// register new user
	// get the email/password off request body

	var body struct {
		Fname    string
		Lname    string
		Phone    string
		Email    string `gorm:"unique"`
		Password string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	// Hash the password (returns as a byte slice of arbitrary length)
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to hash password",
		})
		return
	}

	// Create the new user
	user := models.User{
		Fname:    body.Fname,
		Lname:    body.Lname,
		Phone:    body.Phone,
		Email:    body.Email,
		Password: string(hash),
	}
	result := initializers.DB.Create(&user)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create user",
		})
		return
	}

	// Responde request
	c.JSON(http.StatusOK, user.ID)
}

func Login(c *gin.Context) {
	// user can login
	// Get the email and pass off req body
	var body struct {
		Email    string
		Password string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	// Check requested user
	var user models.User
	initializers.DB.First(&user, "email = ?", body.Email)

	if user.ID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid email or password",
		})
		return
	}

	// Check sent password with user password hash
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid email or password",
		})
		return
	}

	// Generate a jwt token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":         user.ID,
		"expirationDate": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create token",
		})
		return
	}

	// Store the token in tokens table
	t := models.Token{
		Token:  tokenString,
		UserID: user.ID,
		Valid:  true,
	}

	var conn = initializers.DB.Save(&t)
	if conn.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Se fodeu",
		})
		return
	}

	// Send it back
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600*24*30, "", "", false, true)
	c.JSON(http.StatusOK, gin.H{})
}

func Logout(c *gin.Context) {
	// check credentials
	tokenString, err := c.Cookie("Authorization")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid credentials",
		})
		return
	}

	// Get token from database
	var t models.Token = models.Token{}
	var conn = initializers.DB.Table("tokens").Where("token=?", tokenString).First(&t)
	if conn.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Failed to get token from database",
		})
		return
	}

	// Set token as invalid
	t.Valid = false
	initializers.DB.Save(&t)
	if conn.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to store token in database",
		})
		return
	}

	c.Status(http.StatusOK)
}

func VerifyUser(c *gin.Context) {
	// check if someone is logedin and who it is
	user, _ := c.Get("user")

	c.JSON(http.StatusOK, gin.H{
		"message": user,
	})
}

func ListUsers(c *gin.Context) {
	// list all existing users

	//return list of users
	var users []models.User
	result := initializers.DB.Find(&users)

	if result != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": users,
		})
	}
}

func GetUser(c *gin.Context) {

	email := c.Param("email")

	// Get users from DB to check its status after update

	var user models.User = models.User{}
	result := initializers.DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to get user info from database.",
		})
	}

	var responseDTO = struct {
		ID    uint   `json:"id"`
		Fname string `json:"fname"`
		Lname string `json:"lname"`
		Phone string `json:"phone"`
		Email string `json:"email" gorm:"unique"`
	}{
		ID:    user.ID,
		Fname: user.Fname,
		Lname: user.Lname,
		Phone: user.Phone,
		Email: user.Email,
	}

	// Send it back
	c.JSON(http.StatusOK, gin.H{
		"message": responseDTO,
	})
}

func UpdateUser(c *gin.Context) {
	// alter user profile

	var requestDTO = struct {
		ID       uint   `json:"id"`
		Fname    string `json:"fname"`
		Lname    string `json:"lname"`
		Phone    string `json:"phone"`
		Email    string `json:"email" gorm:"unique"`
		Password string `json:"password"`
	}{}

	if c.Bind(&requestDTO) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	//Map to store fields to be changed
	var mp map[string]any = map[string]any{}
	if len(requestDTO.Fname) > 0 {
		mp["fname"] = requestDTO.Fname
	}
	if len(requestDTO.Lname) > 0 {
		mp["lname"] = requestDTO.Lname
	}
	if len(requestDTO.Phone) > 0 {
		mp["phone"] = requestDTO.Phone
	}
	if len(requestDTO.Email) > 0 {
		mp["email"] = requestDTO.Email
	}
	if len(requestDTO.Password) > 0 {
		mp["password"] = requestDTO.Password
	}

	// Update changed fields
	var conn = initializers.DB.Table("users").Where("id = ?", requestDTO.ID).Updates(mp)
	if conn.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to update user fields",
			"error":   conn.Error,
		})
	}

	// Get users from DB to check its status after update
	var usr models.User = models.User{}
	conn = initializers.DB.Where("id = ?", requestDTO.ID).First(&usr)
	if conn.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to get user info from database.",
			"error":   conn.Error,
		})
	}

	// Send it back
	c.JSON(http.StatusOK, gin.H{
		"message": usr,
		"map":     mp,
	})
}

func DeleteUser(c *gin.Context) {
	// delete user profile

	// declare body to send post request
	var requestDTO = struct {
		ID uint `json:"id"`
	}{}

	if c.Bind(&requestDTO) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	// check if the user exists
	var user models.User
	initializers.DB.First(&user, "id = ?", requestDTO.ID)
	if user.ID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"errpr": "Invalid user ID",
		})
		return
	}

	// delete user from DB
	var conn = initializers.DB.Where("id = ?", requestDTO.ID).Delete(&models.User{})
	if conn.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to delete user",
			"error":   conn.Error,
		})
		return
	}

	// Send it back
	c.JSON(http.StatusOK, gin.H{
		"message": "User deleted successfully.",
	})

}

func CreateProductCategory(c *gin.Context) {
	// declare structure
	var body struct {
		Name        string
		Description string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
	}

	// create the new category
	category := models.Category{
		Name:        body.Name,
		Description: body.Description,
	}
	result := initializers.DB.Debug().Create(&category)

	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{})
	}
}

func RegisterProduct(c *gin.Context) {
	// add new product to DB
	var body struct {
		Name        string
		Description string
		Quantity    int
		Price       float64
		CategoryID  uint
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	// Create the new product
	product := models.Product{
		Name:        body.Name,
		Description: body.Description,
		Quantity:    body.Quantity,
		Price:       body.Price,
		CategoryID:  body.CategoryID,
	}
	result := initializers.DB.Debug().Create(&product)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to register new product",
		})
		return
	}
	// Responde request
	c.JSON(http.StatusOK, gin.H{})
}

func ListProducts(c *gin.Context) {
	// query for product list
	var products []models.Product = []models.Product{}

	conn := initializers.DB.Debug().Raw("SELECT * FROM products").Find(&products)
	if conn.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to get stuff from db",
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": products,
	})
}

func GetProduct(c *gin.Context) {
	// get request input

	id := c.Param("id")

	// get product profile
	product := models.Product{}
	conn := initializers.DB.Raw("SELECT * FROM products WHERE id=?", id).First(&product)

	// respond
	if conn != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": product,
		})
	}
}

func UpdateProduct(c *gin.Context) {
	// get request input

	var requestDTO = struct {
		ID          uint    `json:"id"`
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Quantity    int     `json:"quantity"`
		Price       float64 `json:"price"`
		CategoryID  uint    `json:"categoryid"`
	}{}

	if c.Bind(&requestDTO) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	fmt.Println("HERE IS THE ID =", requestDTO.ID)

	// query for product update
	//var product []models.Product = []models.Product{}
	conn := initializers.DB.Raw("UPDATE TABLE products SET name=? SET description=? SET quantity=? SET price=? SET categoryId=? WHERE id=?").Updates(requestDTO)

	if conn != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": requestDTO,
		})
	}
}
