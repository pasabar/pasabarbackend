package pasabarbackend

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/aiteung/atdb"
	"github.com/whatsauth/watoken"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// crud
func GetAllDocs(db *mongo.Database, col string, docs interface{}) interface{} {
	collection := db.Collection(col)
	filter := bson.M{}
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return fmt.Errorf("error GetAllDocs %s: %s", col, err)
	}
	err = cursor.All(context.TODO(), &docs)
	if err != nil {
		return err
	}
	return docs
}

func UpdateOneDoc(id primitive.ObjectID, db *mongo.Database, col string, doc interface{}) (err error) {
	filter := bson.M{"_id": id}
	result, err := db.Collection(col).UpdateOne(context.Background(), filter, bson.M{"$set": doc})
	if err != nil {
		return fmt.Errorf("error update: %v", err)
	}
	if result.ModifiedCount == 0 {
		err = fmt.Errorf("tidak ada data yang diubah")
		return
	}
	return nil
}

func DeleteOneDoc(_id primitive.ObjectID, db *mongo.Database, col string) error {
	collection := db.Collection(col)
	filter := bson.M{"_id": _id}
	result, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return fmt.Errorf("error deleting data for ID %s: %s", _id, err.Error())
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("data with ID %s not found", _id)
	}

	return nil
}

// admin
func CreateNewAdminRole(mongoconn *mongo.Database, collection string, admindata Admin) interface{} {
	// Hash the password before storing it
	hashedPassword, err := HashPass(admindata.Password)
	if err != nil {
		return err
	}
	admindata.Password = hashedPassword

	// Insert the admin data into the database
	return atdb.InsertOneDoc(mongoconn, collection, admindata)
}

func CreateAdminAndAddToken(privateKeyEnv string, mongoconn *mongo.Database, collection string, admindata Admin) error {
	// Hash the password before storing it
	hashedPassword, err := HashPass(admindata.Password)
	if err != nil {
		return err
	}
	admindata.Password = hashedPassword

	// Create a token for the admin
	tokenstring, err := watoken.Encode(admindata.Email, os.Getenv(privateKeyEnv))
	if err != nil {
		return err
	}

	admindata.Token = tokenstring

	// Insert the admin data into the MongoDB collection
	if err := atdb.InsertOneDoc(mongoconn, collection, admindata.Email); err != nil {
		return nil // Mengembalikan kesalahan yang dikembalikan oleh atdb.InsertOneDoc
	}

	// Return nil to indicate success
	return nil
}

func CreateResponse(status bool, message string, data interface{}) Response {
	response := Response{
		Status:  status,
		Message: message,
		Data:    data,
	}
	return response
}

func CreateAdmin(mongoconn *mongo.Database, collection string, admindata Admin) interface{} {
	// Hash the password before storing it
	hashedPassword, err := HashPass(admindata.Password)
	if err != nil {
		return err
	}
	privateKey, publicKey := watoken.GenerateKey()
	adminid := admindata.Email
	tokenstring, err := watoken.Encode(adminid, privateKey)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(tokenstring)
	// decode token to get adminid
	adminidstring := watoken.DecodeGetId(publicKey, tokenstring)
	if adminidstring == "" {
		fmt.Println("expire token")
	}
	fmt.Println(adminidstring)
	admindata.Private = privateKey
	admindata.Public = publicKey
	admindata.Password = hashedPassword

	// Insert the admin data into the database
	return atdb.InsertOneDoc(mongoconn, collection, admindata)
}

// catalog
func CreateNewCatalog(mongoconn *mongo.Database, collection string, catalogdata Catalog) interface{} {
	return atdb.InsertOneDoc(mongoconn, collection, catalogdata)
}

// catalog function
func insertCatalog(mongoconn *mongo.Database, collection string, catalogdata Catalog) interface{} {
	return atdb.InsertOneDoc(mongoconn, collection, catalogdata)
}

func DeleteCatalog(mongoconn *mongo.Database, collection string, catalogdata Catalog) interface{} {
	filter := bson.M{"nomorid": catalogdata.Nomorid}
	return atdb.DeleteOneDoc(mongoconn, collection, filter)
}

func UpdatedCatalog(mongoconn *mongo.Database, collection string, filter bson.M, catalogdata Catalog) interface{} {
	updatedFilter := bson.M{"nomorid": catalogdata.Nomorid}
	return atdb.ReplaceOneDoc(mongoconn, collection, updatedFilter, catalogdata)
}

func GetAllCatalog(mongoconn *mongo.Database, collection string) []Catalog {
	catalog := atdb.GetAllDoc[[]Catalog](mongoconn, collection)
	return catalog
}
func GetAllCatalogs(MongoConn *mongo.Database, colname string, email string) []Admin {
	data := atdb.GetAllDoc[[]Admin](MongoConn, colname)
	return data
}

