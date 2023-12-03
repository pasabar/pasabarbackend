package pasabarbackend

import (
	"encoding/json"

	"fmt"

	"aidanwoods.dev/go-paseto"
)

// <--- paseto encoder --->
// pengecekantoken
func Decoder(publickey, tokenstr string) (payload Payload, err error) {
	var token *paseto.Token
	var pubKey paseto.V4AsymmetricPublicKey
	pubKey, err = paseto.NewV4AsymmetricPublicKeyFromHex(publickey) // this wil fail if given key in an invalid format
	if err != nil {
		fmt.Println("Decode NewV4AsymmetricPublicKeyFromHex : ", err)
	}
	parser := paseto.NewParser()                             // only used because this example token has expired, use NewParser() (which checks expiry by default)
	token, err = parser.ParseV4Public(pubKey, tokenstr, nil) // this will fail if parsing failes, cryptographic checks fail, or validation rules fail
	if err != nil {
		fmt.Println("Decode ParseV4Public : ", err)
	} else {
		json.Unmarshal(token.ClaimsJSON(), &payload)
	}
	return payload, err
}

func DecodeGetCatalog(PublicKey, tokenStr string) (pay string, err error) {
	key, err := Decoder(PublicKey, tokenStr)
	if err != nil {
		fmt.Println("Cannot decode the token", err.Error())
	}
	return key.Catalog, nil
}
