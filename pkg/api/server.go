package api

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/fraenky8/zlr-ca/pkg/domain"
	"github.com/fraenky8/zlr-ca/pkg/storage"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

const DefaultPort = "8080"

type ServerConfig struct {
	Port string
}

func (s *ServerConfig) verify() error {
	if s.Port == "" {
		s.Port = DefaultPort
	}
	return nil
}

type Storage struct {
	Db                   *storage.Database
	IcecreamService      domain.IcecreamService
	IngredientService    domain.IngredientService
	SourcingValueService domain.SourcingValueService
}

func (s *Storage) verify() error {
	if s.IcecreamService == nil {
		return fmt.Errorf("no icecreamer given")
	}
	if s.IngredientService == nil {
		return fmt.Errorf("no ingredienter given")
	}
	if s.SourcingValueService == nil {
		return fmt.Errorf("no sourcingValuer given")
	}
	return nil
}

type Server struct {
	config  *ServerConfig
	storage *Storage
	engine  *gin.Engine
}

func NewServer(config *ServerConfig, storage *Storage) (*Server, error) {

	if err := config.verify(); err != nil {
		return nil, fmt.Errorf("could not create server with config: %v", err)
	}

	if err := storage.verify(); err != nil {
		return nil, fmt.Errorf("could not create server with storage: %v", err)
	}

	engine := gin.Default()
	engine.Use(gzip.Gzip(gzip.DefaultCompression))

	s := &Server{
		config:  config,
		storage: storage,
		engine:  engine,
	}

	return s.setupRoutes(), nil
}

func (s *Server) Run() error {
	// gin.SetMode(gin.ReleaseMode)
	return s.engine.Run(":" + s.config.Port)
}

func (s *Server) setupRoutes() *Server {
	icecreams := s.engine.Group("/icecreams")
	{
		// CREATE
		icecreams.POST("/", s.createIcecreams)
		icecreams.PUT("/", s.createIcecreams)

		// READ
		icecreams.GET("/", s.readIcecream)
		icecreams.GET("/:ids", s.readIcecream)
		icecreams.GET("/:ids/", s.readIcecream)
		icecreams.GET("/:ids/ingredients", s.readIcecreamIngredients)
		icecreams.GET("/:ids/sourcingvalues", s.readIcecreamSourcingValues)

		// UPDATE
		icecreams.PUT("/:ids", s.updateIcecreams)
		icecreams.PUT("/:ids/", s.updateIcecreams)
		icecreams.PATCH("/", func(c *gin.Context) {
			c.JSON(http.StatusMethodNotAllowed, FailStringResponse("updating the entire collections is not allowed"))
		})
		icecreams.PATCH("/:ids", s.updateIcecreams)
		icecreams.PATCH("/:ids/", s.updateIcecreams)

		// DELETE
		icecreams.DELETE("/", func(c *gin.Context) {
			c.JSON(http.StatusMethodNotAllowed, FailStringResponse("deleting the entire collections is not allowed"))
		})
		icecreams.DELETE("/:ids", s.deleteIcecreams)
		icecreams.DELETE("/:ids/", s.deleteIcecreams)
	}

	ingredients := s.engine.Group("/ingredients")
	{
		ingredients.GET("/", s.readIngredients)
	}

	sourcingvalues := s.engine.Group("/sourcingvalues")
	{
		sourcingvalues.GET("/", s.readSourcingValues)
	}

	return s
}

func (s *Server) readIcecream(c *gin.Context) {

	ids, err := convertIdsParam(c.Param("ids"))
	if err != nil {
		c.JSON(http.StatusBadRequest, FailResponse(err))
		return
	}

	icecreams, err := s.storage.IcecreamService.Read(ids)
	if err != nil {
		log.Printf("could not get icecream: %v", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse("a database error occured, please try again later"))
		return
	}

	if len(icecreams) == 1 {
		c.JSON(http.StatusOK, SuccessResponse(
			&IcecreamResponse{Icecream: icecreams[0]}),
		)
		return
	}

	c.JSON(http.StatusOK, SuccessResponse(
		&IcecreamsResponse{Icecreams: icecreams}),
	)
}

func (s *Server) readIcecreamIngredients(c *gin.Context) {

	ids, err := convertIdsParam(c.Param("ids"))
	if err != nil {
		c.JSON(http.StatusBadRequest, FailResponse(err))
		return
	}

	ingredients, err := s.storage.IngredientService.Reads(ids)
	if err != nil {
		log.Printf("could not get ingredients: %v", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse("a database error occured, please try again later"))
	}

	if len(ingredients) == 1 {
		c.JSON(http.StatusOK, SuccessResponse(
			&IngredientResponse{Ingredient: ingredients[0]},
		))
		return
	}

	c.JSON(http.StatusOK, SuccessResponse(
		&IngredientsResponse{Ingredients: ingredients}),
	)
}

