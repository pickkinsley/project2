package models

import "time"

// PackingItem maps to the packing_items table.
// This same struct is used as the API response — the shape is identical.
type PackingItem struct {
	ID          int       `json:"id"`
	TripID      string    `json:"-"`           // Foreign key — omitted from API response
	Name        string    `json:"name"`
	Category    string    `json:"category"`
	IsEssential bool      `json:"is_essential"`
	Reason      string    `json:"reason"`
	IsChecked   bool      `json:"is_checked"`
	SortOrder   int       `json:"sort_order"`
	CreatedAt   time.Time `json:"-"`           // Internal — omitted from API response
}

// CreatePackingItemRequest is the JSON body received by POST /api/trips/{uuid}/items.
type CreatePackingItemRequest struct {
	Name     string `json:"name"`
	Category string `json:"category"`
}

// UpdateItemRequest is the JSON body received by PATCH /api/trips/{uuid}/items/{itemId}.
type UpdateItemRequest struct {
	IsChecked bool `json:"is_checked"`
}

// UpdateItemResponse is returned after a successful PATCH.
type UpdateItemResponse struct {
	ID        int  `json:"id"`
	IsChecked bool `json:"is_checked"`
}
