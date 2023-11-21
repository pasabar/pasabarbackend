package pasabarbackend

import (
	"context"
	"fmt"
	"os"

	"github.com/aiteung/atdb"
	"github.com/whatsauth/watoken"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// mongo
func MongoConnect(MongoString, dbname string) *mongo.Database {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(os.Getenv(MongoString)))
	if err != nil {
		fmt.Printf("MongoConnect: %v\n", err)
	}
	return client.Database(dbname)
}

func SetConnection(MONGOCONNSTRINGENV, dbname string) *mongo.Database {
	var DBmongoinfo = atdb.DBInfo{
		// DBString: "mongodb+srv://ryaasishlah123:ryaas123@ryaas.x1atjad.mongodb.net/", //os.Getenv(MONGOCONNSTRINGENV),
		DBString: os.Getenv(MONGOCONNSTRINGENV),
		DBName:   dbname,
	}
	return atdb.MongoConnect(DBmongoinfo)
}

func IsPasswordValid(mongoconn *mongo.Database, collection string, userdata Admin) bool {
	filter := bson.M{"username": userdata.Username}
	res := atdb.GetOneDoc[Admin](mongoconn, collection, filter)
	return CheckPasswordHash(userdata.Password, res.Password)
}

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

func InsertOneDoc(db *mongo.Database, col string, doc interface{}) (insertedID primitive.ObjectID, err error) {
	result, err := db.Collection(col).InsertOne(context.Background(), doc)
	if err != nil {
		return insertedID, fmt.Errorf("kesalahan server : insert")
	}
	insertedID = result.InsertedID.(primitive.ObjectID)
	return insertedID, nil
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

func CreateAdmin(mongoconn *mongo.Database, collection string, userdata Admin) interface{} {
	// Hash the password before storing it
	hashedPassword, err := HashPassword(userdata.Password)
	if err != nil {
		return err
	}
	privateKey, publicKey := watoken.GenerateKey()
	userid := userdata.Username
	tokenstring, err := watoken.Encode(userid, privateKey)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(tokenstring)
	// decode token to get userid
	useridstring := watoken.DecodeGetId(publicKey, tokenstring)
	if useridstring == "" {
		fmt.Println("expire token")
	}
	fmt.Println(useridstring)
	userdata.Private = privateKey
	userdata.Public = publicKey
	userdata.Password = hashedPassword

	// Insert the user data into the database
	return atdb.InsertOneDoc(mongoconn, collection, userdata)
}

// catalog
func CreateNewCatalog(mongoconn *mongo.Database, collection string, catalogdata Catalog) interface{} {
	return atdb.InsertOneDoc(mongoconn, collection, catalogdata)
}

// catalog function
func CreateCatalog(mongoconn *mongo.Database, collection string, catalogdata Catalog) interface{} {
	return atdb.InsertOneDoc(mongoconn, collection, catalogdata)
}

func DeleteCatalog(mongoconn *mongo.Database, collection string, catalogdata Catalog) interface{} {
	filter := bson.M{"nomorid": catalogdata.Nomorid}
	return atdb.DeleteOneDoc(mongoconn, collection, filter)
}

func UpdatedCatalog(mongoconn *mongo.Database, collection string, filter bson.M, catalogdata Catalog) interface{} {
	filter = bson.M{"nomorid": catalogdata.Nomorid}
	return atdb.ReplaceOneDoc(mongoconn, collection, filter, catalogdata)
}

func GetAllCatalog(mongoconn *mongo.Database, collection string) []Catalog {
	catalog := atdb.GetAllDoc[[]Catalog](mongoconn, collection)
	return catalog
}

func GetAllCatalogID(mongoconn *mongo.Database, collection string, catalogdata Catalog) Catalog {
	filter := bson.M{
		"nomorid":     catalogdata.Nomorid,
		"title":       catalogdata.Title,
		"description": catalogdata.Description,
		"image":       catalogdata.Image,
	}
	catalogID := atdb.GetOneDoc[Catalog](mongoconn, collection, filter)
	return catalogID
}
