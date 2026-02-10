package dtos

import "time"

type CreateMapRequest struct {
	Name        string            `json:"name" binding:"required"`
	Description string            `json:"description" binding:"required"`
	Locations   []LocationInputDTO `json:"locations" binding:"required,dive"`
}

type LocationInputDTO struct {
	PanoId    string  `json:"pano_id" binding:"required"`
	Latitude  float64 `json:"latitude" binding:"required"`
	Longitude float64 `json:"longitude" binding:"required"`
	Heading   float64 `json:"heading"`
	Pitch     float64 `json:"pitch"`
}

type CreateMapResponse struct {
	ID          string               `json:"id"`
	Name        string               `json:"name"`
	Description string               `json:"description"`
	OwnerId     string               `json:"owner_id"`
	Locations   []LocationOutputDTO  `json:"locations"`
	CreatedAt   time.Time            `json:"created_at"`
}

type LocationOutputDTO struct {
	ID        string  `json:"id"`
	PanoId    string  `json:"pano_id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Heading   float64 `json:"heading"`
	Pitch     float64 `json:"pitch"`
}
