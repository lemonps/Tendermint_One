package main
import (
        "fmt"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"

)
func main() {
	db, err := sql.Open("mysql","root:tendermint@tcp(52.36.211.159)/tendermint")
	if err != nil {
		fmt.Println("Connection fail")
	}else{
		fmt.Println("successful connect")
	}
    var usernum int
        stmtOut, err := db.Prepare("SELECT count(*) FROM user WHERE public_key = ?")
	if err != nil {
        panic(err.Error()) 
	}
	defer stmtOut.Close()

	err = stmtOut.QueryRow("public_key1").Scan(&usernum) // WHERE number = 13
	if err != nil {
	panic(err.Error())
	}
	fmt.Println(usernum)
	
    }
