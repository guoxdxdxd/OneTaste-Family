package models

import "time"

const (
	// MealTypeBreakfast 早餐
	MealTypeBreakfast = "breakfast"
	// MealTypeLunch 午餐
	MealTypeLunch = "lunch"
	// MealTypeDinner 晚餐
	MealTypeDinner = "dinner"
)

const (
	// MenuSourceManual 手动创建
	MenuSourceManual = "manual"
	// MenuSourceAI AI生成
	MenuSourceAI = "ai"
)

// CreateMenuRequest 创建菜单请求
type CreateMenuRequest struct {
	Date      string   `json:"date" binding:"required"`      // 日期，格式：YYYY-MM-DD
	MealType  string   `json:"meal_type" binding:"required,oneof=breakfast lunch dinner"` // 餐次
	DishIDs   []string `json:"dish_ids" binding:"required,min=1,dive,len=26"`            // 菜式ID列表，至少1个
}

// UpdateMenuRequest 更新菜单请求
type UpdateMenuRequest struct {
	Date     string   `json:"date" binding:"omitempty"`     // 日期，格式：YYYY-MM-DD
	MealType string   `json:"meal_type" binding:"omitempty,oneof=breakfast lunch dinner"` // 餐次
	DishIDs  []string `json:"dish_ids" binding:"omitempty,min=1,dive,len=26"`           // 菜式ID列表
}

// MenuIDRequest 菜单ID请求
type MenuIDRequest struct {
	ID string `uri:"id" binding:"required,len=26"`
}

// DailyMenuRequest 每日菜单查询请求
type DailyMenuRequest struct {
	Date string `form:"date" binding:"required"` // 日期，格式：YYYY-MM-DD
}

// WeeklyMenuRequest 每周菜单查询请求
type WeeklyMenuRequest struct {
	StartDate string `form:"start_date" binding:"required"` // 开始日期，格式：YYYY-MM-DD
}

// Menu 菜单数据库实体
type Menu struct {
	ID        string    `json:"menu_id"`
	FamilyID string    `json:"family_id"`
	Date     time.Time `json:"date"`
	MealType string    `json:"meal_type"`
	CreatedBy string    `json:"created_by"`
	Source   string    `json:"source"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// MenuDish 菜单菜式关联实体
type MenuDish struct {
	ID        string    `json:"id"`
	MenuID    string    `json:"menu_id"`
	DishID    string    `json:"dish_id"`
	CreatedAt time.Time `json:"created_at"`
}

// MenuDetail 菜单详情（包含菜式信息）
type MenuDetail struct {
	MenuID    string         `json:"menu_id"`
	FamilyID  string         `json:"family_id"`
	Date      string         `json:"date"`      // 格式：YYYY-MM-DD
	MealType  string         `json:"meal_type"` // breakfast, lunch, dinner
	CreatedBy string         `json:"created_by"`
	Source    string         `json:"source"`
	Dishes    []*DishSummary `json:"dishes"`   // 菜式列表
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

// MenuCreateResponse 创建菜单响应
type MenuCreateResponse struct {
	MenuID   string         `json:"menu_id"`
	Date     string         `json:"date"`
	MealType string         `json:"meal_type"`
	Dishes   []*DishSummary `json:"dishes"`
}

// MenuUpdateResponse 更新菜单响应
type MenuUpdateResponse MenuDetail

// DailyMenuResponse 每日菜单响应
type DailyMenuResponse struct {
	Date  string        `json:"date"`
	Menus []*MenuDetail `json:"menus"` // 三餐菜单（早餐、午餐、晚餐）
}

// WeeklyMenuResponse 每周菜单响应
type WeeklyMenuResponse struct {
	StartDate string        `json:"start_date"`
	EndDate   string        `json:"end_date"`
	Menus     []*MenuDetail `json:"menus"` // 一周的菜单列表
}

