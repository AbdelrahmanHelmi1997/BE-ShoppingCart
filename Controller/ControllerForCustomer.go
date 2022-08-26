package Controller

import (
	"AmzonElGalaba/Helper"
	"AmzonElGalaba/dataBase"
	"AmzonElGalaba/module"
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct{}

func (u *User) SignUp(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	var user module.User

	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Println()

	count, err := dataBase.UsersCollection.CountDocuments(ctx, bson.M{"email": user.Email})
	defer cancel()
	if err != nil {
		log.Panic(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while checking for the email"})
		return
	}

	password := Helper.HashPassword(user.Password)
	user.Password = password

	if count > 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "this email  already exists"})
		return
	}


	user.ID = primitive.NewObjectID()
	token, _ := Helper.GenerateAllTokens(user.Email, user.FirstName, user.LastName, user.RoleType)
	user.Token = token

	resultInsertionNumber, insertErr := dataBase.UsersCollection.InsertOne(ctx, user)
	if insertErr != nil {
		msg := fmt.Sprintf("User item was not created")
		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		return
	}
	defer cancel()


	c.JSON(http.StatusOK, resultInsertionNumber)

}

func (u *User) Login(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	var user module.User
	var foundUser module.User

	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := dataBase.UsersCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
	defer cancel()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": " Email is incorrect"})
		return
	}

	passwordIsValid, _ := Helper.VerifyPassword(user.Password, foundUser.Password)
	defer cancel()
	if passwordIsValid != true {
		c.JSON(http.StatusInternalServerError, gin.H{"error": " Password is incorrect"})
		return
	}

	if foundUser.Email == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found"})
		return
	}
	token, _ := Helper.GenerateAllTokens(foundUser.Email, foundUser.FirstName, foundUser.LastName, foundUser.RoleType)

	Helper.UpdateAllTokens(token, foundUser.ID)
	err = dataBase.UsersCollection.FindOne(ctx, bson.M{"_id": foundUser.ID}).Decode(&foundUser)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	userFound := module.LoginData{
		ID:        foundUser.ID,
		FirstName: foundUser.FirstName,
		LastName:  foundUser.LastName,
		RoleType:  foundUser.RoleType,
		Token:     foundUser.Token,
	}

	c.JSON(http.StatusOK, userFound)

}
