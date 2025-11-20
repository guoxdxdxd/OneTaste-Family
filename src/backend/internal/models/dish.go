package models

import "time"

// IngredientInput 食材入参
type IngredientInput struct {
	Name        string  `json:"name" binding:"required,max=100"`
	Amount      float64 `json:"amount" binding:"required,gt=0"`
	Unit        string  `json:"unit" binding:"required,max=20"`
	Category    string  `json:"category" binding:"omitempty,max=50"`
	StorageDays *int    `json:"storage_days" binding:"omitempty,min=0,max=365"`
	SortOrder   int     `json:"sort_order" binding:"omitempty,min=0,max=1000"`
}

// CookingStepInput 烹饪步骤入参
type CookingStepInput struct {
	Order    int    `json:"order" binding:"required,min=1,max=200"`
	Content  string `json:"content" binding:"required,max=2000"`
	ImageURL string `json:"image_url" binding:"omitempty,max=500"`
}

// CreateDishRequest 创建菜式请求
type CreateDishRequest struct {
	Name        string             `json:"name" binding:"required,max=100"`
	Category    string             `json:"category" binding:"omitempty,max=50"`
	Description string             `json:"description" binding:"omitempty,max=2000"`
	ImageURL    string             `json:"image_url" binding:"omitempty,max=500"`
	Ingredients []IngredientInput  `json:"ingredients" binding:"required"`
	Steps       []CookingStepInput `json:"steps" binding:"required"`
}

// UpdateDishRequest 更新菜式请求
type UpdateDishRequest CreateDishRequest

// DishListRequest 菜式列表查询请求
type DishListRequest struct {
	Page     int    `form:"page,default=1" binding:"min=1"`
	PageSize int    `form:"page_size,default=20" binding:"min=1,max=100"`
	Category string `form:"category" binding:"omitempty,max=50"`
	Keyword  string `form:"keyword" binding:"omitempty,max=100"`
}

// DishIDRequest 菜式ID请求
type DishIDRequest struct {
	ID string `uri:"id" binding:"required,len=26"`
}

// Dish 菜式数据库实体
type Dish struct {
	ID          string    `json:"dish_id"`
	FamilyID    string    `json:"family_id"`
	Name        string    `json:"name"`
	Category    string    `json:"category,omitempty"`
	Description string    `json:"description,omitempty"`
	ImageURL    string    `json:"image_url,omitempty"`
	CreatedBy   string    `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Ingredient 食材实体
type Ingredient struct {
	ID          string  `json:"ingredient_id"`
	DishID      string  `json:"-"`
	Name        string  `json:"name"`
	Amount      float64 `json:"amount"`
	Unit        string  `json:"unit"`
	Category    string  `json:"category,omitempty"`
	StorageDays *int    `json:"storage_days,omitempty"`
	SortOrder   int     `json:"sort_order"`
}

// CookingStep 烹饪步骤实体
type CookingStep struct {
	ID       string `json:"step_id"`
	DishID   string `json:"-"`
	Order    int    `json:"order"`
	Content  string `json:"content"`
	ImageURL string `json:"image_url,omitempty"`
}

// DishCreateResponse 创建菜式响应
type DishCreateResponse struct {
	DishID string `json:"dish_id"`
	Name   string `json:"name"`
}

// DishSummary 菜式列表项
type DishSummary struct {
	DishID    string    `json:"dish_id"`
	Name      string    `json:"name"`
	Category  string    `json:"category,omitempty"`
	ImageURL  string    `json:"image_url,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// DishListResponse 菜式列表响应
type DishListResponse struct {
	Dishes   []*DishSummary `json:"dishes"`
	Total    int64          `json:"total"`
	Page     int            `json:"page"`
	PageSize int            `json:"page_size"`
}

// DishDetailResponse 菜式详情响应
type DishDetailResponse struct {
	DishID      string         `json:"dish_id"`
	Name        string         `json:"name"`
	Category    string         `json:"category,omitempty"`
	Description string         `json:"description,omitempty"`
	ImageURL    string         `json:"image_url,omitempty"`
	Ingredients []*Ingredient  `json:"ingredients"`
	Steps       []*CookingStep `json:"steps"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}
