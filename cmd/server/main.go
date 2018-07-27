package main

import (
	"fmt"

	"github.com/fraenky8/zlr-ca/pkg/core/api"
	"github.com/fraenky8/zlr-ca/pkg/infrastructure/storage"
	"github.com/fraenky8/zlr-ca/pkg/infrastructure/storage/repos"
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
			Icecreamer:     repos.NewIcecreamRepo(db),
			Ingredienter:   repos.NewIngredientsRepo(db),
			SourcingValuer: repos.NewSourcingValuesRepo(db),
		},
	)

	// r.POST("/icecream", createIcecream)

	s.Run()
}

// func createIcecream(c *gin.Context) {
// 	var ice icecream
//
// 	err := c.BindJSON(&ice)
// 	if err != nil {
// 		if err.Error() == "EOF" {
// 			c.JSON(400, gin.H{})
// 			return
// 		}
// 		log.Printf("could not bind: %v", err)
// 		c.Writer.WriteString("could not bind: " + err.Error())
// 		return
// 	}
//
// 	ice.ProductId = fmt.Sprintf("%v", len(icecreams)+1)
// 	log.Printf("new icecream: %+v", ice)
// 	icecreams = append(icecreams, ice)
//
// 	c.Writer.WriteString("new icecream created: " + ice.ProductId)
// }
