package main

import (
	"fmt"
	"log"

	"github.com/fraenky8/zlr-ca/pkg/api"
	"github.com/fraenky8/zlr-ca/pkg/storage"
	"github.com/fraenky8/zlr-ca/pkg/storage/repos"
)

func main() {

	db, err := storage.Connect(&storage.Config{
		Host:     "192.168.99.100",
		Username: "postgres",
		Password: "mysecretpassword",
		Database: "postgres",
		Schema:   "zlr_ca",
	})

	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	s, err := api.NewServer(
		&api.ServerConfig{},
		&api.Storage{
			Db:                   db,
			IcecreamService:      repos.NewIcecreamRepo(db),
			IngredientService:    repos.NewIngredientsRepo(db),
			SourcingValueService: repos.NewSourcingValuesRepo(db),
		},
	)

	log.Fatal(s.Run())
}
