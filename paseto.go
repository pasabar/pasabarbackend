package pasabarbackend

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	"github.com/aiteung/atapi"
	"github.com/aiteung/atmessage"
	"github.com/whatsauth/wa"
	"github.com/whatsauth/watoken"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// <--- ini Login & Register Admin --->
func Login(Privatekey, MongoEnv, dbname, Colname string, r *http.Request) string {
	var resp Credential
	mconn := SetConnection(MongoEnv, dbname)
	var dataadmin Admin
	err := json.NewDecoder(r.Body).Decode(&dataadmin)
	if err != nil {
		resp.Message = "error parsing application/json: " + err.Error()
	} else {
		if IsPasswordValid(mconn, Colname, dataadmin) {
			tokenstring, err := watoken.Encode(dataadmin.Email, os.Getenv(Privatekey))
			if err != nil {
				resp.Message = "Gagal Encode Token : " + err.Error()
			} else {
				resp.Status = true
				resp.Message = "Selamat Datang Admin"
				resp.Token = tokenstring
			}
		} else {
			resp.Message = "Password Salah"
		}
	}
	return GCFReturnStruct(resp)
}

func LoginWA(token, Privatekey, MongoEnv, dbname, Colname string, r *http.Request) string {
	var resp Credential
	mconn := SetConnection(MongoEnv, dbname)
	var dataadmin Admin
	err := json.NewDecoder(r.Body).Decode(&dataadmin)
	if err != nil {
		resp.Message = "error parsing application/json: " + err.Error()
	} else {
		if IsPasswordValid(mconn, Colname, dataadmin) {
			tokenstring, err := watoken.Encode(dataadmin.Email, os.Getenv(Privatekey))
			if err != nil {
				resp.Message = "Gagal Encode Token : " + err.Error()
			} else {
				resp.Status = true
				resp.Message = "Selamat Datang SUPERADMIN"
				resp.Token = tokenstring
			}

			var email = dataadmin.Email
			var nohp = dataadmin.No_whatsapp

			dt := &wa.TextMessage{
				To:       nohp,
				IsGroup:  false,
				Messages: "Selamat datang Admin PASABAR anda berhasil Login, anda masuk menggunakan akun: " + email + "\nSelamat menggunakanya ya",
			}

			atapi.PostStructWithToken[atmessage.Response]("Token", os.Getenv(token), dt, "https://api.wa.my.id/api/send/message/text")
		} else {
			resp.Message = "Password Salah"
		}
	}
	return ReturnStringStruct(resp)
}

// return struct
func GCFReturnStruct(DataStuct any) string {
	jsondata, _ := json.Marshal(DataStuct)
	return string(jsondata)
}

func ReturnStringStruct(Data any) string {
	jsonee, _ := json.Marshal(Data)
	return string(jsonee)
}

// func Register(Mongoenv, dbname string, r *http.Request) string {
// 	resp := new(Credential)
// 	admindata := new(Admin)
// 	resp.Status = false
// 	conn := SetConnection(Mongoenv, dbname)
// 	err := json.NewDecoder(r.Body).Decode(&admindata)
// 	if err != nil {
// 		resp.Message = "error parsing application/json: " + err.Error()
// 	} else {
// 		resp.Status = true
// 		hash, err := HashPass(admindata.Password)
// 		if err != nil {
// 			resp.Message = "Gagal Hash Password" + err.Error()
// 		}
// 		InsertAdmindata(conn, admindata.Email, admindata.Role, hash)
// 		resp.Message = "Berhasil Input data"
// 	}
// 	response := ReturnStringStruct(resp)
// 	return response
// }

// <--- ini catalog --->

// catalog post
func GCFInsertCatalog(publickey, MONGOCONNSTRINGENV, dbname, colladmin, collcatalog string, r *http.Request) string {
	var response Credential
	response.Status = false
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
	var admindata Admin
	gettoken := r.Header.Get("Login")
	if gettoken == "" {
		response.Message = "Header Login Not Exist"
	} else {
		// Process the request with the "Login" token
		checktoken := watoken.DecodeGetId(os.Getenv(publickey), gettoken)
		admindata.Email = checktoken
		if checktoken == "" {
			response.Message = "Kamu kayaknya belum punya akun"
		} else {
			admin2 := FindAdmin(mconn, colladmin, admindata)
			if admin2.Role == "admin" {
				var datacatalog Catalog
				err := json.NewDecoder(r.Body).Decode(&datacatalog)
				if err != nil {
					response.Message = "Error parsing application/json: " + err.Error()
				} else {
					insertCatalog(mconn, collcatalog, Catalog{
						Nomorid:     datacatalog.Nomorid,
						Title:       datacatalog.Title,
						Description: datacatalog.Description,
						Lokasi:      datacatalog.Lokasi,
						Image:       datacatalog.Image,
						Status:      datacatalog.Status,
					})
					response.Status = true
					response.Message = "Berhasil Insert Catalog"
				}
			} else {
				response.Message = "Anda tidak dapat Insert data karena bukan admin"
			}
		}
	}
	return GCFReturnStruct(response)
}

// delete catalog
func GCFDeleteCatalog(publickey, MONGOCONNSTRINGENV, dbname, colladmin, collcatalog string, r *http.Request) string {

	var respon Credential
	respon.Status = false
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
	var admindata Admin

	gettoken := r.Header.Get("Login")
	if gettoken == "" {
		respon.Message = "Header Login Not Exist"
	} else {
		// Process the request with the "Login" token
		checktoken := watoken.DecodeGetId(os.Getenv(publickey), gettoken)
		admindata.Email = checktoken
		if checktoken == "" {
			respon.Message = "Kamu kayaknya belum punya akun"
		} else {
			admin2 := FindAdmin(mconn, colladmin, admindata)
			if admin2.Role == "admin" {
				var datacatalog Catalog
				err := json.NewDecoder(r.Body).Decode(&datacatalog)
				if err != nil {
					respon.Message = "Error parsing application/json: " + err.Error()
				} else {
					DeleteCatalog(mconn, collcatalog, datacatalog)
					respon.Status = true
					respon.Message = "Berhasil Delete Catalog"
				}
			} else {
				respon.Message = "Anda tidak dapat Delete data karena bukan admin"
			}
		}
	}
	return GCFReturnStruct(respon)
}

// update catalog
func GCFUpdateCatalog(publickey, MONGOCONNSTRINGENV, dbname, colladmin, collcatalog string, r *http.Request) string {
	var response Credential
	response.Status = false
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
	var admindata Admin

	gettoken := r.Header.Get("Login")
	if gettoken == "" {
		response.Message = "Header Login Not Exist"
	} else {
		checktoken := watoken.DecodeGetId(os.Getenv(publickey), gettoken)
		admindata.Email = checktoken
		if checktoken == "" {
			response.Message = "Kamu kayaknya belum punya akun"
		} else {
			admin2 := FindAdmin(mconn, colladmin, admindata)
			if admin2.Role == "admin" {
				var datacatalog Catalog
				err := json.NewDecoder(r.Body).Decode(&datacatalog)
				if err != nil {
					response.Message = "Error parsing application/json: " + err.Error()

				} else {
					UpdatedCatalog(mconn, collcatalog, bson.M{"id": datacatalog.ID}, datacatalog)
					response.Status = true
					response.Message = "Berhasil Update Catalog"
					GCFReturnStruct(CreateResponse(true, "Success Update Catalog", datacatalog))
				}
			} else {
				response.Message = "Anda tidak dapat Update data karena bukan admin"
			}

		}
	}
	return GCFReturnStruct(response)
}

