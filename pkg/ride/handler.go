package ride

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/yura888840/scootgo/pkg/ride/usecase"
)

type RideHandler struct {
	usecase usecase.RideUseCase
}

func (r *RideHandler) StartRide(c *gin.Context) {
	c.Status(http.StatusOK)
}

func (r *RideHandler) EndRide(c *gin.Context) {
	c.Status(http.StatusOK)
}

func (r *RideHandler) GetRide(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

func (r *RideHandler) ListMyRides(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

func SetupRideHandler(logger *zerolog.Logger) *RideHandler {
	uc := usecase.NewRideUseCase(logger)

	return &RideHandler{
		usecase: uc,
	}
}
