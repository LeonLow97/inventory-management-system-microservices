package main

import (
	"os"
	"testing"

	"github.com/LeonLow97/inventory-management-system-golang-react-postgresql/database/repository/dbrepo"
)

var app application

func TestMain(m *testing.M) {

	app.DB = &dbrepo.TestDBRepo{}

	os.Exit(m.Run())
}