// get all catalog
func GCFGetAllCatalog(MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
	datacatalog := GetAllCatalog(mconn, collectionname)
	if datacatalog != nil {
		return GCFReturnStruct(CreateResponse(true, "success Get All Catalog", datacatalog))
	} else {
		return GCFReturnStruct(CreateResponse(false, "Failed Get All Catalog", datacatalog))
	}
}

// func GCFGetAllCatalogg(publickey, Mongostring, dbname, colname string, r *http.Request) string {
// 	resp := new(Credential)
// 	tokenlogin := r.Header.Get("Login")
// 	if tokenlogin == "" {
// 		resp.Status = false
// 		resp.Message = "Header Login Not Exist"
// 	} else {
// 		existing := IsExist(tokenlogin, os.Getenv(publickey))
// 		if !existing {
// 			resp.Status = false
// 			resp.Message = "Kamu kayaknya belum punya akun"
// 		} else {
// 			koneksyen := SetConnection(Mongostring, dbname)
// 			datacatalog := GetAllCatalog(koneksyen, colname)
// 			yas, _ := json.Marshal(datacatalog)
// 			resp.Status = true
// 			resp.Message = "Data Berhasil diambil"
// 			resp.Token = string(yas)
// 		}
// 	}
// 	return ReturnStringStruct(resp)
// }

func GetAllDataCatalogs(PublicKey, MongoEnv, dbname, colname string, r *http.Request) string {
	req := new(Response)
	conn := SetConnection(MongoEnv, dbname)
	tokenlogin := r.Header.Get("Login")
	if tokenlogin == "" {
		req.Status = false
		req.Message = "Header Login Not Found"
	} else {
		// Dekode token untuk mendapatkan
		_, err := DecodeGetCatalog(os.Getenv(PublicKey), tokenlogin)
		if err != nil {
			req.Status = false
			req.Message = "Data Tersebut tidak ada" + tokenlogin
		} else {
			// Langsung ambil data catalog
			datacatalog := GetAllCatalog(conn, colname)
			if datacatalog == nil {
				req.Status = false
				req.Message = "Data catalog tidak ada"
			} else {
				req.Status = true
				req.Message = "Data Catalog berhasil diambil"
				req.Data = datacatalog
			}
		}
	}
	return ReturnStringStruct(req)
}

// get all catalog by id
func GCFGetAllCatalogID(MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)

	var datacatalog Catalog
	err := json.NewDecoder(r.Body).Decode(&datacatalog)
	if err != nil {
		return err.Error()
	}

	catalog := GetAllCatalogID(mconn, collectionname, datacatalog)
	if catalog != (Catalog{}) {
		return GCFReturnStruct(CreateResponse(true, "Success: Get ID Catalog", datacatalog))
	} else {
		return GCFReturnStruct(CreateResponse(false, "Failed to Get ID Catalog", datacatalog))
	}
}

func GetOneDataCatalog(MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
	resp := new(Credential)
	catalogdata := new(Catalog)
	resp.Status = false
	err := json.NewDecoder(r.Body).Decode(&catalogdata)

	id := r.URL.Query().Get("_id")
	if id == "" {
		resp.Message = "Missing '_id' parameter in the URL"
		return GCFReturnStruct(resp)
	}

	ID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		resp.Message = "Invalid '_id' parameter in the URL"
		return GCFReturnStruct(resp)
	}

	catalogdata.ID = ID

	// Menggunakan fungsi GetProdukFromID untuk mendapatkan data produk berdasarkan ID
	catalogdata, err = GetCatalogFromID(mconn, collectionname, ID)
	if err != nil {
		resp.Message = err.Error()
		return GCFReturnStruct(resp)
	}

	resp.Status = true
	resp.Message = "Get Data Berhasil"
	resp.Data = []Catalog{*catalogdata}

	return GCFReturnStruct(resp)
}

func GetOneDataCatalogs(MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
	resp := new(Credential)
	catalogdata := new(Catalog)
	resp.Status = false

	err := json.NewDecoder(r.Body).Decode(&catalogdata)

	nomorID := r.URL.Query().Get("nomorid")
	if nomorID == "" {
		resp.Message = "Missing 'nomorid' parameter in the URL"
		return GCFReturnStruct(resp)
	}

	ID, err := strconv.Atoi(nomorID)
	if err != nil {
		resp.Message = "Invalid 'nomorid' parameter in the URL"
		return GCFReturnStruct(resp)
	}

	catalogdata.Nomorid = ID // Gunakan casing yang sesuai

	// Menggunakan fungsi GetCatalogFromIDs untuk mendapatkan data produk berdasarkan NomorID
	catalogData, err := GetCatalogFromIDs(mconn, collectionname, ID)
	if err != nil {
		resp.Message = err.Error()
		return GCFReturnStruct(resp)
	}

	resp.Status = true
	resp.Message = "Get Data Berhasil"
	resp.Data = []Catalog{*catalogData}

	return GCFReturnStruct(resp)
}

// <--- ini wisata --->

// wisata post
func GCFInsertWisata(publickey, MONGOCONNSTRINGENV, dbname, colladmin, collwisata string, r *http.Request) string {
	var response Credential
	response.Status = false
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
	var admindata Admin
	gettoken := r.Header.Get("Login")
	if gettoken == "" {
		response.Message = "Header Login Not Exist"
	} else {
		// Process the request with the "Login" token
		checktoken := watoken.DecodeGetId(os.Getenv(publickey), gettoken)
		admindata.Email = checktoken
		if checktoken == "" {
			response.Message = "Kamu kayaknya belum punya akun"
		} else {
			admin2 := FindAdmin(mconn, colladmin, admindata)
			if admin2.Role == "admin" {
				var datawisata Wisata
				err := json.NewDecoder(r.Body).Decode(&datawisata)
				if err != nil {
					response.Message = "Error parsing application/json: " + err.Error()
				} else {
					insertWisata(mconn, collwisata, Wisata{
						Nomorid:     datawisata.Nomorid,
						Title:       datawisata.Title,
						Description: datawisata.Description,
						Lokasi:      datawisata.Lokasi,
						Image:       datawisata.Image,
						Status:      datawisata.Status,
					})
					response.Status = true
					response.Message = "Berhasil Insert Wisata"
				}
			} else {
				response.Message = "Anda tidak dapat Insert data karena bukan admin"
			}
		}
	}
	return GCFReturnStruct(response)
}

