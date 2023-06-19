package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"reflect"
	"runtime"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	IDNumber  int64  `json:"idNumber" gorm:"column:id_number"`
	Name      string `json:"name" gorm:"column:name"`
	Surname   string `json:"surname" gorm:"column:surname"`
	BirthYear int    `json:"birthyear" gorm:"column:birthyear"`
}

type UserMinimal struct {
	IDNumber int64  `json:"idNumber"`
	Name     string `json:"name"`
}

type AfetzedeLokasyon struct {
	gorm.Model
	PhoneNumber string  `json:"phoneNumber" gorm:"column:phone_number"`
	UserMedia   string  `json:"mediaLink" gorm:"column:user_media_link"`
	Latitude    float64 `json:"latitude" gorm:"column:latitude"`
	Longitude   float64 `json:"longitude" gorm:"column:longitude"`
	UserID      uint    `json:"userID" gorm:"column:user_id;index;ForeignKey:id"` // foreign key to User table
}

type JWTToken struct {
	Token string `json:"token"`
}

type CustomClaims struct {
	UserID uint `json:"UserID"`
	jwt.StandardClaims
}

var (
	db     *gorm.DB
	secret = []byte("my-secret-key") // TODO: replace with your own secret key
)

func connection() (*gorm.DB, error) {

	dsn := "" // TODO: Add your dsn connection string here
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	db.AutoMigrate(User{}, &AfetzedeLokasyon{})
	return db, err
}

func main() {
	var err error
	db, err = connection()
	if err != nil {
		log.Fatal(err)
	}
	db.AutoMigrate(&User{})

	router := gin.Default()
	router.POST("/signup", handleSignup)
	router.POST("/login", handleLogin)

	authorized := router.Group("/")
	authorized.Use(ValidateJWT())
	{
		authorized.POST("/yardimaihtiyacimvar", handleYardimaihtiyacimvar)
		authorized.GET("/yardimedebilirim", handleYardimedebilirim)
	}

	router.Run(":8080")
}

