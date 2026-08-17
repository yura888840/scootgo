package ride

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type Router struct {
	logger *zerolog.Logger
}

func (r *Router) CreateNewRouter(logger *zerolog.Logger) *Router {
	return &Router{
		logger: logger,
	}
}

func NewRideUC(c *gin.Context) {

}

func (r *Router) SetupRoutes(routesEngine *gin.Engine) error {
	group := routesEngine.Group("/ride")

	group.GET("/new", NewRideUC)
	return nil
}
