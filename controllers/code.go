package controllers

import (
	"context"
	"fmt"
	"net/http"
	"test/dilaf/structs"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

func DBConnection() (*mongo.Client, context.Context) {
	url := options.Client().ApplyURI("mongodb://localhost:27017/")
	NewCtx, _ := context.WithTimeout(context.Background(), 2*time.Second)

	Client, err := mongo.Connect(NewCtx, url)
	if err != nil {
		fmt.Printf("errors: %v\n", err)
	}
	return Client, NewCtx
}

// ?HashPassword
func HashPassword(password string) (string, error) {
	var passwordBytes = []byte(password)
	hashedPasswordBytes, err := bcrypt.GenerateFromPassword(passwordBytes, bcrypt.MinCost)

	return string(hashedPasswordBytes), err
}

// ?Compare The Password
func CompareHashPasswords(HashedPasswordFromDB, PasswordToCampare string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(HashedPasswordFromDB), []byte(PasswordToCampare))
	return err == nil
}

//!Register
func Register(c *gin.Context) {
	var registerTemp structs.RegStruct
	c.ShouldBindJSON(&registerTemp)

	if registerTemp.Name == "" || registerTemp.Login == "" || registerTemp.Balance == 0 || registerTemp.Surname == "" || registerTemp.Password == "" {

		c.JSON(404, "Empty!")
	} else {
		client, ctx := DBConnection()

		DBConnect := client.Database("E-BookStore").Collection("Users")

		id := primitive.NewObjectID().Hex()
		Hashed, _ := HashPassword(registerTemp.Password)

		DBConnect.InsertOne(ctx, bson.M{
			"_id":      id,
			"name":     registerTemp.Name,
			"surname":  registerTemp.Surname,
			"balance":  registerTemp.Balance,
			"login":    registerTemp.Login,
			"password": Hashed,
		})
		c.JSON(200, "registered!")
	}
}

func Login(c *gin.Context) {
	var loginTemp structs.LoginStruct
	c.ShouldBindJSON(&loginTemp)
	if loginTemp.Login == "" || loginTemp.Password == "" {
		c.JSON(404, "empty!")
	} else {
		client, ctx := DBConnection()

		DBConnect := client.Database("E-BookStore").Collection("Users")

		result := DBConnect.FindOne(ctx, bson.M{
			"login": loginTemp.Login,
		})

		var userdata structs.RegStruct
		result.Decode(&userdata)
		isValidPass := CompareHashPasswords(userdata.Password, loginTemp.Password)
		fmt.Println(isValidPass)

		if isValidPass {
			http.SetCookie(c.Writer, &http.Cookie{
				Name:     "cookie",
				Value:    userdata.Id,
				Expires:  time.Now().Add(60 * time.Hour),
				Secure:   false,
				SameSite: http.SameSiteLaxMode,
			})
			c.JSON(200, "success")
		} else {
			c.JSON(404, "Wrong login or password")
		}
	}
}

func Allbooks(c *gin.Context) {
	_, err := c.Request.Cookie("cookie")
	var booksSlice = []structs.Books{
		{
			Name:   "Wonderful",
			Author: "John",
			Year:   "",
		},
	}
	if err == nil {
		c.JSON(401, "No cookie")
	} else {
		client, ctx := DBConnection()

		DBConnect := client.Database("E-BookStore").Collection("Books")
		result, _ := DBConnect.Find(ctx, bson.M{})

		for result.Next(ctx) {
			var dbDataTemp structs.Books
			result.Decode(&dbDataTemp)
			booksSlice = append(booksSlice, dbDataTemp)

		}
		c.JSON(200, booksSlice)
	}
}