func GetAllCatalogID(mongoconn *mongo.Database, collection string, catalogdata Catalog) Catalog {
	filter := bson.M{
		"nomorid":     catalogdata.Nomorid,
		"title":       catalogdata.Title,
		"description": catalogdata.Description,
		"image":       catalogdata.Image,
		"lokasi":      catalogdata.Lokasi,
	}
	catalogID := atdb.GetOneDoc[Catalog](mongoconn, collection, filter)
	return catalogID
}

// wisata function
func insertWisata(mongoconn *mongo.Database, collection string, wisatadata Wisata) interface{} {
	return atdb.InsertOneDoc(mongoconn, collection, wisatadata)
}

func DeleteWisata(mongoconn *mongo.Database, collection string, wisatadata Wisata) interface{} {
	filter := bson.M{"nomorid": wisatadata.Nomorid}
	return atdb.DeleteOneDoc(mongoconn, collection, filter)
}

func UpdatedWisata(mongoconn *mongo.Database, collection string, filter bson.M, wisatadata Wisata) interface{} {
	updatedFilter := bson.M{"nomorid": wisatadata.Nomorid}
	return atdb.ReplaceOneDoc(mongoconn, collection, updatedFilter, wisatadata)
}

func GetAllWisata(mongoconn *mongo.Database, collection string) []Wisata {
	wisata := atdb.GetAllDoc[[]Wisata](mongoconn, collection)
	return wisata
}

// func GetAllWisatas(MongoConn *mongo.Database, colname string, email string) []Admin {
// 	data := atdb.GetAllDoc[[]Admin](MongoConn, colname)
// 	return data
// }

func GetAllWisataID(mongoconn *mongo.Database, collection string, wisatadata Wisata) Wisata {
	filter := bson.M{
		"nomorid":     wisatadata.Nomorid,
		"title":       wisatadata.Title,
		"description": wisatadata.Description,
		"lokasi":      wisatadata.Lokasi,
		"image":       wisatadata.Image,
	}
	wisataID := atdb.GetOneDoc[Wisata](mongoconn, collection, filter)
	return wisataID
}

func GetCatalogFromID(db *mongo.Database, col string, _id primitive.ObjectID) (*Catalog, error) {
	cols := db.Collection(col)
	filter := bson.M{"_id": _id}

	cataloglist := new(Catalog)

	err := cols.FindOne(context.Background(), filter).Decode(cataloglist)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("no data found for ID %s", _id.Hex())
		}
		return nil, fmt.Errorf("error retrieving data for ID %s: %s", _id.Hex(), err.Error())
	}

	return cataloglist, nil
}

// hotel function
func insertHotel(mongoconn *mongo.Database, collection string, hoteldata Hotel) interface{} {
	return atdb.InsertOneDoc(mongoconn, collection, hoteldata)
}

func DeleteHotel(mongoconn *mongo.Database, collection string, hoteldata Hotel) interface{} {
	filter := bson.M{"nomorid": hoteldata.Nomorid}
	return atdb.DeleteOneDoc(mongoconn, collection, filter)
}

func UpdatedHotel(mongoconn *mongo.Database, collection string, filter bson.M, hoteldata Hotel) interface{} {
	updatedFilter := bson.M{"nomorid": hoteldata.Nomorid}
	return atdb.ReplaceOneDoc(mongoconn, collection, updatedFilter, hoteldata)
}

func GetAllHotel(mongoconn *mongo.Database, collection string) []Hotel {
	hotel := atdb.GetAllDoc[[]Hotel](mongoconn, collection)
	return hotel
}

// func GetAllHotels(MongoConn *mongo.Database, colname string, email string) []Admin {
// 	data := atdb.GetAllDoc[[]Admin](MongoConn, colname)
// 	return data
// }

func GetAllHotelID(mongoconn *mongo.Database, collection string, hoteldata Hotel) Hotel {
	filter := bson.M{
		"nomorid":     hoteldata.Nomorid,
		"title":       hoteldata.Title,
		"description": hoteldata.Description,
		"lokasi":      hoteldata.Lokasi,
		"image":       hoteldata.Image,
	}
	hotelID := atdb.GetOneDoc[Hotel](mongoconn, collection, filter)
	return hotelID
}

// restoran function
func insertRestoran(mongoconn *mongo.Database, collection string, restorandata Restoran) interface{} {
	return atdb.InsertOneDoc(mongoconn, collection, restorandata)
}

func DeleteRestoran(mongoconn *mongo.Database, collection string, restorandata Restoran) interface{} {
	filter := bson.M{"nomorid": restorandata.Nomorid}
	return atdb.DeleteOneDoc(mongoconn, collection, filter)
}

