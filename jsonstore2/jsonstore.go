package jsonstore

import (
//	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
//	"math"
//	"net/url"
//	"regexp"
	"strconv"
//	"strings"
	"time"
	"github.com/mint/code"
	"github.com/tendermint/tendermint/abci/types"
	"golang.org/x/crypto/ed25519"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
//	"gopkg.in/mgo.v2/bson"
)
var (
	ErrBadAddress     = &appError{"Ticket must have an address"}
	ErrBadNonce       = &appError{"Ticket nonce must increase on resale"}
	ErrBadSignature   = &appError{"Resale must be signed by the previous owner"}
	ErrTicketNotFound = &appError{"Ticket could not be found"}
)

var _ types.Application = (*JSONStoreApplication)(nil)
var db *sql.DB
var count int64 = 0

type appError struct{ msg string }

func (err appError) Error() string { return err.msg }

// User ...
type User struct {
	ID        int
	PublicKey string    
	Role      int
}

// JSONStoreApplication ...
type JSONStoreApplication struct {
	types.BaseApplication
}

func byteToHex(input []byte) string {
	var hexValue string
	for _, v := range input {
		hexValue += fmt.Sprintf("%02x", v)
	}
	return hexValue
}

func checkUserPublic(db *sql.DB,pub string) int64 {
    var usernum int
    stmtOut, err := db.Prepare("SELECT count(*) FROM user WHERE public_key = ?")
	if err != nil {
        panic(err.Error()) 
	}
	defer stmtOut.Close()

	err = stmtOut.QueryRow(pub).Scan(&usernum) // WHERE number = 13
	if err != nil {
		panic(err.Error())
	}
	if usernum == 1 {
        return 1
	}else{
		return 0
	}
}

func signatureValidate(pub string,sig string,msg string) bool {
  	pub_b, _ := hex.DecodeString(pub)
 	sig_b, _ := hex.DecodeString(sig)
 	buffer := []byte(msg)
 	if ed25519.Verify(pub_b, buffer, sig_b) {
 		return true
 	}
        return false
}

// FindTimeFromObjectID ... Convert ObjectID string to Time
func FindTimeFromObjectID(id string) time.Time {
ts, _ := strconv.ParseInt(id[0:8], 16, 64)
	return time.Unix(ts, 0)}

// NewJSONStoreApplication ...
func NewJSONStoreApplication(dbCopy *sql.DB) *JSONStoreApplication {
	db = dbCopy
    return &JSONStoreApplication{}
}

// Info ...
func (app *JSONStoreApplication) Info(req types.RequestInfo) (resInfo types.ResponseInfo) {
	return types.ResponseInfo{Data: fmt.Sprintf("{\"size\":%v}", 0)}
}

// DeliverTx ... Update the global state
func (app *JSONStoreApplication) DeliverTx(tx types.RequestDeliverTx) types.ResponseDeliverTx {
	 return types.ResponseDeliverTx{Code: code.CodeTypeOK}

	 var temp interface{}
	 err := json.Unmarshal(tx.Tx, &temp)

	 if err != nil {
		 return types.ResponseDeliverTx{Code: code.CodeTypeEncodingError,Log: fmt.Sprint(err)}
	 }

	 message := temp.(map[string]interface{})

	 PublicKey := message["publicKey"].(string)

	 count := checkUserPublic(db,PublicKey)
        
	 if count != 0 {
                //var temp2 interface{}
		//userInfo := message["userInfo"].(map[string]interface{})
		// err2 := json.Unmarshal([]byte(message["userInfo"].(string)), &temp2)
                // message2 := temp2.(map[string]interface{})
		//if err2 != nil {
		//	panic(err.Error)
		//}
             
		var user User
		user.ID = message["id"].(int)
		user.PublicKey = message["public_key"].(string)
		user.Role = message["role"].(int)

		fmt.Printf(user.PublicKey)
               
		// log.PrintIn("id: ", user.ID, "public_key: ", user.PublicKey, "role: ", user.Role)

		stmt, err := db.Prepare("INSERT INTO user(id, public_key, role) VALUES(?,?,?)")

		if err != nil {
			panic(err.Error)
		}
			
		stmt.Exec(user.ID, user.PublicKey, user.Role)

		// log.PrintIn("insert result: ", res.LastInsertId())

		return types.ResponseDeliverTx{Code: code.CodeTypeOK}
	 } else {
		return types.ResponseDeliverTx{Code: code.CodeTypeBadData}
	 }
	 
	//  var types interface{}
	//  errType := json.Unmarshall(message["types"].(string), &types)
	 
	//  if errType != nil {
	// 	 panic(err.Error)
	//  }

	// switch types["types"] {
	// 	case "createUser":
	// 		entity := types["entity"].(map[string]interface{})

	// 		var user User
	// 		user.ID = entity["id"].(int)
	// 		user.PublicKey = entity["publicKey"].(string)
	// 		user.Role = entity["role"].(int)
	// }
}

// CheckTx ... Verify the transaction
func (app *JSONStoreApplication) CheckTx(tx types.RequestCheckTx) types.ResponseCheckTx {

	var temp interface{}
	err := json.Unmarshal(tx.Tx, &temp)

	if err != nil {
		return types.ResponseCheckTx{Code: code.CodeTypeEncodingError,Log: fmt.Sprint(err)}
	}

	message := temp.(map[string]interface{})
	PublicKey:= message["publicKey"].(string)
	ByteString:= message["msg"].(string)
	SignaTure:= message["sig"].(string)   
	count := checkUserPublic(db,PublicKey)

    if signatureValidate(PublicKey,SignaTure,ByteString) && count !=0 {
    return types.ResponseCheckTx{Code: code.CodeTypeOK}
}
    return types.ResponseCheckTx{Code: code.CodeTypeBadData}
}

// Commit ...Commit the block. Calculate the appHash
func (app *JSONStoreApplication) Commit() types.ResponseCommit {
	appHash := make([]byte, 8)
	count++
	binary.PutVarint(appHash, count)
	return types.ResponseCommit{Data: appHash}
}

// Query ... Query the blockchain. Unimplemented as of now.
func (app *JSONStoreApplication) Query(reqQuery types.RequestQuery) (resQuery types.ResponseQuery) {
	return types.ResponseQuery{Code: code.CodeTypeOK}
}
