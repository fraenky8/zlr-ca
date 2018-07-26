package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/fraenky8/zlr-ca/pkg/core/domain"
	"github.com/fraenky8/zlr-ca/pkg/infrastructure/storage"
	"github.com/fraenky8/zlr-ca/pkg/infrastructure/storage/repos"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

type response struct {
	Message   string             `json:"message"`
	Icecreams []*domain.Icecream `json:"icecreams"`
}

var (
	db  *storage.Database
	err error
)

// TODO dont expose errors to users

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

	read := r.Group("/icecream")
	{
		read.GET("/", readIcecream)
		read.GET("/:id", readIcecream)
		read.GET("/:id/", readIcecream)
	}

	// reads := r.Group("/icecreams")
	// {
	// 	reads.GET("/", readIcecreams)
	// 	reads.GET("/:ids", readIcecreams)
	// 	reads.GET("/:ids/", readIcecreams)
	// }
	//
	// r.POST("/icecream", createIcecream)

	// listen and serve on 0.0.0.0:8080
	r.Run()
}

func readIcecream(c *gin.Context) {
	sid := c.Param("id")

	if sid == "" {
		c.JSON(http.StatusBadRequest, &response{
			Message:   "no id provided",
			Icecreams: []*domain.Icecream{},
		})
		return
	}

	id, err := strconv.ParseInt(sid, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, &response{
			Message:   fmt.Sprintf("faulty id: %v", err),
			Icecreams: []*domain.Icecream{},
		})
		return
	}

	icecream, err := repos.NewIcecreamRepo(db).Read(id)
	if err != nil {
		log.Printf("could not get icecream: %v", err)
		c.JSON(http.StatusInternalServerError, &response{
			Message:   "a database error occured, please try again later",
			Icecreams: []*domain.Icecream{},
		})
		return
	}

	icecreams := []*domain.Icecream{}
	if icecream != nil {
		icecreams = append(icecreams, icecream)
	}

	c.JSON(http.StatusOK, &response{
		Message:   "",
		Icecreams: icecreams,
	})
}

// func readIcecreams(c *gin.Context) {
// 	rawIds := c.Param("ids")
// 	rawIds = strings.TrimSpace(rawIds)
//
// 	if rawIds == "" {
// 		c.JSON(200, []string{})
// 		return
// 	}
//
// 	var ids []string
// 	for _, id := range strings.Split(rawIds, ",") {
// 		tid := strings.TrimSpace(id)
// 		if tid != "" {
// 			ids = append(ids, tid)
// 		}
// 	}
//
// 	if len(ids) == 0 {
// 		c.JSON(200, []string{})
// 		return
// 	}
//
// 	var ices []icecream
// 	for _, icecream := range icecreams {
// 		for _, id := range ids {
// 			if icecream.ProductId == id {
// 				ices = append(ices, icecream)
// 			}
// 		}
// 	}
//
// 	if len(ices) == 0 {
// 		c.JSON(200, []string{})
// 		return
// 	}
//
// 	c.JSON(200, ices)
// }
//
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
