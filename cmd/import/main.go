package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/fraenky8/zlr-ca/pkg/domain"
	"github.com/fraenky8/zlr-ca/pkg/storage"
	"github.com/fraenky8/zlr-ca/pkg/storage/repos"
)

func main() {
	start := time.Now()
	fmt.Println("starting import of icecream.json")

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

	b, err := ioutil.ReadFile("icecream.json")
	if err != nil {
		fmt.Println(err)
		return
	}

	var icecreams []domain.Icecream

	if err = json.Unmarshal(b, &icecreams); err != nil {
		fmt.Println(err)
		return
	}

	if _, err = repos.NewIcecreamRepo(db).Creates(icecreams); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("done (%s)", time.Since(start))
}