func (s *Server) readIngredients(c *gin.Context) {
	ingredients, err := s.storage.IngredientService.ReadAll()
	if err != nil {
		log.Printf("could not get ingredients: %v", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse("a database error occured, please try again later"))
	}

	c.JSON(http.StatusOK, SuccessResponse(
		&IngredientResponse{Ingredient: ingredients}),
	)
}

func (s *Server) readIcecreamSourcingValues(c *gin.Context) {

	ids, err := convertIdsParam(c.Param("ids"))
	if err != nil {
		c.JSON(http.StatusBadRequest, FailResponse(err))
		return
	}

	sourcingValues, err := s.storage.SourcingValueService.Reads(ids)
	if err != nil {
		log.Printf("could not get sourcing values: %v", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse("a database error occured, please try again later"))
	}

	if len(sourcingValues) == 1 {
		c.JSON(http.StatusOK, SuccessResponse(
			&SourcingValueResponse{SourcingValue: sourcingValues[0]},
		))
		return
	}

	c.JSON(http.StatusOK, SuccessResponse(
		&SourcingValuesResponse{SourcingValues: sourcingValues}),
	)
}

func (s *Server) readSourcingValues(c *gin.Context) {

	sourcingValues, err := s.storage.SourcingValueService.ReadAll()
	if err != nil {
		log.Printf("could not get sourcing values: %v", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse("a database error occured, please try again later"))
	}

	c.JSON(http.StatusOK, SuccessResponse(
		&SourcingValueResponse{SourcingValue: sourcingValues}),
	)
}

func (s *Server) createIcecreams(c *gin.Context) {

	if !strings.Contains(c.ContentType(), "application/json") {
		c.JSON(http.StatusBadRequest, FailStringResponse("only Content-Type: application/json is supported"))
		return
	}

	var icecreams []*domain.Icecream
	if err := c.ShouldBind(&icecreams); err != nil {
		if err.Error() == "EOF" {
			c.JSON(http.StatusBadRequest, FailStringResponse("no icecream data provided"))
			return
		}
		c.JSON(http.StatusBadRequest, FailResponse(fmt.Errorf("faulty data provided: %v (forgot to wrap in []?)", err)))
		return
	}

	// if one icecream fails, all icecreams fail - "all or nothing"
	for k, icecream := range icecreams {
		if err := icecream.Verify(); err != nil {
			c.JSON(http.StatusBadRequest, FailResponse(fmt.Errorf("icecream #%d: %v", k, err)))
			return
		}

		productId, err := strconv.Atoi(icecream.ProductID)
		if err != nil {
			log.Printf("faulty id: %v", err)
			c.JSON(http.StatusBadRequest, FailResponse(fmt.Errorf("icecream #%d: faulty productId provided: %s", k, icecream.ProductID)))
			return
		}

		if existingIcecream, _ := s.storage.IcecreamService.Read([]int{productId}); existingIcecream != nil {
			c.JSON(http.StatusBadRequest, FailStringResponse("icecream with productId = "+icecream.ProductID+" already exists"))
			return
		}
	}

	for _, icecream := range icecreams {
		if _, err := s.storage.IcecreamService.Create(*icecream); err != nil {
			log.Printf("could not create icecream: %v", err)
			c.JSON(http.StatusInternalServerError, ErrorResponse("a database error occured, please try again later"))
			return
		}
	}

	c.JSON(http.StatusCreated, SuccessResponse(
		&IcecreamsResponse{Icecreams: icecreams},
	))
}

func (s *Server) updateIcecreams(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, ErrorResponse("not implemented yet"))
}

func (s *Server) deleteIcecreams(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, ErrorResponse("not implemented yet"))
}

func convertIdsParam(sids string) (ids []int, err error) {

	sids = strings.TrimSpace(sids)
	if sids == "" {
		return []int{}, fmt.Errorf("no id(s) provided")
	}

	for _, id := range strings.Split(sids, ",") {
		tid := strings.TrimSpace(id)

		// ignore empty ids
		if tid == "" {
			continue
		}

		id, parseErr := strconv.Atoi(tid)
		if parseErr != nil {
			log.Printf("faulty id: %v", parseErr)
			err = parseErr
			continue
		}

		ids = append(ids, id)
	}

	if len(ids) == 0 {
		return []int{}, fmt.Errorf("no valid id(s) provided")
	}

	if err != nil {
		return []int{}, fmt.Errorf("at least one invalid id detected: %v", err)
	}

	return ids, nil
}
