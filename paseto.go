package pasabarbackend

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/whatsauth/watoken"
	"go.mongodb.org/mongo-driver/bson"
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
				resp.Message = "Selamat Datang SUPERADMIN"
				resp.Token = tokenstring
			}
		} else {
			resp.Message = "Password Salah"
		}
	}
	return GCFReturnStruct(resp)
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

func Register(Mongoenv, dbname string, r *http.Request) string {
	resp := new(Credential)
	admindata := new(Admin)
	resp.Status = false
	conn := SetConnection(Mongoenv, dbname)
	err := json.NewDecoder(r.Body).Decode(&admindata)
	if err != nil {
		resp.Message = "error parsing application/json: " + err.Error()
	} else {
		resp.Status = true
		hash, err := HashPass(admindata.Password)
		if err != nil {
			resp.Message = "Gagal Hash Password" + err.Error()
		}
		InsertAdmindata(conn, admindata.Email, admindata.Role, hash)
		resp.Message = "Berhasil Input data"
	}
	response := ReturnStringStruct(resp)
	return response
}

// <--- ini catalog --->

// catalog post
func GCFInsertCatalog(publickey, MONGOCONNSTRINGENV, dbname, colladmin, collcatalog string, r *http.Request) string {
	var response Credential
	response.Status = false
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
	var admindata Admin
	gettoken := r.Header.Get("token")
	if gettoken == "" {
		response.Message = "Missing token in headers"
	} else {
		// Process the request with the "Login" token
		checktoken := watoken.DecodeGetId(os.Getenv(publickey), gettoken)
		admindata.Email = checktoken
		if checktoken == "" {
			response.Message = "Invalid token"
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

	gettoken := r.Header.Get("token")
	if gettoken == "" {
		respon.Message = "Missing token in headers"
	} else {
		// Process the request with the "Login" token
		checktoken := watoken.DecodeGetId(os.Getenv(publickey), gettoken)
		admindata.Email = checktoken
		if checktoken == "" {
			respon.Message = "Invalid token"
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

	gettoken := r.Header.Get("token")
	if gettoken == "" {
		response.Message = "Missing token in Headers"
	} else {
		checktoken := watoken.DecodeGetId(os.Getenv(publickey), gettoken)
		admindata.Email = checktoken
		if checktoken == "" {
			response.Message = "Invalid token"
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

func GCFGetAllCatalogs(publickey, MONGOCONNSTRINGENV, dbname, colladmin, collcatalog string, r *http.Request) string {
	var response Credential
	response.Status = false
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
	var admindata Admin

	gettoken := r.Header.Get("token")
	if gettoken == "" {
		response.Message = "Missing token in Headers"
	} else {
		checktoken := watoken.DecodeGetId(os.Getenv(publickey), gettoken)
		admindata.Email = checktoken
		if checktoken == "" {
			response.Message = "Invalid token"
		} else {
			admin2 := FindAdmin(mconn, colladmin, admindata)
			if admin2.Role == "admin" {
				var datacatalog Catalog
				err := json.NewDecoder(r.Body).Decode(&datacatalog)
				if err != nil {
					response.Message = "Error parsing application/json: " + err.Error()

				} else {
					GetAllCatalog(mconn, collcatalog)
					response.Status = true
					response.Message = "Berhasil Ambil data"
					GCFReturnStruct(CreateResponse(true, "Success Get Catalog", datacatalog))
				}
			} else {
				response.Message = "Anda tidak dapat Get data karena bukan admin"
			}

		}
	}
	return GCFReturnStruct(response)
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

// <--- ini about --->

// about post
func GCFInsertAbout(publickey, MONGOCONNSTRINGENV, dbname, colladmin, collabout string, r *http.Request) string {
	var response Credential
	response.Status = false
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
	var authdata Admin

	gettoken := r.Header.Get("token")

	if gettoken == "" {
		response.Message = "Missing token in headers"
	} else {
		// Process the request with the "Login" token
		checktoken := watoken.DecodeGetId(os.Getenv(publickey), gettoken)
		authdata.Email = checktoken
		if checktoken == "" {
			response.Message = "Invalid token"
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

	gettoken := r.Header.Get("token")

	if gettoken == "" {
		response.Message = "Missing token in headers"
	} else {
		// Process the request with the "Login" token
		checktoken := watoken.DecodeGetId(os.Getenv(publickey), gettoken)
		authdata.Email = checktoken
		if checktoken == "" {
			response.Message = "Invalid token"
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

	gettoken := r.Header.Get("token")

	if gettoken == "" {
		response.Message = "Missing token in headers"
	} else {
		// Process the request with the "Login" token
		checktoken := watoken.DecodeGetId(os.Getenv(publickey), gettoken)
		authdata.Email = checktoken
		if checktoken == "" {
			response.Message = "Invalid token"
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
func GCFGetAllAbout(MONGOCONNSTRINGENV, dbname, collectionname string) string {
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
	dataabout := GetAllAbout(mconn, collectionname)
	if dataabout != nil {
		return GCFReturnStruct(CreateResponse(true, "Berhasil Get All About", dataabout))
	} else {
		return GCFReturnStruct(CreateResponse(false, "Gagal Get All About", dataabout))
	}
}

// <--- ini contact --->

// contact post
func GCFCreateContact(publickey, MONGOCONNSTRINGENV, dbname, colladmin, collcontact string, r *http.Request) string {
	var response Credential
	response.Status = false
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
	var authdata Admin

	gettoken := r.Header.Get("token")

	if gettoken == "" {
		response.Message = "Missing token in headers"
	} else {
		// Process the request with the "Login" token
		checktoken := watoken.DecodeGetId(os.Getenv(publickey), gettoken)
		authdata.Email = checktoken
		if checktoken == "" {
			response.Message = "Invalid token"
		} else {
			auth2 := FindAdmin(mconn, colladmin, authdata)
			if auth2.Role == "admin" {
				var datacontact Contact
				err := json.NewDecoder(r.Body).Decode(&datacontact)
				if err != nil {
					response.Message = "Error parsing application/json: " + err.Error()
				} else {
					CreateContact(mconn, collcontact, Contact{
						ID:       datacontact.ID,
						FullName: datacontact.FullName,
						Email:    datacontact.Email,
						Phone:    datacontact.Phone,
						Message:  datacontact.Message,
						Status:   datacontact.Status,
					})
					response.Status = true
					response.Message = "Berhasil Insert Contact"
				}
			} else {
				response.Message = "Anda tidak dapat Insert data karena bukan admin"
			}
		}
	}
	return GCFReturnStruct(response)
}

// delete contact
func GCFDeleteContact(publickey, MONGOCONNSTRINGENV, dbname, colladmin, collcontact string, r *http.Request) string {
	var response Credential
	response.Status = false
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
	var authdata Admin

	gettoken := r.Header.Get("token")

	if gettoken == "" {
		response.Message = "Missing token in headers"
	} else {
		// Process the request with the "Login" token
		checktoken := watoken.DecodeGetId(os.Getenv(publickey), gettoken)
		authdata.Email = checktoken
		if checktoken == "" {
			response.Message = "Invalid token"
		} else {
			auth2 := FindAdmin(mconn, colladmin, authdata)
			if auth2.Role == "admin" {
				var datacontact Contact
				err := json.NewDecoder(r.Body).Decode(&datacontact)
				if err != nil {
					response.Message = "Error parsing application/json: " + err.Error()
				} else {
					DeleteContact(mconn, colladmin, datacontact)
					response.Status = true
					response.Message = "Berhasil Delete Contact"
				}
			} else {
				response.Message = "Anda tidak dapat Delete data karena bukan admin"
			}
		}
	}
	return GCFReturnStruct(response)
}

// update contact
func GCFUpdateContact(publickey, MONGOCONNSTRINGENV, dbname, colladmin, collcontact string, r *http.Request) string {
	var respon Credential
	respon.Status = false
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
	var authdata Admin

	gettoken := r.Header.Get("token")
	if gettoken == "" {
		respon.Message = "Missing token in headers"
	} else {
		// Process the request with the "Login" token
		checktoken := watoken.DecodeGetId(os.Getenv(publickey), gettoken)
		authdata.Email = checktoken
		if checktoken == "" {
			respon.Message = "Invalid token"
		} else {
			auth2 := FindAdmin(mconn, colladmin, authdata)
			if auth2.Role == "admin" {
				var datacontact Contact
				err := json.NewDecoder(r.Body).Decode(&datacontact)
				if err != nil {
					respon.Message = "Error parsing application/json: " + err.Error()
				} else {
					UpdatedContact(mconn, colladmin, bson.M{"id": datacontact.ID}, datacontact)
					respon.Status = true
					respon.Message = "Berhasil Updated Contact"
					GCFReturnStruct(CreateResponse(true, "Success Update Product", datacontact))
				}
			} else {
				respon.Message = "Anda tidak dapat Update data karena bukan admin"
			}

		}
	}
	return GCFReturnStruct(respon)
}

// get all contact
func GCFGetAllContact(MONGOCONNSTRINGENV, dbname, collectionname string) string {
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
	datacontact := GetAllContact(mconn, collectionname)
	if datacontact != nil {
		return GCFReturnStruct(CreateResponse(true, "success Get All Contact", datacontact))
	} else {
		return GCFReturnStruct(CreateResponse(false, "Failed Get All Contact", datacontact))
	}
}

// get all contact by id
func GCFGetAllContactID(MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
	var datacontact Contact
	err := json.NewDecoder(r.Body).Decode(&datacontact)
	if err != nil {
		return err.Error()
	}

	contact := GetIdContact(mconn, collectionname, datacontact)
	if contact != (Contact{}) {
		return GCFReturnStruct(CreateResponse(true, "Success: Get ID Contact", datacontact))
	} else {
		return GCFReturnStruct(CreateResponse(false, "Failed to Get ID Contact", datacontact))
	}
}
