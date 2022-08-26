package routes

import (
	"AmzonElGalaba/Controller"
	"AmzonElGalaba/middleware"

	"github.com/gin-gonic/gin"
)

var Shop Controller.Shop
var Cart Controller.ShopCart
var auth middleware.Auth

func UserRoute(router *gin.Engine) {
	router.POST("/addItemToShop", auth.Authentication, Shop.AddToShop)
	router.POST("/createCart", Cart.CreateCart)
	router.GET("/getAllShopItems", Shop.GetAllShopItems)
	router.GET("getItem/:itemId", Shop.GetItemById)
	router.GET("getCart/:cartId", Cart.GetCart)
	router.POST("deleteItemFromCart", Cart.DeleteFromCart)
	router.DELETE("deleteFromShop/:itemId", Shop.DeleteItemFromShop)
	router.POST("addItemToCart", Cart.AddItemToCart)

}
