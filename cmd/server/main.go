package main

import (
	"fmt"
	"log"

	"github.com/fraenky8/zlr-ca/pkg/api"
	"github.com/fraenky8/zlr-ca/pkg/storage"
	"github.com/fraenky8/zlr-ca/pkg/storage/repos"
)

// go run main.go -h 192.168.99.100 -pt 5432 -u postgres -p mysecretpassword -d postgres -s zlr_ca
func main() {

	db, err := storage.NewPostgres(
		storage.NewConfigByCmdArgs(),
	)

	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	repository, err := repos.NewRepository(db)
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
