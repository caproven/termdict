package main

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/caproven/termdict/cmd"
	"github.com/caproven/termdict/dictionary"
	"github.com/caproven/termdict/storage/sqlite"
	_ "github.com/mattn/go-sqlite3"
)

const (
	appName    = "termdict"
	dbFilename = "termdict.sqlite"
)

func main() {
	dataDir := filepath.Join(xdg.DataHome, appName)
	if err := os.MkdirAll(dataDir, os.ModePerm); err != nil {
		fmt.Printf("Could not create data directory: %v\n", err)
		os.Exit(1)
	}

	db, err := sql.Open("sqlite3", filepath.Join(dataDir, dbFilename))
	if err != nil {
		fmt.Println("Failed to open database")
		os.Exit(1)
	}
	defer func(db *sql.DB) {
		if err := db.Close(); err != nil {
			slog.Warn("Failed to close database", "error", err)
		}
	}(db)
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