// delete wisata
func GCFDeleteWisata(publickey, MONGOCONNSTRINGENV, dbname, colladmin, collwisata string, r *http.Request) string {

	var respon Credential
	respon.Status = false
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
	var admindata Admin

	gettoken := r.Header.Get("Login")
	if gettoken == "" {
		respon.Message = "Header Login Not Exist"
	} else {
		// Process the request with the "Login" token
		checktoken := watoken.DecodeGetId(os.Getenv(publickey), gettoken)
		admindata.Email = checktoken
		if checktoken == "" {
			respon.Message = "Kamu kayaknya belum punya akun"
		} else {
			admin2 := FindAdmin(mconn, colladmin, admindata)
			if admin2.Role == "admin" {
				var datawisata Wisata
				err := json.NewDecoder(r.Body).Decode(&datawisata)
				if err != nil {
					respon.Message = "Error parsing application/json: " + err.Error()
				} else {
					DeleteWisata(mconn, collwisata, datawisata)
					respon.Status = true
					respon.Message = "Berhasil Delete Wisata"
				}
			} else {
				respon.Message = "Anda tidak dapat Delete data karena bukan admin"
			}
		}
	}
	return GCFReturnStruct(respon)
}

// update wisata
func GCFUpdateWisata(publickey, MONGOCONNSTRINGENV, dbname, colladmin, collwisata string, r *http.Request) string {
	var response Credential
	response.Status = false
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
	var admindata Admin

	gettoken := r.Header.Get("Login")
	if gettoken == "" {
		response.Message = "Header Login Not Exist"
	} else {
		checktoken := watoken.DecodeGetId(os.Getenv(publickey), gettoken)
		admindata.Email = checktoken
		if checktoken == "" {
			response.Message = "Kamu kayaknya belum punya akun"
		} else {
			admin2 := FindAdmin(mconn, colladmin, admindata)
			if admin2.Role == "admin" {
				var datawisata Wisata
				err := json.NewDecoder(r.Body).Decode(&datawisata)
				if err != nil {
					response.Message = "Error parsing application/json: " + err.Error()

				} else {
					UpdatedWisata(mconn, collwisata, bson.M{"id": datawisata.ID}, datawisata)
					response.Status = true
					response.Message = "Berhasil Update Wisata"
					GCFReturnStruct(CreateResponse(true, "Success Update Wisata", datawisata))
				}
			} else {
				response.Message = "Anda tidak dapat Update data karena bukan admin"
			}

		}
	}
	return GCFReturnStruct(response)
}

// get all wisata
func GCFGetAllWisata(MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
	datawisata := GetAllWisata(mconn, collectionname)
	if datawisata != nil {
		return GCFReturnStruct(CreateResponse(true, "success Get All Wisata", datawisata))
	} else {
		return GCFReturnStruct(CreateResponse(false, "Failed Get All Wisata", datawisata))
	}
}

// func GCFGetAllWisataa(publickey, Mongostring, dbname, colname string, r *http.Request) string {
// 	resp := new(Credential)
// 	tokenlogin := r.Header.Get("Login")
// 	if tokenlogin == "" {
// 		resp.Status = false
// 		resp.Message = "Header Login Not Exist"
// 	} else {
// 		existing := IsExist(tokenlogin, os.Getenv(publickey))
// 		if !existing {
// 			resp.Status = false
// 			resp.Message = "Kamu kayaknya belum punya akun"
// 		} else {
// 			koneksyen := SetConnection(Mongostring, dbname)
// 			datawisata := GetAllWisata(koneksyen, colname)
// 			yas, _ := json.Marshal(datawisata)
// 			resp.Status = true
// 			resp.Message = "Data Berhasil diambil"
// 			resp.Token = string(yas)
// 		}
// 	}
// 	return ReturnStringStruct(resp)
// }

func GetAllDataWisataa(PublicKey, MongoEnv, dbname, colname string, r *http.Request) string {
	req := new(Response)
	conn := SetConnection(MongoEnv, dbname)
	tokenlogin := r.Header.Get("Login")
	if tokenlogin == "" {
		req.Status = false
		req.Message = "Header Login Not Found"
	} else {
		// Dekode token untuk mendapatkan
		_, err := DecodeGetWisata(os.Getenv(PublicKey), tokenlogin)
		if err != nil {
			req.Status = false
			req.Message = "Data Tersebut tidak ada" + tokenlogin
		} else {
			// Langsung ambil data wisata
			datawisata := GetAllWisata(conn, colname)
			if datawisata == nil {
				req.Status = false
				req.Message = "Data wisata tidak ada"
			} else {
				req.Status = true
				req.Message = "Data Wisata berhasil diambil"
				req.Data = datawisata
			}
		}
	}
	return ReturnStringStruct(req)
}

// get all wisata by id
func GCFGetAllWisataID(MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)

	var datawisata Wisata
	err := json.NewDecoder(r.Body).Decode(&datawisata)
	if err != nil {
		return err.Error()
	}

	wisata := GetAllWisataID(mconn, collectionname, datawisata)
	if wisata != (Wisata{}) {
		return GCFReturnStruct(CreateResponse(true, "Success: Get ID Wisata", datawisata))
	} else {
		return GCFReturnStruct(CreateResponse(false, "Failed to Get ID Wisata", datawisata))
	}
}

func GetOneDataWisatas(MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
	resp := new(Credential)
	wisatadata := new(Wisata)
	resp.Status = false

	err := json.NewDecoder(r.Body).Decode(&wisatadata)

	nomorID := r.URL.Query().Get("nomorid")
	if nomorID == "" {
		resp.Message = "Missing 'nomorid' parameter in the URL"
		return GCFReturnStruct(resp)
	}

	ID, err := strconv.Atoi(nomorID)
	if err != nil {
		resp.Message = "Invalid 'nomorid' parameter in the URL"
		return GCFReturnStruct(resp)
	}

	wisatadata.Nomorid = ID // Gunakan casing yang sesuai

	// Menggunakan fungsi GetWisataFromIDs untuk mendapatkan data produk berdasarkan NomorID
	wisataData, err := GetWisataFromIDs(mconn, collectionname, ID)
	if err != nil {
		resp.Message = err.Error()
		return GCFReturnStruct(resp)
	}

	resp.Status = true
	resp.Message = "Get Data Berhasil"
	resp.Dataw = []Wisata{*wisataData}

	return GCFReturnStruct(resp)
}

func GetOneDataWisata(MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
	resp := new(Credential)
	wisatadata := new(Wisata)
	resp.Status = false
	err := json.NewDecoder(r.Body).Decode(&wisatadata)

	id := r.URL.Query().Get("_id")
	if id == "" {
		resp.Message = "Missing '_id' parameter in the URL"
		return GCFReturnStruct(resp)
	}

	ID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		resp.Message = "Invalid '_id' parameter in the URL"
		return GCFReturnStruct(resp)
	}

	wisatadata.ID = ID

	// Menggunakan fungsi GetProdukFromID untuk mendapatkan data produk berdasarkan ID
	wisatadata, err = GetWisataFromID(mconn, collectionname, ID)
	if err != nil {
		resp.Message = err.Error()
		return GCFReturnStruct(resp)
	}

	resp.Status = true
	resp.Message = "Get Data Berhasil"
	resp.Dataw = []Wisata{*wisatadata}

	return GCFReturnStruct(resp)
}