func handleYardimaihtiyacimvar(c *gin.Context) {
	var afetzedelokasyon AfetzedeLokasyon
	if err := c.ShouldBindJSON(&afetzedelokasyon); err != nil {
		// Get the name of the current function
		funcName := runtime.FuncForPC(reflect.ValueOf(handleYardimaihtiyacimvar).Pointer()).Name()

		// Print the function name along with a message
		fmt.Printf("%s, afetzedelokasyon json to model: %s\n", funcName, err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// get the current user ID from the JWT token
	userID, ok := c.Get("UserID")
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to get userID from token"})
		return
	}

	// set the userID field in the AfetzedeLokasyon object
	afetzedelokasyon.UserID = userID.(uint)
	db.AutoMigrate(&AfetzedeLokasyon{})
	// save the AfetzedeLokasyon object to the database
	if err := db.Create(&afetzedelokasyon).Error; err != nil {
		// Get the name of the current function
		funcName := runtime.FuncForPC(reflect.ValueOf(handleYardimaihtiyacimvar).Pointer()).Name()

		// Print the function name along with a message
		fmt.Printf("%s, afetzedelokasyon create db element: %s\n", funcName, err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, afetzedelokasyon)

}

func handleYardimedebilirim(c *gin.Context) {
	afetzedeLokasyons := []AfetzedeLokasyon{}
	if err := db.Find(&afetzedeLokasyons).Error; err != nil {
		// Get the name of the current function
		funcName := runtime.FuncForPC(reflect.ValueOf(handleYardimedebilirim).Pointer()).Name()

		// Print the function name along with a message
		fmt.Printf("%s: %s\n", funcName, err.Error())

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, afetzedeLokasyons)
}

func handleSignup(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		// Get the name of the current function
		funcName := runtime.FuncForPC(reflect.ValueOf(handleSignup).Pointer()).Name()

		// Print the function name along with a message
		fmt.Printf("%s, user json to model: %s\n", funcName, err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if the input is valid or not
	if !checkCitizenship(user.IDNumber, user.Name, user.Surname, user.BirthYear) {
		// Get the name of the current function
		funcName := runtime.FuncForPC(reflect.ValueOf(handleSignup).Pointer()).Name()

		// Print the function name along with a message
		fmt.Printf("%s, user json to model: %s\n", funcName, "User identity information is not valid")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "User identity information is not valid"})
		return
	}

	// Check if the user already exists in the database
	var count int64
	if err := db.Model(&User{}).Where("id_number = ?", user.IDNumber).Count(&count).Error; err != nil {
		// Get the name of the current function
		funcName := runtime.FuncForPC(reflect.ValueOf(handleSignup).Pointer()).Name()

		// Print the function name along with a message
		fmt.Printf("%s, user count error: %s\n", funcName, err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if count > 0 {
		// Get the name of the current function
		funcName := runtime.FuncForPC(reflect.ValueOf(handleSignup).Pointer()).Name()

		// Print the function name along with a message
		fmt.Printf("%s: %s\n", funcName, "User already exists")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "User already exists"})
		return
	}

	// Insert the user into the database
	if err := db.Create(&user).Error; err != nil {
		// Get the name of the current function
		funcName := runtime.FuncForPC(reflect.ValueOf(handleSignup).Pointer()).Name()

		// Print the function name along with a message
		fmt.Printf("%s, user db element create: %s\n", funcName, err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// JWT token Creation Segment
	tokenString, err := createJWTToken(user.ID)
	if err != nil {
		// Get the name of the current function
		funcName := runtime.FuncForPC(reflect.ValueOf(handleSignup).Pointer()).Name()

		// Print the function name along with a message
		fmt.Printf("%s, token string create: %s\n", funcName, err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user":  user,
		"token": tokenString,
	})

}

func handleLogin(c *gin.Context) {
	var userMinimal UserMinimal
	if err := c.ShouldBindJSON(&userMinimal); err != nil {
		// Get the name of the current function
		funcName := runtime.FuncForPC(reflect.ValueOf(handleLogin).Pointer()).Name()

		// Print the function name along with a message
		fmt.Printf("%s, userMinimal json to model: %s\n", funcName, err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find the user in the database
	var user User
	if err := db.Where("id_number = ? AND name = ?", userMinimal.IDNumber, userMinimal.Name).First(&user).Error; err != nil {
		// Get the name of the current function
		funcName := runtime.FuncForPC(reflect.ValueOf(handleLogin).Pointer()).Name()

		// Print the function name along with a message
		fmt.Printf("%s, search user in database: %s\n", funcName, err.Error())
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// JWT token Creation Segment
	tokenString, err := createJWTToken(user.ID)
	if err != nil {
		// Get the name of the current function
		funcName := runtime.FuncForPC(reflect.ValueOf(handleLogin).Pointer()).Name()

		// Print the function name along with a message
		fmt.Printf("%s, create token string: %s\n", funcName, err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user":  user,
		"token": tokenString,
	})
}

func createJWTToken(id uint) (string, error) {
	// JWT token Creation Segment
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"UserID": id,
		"exp":    time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString(secret)
	return tokenString, err
}

func checkCitizenship(tckimlikNo int64, ad string, soyad string, dogumYili int) bool {
	// Create the SOAP request body
	body := fmt.Sprintf(`<?xml version="1.0" encoding="utf-8"?>
        <soap12:Envelope xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:soap12="http://www.w3.org/2003/05/soap-envelope">
          <soap12:Body>
            <TCKimlikNoDogrula xmlns="http://tckimlik.nvi.gov.tr/WS">
              <TCKimlikNo>%d</TCKimlikNo>
              <Ad>%s</Ad>
              <Soyad>%s</Soyad>
              <DogumYili>%d</DogumYili>
            </TCKimlikNoDogrula>
          </soap12:Body>
        </soap12:Envelope>`, tckimlikNo, ad, soyad, dogumYili)

	// Create a new HTTP request with the SOAP body
	req, err := http.NewRequest("POST", "https://tckimlik.nvi.gov.tr/Service/KPSPublic.asmx", bytes.NewBufferString(body))
	if err != nil {
		return false
	}

	// Set the SOAP action header
	req.Header.Set("SOAPAction", "http://tckimlik.nvi.gov.tr/WS/TCKimlikNoDogrula")

	// Set the content type header
	req.Header.Set("Content-Type", "application/soap+xml; charset=utf-8")

	// Send the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	// Read the SOAP response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return false
	}

	// Unmarshal the SOAP response into a struct
	type Envelope struct {
		Body struct {
			Response struct {
				Result bool `xml:"TCKimlikNoDogrulaResult"`
			} `xml:"TCKimlikNoDogrulaResponse"`
		} `xml:"Body"`
	}
	var env Envelope
	if err := xml.Unmarshal(bodyBytes, &env); err != nil {
		return false
	}

	// Return the boolean value of TCKimlikNoDogrulaResult
	return env.Body.Response.Result
}

func ValidateJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve JWT token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		// Remove the "Bearer " prefix from the token string
		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

		// Parse JWT token and validate signing method
		token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			}
			return secret, nil
		})

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized", "message": err.Error()})
			return
		}
		claims, ok := token.Claims.(*CustomClaims)
		if !ok || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Set("UserID", claims.UserID)
		c.Next()
	}
}
