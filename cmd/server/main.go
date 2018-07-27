package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/fraenky8/zlr-ca/pkg/core/api"
	"github.com/fraenky8/zlr-ca/pkg/infrastructure/storage"
	"github.com/fraenky8/zlr-ca/pkg/infrastructure/storage/repos"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

var (
	db  *storage.Database
	err error
)

func main() {

	db, err = storage.Connect(&storage.Config{
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

	// gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(gzip.Gzip(gzip.DefaultCompression))

	icecreams := r.Group("/icecreams")
	{
		icecreams.GET("/", readIcecream)
		icecreams.GET("/:ids", readIcecream)
		icecreams.GET("/:ids/", readIcecream)
		icecreams.GET("/:ids/ingredients", readIcecreamIngredients)
		icecreams.GET("/:ids/sourcingvalues", readIcecreamSourcingValues)
	}

	ingredients := r.Group("/ingredients")
	{
		ingredients.GET("/", readIngredients)
	}

	sourcingvalues := r.Group("/sourcingvalues")
	{
		sourcingvalues.GET("/", readSourcingValues)
	}

	// r.POST("/icecream", createIcecream)

	// listen and serve on 0.0.0.0:8080
	r.Run()
}

func readIcecream(c *gin.Context) {

	ids, err := api.ConvertIdsParam(c.Param("ids"))
	if err != nil {
		c.JSON(http.StatusBadRequest, api.Fail(err))
		return
	}

	icecreams, err := repos.NewIcecreamRepo(db).Read(ids)
	if err != nil {
		log.Printf("could not get icecream: %v", err)
		c.JSON(http.StatusInternalServerError, api.Error("a database error occured, please try again later"))
		return
	}

	if len(icecreams) == 1 {
		c.JSON(http.StatusOK, api.Success(
			&api.IcecreamResponse{Icecream: icecreams[0]}),
		)
		return
	}

	c.JSON(http.StatusOK, api.Success(
		&api.IcecreamsResponse{Icecreams: icecreams}),
	)
}

func readIcecreamIngredients(c *gin.Context) {

	ids, err := api.ConvertIdsParam(c.Param("ids"))
	if err != nil {
		c.JSON(http.StatusBadRequest, api.Fail(err))
		return
	}

	ingredients, err := repos.NewIngredientsRepo(db).Reads(ids)
	if err != nil {
		log.Printf("could not get ingredients: %v", err)
		c.JSON(http.StatusInternalServerError, api.Error("a database error occured, please try again later"))
	}

	if len(ingredients) == 1 {
		c.JSON(http.StatusOK, api.Success(
			&api.IngredientResponse{Ingredient: ingredients[0]},
		))
		return
	}

	c.JSON(http.StatusOK, api.Success(
		&api.IngredientsResponse{Ingredients: ingredients}),
	)
}

func readIngredients(c *gin.Context) {
	ingredients, err := repos.NewIngredientsRepo(db).ReadAll()
	if err != nil {
		log.Printf("could not get ingredients: %v", err)
		c.JSON(http.StatusInternalServerError, api.Error("a database error occured, please try again later"))
	}

	c.JSON(http.StatusOK, api.Success(
		&api.IngredientResponse{Ingredient: ingredients}),
	)
}

func readIcecreamSourcingValues(c *gin.Context) {

	ids, err := api.ConvertIdsParam(c.Param("ids"))
	if err != nil {
		c.JSON(http.StatusBadRequest, api.Fail(err))
		return
	}

	sourcingValues, err := repos.NewSourcingValuesRepo(db).Reads(ids)
	if err != nil {
		log.Printf("could not get sourcing values: %v", err)
		c.JSON(http.StatusInternalServerError, api.Error("a database error occured, please try again later"))
	}

	if len(sourcingValues) == 1 {
		c.JSON(http.StatusOK, api.Success(
			&api.SourcingValueResponse{SourcingValue: sourcingValues[0]},
		))
		return
	}

	c.JSON(http.StatusOK, api.Success(
		&api.SourcingValuesResponse{SourcingValues: sourcingValues}),
	)
}

func readSourcingValues(c *gin.Context) {

	sourcingValues, err := repos.NewSourcingValuesRepo(db).ReadAll()
	if err != nil {
		log.Printf("could not get sourcing values: %v", err)
		c.JSON(http.StatusInternalServerError, api.Error("a database error occured, please try again later"))
	}

	c.JSON(http.StatusOK, api.Success(
		&api.SourcingValueResponse{SourcingValue: sourcingValues}),
	)
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
