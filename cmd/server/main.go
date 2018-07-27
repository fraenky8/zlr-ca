package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

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

	read := r.Group("/icecreams")
	{
		read.GET("/", readIcecream)
		read.GET("/:ids", readIcecream)
		read.GET("/:ids/", readIcecream)
	}

	// r.POST("/icecream", createIcecream)

	// listen and serve on 0.0.0.0:8080
	r.Run()
}

func readIcecream(c *gin.Context) {
	rawIds := c.Param("ids")
	rawIds = strings.TrimSpace(rawIds)

	if rawIds == "" {
		c.JSON(http.StatusBadRequest, api.FailString("no id(s) provided"))
		return
	}

	var ids []int64
	for _, id := range strings.Split(rawIds, ",") {
		tid := strings.TrimSpace(id)
		if tid != "" {
			id, err := strconv.ParseInt(tid, 10, 64)
			if err != nil {
				// ignore faulty ids - but give message
				log.Printf("faulty id: %v", err)
				// errors = append(errors, fmt.Errorf("faulty id: %v", err))
				continue
			}
			ids = append(ids, id)
		}
	}

	if len(ids) == 0 {
		c.JSON(http.StatusBadRequest, api.FailString("no valid id(s) provided"))
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
