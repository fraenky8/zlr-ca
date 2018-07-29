package api

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/fraenky8/zlr-ca/pkg/domain"
	"github.com/fraenky8/zlr-ca/pkg/storage/repos"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

const (
	DefaultPort = "8080"

	RequestIcecreamKey = "icecreams"
)

type ServerConfig struct {
	Port string
	Mode string
}

func (s *ServerConfig) Verify() error {
	if s.Port == "" {
		s.Port = DefaultPort
	}
	if s.Mode == "" {
		s.Mode = gin.DebugMode
	}
	return nil
}

type Server struct {
	config *ServerConfig
	repo   *repos.Repository
	engine *gin.Engine
}

func NewServer(config *ServerConfig, repo *repos.Repository) (*Server, error) {

	if err := config.Verify(); err != nil {
		return nil, fmt.Errorf("could not create server with config: %v", err)
	}

	if err := repo.Verify(); err != nil {
		return nil, fmt.Errorf("could not create server with repo: %v", err)
	}

	gin.SetMode(config.Mode)
	engine := gin.Default()
	engine.Use(gzip.Gzip(gzip.DefaultCompression))

	s := &Server{
		config: config,
		repo:   repo,
		engine: engine,
	}

	return s.setupRoutes(), nil
}

func (s *Server) Run() error {
	return s.engine.Run(":" + s.config.Port)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	s.engine.ServeHTTP(w, req)
}

func (s *Server) setupRoutes() *Server {

	// empty "" routes here  to avoid a 307 redirect response
	// so it only occures when query has trailing slash

	icecreams := s.engine.Group("/icecreams")
	{
		create := icecreams.Group("").Use(s.icecreamRequest)
		{
			create.POST("", s.createIcecreams)
			create.PUT("", s.createIcecreams)
		}

		read := icecreams.Group("")
		{
			read.GET("", s.readIcecreams)
			read.GET("/:ids", s.readIcecreams)
			read.GET("/:ids/", s.readIcecreams)
			read.GET("/:ids/ingredients", s.readIcecreamIngredients)
			read.GET("/:ids/sourcingvalues", s.readIcecreamSourcingValues)
		}

		update := icecreams.Group("").Use(s.icecreamRequest)
		{
			update.PATCH("", s.updateIcecreams)
		}

		// delete collides with an in-built function
		del := icecreams.Group("")
		{
			del.DELETE("", func(c *gin.Context) {
				c.JSON(http.StatusMethodNotAllowed, FailStringResponse("deleting the entire collection is not allowed"))
			})
			del.DELETE("/:ids", s.deleteIcecreams)
			del.DELETE("/:ids/", s.deleteIcecreams)
			del.DELETE("/:ids/sourcingvalues", s.deleteIcecreamSourcingValues)

			del.DELETE("/:ids/ingredients", func(c *gin.Context) {
				c.JSON(http.StatusMethodNotAllowed, FailStringResponse("deleting all ingredients of an icecream is not allowed"))
			})
			del.DELETE("/:ids/ingredients/:name", func(c *gin.Context) {
				c.JSON(http.StatusNotImplemented, ErrorResponse("not implemented yet"))
			})
			del.DELETE("/:ids/sourcingvalues/:name", func(c *gin.Context) {
				c.JSON(http.StatusNotImplemented, ErrorResponse("not implemented yet"))
			})
		}
	}

	ingredients := s.engine.Group("/ingredients")
	{
		ingredients.GET("", s.readIngredients)
	}

	sourcingvalues := s.engine.Group("/sourcingvalues")
	{
		sourcingvalues.GET("", s.readSourcingValues)
	}

	return s
}

func (s *Server) icecreamRequest(c *gin.Context) {

	if !strings.Contains(c.ContentType(), "application/json") {
		c.AbortWithStatusJSON(http.StatusBadRequest, FailStringResponse("only Content-Type: application/json is supported"))
		return
	}

	var icecreams []*domain.Icecream
	if err := c.ShouldBind(&icecreams); err != nil {
		if err.Error() == "EOF" {
			c.AbortWithStatusJSON(http.StatusBadRequest, FailStringResponse("no icecream data provided"))
			return
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, FailResponse(fmt.Errorf("faulty data provided: %v (forgot to wrap in []?)", err)))
		return
	}

	c.Set(RequestIcecreamKey, icecreams)
	c.Next()
}

