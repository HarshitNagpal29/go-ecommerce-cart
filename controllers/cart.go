package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/HarshitNagpal29/go-ecommerce-cart/database"
	"github.com/HarshitNagpal29/go-ecommerce-cart/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Application struct {
	prodCollection *mongo.Collection
	userCollection *mongo.Collection
}

func NewApplication(prodCollection *mongo.Collection, userCollection *mongo.Collection) *Application {
	return &Application{
		prodCollection: prodCollection,
		userCollection: userCollection,
	}
}

func (app *Application) AddToCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		productQueryID := c.Query("id")
		if productQueryID == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Product ID not provided",
			})
			return
		}
		userQueryID := c.Query("userID")
		if userQueryID == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "User ID not provided",
			})
			return
		}
		productId, err := primitive.ObjectIDFromHex(productQueryID)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid Product ID",
			})
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err = database.AddProductToCart(ctx, app.prodCollection, app.userCollection, productId, userQueryID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "Product added to cart",
		})

	}
}

func (app *Application) RemoveItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		productQueryID := c.Query("id")
		if productQueryID == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Product ID not provided",
			})
			return
		}
		userQueryID := c.Query("userID")
		if userQueryID == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "User ID not provided",
			})
			return
		}
		productId, err := primitive.ObjectIDFromHex(productQueryID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid Product ID",
			})
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err = database.RemoveProductFromCart(ctx, app.prodCollection, app.userCollection, productId, userQueryID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "Product removed from cart",
		})
	}
}

func GetItemsFromCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id := c.Query("id")

		if user_id == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "User ID not provided",
			})
			return
		}
		usert_id, _ := primitive.ObjectIDFromHex(user_id)
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var filledcart models.User
		err := userCollection.FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: usert_id}}).Decode(&filledcart)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		filter_match := bson.D{{Key: "$match", Value: bson.D{primitive.E{Key: "_id", Value: usert_id}}}}
		unwind := bson.D{{Key: "$unwind", Value: bson.D{primitive.E{Key: "path", Value: "$usercart"}}}}
		grouping := bson.D{{Key: "$group", Value: bson.D{primitive.E{Key: "_id", Value: "$_id"}, {Key: "total", Value: bson.D{primitive.E{Key: "$sum", Value: "$usercart.price"}}}}}}

		pointcursor, err := userCollection.Aggregate(ctx, mongo.Pipeline{filter_match, unwind, grouping})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		var listing []bson.M
		if err = pointcursor.All(ctx, &listing); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		for _, json := range listing {
			c.JSON(http.StatusOK, json["total"])
		}
		ctx.Done()
	}
}

func (app *Application) BuyFromCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		userQueryID := c.Query("userID")
		if userQueryID == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "User ID not provided",
			})
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err := database.BuyItemFromCart(ctx, app.userCollection, userQueryID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "Items bought successfully",
		})
	}
}

func (app *Application) InstantBuy() gin.HandlerFunc {
	return func(c *gin.Context) {
		productQueryID := c.Query("id")
		if productQueryID == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Product ID not provided",
			})
			return
		}
		userQueryID := c.Query("userID")
		if userQueryID == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "User ID not provided",
			})
			return
		}
		productId, err := primitive.ObjectIDFromHex(productQueryID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid Product ID",
			})
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err = database.InstantBuyer(ctx, app.prodCollection, app.userCollection, productId, userQueryID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "Product bought successfully",
		})

	}
}
