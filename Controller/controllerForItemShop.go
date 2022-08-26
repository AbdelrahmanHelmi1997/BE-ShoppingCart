package Controller

import (
	"AmzonElGalaba/dataBase"
	"AmzonElGalaba/module"
	"AmzonElGalaba/responses"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var validate = validator.New()

type Shop struct {
}

func (s *Shop) AddToShop(c *gin.Context) {

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	var item module.ShopItems

	if err := c.BindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
		return
	}

	if validationErr := validate.Struct(&item); validationErr != nil {
		c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
		return
	}

	newItem := module.ShopItems{
		ID:       primitive.NewObjectID(),
		ItemName: item.ItemName,
	}

	result, err := dataBase.Collection.InsertOne(ctx, newItem)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
		return
	}

	c.JSON(http.StatusCreated, responses.UserResponse{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"data": result}})
}

func (s *Shop) GetAllShopItems(c *gin.Context) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var items []module.ShopItems
	defer cancel()

	results, err := dataBase.Collection.Find(ctx, bson.M{})

	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
		return
	}

	//reading from the db in an optimal way
	defer results.Close(ctx)
	for results.Next(ctx) {
		var singleItem module.ShopItems
		if err = results.Decode(&singleItem); err != nil {
			c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
		}

		items = append(items, singleItem)
	}

	c.JSON(http.StatusOK,
		responses.UserResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": items}},
	)
}

func (s *Shop) GetItemById(c *gin.Context) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	userId := c.Param("itemId")
	var item module.ShopItems
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(userId)
	fmt.Println(userId)

	err := dataBase.Collection.FindOne(ctx, bson.M{"_id": objId}).Decode(&item)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
		return
	}

	c.JSON(http.StatusOK, responses.UserResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": item}})
}

func (s *Shop) DeleteItemFromShop(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	userId := c.Param("itemId")
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(userId)

	result, err := dataBase.Collection.DeleteOne(ctx, bson.M{"_id": objId})
	fmt.Println(objId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
		return
	}

	if result.DeletedCount < 1 {
		c.JSON(http.StatusNotFound,
			responses.UserResponse{Status: http.StatusNotFound, Message: "error", Data: map[string]interface{}{"data": "item with specified ID not found!"}},
		)
		return
	}

	c.JSON(http.StatusOK,
		responses.UserResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": "item successfully deleted!"}},
	)
}
