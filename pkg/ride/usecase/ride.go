package usecase

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type StartRideRequest struct {
	UserID int64
}

type EndRideRequest struct {
}

type GetRideRequest struct {
}
type GetRideResponse struct {
}

type ListMyRidesRequest struct {
}
type ListMyRidesResponse struct {
}

type RideUseCase interface {
	StartRide(c *gin.Context, request StartRideRequest) error
	EndRide(c *gin.Context, request EndRideRequest) error
	GetRide(c *gin.Context, req GetRideRequest) (*GetRideResponse, error)
	ListMyRides(c *gin.Context, req ListMyRidesRequest) (*ListMyRidesResponse, error)
}

type rideUseCase struct {
	// repository here
	logger *zerolog.Logger
}

func NewRideUseCase(logger *zerolog.Logger) RideUseCase {
	return &rideUseCase{
		logger: logger,
	}
}

func (uc *rideUseCase) StartRide(c *gin.Context, request StartRideRequest) error {
	return nil
}

func (uc *rideUseCase) EndRide(c *gin.Context, request EndRideRequest) error {
	return nil
}

func (uc *rideUseCase) GetRide(c *gin.Context, req GetRideRequest) (*GetRideResponse, error) {
	return nil, nil
}

func (uc *rideUseCase) ListMyRides(c *gin.Context, req ListMyRidesRequest) (*ListMyRidesResponse, error) {
	return nil, nil
}
