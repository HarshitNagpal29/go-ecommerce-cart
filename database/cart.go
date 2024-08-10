package database

import (
	"context"
	"errors"
	"time"

	"github.com/HarshitNagpal29/go-ecommerce-cart/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrCantFindProduct    = errors.New("product not found")
	ErrCantDecodeProduct  = errors.New("can't decode product")
	ErrUserIdIsNotValid   = errors.New("user id is not valid")
	ErrCantUpdateUser     = errors.New("can't update user")
	ErrCantRemoveItemCart = errors.New("can't remove item from cart")
	ErrCantGetItem        = errors.New("can't get item")
	ErrCantBuyItem        = errors.New("can't buy item")
)

func AddProductToCart(ctx context.Context, prodCollection, userCollection *mongo.Collection, productID primitive.ObjectID, userID string) error {
	searchfromdb, err := prodCollection.Find(ctx, bson.M{"_id": productID})
	if err != nil {
		return ErrCantFindProduct
	}
	var productCart []models.ProductUser
	if err = searchfromdb.All(ctx, &productCart); err != nil {
		return ErrCantDecodeProduct
	}
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return ErrUserIdIsNotValid
	}
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "usercart", Value: bson.D{{Key: "$each", Value: productCart}}}}}}

	_, err = userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return ErrCantUpdateUser
	}
	return nil
}

func RemoveProductFromCart(ctx context.Context, prodCollection, userCollection *mongo.Collection, productID primitive.ObjectID, userID string) error {
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return ErrUserIdIsNotValid
	}
	filter := bson.M{"_id": id}
	update := bson.M{"$pull": bson.M{"usercart": bson.M{"_id": productID}}}
	_, err = userCollection.UpdateMany(ctx, filter, update)
	if err != nil {
		return ErrCantRemoveItemCart
	}
	return nil

}

func BuyItemFromCart(ctx context.Context, userCollection *mongo.Collection, userID string) error {
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return ErrUserIdIsNotValid
	}
	var getcartitems models.User
	var ordercart models.Order

	ordercart.Order_ID = primitive.NewObjectID()
	ordercart.Order_at = time.Now()
	ordercart.Order_Cart = make([]models.ProductUser, 0)
	ordercart.Payment_Method.COD = true

	unwind := bson.D{{Key: "$unwind", Value: bson.D{primitive.E{Key: "path", Value: "$usercart"}}}}
	grouping := bson.D{{Key: "$group", Value: bson.D{primitive.E{Key: "_id", Value: "$_id"}, {Key: "total", Value: bson.D{primitive.E{Key: "$sum", Value: "$usercart.price"}}}}}}
	userCollection.Aggregate(ctx, mongo.Pipeline{unwind, grouping})
	currentresults, err := userCollection.Aggregate(ctx, mongo.Pipeline{unwind, grouping})
	if err != nil {
		return ErrCantGetItem
	}
	ctx.Done()
	var getusercart []bson.M
	if err = currentresults.All(ctx, &getusercart); err != nil {
		return ErrCantGetItem
	}
	var total_price int32
	for _, user_item := range getusercart {
		price := user_item["total"]
		total_price = price.(int32)
	}
	ordercart.Price = uint64(total_price)

	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "orders", Value: ordercart}}}}
	_, err = userCollection.UpdateMany(ctx, filter, update)
	if err != nil {
		return ErrCantBuyItem
	}
	err = userCollection.FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: id}}).Decode(&getcartitems)
	if err != nil {
		return ErrCantBuyItem
	}
	filter2 := bson.D{primitive.E{Key: "_id", Value: id}}
	update2 := bson.M{"$push": bson.M{"order.$[].order_list": bson.M{"$each": getcartitems.UserCart}}}
	_, err = userCollection.UpdateOne(ctx, filter2, update2)
	if err != nil {
		return ErrCantBuyItem
	}
	usercart_empty := make([]models.ProductUser, 0)
	filter3 := bson.D{primitive.E{Key: "_id", Value: id}}
	update3 := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "usercart", Value: usercart_empty}}}}
	_, err = userCollection.UpdateOne(ctx, filter3, update3)
	if err != nil {
		return ErrCantBuyItem
	}
	return nil
}

func InstantBuyer(ctx context.Context, prodCollection, userCollection *mongo.Collection, productID primitive.ObjectID, UserID string) error {
	id, err := primitive.ObjectIDFromHex(UserID)
	if err != nil {
		return ErrUserIdIsNotValid
	}

	var product_details models.ProductUser
	var orders_details models.Order

	orders_details.Order_ID = primitive.NewObjectID()
	orders_details.Order_at = time.Now()
	orders_details.Order_Cart = make([]models.ProductUser, 0)
	orders_details.Payment_Method.COD = true

	err = prodCollection.FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: productID}}).Decode(&product_details)
	if err != nil {
		return ErrCantGetItem
	}
	orders_details.Price = product_details.Price
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "orders", Value: orders_details}}}}
	_, err = userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return ErrCantBuyItem
	}
	filter2 := bson.D{primitive.E{Key: "_id", Value: id}}
	update2 := bson.M{"$push": bson.M{"order.$[].order_list": product_details}}
	_, err = userCollection.UpdateOne(ctx, filter2, update2)
	if err != nil {
		return ErrCantBuyItem
	}
	return nil
}
