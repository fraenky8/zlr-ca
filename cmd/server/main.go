package main

import (
	"fmt"
	"log"

	"github.com/fraenky8/zlr-ca/pkg/api"
	"github.com/fraenky8/zlr-ca/pkg/storage"
	"github.com/fraenky8/zlr-ca/pkg/storage/repos"
)

func main() {

	db, err := storage.NewPostgres(&storage.Config{
		Host:     "192.168.99.100",
		Port:     "5432",
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

	repository, err := repos.NewService(db)
	if err != nil {
		fmt.Println(err)
		return
	}

	s, err := api.NewServer(
		&api.ServerConfig{},
		repository,
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	log.Fatal(s.Run())
}
