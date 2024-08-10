package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/HarshitNagpal29/go-ecommerce-cart/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func AddAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id := c.Query("id")

		if user_id == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "User ID not provided",
			})
			return
		}
		address, err := primitive.ObjectIDFromHex(user_id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid User ID",
			})
			return
		}

		var addresses models.Address

		addresses.Address_ID = primitive.NewObjectID()

		if err := c.BindJSON(&addresses); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		match_filter := bson.D{{Key: "$match", Value: address}}
		unwind := bson.D{{Key: "$unwind", Value: bson.D{primitive.E{Key: "path", Value: "$address"}}}}
		group := bson.D{{Key: "$group", Value: bson.D{primitive.E{Key: "_id", Value: "$address_id"}, {Key: "count", Value: bson.D{primitive.E{Key: "$sum", Value: 1}}}}}}
		pointcursor, err := userCollection.Aggregate(ctx, mongo.Pipeline{match_filter, unwind, group})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		var addressinfo []bson.M
		if err = pointcursor.All(ctx, &addressinfo); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		var size int32
		for _, address_no := range addressinfo {
			count := address_no["count"]
			size = count.(int32)
		}
		if size < 2 {
			filter := bson.D{primitive.E{Key: "_id", Value: address}}
			update := bson.D{primitive.E{Key: "$push", Value: bson.D{primitive.E{Key: "address", Value: addresses}}}}
			_, err = userCollection.UpdateOne(ctx, filter, update)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Address limit exceeded",
			})
			return
		}
		defer cancel()
		ctx.Done()
	}
}
func RemoveAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id := c.Query("id")

		if user_id == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "User ID not provided",
			})
			return
		}
		address_id := c.Query("addressid")

		if address_id == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Address ID not provided",
			})
			return
		}
		usert_id, err := primitive.ObjectIDFromHex(user_id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid User ID",
			})
			return
		}
		addressid, err := primitive.ObjectIDFromHex(address_id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid Address ID",
			})
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		filter := bson.D{primitive.E{Key: "_id", Value: usert_id}}
		update := bson.D{{Key: "$pull", Value: bson.D{primitive.E{Key: "address", Value: bson.D{primitive.E{Key: "address_id", Value: addressid}}}}}}
		_, err = userCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "Address deleted successfully",
		})
	}
}

func EditHomeAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id := c.Query("id")

		if user_id == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "User ID not provided",
			})
			return
		}
		usert_id, err := primitive.ObjectIDFromHex(user_id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid User ID",
			})
			return
		}
		var editaddress models.Address
		if err := c.BindJSON(&editaddress); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		filter := bson.D{primitive.E{Key: "_id", Value: usert_id}}
		update := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "address.0.house_name", Value: editaddress.House}, {Key: "address.0.street_name", Value: editaddress.Street}, {Key: "address.0.city_name", Value: editaddress.City}, {Key: "address.0.pin_code", Value: editaddress.Pincode}}}}

		_, err = userCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "Home address updated successfully",
		})
		defer cancel()
		ctx.Done()
	}
}

func EditWorkAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id := c.Query("id")

		if user_id == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "User ID not provided",
			})
			return
		}
		usert_id, err := primitive.ObjectIDFromHex(user_id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid User ID",
			})
			return
		}
		var editaddress models.Address
		if err := c.BindJSON(&editaddress); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		filter := bson.D{primitive.E{Key: "_id", Value: usert_id}}
		update := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "address.1.house_name", Value: editaddress.House}, {Key: "address.1.street_name", Value: editaddress.Street}, {Key: "address.1.city_name", Value: editaddress.City}, {Key: "address.1.pin_code", Value: editaddress.Pincode}}}}

		_, err = userCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "Work address updated successfully",
		})
		defer cancel()
		ctx.Done()
	}
}

func DeleteAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id := c.Query("id")

		if user_id == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "User ID not provided",
			})
			return
		}
		addresses := make([]models.Address, 0)
		usert_id, err := primitive.ObjectIDFromHex(user_id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid User ID",
			})
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		filter := bson.D{{Key: "_id", Value: usert_id}}
		update := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "address_details", Value: addresses}}}}
		_, err = userCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "Address deleted successfully",
		})

	}
}
