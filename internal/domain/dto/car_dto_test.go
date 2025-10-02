package dto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPaginationRequest_SetDefaults(t *testing.T) {
	tests := []struct {
		name     string
		input    PaginationRequest
		expected PaginationRequest
	}{
		{
			name:  "Empty pagination request should set defaults",
			input: PaginationRequest{},
			expected: PaginationRequest{
				Page:     DefaultPage,
				PageSize: DefaultPageSize,
				SortBy:   "created_at",
				SortDir:  "desc",
			},
		},
		{
			name: "Invalid page number should be corrected",
			input: PaginationRequest{
				Page:     0,
				PageSize: 20,
			},
			expected: PaginationRequest{
				Page:     DefaultPage,
				PageSize: 20,
				SortBy:   "created_at",
				SortDir:  "desc",
			},
		},
		{
			name: "Page size exceeding max should be limited",
			input: PaginationRequest{
				Page:     1,
				PageSize: 150,
			},
			expected: PaginationRequest{
				Page:     1,
				PageSize: MaxPageSize,
				SortBy:   "created_at",
				SortDir:  "desc",
			},
		},
		{
			name: "Page size below 1 should use default",
			input: PaginationRequest{
				Page:     1,
				PageSize: 0,
			},
			expected: PaginationRequest{
				Page:     1,
				PageSize: DefaultPageSize,
				SortBy:   "created_at",
				SortDir:  "desc",
			},
		},
		{
			name: "Valid custom values should be preserved",
			input: PaginationRequest{
				Page:     2,
				PageSize: 25,
				SortBy:   "name",
				SortDir:  "asc",
			},
			expected: PaginationRequest{
				Page:     2,
				PageSize: 25,
				SortBy:   "name",
				SortDir:  "asc",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.input.SetDefaults()
			assert.Equal(t, tt.expected, tt.input)
		})
	}
}

func TestPaginationRequest_GetOffset(t *testing.T) {
	tests := []struct {
		name     string
		page     int
		pageSize int
		expected int
	}{
		{
			name:     "First page should have offset 0",
			page:     1,
			pageSize: 10,
			expected: 0,
		},
		{
			name:     "Second page with page size 10",
			page:     2,
			pageSize: 10,
			expected: 10,
		},
		{
			name:     "Third page with page size 25",
			page:     3,
			pageSize: 25,
			expected: 50,
		},
		{
			name:     "Page 5 with page size 20",
			page:     5,
			pageSize: 20,
			expected: 80,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := PaginationRequest{
				Page:     tt.page,
				PageSize: tt.pageSize,
			}
			assert.Equal(t, tt.expected, p.GetOffset())
		})
	}
}

func TestPaginationRequest_GetOrderBy(t *testing.T) {
	tests := []struct {
		name     string
		sortBy   string
		sortDir  string
		expected string
	}{
		{
			name:     "Sort by name ascending",
			sortBy:   "name",
			sortDir:  "asc",
			expected: "name ASC",
		},
		{
			name:     "Sort by engine_version descending",
			sortBy:   "engine_version",
			sortDir:  "desc",
			expected: "engine_version DESC",
		},
		{
			name:     "Sort by created_at descending",
			sortBy:   "created_at",
			sortDir:  "desc",
			expected: "created_at DESC",
		},
		{
			name:     "Invalid sort field should default to created_at",
			sortBy:   "invalid_field",
			sortDir:  "asc",
			expected: "created_at ASC",
		},
		{
			name:     "Invalid sort direction should default to DESC",
			sortBy:   "name",
			sortDir:  "invalid",
			expected: "name DESC",
		},
		{
			name:     "SQL injection attempt should be blocked",
			sortBy:   "name; DROP TABLE cars;",
			sortDir:  "desc",
			expected: "created_at DESC",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := PaginationRequest{
				SortBy:  tt.sortBy,
				SortDir: tt.sortDir,
			}
			assert.Equal(t, tt.expected, p.GetOrderBy())
		})
	}
}

func TestPaginationConstants(t *testing.T) {
	assert.Equal(t, 1, DefaultPage)
	assert.Equal(t, 10, DefaultPageSize)
	assert.Equal(t, 100, MaxPageSize)
}
