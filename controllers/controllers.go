package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/HarshitNagpal29/go-ecommerce-cart/database"
	"github.com/HarshitNagpal29/go-ecommerce-cart/models"
	"github.com/HarshitNagpal29/go-ecommerce-cart/tokens"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection = database.UserData(database.Client, "Users")
var productCollection *mongo.Collection = database.ProductData(database.Client, "Products")
var Validate = validator.New()

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Fatal(err)
	}
	return string(bytes)
}

func VerifyPassword(userPassword string, givenPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(givenPassword), []byte(userPassword))
	valid := true
	msg := "Password is correct"
	if err != nil {
		valid = false
		msg = "Password is incorrect"
	}
	return valid, msg

}

func Signup() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		validateErr := Validate.Struct(user)
		if validateErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": validateErr.Error(),
			})
			return
		}
		count, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error while checking for existing user",
			})
			return
		}
		if count > 0 {
			c.JSON(http.StatusConflict, gin.H{
				"error": "User already exists",
			})
			return
		}
		count, err = userCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error while checking for existing user",
			})
			return
		}
		if count > 0 {
			c.JSON(http.StatusConflict, gin.H{
				"error": "User already exists",
			})
			return
		}
		password := HashPassword(user.Password)
		user.Password = password
		user.Created_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.User_ID = user.ID.Hex()
		token, refreshtoken, _ := tokens.TokenGenerator(user.Email, user.First_Name, user.Last_Name, user.User_ID)
		user.Token = token
		user.Refresh_Token = refreshtoken
		user.UserCart = []models.ProductUser{}
		user.Address_Details = []models.Address{}
		user.Order_Status = []models.Order{}
		_, err = userCollection.InsertOne(ctx, user)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error while inserting the user",
			})
			return
		}
		defer cancel()
		c.JSON(http.StatusCreated, gin.H{
			"message": "User created successfully",
		})
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var founduser models.User
		var user models.User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		cur, err := userCollection.Find(ctx, bson.M{"email": user.Email})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		err = cur.Decode(&founduser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		match, msg := VerifyPassword(user.Password, founduser.Password)
		defer cancel()

		if !match {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": msg,
			})
			return
		}
		token, refreshtoken, _ := tokens.TokenGenerator(founduser.Email, founduser.First_Name, founduser.Last_Name, founduser.User_ID)
		defer cancel()

		tokens.UpdateAllTokens(token, refreshtoken, founduser.User_ID)
		c.JSON(http.StatusOK, founduser)
	}
}

func ProductViewerAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var product models.Product
		if err := c.BindJSON(&product); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		validateErr := Validate.Struct(product)
		if validateErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": validateErr.Error(),
			})
			return
		}
		product.Product_ID = primitive.NewObjectID()
		_, err := productCollection.InsertOne(ctx, product)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		defer cancel()
		c.JSON(http.StatusCreated, gin.H{
			"message": "Product added successfully",
		})
	}
}

func SearchProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		var productList []models.Product
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		cursor, err := productCollection.Find(ctx, bson.D{{}})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		err = cursor.All(ctx, &productList)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		defer cursor.Close(ctx)

		if err := cursor.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, productList)
	}
}

func SearchProductByQuery() gin.HandlerFunc {
	return func(c *gin.Context) {
		var productList []models.Product
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		queryParam := c.Query("name")
		if queryParam == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "No query provided",
			})
			return
		}
		cursor, err := productCollection.Find(ctx, bson.M{"product_name": primitive.Regex{Pattern: queryParam, Options: "i"}})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		err = cursor.All(ctx, &productList)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		defer cursor.Close(ctx)
		if err := cursor.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, productList)
	}
}
