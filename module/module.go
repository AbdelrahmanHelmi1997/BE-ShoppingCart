package module

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID        primitive.ObjectID `bson:"_id"`
	FirstName string             `bson:"firstName"`
	LastName  string             `bson:"lastName"`
	Email     string             `bson:"email"`
	Password  string             `bson:"password"`
	RoleType  string             `bson:"roleType"`
	Token     string             `bson:"token"`
}
type LoginData struct {
	ID        primitive.ObjectID `bson:"_id"`
	FirstName string             `bson:"firstName"`
	LastName  string             `bson:"lastName"`
	RoleType  string             `bson: "roleType"`
	Token     string             `bson:"token"`
}

type ShopItems struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	ItemName string             `json:"itemName,omitempty"`
}

type Cart struct {
	ID         primitive.ObjectID `json:"_id,omitempty" bson:"_id"`
	CustomerId string             `json:"customerId,omitempty" bson:"customerId"`
	CartItems  []ShopItems        `json:"cartItems,omitempty"bson:"cartItems"`
}

type CartItems struct {
	ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
}

type RemoveItems struct {
	CartID string `json:"_id" bson:"_id"`
	ItemID string `json:"_id" bson:"_id"`
}

type AddItems struct {
	CustomerId string             `json:"CustomerId,omitempty" bson:"_id,omitempty"`
	ID         primitive.ObjectID `json:"_id,omitempty"bson:"_id,omitempty"`
	ItemName   string             `json:"itemName,omitempty"`
}
