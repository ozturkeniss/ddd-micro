package domain

import (
	"time"

	"gorm.io/gorm"
)

// Category represents the product category domain entity
type Category struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"not null;size:255" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	Slug        string         `gorm:"uniqueIndex;not null;size:255" json:"slug"`
	ParentID    *uint          `gorm:"index" json:"parent_id"`
	Parent      *Category      `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Children    []Category     `gorm:"foreignKey:ParentID" json:"children,omitempty"`
	Image       string         `gorm:"size:500" json:"image"`
	Icon        string         `gorm:"size:100" json:"icon"`
	SortOrder   int            `gorm:"default:0" json:"sort_order"`
	IsActive    bool           `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName specifies the table name for Category entity
func (Category) TableName() string {
	return "categories"
}

// IsValidName checks if the category name is valid
func (c *Category) IsValidName() bool {
	return len(c.Name) > 0 && len(c.Name) <= 255
}

// IsValidSlug checks if the slug is valid
func (c *Category) IsValidSlug() bool {
	return len(c.Slug) > 0 && len(c.Slug) <= 255
}

// IsRootCategory checks if this is a root category
func (c *Category) IsRootCategory() bool {
	return c.ParentID == nil
}

// IsChildCategory checks if this is a child category
func (c *Category) IsChildCategory() bool {
	return c.ParentID != nil
}

// HasChildren checks if this category has child categories
func (c *Category) HasChildren() bool {
	return len(c.Children) > 0
}

// Activate activates the category
func (c *Category) Activate() {
	c.IsActive = true
}

// Deactivate deactivates the category
func (c *Category) Deactivate() {
	c.IsActive = false
}

// SetSortOrder sets the sort order for the category
func (c *Category) SetSortOrder(order int) {
	c.SortOrder = order
}
