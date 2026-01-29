package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/caproven/termdict/cmd"
	"github.com/caproven/termdict/dictionary"
	"github.com/caproven/termdict/storage/sqlite"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// TODO give proper location for this
	db, err := sql.Open("sqlite3", "data.sqlite")
	if err != nil {
		fmt.Println("Failed to open database")
		os.Exit(1)
	}
	defer db.Close()
	store, err := sqlite.NewStore(context.Background(), db)
	if err != nil {
		fmt.Println("Failed to instantiate cache")
		os.Exit(1)
	}
	api := dictionary.NewDefaultWebAPI()
	dict := dictionary.NewCachedDefiner(store, api)

	cfg := &cmd.Config{
		Out:   os.Stdout,
		Vocab: store,
		Dict:  dict,
	}
	if err := cmd.NewRootCmd(cfg).Execute(); err != nil {
		os.Exit(1)
	}
}