func (s *Server) readIcecreams(c *gin.Context) {

	ids, err := convertIdsParam(c.Param("ids"))
	if err != nil {
		c.JSON(http.StatusBadRequest, FailResponse(err))
		return
	}

	icecreams, err := s.repo.IcecreamService.Reads(ids)
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

	if icecreams == nil {
		c.JSON(http.StatusNotFound, FailStringResponse("no icecream found"))
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

	ingredients, err := s.repo.IngredientService.Reads(ids)
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

func (s *Server) readIcecreamSourcingValues(c *gin.Context) {

	ids, err := convertIdsParam(c.Param("ids"))
	if err != nil {
		c.JSON(http.StatusBadRequest, FailResponse(err))
		return
	}

	sourcingValues, err := s.repo.SourcingValueService.Reads(ids)
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

func (s *Server) createIcecreams(c *gin.Context) {

	icecreams := c.MustGet(RequestIcecreamKey).([]*domain.Icecream)

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

		if existingIcecream, _ := s.repo.IcecreamService.Reads([]int{productId}); existingIcecream != nil {
			c.JSON(http.StatusBadRequest, FailStringResponse("icecream with productId = "+icecream.ProductID+" already exists"))
			return
		}
	}

	if _, err := s.repo.IcecreamService.Creates(icecreams); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse("a database error occured, please try again later"))
		return
	}

	c.JSON(http.StatusCreated, SuccessResponse(
		&IcecreamsResponse{Icecreams: icecreams},
	))
}

func (s *Server) updateIcecreams(c *gin.Context) {

	icecreams := c.MustGet(RequestIcecreamKey).([]*domain.Icecream)

	for k, icecream := range icecreams {
		if err := icecream.Verify(); err != nil {
			c.JSON(http.StatusBadRequest, FailResponse(fmt.Errorf("icecream #%d: %v", k, err)))
			return
		}
	}

	if err := s.repo.IcecreamService.Updates(icecreams); err != nil {
		c.JSON(http.StatusInternalServerError, FailResponse(err))
		return
	}

	c.JSON(http.StatusOK, SuccessResponse(
		&IcecreamsResponse{Icecreams: icecreams},
	))
}

func (s *Server) deleteIcecreams(c *gin.Context) {

	ids, err := convertIdsParam(c.Param("ids"))
	if err != nil {
		c.JSON(http.StatusBadRequest, FailResponse(err))
		return
	}

	if err := s.repo.IcecreamService.Deletes(ids); err != nil {
		c.JSON(http.StatusInternalServerError, FailResponse(err))
		return
	}

	c.JSON(http.StatusOK, SuccessResponse(nil))
}

func (s *Server) deleteIcecreamSourcingValues(c *gin.Context) {

	ids, err := convertIdsParam(c.Param("ids"))
	if err != nil {
		c.JSON(http.StatusBadRequest, FailResponse(err))
		return
	}

	if err := s.repo.SourcingValueService.Deletes(ids); err != nil {
		c.JSON(http.StatusInternalServerError, FailResponse(err))
		return
	}

	c.JSON(http.StatusOK, SuccessResponse(nil))
}

func (s *Server) readIngredients(c *gin.Context) {
	ingredients, err := s.repo.IngredientService.ReadAll()
	if err != nil {
		log.Printf("could not get ingredients: %v", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse("a database error occured, please try again later"))
	}
	c.JSON(http.StatusOK, SuccessResponse(
		&IngredientResponse{Ingredient: ingredients}),
	)
}

func (s *Server) readSourcingValues(c *gin.Context) {
	sourcingValues, err := s.repo.SourcingValueService.ReadAll()
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
