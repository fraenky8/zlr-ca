package api

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/fraenky8/zlr-ca/pkg/domain"
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

func (s *Server) Run() {
	// gin.SetMode(gin.ReleaseMode)
	log.Fatal(s.engine.Run(":" + s.config.Port))
}

func (s *Server) setupRoutes() *Server {
	icecreams := s.engine.Group("/icecreams")
	{
		icecreams.GET("/", s.readIcecream)
		icecreams.GET("/:ids", s.readIcecream)
		icecreams.GET("/:ids/", s.readIcecream)
		icecreams.GET("/:ids/ingredients", s.readIcecreamIngredients)
		icecreams.GET("/:ids/sourcingvalues", s.readIcecreamSourcingValues)
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
