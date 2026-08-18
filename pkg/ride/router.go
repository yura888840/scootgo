package ride

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type Router struct {
	logger *zerolog.Logger
}

func CreateNewRouter(logger *zerolog.Logger) *Router {
	return &Router{
		logger: logger,
	}
}
func (r *Router) SetupRoutes(routesEngine *gin.Engine) error {
	group := routesEngine.Group("/api/v1/ride")
	h := SetupRideHandler(r.logger)

	group.POST("/start", h.StartRide)
	group.POST("/end", h.EndRide)
	group.GET("/current", h.GetRide)
	group.GET("/list", h.ListMyRides)

	return nil
}
