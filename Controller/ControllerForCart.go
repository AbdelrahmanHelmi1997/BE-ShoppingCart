package Controller

import (
	"AmzonElGalaba/dataBase"
	"AmzonElGalaba/module"
	"AmzonElGalaba/responses"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
)

type ShopCart struct {
}

func (s *ShopCart) CreateCart(c *gin.Context) {

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	var Cart module.Cart

	if err := c.BindJSON(&Cart); err != nil {
		c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
		return
	}

	if validationErr := validate.Struct(&Cart); validationErr != nil {
		c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
		return
	}

	newItem := module.Cart{
		ID:         primitive.NewObjectID(),
		CustomerId: Cart.CustomerId,
		CartItems:  Cart.CartItems,
	}

	result, err := dataBase.ShoppingCartColliction.InsertOne(ctx, newItem)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
		return
	}

	c.JSON(http.StatusCreated, responses.UserResponse{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"data": result}})
}

func (s *ShopCart) GetCart(c *gin.Context) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	userId := c.Param("cartId")
	var item module.Cart
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(userId)

	err := dataBase.ShoppingCartColliction.FindOne(ctx, bson.M{"_id": objId}).Decode(&item)
	fmt.Println(item)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
		return
	}

	c.JSON(http.StatusOK, responses.UserResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": item}})
}

func (s *ShopCart) DeleteFromCart(c *gin.Context) {

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	var req module.AddItems

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
		return
	}

	if validationErr := validate.Struct(&req); validationErr != nil {
		c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
		return
	}

	filter := bson.M{"customerId": req.CustomerId}
	update := bson.M{
		"$pull": bson.M{
			"cartItems": bson.M{"_id": req.ID, "itemname": req.ItemName},
		},
	}

	result := dataBase.ShoppingCartColliction.FindOneAndUpdate(ctx, filter, update)
	if result.Err() != nil {
		fmt.Println(result.Err())
		c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": "Item Not Found"}})
		return
	}

	c.JSON(http.StatusCreated, responses.UserResponse{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"data": "Item Removed "}})
}

func (s *ShopCart) AddItemToCart(c *gin.Context) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	var req module.AddItems
	var cart module.Cart

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
		return
	}

	filter := bson.M{"customerId": req.CustomerId}
	update := bson.M{
		"$push": bson.M{
			"cartItems": bson.M{"_id": req.ID, "itemname": req.ItemName},
		},
	}

	result := dataBase.ShoppingCartColliction.FindOneAndUpdate(ctx, filter, update)

	if result.Err() != nil {
		fmt.Println(result.Err())
		c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": "Item Not Added"}})
		return
	}

	c.JSON(http.StatusCreated, responses.UserResponse{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"data": cart}})
}
