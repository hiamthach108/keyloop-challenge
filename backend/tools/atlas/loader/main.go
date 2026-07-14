package main

import (
	"fmt"
	"io"
	"os"

	"ariga.io/atlas-provider-gorm/gormschema"
	"github.com/hiamthach108/keyloop-challenge/backend/internal/model"
)

func main() {
	statements, err := gormschema.New("postgres").Load(
		&model.Dealership{},
		&model.Vehicle{},
		&model.InventoryStock{},
		&model.StockMovement{},
		&model.InventoryAction{},
	)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "load GORM schema: %v\n", err)
		os.Exit(1)
	}
	_, _ = io.WriteString(os.Stdout, statements)
}