// <--- ini hotel --->

// hotel post
func GCFInsertHotel(publickey, MONGOCONNSTRINGENV, dbname, colladmin, collhotel string, r *http.Request) string {
	var response Credential
	response.Status = false
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
	var admindata Admin
	gettoken := r.Header.Get("Login")
	if gettoken == "" {
		response.Message = "Header Login Not Exist"
	} else {
		// Process the request with the "Login" token
		checktoken := watoken.DecodeGetId(os.Getenv(publickey), gettoken)
		admindata.Email = checktoken
		if checktoken == "" {
			response.Message = "Kamu kayaknya belum punya akun"
		} else {
			admin2 := FindAdmin(mconn, colladmin, admindata)
			if admin2.Role == "admin" {
				var datahotel Hotel
				err := json.NewDecoder(r.Body).Decode(&datahotel)
				if err != nil {
					response.Message = "Error parsing application/json: " + err.Error()
				} else {
					insertHotel(mconn, collhotel, Hotel{
						Nomorid:     datahotel.Nomorid,
						Title:       datahotel.Title,
						Description: datahotel.Description,
						Lokasi:      datahotel.Lokasi,
						Image:       datahotel.Image,
						Status:      datahotel.Status,
					})
					response.Status = true
					response.Message = "Berhasil Insert Hotel"
				}
			} else {
				response.Message = "Anda tidak dapat Insert data karena bukan admin"
			}
		}
	}
	return GCFReturnStruct(response)
}

// delete hotel
func GCFDeleteHotel(publickey, MONGOCONNSTRINGENV, dbname, colladmin, collhotel string, r *http.Request) string {

	var respon Credential
	respon.Status = false
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
	var admindata Admin

	gettoken := r.Header.Get("Login")
	if gettoken == "" {
		respon.Message = "Header Login Not Exist"
	} else {
		// Process the request with the "Login" token
		checktoken := watoken.DecodeGetId(os.Getenv(publickey), gettoken)
		admindata.Email = checktoken
		if checktoken == "" {
			respon.Message = "Kamu kayaknya belum punya akun"
		} else {
			admin2 := FindAdmin(mconn, colladmin, admindata)
			if admin2.Role == "admin" {
				var datahotel Hotel
				err := json.NewDecoder(r.Body).Decode(&datahotel)
				if err != nil {
					respon.Message = "Error parsing application/json: " + err.Error()
				} else {
					DeleteHotel(mconn, collhotel, datahotel)
					respon.Status = true
					respon.Message = "Berhasil Delete Hotel"
				}
			} else {
				respon.Message = "Anda tidak dapat Delete data karena bukan admin"
			}
		}
	}
	return GCFReturnStruct(respon)
}

// update hotel
func GCFUpdateHotel(publickey, MONGOCONNSTRINGENV, dbname, colladmin, collhotel string, r *http.Request) string {
	var response Credential
	response.Status = false
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
	var admindata Admin

	gettoken := r.Header.Get("Login")
	if gettoken == "" {
		response.Message = "Header Login Not Exist"
	} else {
		checktoken := watoken.DecodeGetId(os.Getenv(publickey), gettoken)
		admindata.Email = checktoken
		if checktoken == "" {
			response.Message = "Kamu kayaknya belum punya akun"
		} else {
			admin2 := FindAdmin(mconn, colladmin, admindata)
			if admin2.Role == "admin" {
				var datahotel Hotel
				err := json.NewDecoder(r.Body).Decode(&datahotel)
				if err != nil {
					response.Message = "Error parsing application/json: " + err.Error()

				} else {
					UpdatedHotel(mconn, collhotel, bson.M{"id": datahotel.ID}, datahotel)
					response.Status = true
					response.Message = "Berhasil Update Hotel"
					GCFReturnStruct(CreateResponse(true, "Success Update Hotel", datahotel))
				}
			} else {
				response.Message = "Anda tidak dapat Update data karena bukan admin"
			}

		}
	}
	return GCFReturnStruct(response)
}

// get all hotel
func GCFGetAllHotel(MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
	datahotel := GetAllHotel(mconn, collectionname)
	if datahotel != nil {
		return GCFReturnStruct(CreateResponse(true, "success Get All Hotel", datahotel))
	} else {
		return GCFReturnStruct(CreateResponse(false, "Failed Get All Hotel", datahotel))
	}
}

// func GCFGetAllHotell(publickey, Mongostring, dbname, colname string, r *http.Request) string {
// 	resp := new(Credential)
// 	tokenlogin := r.Header.Get("Login")
// 	if tokenlogin == "" {
// 		resp.Status = false
// 		resp.Message = "Header Login Not Exist"
// 	} else {
// 		existing := IsExist(tokenlogin, os.Getenv(publickey))
// 		if !existing {
// 			resp.Status = false
// 			resp.Message = "Kamu kayaknya belum punya akun"
// 		} else {
// 			koneksyen := SetConnection(Mongostring, dbname)
// 			datahotel := GetAllHotel(koneksyen, colname)
// 			yas, _ := json.Marshal(datahotel)
// 			resp.Status = true
// 			resp.Message = "Data Berhasil diambil"
// 			resp.Token = string(yas)
// 		}
// 	}
// 	return ReturnStringStruct(resp)
// }

func GetAllDataHotels(PublicKey, MongoEnv, dbname, colname string, r *http.Request) string {
	req := new(Response)
	conn := SetConnection(MongoEnv, dbname)
	tokenlogin := r.Header.Get("Login")
	if tokenlogin == "" {
		req.Status = false
		req.Message = "Header Login Not Found"
	} else {
		// Dekode token untuk mendapatkan
		_, err := DecodeGetHotel(os.Getenv(PublicKey), tokenlogin)
		if err != nil {
			req.Status = false
			req.Message = "Data Tersebut tidak ada" + tokenlogin
		} else {
			// Langsung ambil data hotel
			datahotel := GetAllHotel(conn, colname)
			if datahotel == nil {
				req.Status = false
				req.Message = "Data hotel tidak ada"
			} else {
				req.Status = true
				req.Message = "Data Hotel berhasil diambil"
				req.Data = datahotel
			}
		}
	}
	return ReturnStringStruct(req)
}

// get all hotel by id
func GCFGetAllHotelID(MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)

	var datahotel Hotel
	err := json.NewDecoder(r.Body).Decode(&datahotel)
	if err != nil {
		return err.Error()
	}

	hotel := GetAllHotelID(mconn, collectionname, datahotel)
	if hotel != (Hotel{}) {
		return GCFReturnStruct(CreateResponse(true, "Success: Get ID Hotel", datahotel))
	} else {
		return GCFReturnStruct(CreateResponse(false, "Failed to Get ID Hotel", datahotel))
	}
}