func UpdatedRestoran(mongoconn *mongo.Database, collection string, filter bson.M, restorandata Restoran) interface{} {
	updatedFilter := bson.M{"nomorid": restorandata.Nomorid}
	return atdb.ReplaceOneDoc(mongoconn, collection, updatedFilter, restorandata)
}

func GetAllRestoran(mongoconn *mongo.Database, collection string) []Restoran {
	restoran := atdb.GetAllDoc[[]Restoran](mongoconn, collection)
	return restoran
}

// func GetAllRestorans(MongoConn *mongo.Database, colname string, email string) []Admin {
// 	data := atdb.GetAllDoc[[]Admin](MongoConn, colname)
// 	return data
// }

func GetAllRestoranID(mongoconn *mongo.Database, collection string, restorandata Restoran) Restoran {
	filter := bson.M{
		"nomorid":     restorandata.Nomorid,
		"title":       restorandata.Title,
		"description": restorandata.Description,
		"lokasi":      restorandata.Lokasi,
		"image":       restorandata.Image,
	}
	restoranID := atdb.GetOneDoc[Restoran](mongoconn, collection, filter)
	return restoranID
}

// about function

func InsertAbout(mongoconn *mongo.Database, collection string, aboutdata About) interface{} {
	return atdb.InsertOneDoc(mongoconn, collection, aboutdata)
}

func DeleteAbout(mongoconn *mongo.Database, collection string, aboutdata About) interface{} {
	filter := bson.M{"id": aboutdata.ID}
	return atdb.DeleteOneDoc(mongoconn, collection, filter)
}

func UpdatedAbout(mongoconn *mongo.Database, collection string, filter bson.M, aboutdata About) interface{} {
	updatedFilter := bson.M{"id": aboutdata.ID}
	return atdb.ReplaceOneDoc(mongoconn, collection, updatedFilter, aboutdata)
}

func GetAllAbout(mongoconn *mongo.Database, collection string) []About {
	about := atdb.GetAllDoc[[]About](mongoconn, collection)
	return about
}

// contact function

func InsertContact(mongoconn *mongo.Database, collection string, contactdata Contact) interface{} {
	return atdb.InsertOneDoc(mongoconn, collection, contactdata)
}

func GetAllContact(mongoconn *mongo.Database, collection string) []Contact {
	contact := atdb.GetAllDoc[[]Contact](mongoconn, collection)
	return contact
}

func GetIdContact(mongoconn *mongo.Database, collection string, contactdata Contact) Contact {
	filter := bson.M{"id": contactdata.ID}
	return atdb.GetOneDoc[Contact](mongoconn, collection, filter)
}

//crawling function

func GetAllCrawling(mongoconn *mongo.Database, collection string) []Crawling {
	crawling := atdb.GetAllDoc[[]Crawling](mongoconn, collection)
	return crawling
}

// kesimpulan function
func insertKesimpulan(mongoconn *mongo.Database, collection string, kesimpulandata Kesimpulan) interface{} {
	return atdb.InsertOneDoc(mongoconn, collection, kesimpulandata)
}

func DeleteKesimpulan(mongoconn *mongo.Database, collection string, kesimpulandata Kesimpulan) interface{} {
	filter := bson.M{"nomorid": kesimpulandata.Nomorid}
	return atdb.DeleteOneDoc(mongoconn, collection, filter)
}

func UpdatedKesimpulan(mongoconn *mongo.Database, collection string, filter bson.M, kesimpulandata Kesimpulan) interface{} {
	updatedFilter := bson.M{"nomorid": kesimpulandata.Nomorid}
	return atdb.ReplaceOneDoc(mongoconn, collection, updatedFilter, kesimpulandata)
}

func GetAllKesimpulan(mongoconn *mongo.Database, collection string) []Kesimpulan {
	kesimpulan := atdb.GetAllDoc[[]Kesimpulan](mongoconn, collection)
	return kesimpulan
}

// func GetAllKesimpulans(MongoConn *mongo.Database, colname string, email string) []Admin {
// 	data := atdb.GetAllDoc[[]Admin](MongoConn, colname)
// 	return data
// }

func GetAllKesimpulanID(mongoconn *mongo.Database, collection string, kesimpulandata Kesimpulan) Kesimpulan {
	filter := bson.M{
		"nomorid":     kesimpulandata.Nomorid,
		"ticket":      kesimpulandata.Ticket,
		"parkir":      kesimpulandata.Parkir,
		"jarak":       kesimpulandata.Jarak,
		"pemandangan": kesimpulandata.Pemandangan,
		"kelebihan":   kesimpulandata.Kelebihan,
		"kekurangan":  kesimpulandata.Kekurangan,
		"status":      kesimpulandata.Status,
	}
	kesimpulanID := atdb.GetOneDoc[Kesimpulan](mongoconn, collection, filter)
	return kesimpulanID
}
