package dto

import (
	"fmt"

	"github.com/google/uuid"
)

const (
	DefaultPage     = 1
	DefaultPageSize = 10
	MaxPageSize     = 100
)

// CreateCarRequest represents the request body for creating a car
type CreateCarRequest struct {
	Name          string `json:"name" binding:"required,min=2,max=100" example:"Honda Civic"`
	EngineVersion string `json:"engine_version" binding:"required,oneof=1.0 1.4 1.5 1.6 1.8 2.0 2.4 2.5 3.0 3.5 4.0" example:"2.0"`
}

// UpdateCarRequest represents the request body for updating a car
type UpdateCarRequest struct {
	Name          string `json:"name" binding:"omitempty,min=2,max=100" example:"Honda Civic Sport"`
	EngineVersion string `json:"engine_version" binding:"omitempty,oneof=1.0 1.4 1.5 1.6 1.8 2.0 2.4 2.5 3.0 3.5 4.0" example:"2.0"`
}

// CarResponse represents the response body for a car
type CarResponse struct {
	ID            uuid.UUID `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name          string    `json:"name" example:"Honda Civic"`
	EngineVersion string    `json:"engine_version" example:"2.0"`
	CreatedAt     string    `json:"created_at" example:"2024-01-01T10:00:00Z"`
	UpdatedAt     string    `json:"updated_at" example:"2024-01-01T10:00:00Z"`
}

// PaginationRequest represents pagination parameters
type PaginationRequest struct {
	Page     int    `form:"page" binding:"omitempty,min=1" example:"1"`
	PageSize int    `form:"page_size" binding:"omitempty,min=1,max=100" example:"10"`
	SortBy   string `form:"sort_by" binding:"omitempty,oneof=name engine_version created_at" example:"created_at"`
	SortDir  string `form:"sort_dir" binding:"omitempty,oneof=asc desc" example:"desc"`
}

// PaginatedResponse represents a paginated response
type PaginatedResponse struct {
	Data       []CarResponse   `json:"data"`
	Pagination PaginationMeta `json:"pagination"`
}

// PaginationMeta represents pagination metadata
type PaginationMeta struct {
	CurrentPage  int   `json:"current_page" example:"1"`
	PageSize     int   `json:"page_size" example:"10"`
	TotalPages   int   `json:"total_pages" example:"5"`
	TotalRecords int64 `json:"total_records" example:"50"`
}

// SetDefaults sets default values for pagination
func (p *PaginationRequest) SetDefaults() {
	if p.Page < 1 {
		p.Page = DefaultPage
	}
	if p.PageSize < 1 {
		p.PageSize = DefaultPageSize
	}
	if p.PageSize > MaxPageSize {
		p.PageSize = MaxPageSize
	}
	if p.SortBy == "" {
		p.SortBy = "created_at"
	}
	if p.SortDir == "" {
		p.SortDir = "desc"
	}
}

// GetOffset calculates the offset for pagination
func (p *PaginationRequest) GetOffset() int {
	return (p.Page - 1) * p.PageSize
}

// GetOrderBy returns the ORDER BY clause with proper sanitization
func (p *PaginationRequest) GetOrderBy() string {
	// Whitelist of allowed columns to prevent SQL injection
	allowedColumns := map[string]string{
		"name":           "name",
		"engine_version": "engine_version",
		"created_at":     "created_at",
	}

	column, ok := allowedColumns[p.SortBy]
	if !ok {
		column = "created_at"
	}

	// Sanitize direction
	direction := "DESC"
	if p.SortDir == "asc" {
		direction = "ASC"
	}

	return fmt.Sprintf("%s %s", column, direction)
}