func GetOneDataHotels(MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
	resp := new(Credential)
	restorandata := new(Hotel)
	resp.Status = false

	err := json.NewDecoder(r.Body).Decode(&restorandata)

	nomorID := r.URL.Query().Get("nomorid")
	if nomorID == "" {
		resp.Message = "Missing 'nomorid' parameter in the URL"
		return GCFReturnStruct(resp)
	}

	ID, err := strconv.Atoi(nomorID)
	if err != nil {
		resp.Message = "Invalid 'nomorid' parameter in the URL"
		return GCFReturnStruct(resp)
	}

	restorandata.Nomorid = ID // Gunakan casing yang sesuai

	// Menggunakan fungsi GetHotelFromIDs untuk mendapatkan data produk berdasarkan NomorID
	restoranData, err := GetHotelFromIDs(mconn, collectionname, ID)
	if err != nil {
		resp.Message = err.Error()
		return GCFReturnStruct(resp)
	}

	resp.Status = true
	resp.Message = "Get Data Berhasil"
	resp.Datah = []Hotel{*restoranData}

	return GCFReturnStruct(resp)
}

func GetOneDataHotel(MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
	resp := new(Credential)
	hoteldata := new(Hotel)
	resp.Status = false
	err := json.NewDecoder(r.Body).Decode(&hoteldata)

	id := r.URL.Query().Get("_id")
	if id == "" {
		resp.Message = "Missing '_id' parameter in the URL"
		return GCFReturnStruct(resp)
	}

	ID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		resp.Message = "Invalid '_id' parameter in the URL"
		return GCFReturnStruct(resp)
	}

	hoteldata.ID = ID

	// Menggunakan fungsi GetProdukFromID untuk mendapatkan data produk berdasarkan ID
	hoteldata, err = GetHotelFromID(mconn, collectionname, ID)
	if err != nil {
		resp.Message = err.Error()
		return GCFReturnStruct(resp)
	}

	resp.Status = true
	resp.Message = "Get Data Berhasil"
	resp.Datah = []Hotel{*hoteldata}

	return GCFReturnStruct(resp)
}

// <--- ini restoran --->

// restoran post
func GCFInsertRestoran(publickey, MONGOCONNSTRINGENV, dbname, colladmin, collrestoran string, r *http.Request) string {
	var response Credential
	response.Status = false
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
	var admindata Admin
	gettoken := r.Header.Get("Login")
	if gettoken == "" {
		response.Message = "Header Login Not Exist"
	} else {
		// Process the request with the "Login" token
		checktoken := watoken.DecodeGetId(os.Getenv(publickey), gettoken)
		admindata.Email = checktoken
		if checktoken == "" {
			response.Message = "Kamu kayaknya belum punya akun"
		} else {
			admin2 := FindAdmin(mconn, colladmin, admindata)
			if admin2.Role == "admin" {
				var datarestoran Restoran
				err := json.NewDecoder(r.Body).Decode(&datarestoran)
				if err != nil {
					response.Message = "Error parsing application/json: " + err.Error()
				} else {
					insertRestoran(mconn, collrestoran, Restoran{
						Nomorid:     datarestoran.Nomorid,
						Title:       datarestoran.Title,
						Description: datarestoran.Description,
						Lokasi:      datarestoran.Lokasi,
						Image:       datarestoran.Image,
						Status:      datarestoran.Status,
					})
					response.Status = true
					response.Message = "Berhasil Insert Restoran"
				}
			} else {
				response.Message = "Anda tidak dapat Insert data karena bukan admin"
			}
		}
	}
	return GCFReturnStruct(response)
}

// delete restoran
func GCFDeleteRestoran(publickey, MONGOCONNSTRINGENV, dbname, colladmin, collrestoran string, r *http.Request) string {

	var respon Credential
	respon.Status = false
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
	var admindata Admin

	gettoken := r.Header.Get("Login")
	if gettoken == "" {
		respon.Message = "Header Login Not Exist"
	} else {
		// Process the request with the "Login" token
		checktoken := watoken.DecodeGetId(os.Getenv(publickey), gettoken)
		admindata.Email = checktoken
		if checktoken == "" {
			respon.Message = "Kamu kayaknya belum punya akun"
		} else {
			admin2 := FindAdmin(mconn, colladmin, admindata)
			if admin2.Role == "admin" {
				var datarestoran Restoran
				err := json.NewDecoder(r.Body).Decode(&datarestoran)
				if err != nil {
					respon.Message = "Error parsing application/json: " + err.Error()
				} else {
					DeleteRestoran(mconn, collrestoran, datarestoran)
					respon.Status = true
					respon.Message = "Berhasil Delete Restoran"
				}
			} else {
				respon.Message = "Anda tidak dapat Delete data karena bukan admin"
			}
		}
	}
	return GCFReturnStruct(respon)
}

// update restoran
func GCFUpdateRestoran(publickey, MONGOCONNSTRINGENV, dbname, colladmin, collrestoran string, r *http.Request) string {
	var response Credential
	response.Status = false
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
	var admindata Admin

	gettoken := r.Header.Get("Login")
	if gettoken == "" {
		response.Message = "Header Login Not Exist"
	} else {
		checktoken := watoken.DecodeGetId(os.Getenv(publickey), gettoken)
		admindata.Email = checktoken
		if checktoken == "" {
			response.Message = "Kamu kayaknya belum punya akun"
		} else {
			admin2 := FindAdmin(mconn, colladmin, admindata)
			if admin2.Role == "admin" {
				var datarestoran Restoran
				err := json.NewDecoder(r.Body).Decode(&datarestoran)
				if err != nil {
					response.Message = "Error parsing application/json: " + err.Error()

				} else {
					UpdatedRestoran(mconn, collrestoran, bson.M{"id": datarestoran.ID}, datarestoran)
					response.Status = true
					response.Message = "Berhasil Update Restoran"
					GCFReturnStruct(CreateResponse(true, "Success Update Restoran", datarestoran))
				}
			} else {
				response.Message = "Anda tidak dapat Update data karena bukan admin"
			}

		}
	}
	return GCFReturnStruct(response)
}

// get all restoran
func GCFGetAllRestoran(MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
	datarestoran := GetAllRestoran(mconn, collectionname)
	if datarestoran != nil {
		return GCFReturnStruct(CreateResponse(true, "success Get All Restoran", datarestoran))
	} else {
		return GCFReturnStruct(CreateResponse(false, "Failed Get All Restoran", datarestoran))
	}
}

// func GCFGetAllRestorann(publickey, Mongostring, dbname, colname string, r *http.Request) string {
// 	resp := new(Credential)
// 	tokenlogin := r.Header.Get("Login")
// 	if tokenlogin == "" {
// 		resp.Status = false
// 		resp.Message = "Header Login Not Exist"
// 	} else {
// 		existing := IsExist(tokenlogin, os.Getenv(publickey))
// 		if !existing {
// 			resp.Status = false
// 			resp.Message = "Kamu kayaknya belum punya akun"
// 		} else {
// 			koneksyen := SetConnection(Mongostring, dbname)
// 			datarestoran := GetAllRestoran(koneksyen, colname)
// 			yas, _ := json.Marshal(datarestoran)
// 			resp.Status = true
// 			resp.Message = "Data Berhasil diambil"
// 			resp.Token = string(yas)
// 		}
// 	}
// 	return ReturnStringStruct(resp)
// }

