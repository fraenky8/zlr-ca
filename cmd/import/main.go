package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/fraenky8/zlr-ca/pkg/core/domain"
	"github.com/fraenky8/zlr-ca/pkg/infrastructure/storage"
	"github.com/fraenky8/zlr-ca/pkg/infrastructure/storage/repos"
)

func main() {

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

	err = json.Unmarshal(b, &icecreams)

	if err != nil {
		fmt.Println(err)
		return
	}

	for _, icecream := range icecreams {
		fmt.Println("\t" + icecream.Name)
		_, err := repos.NewIcecreamRepo(db).Create(icecream)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	fmt.Println("done")
}
