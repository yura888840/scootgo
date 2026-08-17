package usecase

import "github.com/gin-gonic/gin"

type StartRideRequest struct {
	UserID int64
}
type EndRideRequest struct {
}
type GetRideRequest struct {
}
type ListMyRidesRequest struct {
}

type RideUseCase interface {
	StartRide(c *gin.Context, request StartRideRequest) error
	EndRide(c *gin.Context, request EndRideRequest)
	GetRide(c *gin.Context, req GetRideRequest)
	ListMyRides(c *gin.Context, req ListMyRidesRequest)
}