func GetAllDataRestorans(PublicKey, MongoEnv, dbname, colname string, r *http.Request) string {
	req := new(Response)
	conn := SetConnection(MongoEnv, dbname)
	tokenlogin := r.Header.Get("Login")
	if tokenlogin == "" {
		req.Status = false
		req.Message = "Header Login Not Found"
	} else {
		// Dekode token untuk mendapatkan
		_, err := DecodeGetRestoran(os.Getenv(PublicKey), tokenlogin)
		if err != nil {
			req.Status = false
			req.Message = "Data Tersebut tidak ada" + tokenlogin
		} else {
			// Langsung ambil data restoran
			datarestoran := GetAllRestoran(conn, colname)
			if datarestoran == nil {
				req.Status = false
				req.Message = "Data restoran tidak ada"
			} else {
				req.Status = true
				req.Message = "Data Restoran berhasil diambil"
				req.Data = datarestoran
			}
		}
	}
	return ReturnStringStruct(req)
}

// get all restoran by id
func GCFGetAllRestoranID(MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)

	var datarestoran Restoran
	err := json.NewDecoder(r.Body).Decode(&datarestoran)
	if err != nil {
		return err.Error()
	}

	restoran := GetAllRestoranID(mconn, collectionname, datarestoran)
	if restoran != (Restoran{}) {
		return GCFReturnStruct(CreateResponse(true, "Success: Get ID Restoran", datarestoran))
	} else {
		return GCFReturnStruct(CreateResponse(false, "Failed to Get ID Restoran", datarestoran))
	}
}

func GetOneDataRestorans(MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
	resp := new(Credential)
	restorandata := new(Restoran)
	resp.Status = false

	err := json.NewDecoder(r.Body).Decode(&restorandata)

	nomorID := r.URL.Query().Get("nomorid")
	if nomorID == "" {
		resp.Message = "Missing 'nomorid' parameter in the URL"
		return GCFReturnStruct(resp)
	}

	ID, err := strconv.Atoi(nomorID)
	if err != nil {
		resp.Message = "Invalid 'nomorid' parameter in the URL"
		return GCFReturnStruct(resp)
	}

	restorandata.Nomorid = ID // Gunakan casing yang sesuai

	// Menggunakan fungsi GetRestoranFromIDs untuk mendapatkan data produk berdasarkan NomorID
	restoranData, err := GetRestoranFromIDs(mconn, collectionname, ID)
	if err != nil {
		resp.Message = err.Error()
		return GCFReturnStruct(resp)
	}

	resp.Status = true
	resp.Message = "Get Data Berhasil"
	resp.Datar = []Restoran{*restoranData}

	return GCFReturnStruct(resp)
}

func GetOneDataRestoran(MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
	resp := new(Credential)
	restorandata := new(Restoran)
	resp.Status = false
	err := json.NewDecoder(r.Body).Decode(&restorandata)

	id := r.URL.Query().Get("_id")
	if id == "" {
		resp.Message = "Missing '_id' parameter in the URL"
		return GCFReturnStruct(resp)
	}

	ID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		resp.Message = "Invalid '_id' parameter in the URL"
		return GCFReturnStruct(resp)
	}

	restorandata.ID = ID

	// Menggunakan fungsi GetProdukFromID untuk mendapatkan data produk berdasarkan ID
	restorandata, err = GetRestoranFromID(mconn, collectionname, ID)
	if err != nil {
		resp.Message = err.Error()
		return GCFReturnStruct(resp)
	}

	resp.Status = true
	resp.Message = "Get Data Berhasil"
	resp.Datar = []Restoran{*restorandata}

	return GCFReturnStruct(resp)
}

// <--- ini about --->

// about post
func GCFInsertAbout(publickey, MONGOCONNSTRINGENV, dbname, colladmin, collabout string, r *http.Request) string {
	var response Credential
	response.Status = false
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
	var authdata Admin

	gettoken := r.Header.Get("Login")

	if gettoken == "" {
		response.Message = "Header Login Not Exist"
	} else {
		// Process the request with the "Login" token
		checktoken := watoken.DecodeGetId(os.Getenv(publickey), gettoken)
		authdata.Email = checktoken
		if checktoken == "" {
			response.Message = "Kamu kayaknya belum punya akun"
		} else {
			auth2 := FindAdmin(mconn, colladmin, authdata)
			if auth2.Role == "admin" {

				var dataabout About
				err := json.NewDecoder(r.Body).Decode(&dataabout)
				if err != nil {
					response.Message = "Error parsing application/json: " + err.Error()
				} else {
					InsertAbout(mconn, collabout, About{
						ID:          dataabout.ID,
						Title:       dataabout.Title,
						Description: dataabout.Description,
						Image:       dataabout.Image,
						Status:      dataabout.Status,
					})
					response.Status = true
					response.Message = "Berhasil Insert About"
				}
			} else {
				response.Message = "Anda tidak dapat Insert data karena bukan admin"
			}
		}
	}
	return GCFReturnStruct(response)

}

// delete about
func GCFDeleteAbout(publickey, MONGOCONNSTRINGENV, dbname, colladmin, collabout string, r *http.Request) string {
	var response Credential
	response.Status = false
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
	var authdata Admin

	gettoken := r.Header.Get("Login")

	if gettoken == "" {
		response.Message = "Header Login Not Exist"
	} else {
		// Process the request with the "Login" token
		checktoken := watoken.DecodeGetId(os.Getenv(publickey), gettoken)
		authdata.Email = checktoken
		if checktoken == "" {
			response.Message = "Kamu kayaknya belum punya akun"
		} else {
			auth2 := FindAdmin(mconn, colladmin, authdata)
			if auth2.Role == "admin" {

				var dataabout About
				err := json.NewDecoder(r.Body).Decode(&dataabout)
				if err != nil {
					response.Message = "Error parsing application/json: " + err.Error()
				} else {
					DeleteAbout(mconn, collabout, dataabout)
					response.Status = true
					response.Message = "Berhasil Delete About"
					CreateResponse(true, "Success Delete About", dataabout)
				}
			} else {
				response.Message = "Anda tidak dapat Delete data karena bukan admin"
			}
		}
	}
	return GCFReturnStruct(response)
}

// update about
func GCFUpdateAbout(publickey, MONGOCONNSTRINGENV, dbname, colladmin, collabout string, r *http.Request) string {
	var response Credential
	response.Status = false
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
	var authdata Admin

	gettoken := r.Header.Get("Login")

	if gettoken == "" {
		response.Message = "Header Login Not Exist"
	} else {
		// Process the request with the "Login" token
		checktoken := watoken.DecodeGetId(os.Getenv(publickey), gettoken)
		authdata.Email = checktoken
		if checktoken == "" {
			response.Message = "Kamu kayaknya belum punya akun"
		} else {
			auth2 := FindAdmin(mconn, colladmin, authdata)
			if auth2.Role == "admin" {
				var dataabout About
				err := json.NewDecoder(r.Body).Decode(&dataabout)
				if err != nil {
					response.Message = "Error parsing application/json: " + err.Error()
				} else {
					UpdatedAbout(mconn, collabout, bson.M{"id": dataabout.ID}, dataabout)
					response.Status = true
					response.Message = "Berhasil Update Catalog"
					CreateResponse(true, "Success Update About", dataabout)
				}
			} else {
				response.Message = "Anda tidak dapat Update data karena bukan admin"
			}
		}
	}
	return GCFReturnStruct(response)
}

