package pasabarbackend

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Admin struct {
	Email    string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
	Role     string `json:"role,omitempty" bson:"role,omitempty"`
	Token    string `json:"token,omitempty" bson:"token,omitempty"`
	Private  string `json:"private,omitempty" bson:"private,omitempty"`
	Public   string `json:"public,omitempty" bson:"public,omitempty"`
}

type Credential struct {
	Status  bool   `json:"status" bson:"status"`
	Token   string `json:"token,omitempty" bson:"token,omitempty"`
	Message string `json:"message,omitempty" bson:"message,omitempty"`
}

type Response struct {
	Status  bool        `json:"status" bson:"status"`
	Message string      `json:"message" bson:"message"`
	Data    interface{} `json:"data" bson:"data"`
}

type Payload struct {
	Id   primitive.ObjectID `json:"id"`
	Role string             `json:"role"`
	Exp  time.Time          `json:"exp"`
	Iat  time.Time          `json:"iat"`
	Nbf  time.Time          `json:"nbf"`
}

type Catalog struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" `
	Nomorid     int                `json:"nomorid" bson:"nomorid"`
	Title       string             `json:"title" bson:"title"`
	Description string             `json:"description" bson:"description"`
	Image       string             `json:"image" bson:"image"`
	Status      bool               `json:"status" bson:"status"`
}

type About struct {
	ID          int    `json:"id" bson:"id"`
	Title       string `json:"title" bson:"title"`
	Description string `json:"description" bson:"description"`
	Image       string `json:"image" bson:"image"`
	Status      bool   `json:"status" bson:"status"`
}

type Tour struct {
	ID          int       `json:"id" bson:"id"`
	Title       string    `json:"title" bson:"title"`
	Description string    `json:"description" bson:"description"`
	Cari        string    `json:"cari" bson:"cari"`
	Tanggal     string    `json:"tanggal" bson:"tanggal"`
	Image       string    `json:"image" bson:"image"`
	Harga       int       `json:"harga" bson:"harga"`
	Catalog     []Catalog `json:"catalog" bson:"catalog"`
	Rating      string    `json:"rating" bson:"rating"`
	Status      bool      `json:"status" bson:"status"`
}

type Hotel struct {
	ID          int       `json:"id" bson:"id"`
	Title       string    `json:"title" bson:"title"`
	Description string    `json:"description" bson:"description"`
	Cari        string    `json:"cari" bson:"cari"`
	Tanggal     string    `json:"tanggal" bson:"tanggal"`
	Image       string    `json:"image" bson:"image"`
	Harga       int       `json:"harga" bson:"harga"`
	Catalog     []Catalog `json:"catalog" bson:"catalog"`
	Rating      string    `json:"rating" bson:"rating"`
	Status      bool      `json:"status" bson:"status"`
}

type Restoran struct {
	ID          int       `json:"id" bson:"id"`
	Title       string    `json:"title" bson:"title"`
	Description string    `json:"description" bson:"description"`
	Cari        string    `json:"cari" bson:"cari"`
	Tanggal     string    `json:"tanggal" bson:"tanggal"`
	Image       string    `json:"image" bson:"image"`
	Harga       int       `json:"harga" bson:"harga"`
	Catalog     []Catalog `json:"catalog" bson:"catalog"`
	Rating      string    `json:"rating" bson:"rating"`
	Status      bool      `json:"status" bson:"status"`
}

type Contact struct {
	ID      int    `json:"id" bson:"id"`
	Name    string `json:"title" bson:"title"`
	Subject string `json:"description" bson:"description"`
	Alamat  string `json:"alamat" bson:"alamat"`
	Website string `json:"website" bson:"website"`
	Message string `json:"image" bson:"image"`
	Email   string `json:"email" bson:"email"`
	Phone   string `json:"phone" bson:"phone"`
	Status  bool   `json:"status" bson:"status"`
}
