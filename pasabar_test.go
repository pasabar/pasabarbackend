package pasabarbackend

import (
	"fmt"
	"testing"

	"github.com/aiteung/atdb"
	"github.com/whatsauth/watoken"
	"go.mongodb.org/mongo-driver/bson"
)

// PASETO
func TestGeneratePrivateKeyPaseto(t *testing.T) {
	privateKey, publicKey := watoken.GenerateKey()
	fmt.Println(privateKey)
	fmt.Println(publicKey)
	hasil, err := watoken.Encode("admin", privateKey)
	fmt.Println(hasil, err)
}

// Hash Pass
func TestGeneratePasswordHash(t *testing.T) {
	password := "adminpasabar"
	hash, _ := HashPassword(password) // ignore error for the sake of simplicity
	fmt.Println("Password:", password)
	fmt.Println("Hash:    ", hash)

	match := CheckPasswordHash(password, hash)
	fmt.Println("Match:   ", match)
}

func TestHashFunction(t *testing.T) {
	mconn := SetConnection("MONGOSTRING", "pasabar13")
	var userdata Admin
	userdata.Username = "pasabaradmin"
	userdata.Password = "adminpasabar"

	filter := bson.M{"username": userdata.Username}
	res := atdb.GetOneDoc[Admin](mconn, "admin", filter)
	fmt.Println("Mongo User Result: ", res)
	hash, _ := HashPassword(userdata.Password)
	fmt.Println("Hash Password : ", hash)
	match := CheckPasswordHash(userdata.Password, res.Password)
	fmt.Println("Match:   ", match)

}

func TestIsPasswordValid(t *testing.T) {
	mconn := SetConnection("MONGOSTRING", "pasabar13")
	var userdata Admin
	userdata.Username = "pasabaradmin"
	userdata.Password = "adminpasabar"

	anu := IsPasswordValid(mconn, "admin", userdata)
	fmt.Println(anu)
}

func TestCatalog(t *testing.T) {
	mconn := SetConnection("MONGOSTRING", "pasabar13")
	var catalogdata Catalog
	catalogdata.Nomorid = 1
	catalogdata.Title = "Pantai Santolo "
	catalogdata.Description = "Rupawan"
	catalogdata.Image = "https://images3.alphacoders.com/165/thumb-1920-165265.jpg"
	CreateNewCatalog(mconn, "catalog", catalogdata)
}

func TestAllCatalog(t *testing.T) {
	mconn := SetConnection("MONGOSTRING", "pasabar13")
	catalog := GetAllCatalog(mconn, "catalog")
	fmt.Println(catalog)
}