// get all about
func GCFGetAllAbout(MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
	dataabout := GetAllAbout(mconn, collectionname)
	if dataabout != nil {
		return GCFReturnStruct(CreateResponse(true, "Berhasil Get All About", dataabout))
	} else {
		return GCFReturnStruct(CreateResponse(false, "Gagal Get All About", dataabout))
	}
}

// func GCFGetAllAboutt(publickey, Mongostring, dbname, colname string, r *http.Request) string {
// 	resp := new(Credential)
// 	tokenlogin := r.Header.Get("Login")
// 	if tokenlogin == "" {
// 		resp.Status = false
// 		resp.Message = "Header Login Not Exist"
// 	} else {
// 		existing := IsExist(tokenlogin, os.Getenv(publickey))
// 		if !existing {
// 			resp.Status = false
// 			resp.Message = "Kamu kayaknya belum punya akun"
// 		} else {
// 			koneksyen := SetConnection(Mongostring, dbname)
// 			dataabaout := GetAllAbout(koneksyen, colname)
// 			yas, _ := json.Marshal(dataabaout)
// 			resp.Status = true
// 			resp.Message = "Data About Berhasil diambil"
// 			resp.Token = string(yas)
// 		}
// 	}
// 	return ReturnStringStruct(resp)
// }

func GetAllDataAbouts(PublicKey, MongoEnv, dbname, colname string, r *http.Request) string {
	req := new(Response)
	conn := SetConnection(MongoEnv, dbname)
	tokenlogin := r.Header.Get("Login")
	if tokenlogin == "" {
		req.Status = false
		req.Message = "Header Login Not Found"
	} else {
		// Dekode token untuk mendapatkan
		_, err := DecodeGetAbout(os.Getenv(PublicKey), tokenlogin)
		if err != nil {
			req.Status = false
			req.Message = "Data Tersebut tidak ada" + tokenlogin
		} else {
			// Langsung ambil data about
			dataabout := GetAllAbout(conn, colname)
			if dataabout == nil {
				req.Status = false
				req.Message = "Data about tidak ada"
			} else {
				req.Status = true
				req.Message = "Data about berhasil diambil"
				req.Data = dataabout
			}
		}
	}
	return ReturnStringStruct(req)
}

// <--- ini contact --->

// contact post
func GCFInsertContact(MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
	var datacontact Contact
	err := json.NewDecoder(r.Body).Decode(&datacontact)
	if err != nil {
		return err.Error()
	}

	if err := InsertContact(mconn, collectionname, datacontact); err != nil {
		return GCFReturnStruct(CreateResponse(true, "Success Create Contact", datacontact))
	} else {
		return GCFReturnStruct(CreateResponse(false, "Failed Create Contact", datacontact))
	}
}

// get all contact
func GCFGetAllContacts(MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
	datacontact := GetAllContact(mconn, collectionname)
	if datacontact != nil {
		return GCFReturnStruct(CreateResponse(true, "success Get All Contact", datacontact))
	} else {
		return GCFReturnStruct(CreateResponse(false, "Failed Get All Contact", datacontact))
	}
}

func GCFGetAllContactt(publickey, Mongostring, dbname, colname string, r *http.Request) string {
	resp := new(Credential)
	tokenlogin := r.Header.Get("Login")
	if tokenlogin == "" {
		resp.Status = false
		resp.Message = "Header Login Not Exist"
	} else {
		existing := IsExist(tokenlogin, os.Getenv(publickey))
		if !existing {
			resp.Status = false
			resp.Message = "Kamu kayaknya belum punya akun"
		} else {
			koneksyen := SetConnection(Mongostring, dbname)
			datacontact := GetAllContact(koneksyen, colname)
			yas, _ := json.Marshal(datacontact)
			resp.Status = true
			resp.Message = "Data Contact Berhasil diambil"
			resp.Token = string(yas)
		}
	}
	return ReturnStringStruct(resp)
}

// <--- ini crawling --->

// get all crawling
func GCFGetAllCrawling(MONGOCONNSTRINGENV, dbname, collectionname string) string {
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
	datacrawling := GetAllCrawling(mconn, collectionname)
	if datacrawling != nil {
		return GCFReturnStruct(CreateResponse(true, "success Get All Crawling", datacrawling))
	} else {
		return GCFReturnStruct(CreateResponse(false, "Failed Get All Crawling", datacrawling))
	}
}

// <--- ini kesimpulan --->

// kesimpulan post
func GCFInsertKesimpulan(publickey, MONGOCONNSTRINGENV, dbname, colladmin, collkesimpulan string, r *http.Request) string {
	var response Credential
	response.Status = false
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
	var admindata Admin
	gettoken := r.Header.Get("Login")
	if gettoken == "" {
		response.Message = "Header Login Not Exist"
	} else {
		// Process the request with the "Login" token
		checktoken := watoken.DecodeGetId(os.Getenv(publickey), gettoken)
		admindata.Email = checktoken
		if checktoken == "" {
			response.Message = "Kamu kayaknya belum punya akun"
		} else {
			admin2 := FindAdmin(mconn, colladmin, admindata)
			if admin2.Role == "admin" {
				var datakesimpulan Kesimpulan
				err := json.NewDecoder(r.Body).Decode(&datakesimpulan)
				if err != nil {
					response.Message = "Error parsing application/json: " + err.Error()
				} else {
					insertKesimpulan(mconn, collkesimpulan, Kesimpulan{
						Nomorid:     datakesimpulan.Nomorid,
						Ticket:      datakesimpulan.Ticket,
						Parkir:      datakesimpulan.Parkir,
						Jarak:       datakesimpulan.Jarak,
						Pemandangan: datakesimpulan.Pemandangan,
						Kelebihan:   datakesimpulan.Kelebihan,
						Kekurangan:  datakesimpulan.Kekurangan,
						Status:      datakesimpulan.Status,
					})
					response.Status = true
					response.Message = "Berhasil Insert Kesimpulan"
				}
			} else {
				response.Message = "Anda tidak dapat Insert data karena bukan admin"
			}
		}
	}
	return GCFReturnStruct(response)
}

// delete kesimpulan
func GCFDeleteKesimpulan(publickey, MONGOCONNSTRINGENV, dbname, colladmin, collkesimpulan string, r *http.Request) string {

	var respon Credential
	respon.Status = false
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
	var admindata Admin

	gettoken := r.Header.Get("Login")
	if gettoken == "" {
		respon.Message = "Header Login Not Exist"
	} else {
		// Process the request with the "Login" token
		checktoken := watoken.DecodeGetId(os.Getenv(publickey), gettoken)
		admindata.Email = checktoken
		if checktoken == "" {
			respon.Message = "Kamu kayaknya belum punya akun"
		} else {
			admin2 := FindAdmin(mconn, colladmin, admindata)
			if admin2.Role == "admin" {
				var datakesimpulan Kesimpulan
				err := json.NewDecoder(r.Body).Decode(&datakesimpulan)
				if err != nil {
					respon.Message = "Error parsing application/json: " + err.Error()
				} else {
					DeleteKesimpulan(mconn, collkesimpulan, datakesimpulan)
					respon.Status = true
					respon.Message = "Berhasil Delete Kesimpulan"
				}
			} else {
				respon.Message = "Anda tidak dapat Delete data karena bukan admin"
			}
		}
	}
	return GCFReturnStruct(respon)
}

