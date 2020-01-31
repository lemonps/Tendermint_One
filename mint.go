package main

import (
	"github.com/mint/jsonstore2"
	"os"
    "fmt"
	"github.com/tendermint/tendermint/abci/server"
	"github.com/tendermint/tendermint/abci/types"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"

	cmn "github.com/tendermint/tendermint/libs/os"
	"github.com/tendermint/tendermint/libs/log"
)

func main() {
	initJSONStore()
	//createNewUser(3,"FC35EC86ADDD2E90DC158C4BD0FCA8E66617BB68B6AC29E0051743816E40263A",0)
}

func initJSONStore() error {
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout))

	// Create the application
	var app types.Application

	db, err := sql.Open("mysql","root:tendermint@tcp(45.77.170.153:3306)/tendermint")
	if err != nil {
		fmt.Println("Connection fail")
	}else{
		fmt.Println("successful connect")
	}
	//db := session.DB("tendermintdb")

	// Clean the DB on each reboot
	//collections := [5]string{"posts", "comments", "users", "userpostvotes", "usercommentvotes"}

//	for _, collection := range collections {
//		db.C(collection).RemoveAll(nil)
//	}

//	createNewUser(3,"FC35EC86ADDD2E90DC158C4BD0FCA8E66617BB68B6AC29E0051743816E40263A",0, db)

	app = jsonstore.NewJSONStoreApplication(db)

	// Start the listener
	srv, err := server.NewServer("tcp://127.0.0.1:46658", "socket", app)
	if err != nil {
		return err
	}
	srv.SetLogger(logger.With("module", "abci-server"))
	if err := srv.Start(); err != nil {
		return err
	}

	// Wait forever
	cmn.TrapSignal(logger, func(){
		// Cleanup
		srv.Stop()
	})
	select {}
}

func createNewUser(id int, public_key string, role int, db *sql.DB) {
	stmt, err := db.Prepare("INSERT INTO user(id, public_key, role) VALUES(?,?,?)")

	if err != nil {
	   panic(err.Error)
	}
			
	stmt.Exec(id, public_key, role)
}