// update kesimpulan
func GCFUpdateKesimpulan(publickey, MONGOCONNSTRINGENV, dbname, colladmin, collkesimpulan string, r *http.Request) string {
	var response Credential
	response.Status = false
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
	var admindata Admin

	gettoken := r.Header.Get("Login")
	if gettoken == "" {
		response.Message = "Header Login Not Exist"
	} else {
		checktoken := watoken.DecodeGetId(os.Getenv(publickey), gettoken)
		admindata.Email = checktoken
		if checktoken == "" {
			response.Message = "Kamu kayaknya belum punya akun"
		} else {
			admin2 := FindAdmin(mconn, colladmin, admindata)
			if admin2.Role == "admin" {
				var datakesimpulan Kesimpulan
				err := json.NewDecoder(r.Body).Decode(&datakesimpulan)
				if err != nil {
					response.Message = "Error parsing application/json: " + err.Error()

				} else {
					UpdatedKesimpulan(mconn, collkesimpulan, bson.M{"id": datakesimpulan.ID}, datakesimpulan)
					response.Status = true
					response.Message = "Berhasil Update Kesimpulan"
					GCFReturnStruct(CreateResponse(true, "Success Update Kesimpulan", datakesimpulan))
				}
			} else {
				response.Message = "Anda tidak dapat Update data karena bukan admin"
			}

		}
	}
	return GCFReturnStruct(response)
}

// get all kesimpulan
func GCFGetAllKesimpulan(MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
	datakesimpulan := GetAllKesimpulan(mconn, collectionname)
	if datakesimpulan != nil {
		return GCFReturnStruct(CreateResponse(true, "success Get All Kesimpulan", datakesimpulan))
	} else {
		return GCFReturnStruct(CreateResponse(false, "Failed Get All Kesimpulan", datakesimpulan))
	}
}

// func GCFGetAllKesimpulann(publickey, Mongostring, dbname, colname string, r *http.Request) string {
// 	resp := new(Credential)
// 	tokenlogin := r.Header.Get("Login")
// 	if tokenlogin == "" {
// 		resp.Status = false
// 		resp.Message = "Header Login Not Exist"
// 	} else {
// 		existing := IsExist(tokenlogin, os.Getenv(publickey))
// 		if !existing {
// 			resp.Status = false
// 			resp.Message = "Kamu kayaknya belum punya akun"
// 		} else {
// 			koneksyen := SetConnection(Mongostring, dbname)
// 			datakesimpulan := GetAllKesimpulan(koneksyen, colname)
// 			yas, _ := json.Marshal(datakesimpulan)
// 			resp.Status = true
// 			resp.Message = "Data Berhasil diambil"
// 			resp.Token = string(yas)
// 		}
// 	}
// 	return ReturnStringStruct(resp)
// }

func GetAllDataKesimpulans(PublicKey, MongoEnv, dbname, colname string, r *http.Request) string {
	req := new(Response)
	conn := SetConnection(MongoEnv, dbname)
	tokenlogin := r.Header.Get("Login")
	if tokenlogin == "" {
		req.Status = false
		req.Message = "Header Login Not Found"
	} else {
		// Dekode token untuk mendapatkan
		_, err := DecodeGetKesimpulan(os.Getenv(PublicKey), tokenlogin)
		if err != nil {
			req.Status = false
			req.Message = "Data Tersebut tidak ada" + tokenlogin
		} else {
			// Langsung ambil data kesimpulan
			datakesimpulan := GetAllKesimpulan(conn, colname)
			if datakesimpulan == nil {
				req.Status = false
				req.Message = "Data kesimpulan tidak ada"
			} else {
				req.Status = true
				req.Message = "Data Kesimpulan berhasil diambil"
				req.Data = datakesimpulan
			}
		}
	}
	return ReturnStringStruct(req)
}

// get all kesimpulan by id
func GCFGetAllKesimpulanID(MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)

	var datakesimpulan Kesimpulan
	err := json.NewDecoder(r.Body).Decode(&datakesimpulan)
	if err != nil {
		return err.Error()
	}

	kesimpulan := GetAllKesimpulanID(mconn, collectionname, datakesimpulan)
	if kesimpulan != (Kesimpulan{}) {
		return GCFReturnStruct(CreateResponse(true, "Success: Get ID Kesimpulan", datakesimpulan))
	} else {
		return GCFReturnStruct(CreateResponse(false, "Failed to Get ID Kesimpulan", datakesimpulan))
	}
}

func GetOneDataKesimpulans(MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
	resp := new(Credential)
	kesimpulandata := new(Kesimpulan)
	resp.Status = false

	err := json.NewDecoder(r.Body).Decode(&kesimpulandata)

	nomorID := r.URL.Query().Get("nomorid")
	if nomorID == "" {
		resp.Message = "Missing 'nomorid' parameter in the URL"
		return GCFReturnStruct(resp)
	}

	ID, err := strconv.Atoi(nomorID)
	if err != nil {
		resp.Message = "Invalid 'nomorid' parameter in the URL"
		return GCFReturnStruct(resp)
	}

	kesimpulandata.Nomorid = ID // Gunakan casing yang sesuai

	// Menggunakan fungsi GetKesimpulanFromIDs untuk mendapatkan data produk berdasarkan NomorID
	kesimpulanData, err := GetKesimpulanFromIDs(mconn, collectionname, ID)
	if err != nil {
		resp.Message = err.Error()
		return GCFReturnStruct(resp)
	}

	resp.Status = true
	resp.Message = "Get Data Berhasil"
	resp.Datak = []Kesimpulan{*kesimpulanData}

	return GCFReturnStruct(resp)
}
func GetOneDataKesimpulan(MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
	resp := new(Credential)
	kesimpulandata := new(Kesimpulan)
	resp.Status = false
	err := json.NewDecoder(r.Body).Decode(&kesimpulandata)

	id := r.URL.Query().Get("_id")
	if id == "" {
		resp.Message = "Missing '_id' parameter in the URL"
		return GCFReturnStruct(resp)
	}

	ID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		resp.Message = "Invalid '_id' parameter in the URL"
		return GCFReturnStruct(resp)
	}

	kesimpulandata.ID = ID

	// Menggunakan fungsi GetProdukFromID untuk mendapatkan data produk berdasarkan ID
	kesimpulandata, err = GetKesimpulanFromID(mconn, collectionname, ID)
	if err != nil {
		resp.Message = err.Error()
		return GCFReturnStruct(resp)
	}

	resp.Status = true
	resp.Message = "Get Data Berhasil"
	resp.Datak = []Kesimpulan{*kesimpulandata}

	return GCFReturnStruct(resp)
}
