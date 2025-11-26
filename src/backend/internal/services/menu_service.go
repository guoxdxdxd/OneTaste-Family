package services

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"onetaste-family/backend/internal/models"
	"onetaste-family/backend/internal/repositories"
	"onetaste-family/backend/internal/utils"
)

var (
	// ErrMenuNotFound 菜单不存在
	ErrMenuNotFound = errors.New("menu not found")
	// ErrInvalidMenuDate 菜单日期非法
	ErrInvalidMenuDate = errors.New("invalid menu date")
	// ErrInvalidMealType 餐次类型非法
	ErrInvalidMealType = errors.New("invalid meal type")
	// ErrInvalidDishIDs 菜式ID列表非法
	ErrInvalidDishIDs = errors.New("invalid dish ids")
	// ErrDishNotInFamily 菜式不属于该家庭
	ErrDishNotInFamily = errors.New("dish not in family")
)

// MenuService 菜单业务逻辑层
type MenuService struct {
	menuRepo   *repositories.MenuRepository
	dishRepo   *repositories.DishRepository
	familyRepo *repositories.FamilyRepository
}

// NewMenuService 创建MenuService
func NewMenuService() *MenuService {
	return &MenuService{
		menuRepo:   repositories.NewMenuRepository(),
		dishRepo:   repositories.NewDishRepository(),
		familyRepo: repositories.NewFamilyRepository(),
	}
}

// CreateMenu 创建菜单
func (s *MenuService) CreateMenu(userID string, req *models.CreateMenuRequest) (*models.MenuCreateResponse, error) {
	family, err := s.getFamilyForUser(userID)
	if err != nil {
		return nil, err
	}

	// 解析日期
	date, err := parseDate(req.Date)
	if err != nil {
		return nil, ErrInvalidMenuDate
	}

	// 验证餐次
	if !isValidMealType(req.MealType) {
		return nil, ErrInvalidMealType
	}

	// 验证菜式ID列表
	if len(req.DishIDs) == 0 {
		return nil, ErrInvalidDishIDs
	}

	// 验证所有菜式都属于该家庭
	if err = s.validateDishesInFamily(family.ID, req.DishIDs); err != nil {
		return nil, err
	}

	// 检查是否已存在相同日期和餐次的菜单
	existingMenu, err := s.menuRepo.GetMenuByDateAndMealType(family.ID, date, req.MealType)
	if err != nil && !errors.Is(err, repositories.ErrMenuNotFound) {
		return nil, fmt.Errorf("failed to check existing menu: %w", err)
	}

	menu := &models.Menu{
		ID:        utils.GenerateULID(),
		FamilyID:  family.ID,
		Date:      date,
		MealType:  req.MealType,
		CreatedBy: userID,
		Source:    models.MenuSourceManual,
	}

	// 如果已存在，则更新；否则创建
	if existingMenu != nil {
		menu.ID = existingMenu.ID
		if err = s.menuRepo.UpdateMenuWithDishes(menu, req.DishIDs); err != nil {
			if errors.Is(err, repositories.ErrMenuNotFound) {
				return nil, ErrMenuNotFound
			}
			return nil, fmt.Errorf("failed to update menu: %w", err)
		}
	} else {
		if err = s.menuRepo.CreateMenuWithDishes(menu, req.DishIDs); err != nil {
			return nil, fmt.Errorf("failed to create menu: %w", err)
		}
	}

	// 获取菜式详情
	dishes, err := s.getDishSummaries(req.DishIDs, family.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get dish summaries: %w", err)
	}

	return &models.MenuCreateResponse{
		MenuID:   menu.ID,
		Date:     formatDate(menu.Date),
		MealType: menu.MealType,
		Dishes:   dishes,
	}, nil
}

// GetDailyMenu 获取每日菜单
func (s *MenuService) GetDailyMenu(userID string, req *models.DailyMenuRequest) (*models.DailyMenuResponse, error) {
	family, err := s.getFamilyForUser(userID)
	if err != nil {
		return nil, err
	}

	// 解析日期
	date, err := parseDate(req.Date)
	if err != nil {
		return nil, ErrInvalidMenuDate
	}

	// 获取该日期的所有菜单
	menus, err := s.menuRepo.GetMenusByDateRange(family.ID, date, date)
	if err != nil {
		return nil, fmt.Errorf("failed to get menus: %w", err)
	}

	// 构建三餐菜单
	menuMap := make(map[string]*models.Menu)
	for _, menu := range menus {
		menuMap[menu.MealType] = menu
	}

	// 按餐次顺序构建响应
	mealTypes := []string{models.MealTypeBreakfast, models.MealTypeLunch, models.MealTypeDinner}
	menuDetails := make([]*models.MenuDetail, 0, 3)

	for _, mealType := range mealTypes {
		if menu, exists := menuMap[mealType]; exists {
			detail, err := s.buildMenuDetail(menu)
			if err != nil {
				return nil, fmt.Errorf("failed to build menu detail: %w", err)
			}
			menuDetails = append(menuDetails, detail)
		}
	}

	return &models.DailyMenuResponse{
		Date:  formatDate(date),
		Menus: menuDetails,
	}, nil
}

// GetWeeklyMenu 获取每周菜单
func (s *MenuService) GetWeeklyMenu(userID string, req *models.WeeklyMenuRequest) (*models.WeeklyMenuResponse, error) {
	family, err := s.getFamilyForUser(userID)
	if err != nil {
		return nil, err
	}

	// 解析开始日期
	startDate, err := parseDate(req.StartDate)
	if err != nil {
		return nil, ErrInvalidMenuDate
	}

	// 计算结束日期（开始日期+6天）
	endDate := startDate.AddDate(0, 0, 6)

	// 获取该日期范围的所有菜单
	menus, err := s.menuRepo.GetMenusByDateRange(family.ID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get menus: %w", err)
	}

	// 构建菜单详情列表
	menuDetails := make([]*models.MenuDetail, 0, len(menus))
	for _, menu := range menus {
		detail, err := s.buildMenuDetail(menu)
		if err != nil {
			return nil, fmt.Errorf("failed to build menu detail: %w", err)
		}
		menuDetails = append(menuDetails, detail)
	}

	return &models.WeeklyMenuResponse{
		StartDate: formatDate(startDate),
		EndDate:   formatDate(endDate),
		Menus:     menuDetails,
	}, nil
}

// UpdateMenu 更新菜单
func (s *MenuService) UpdateMenu(userID, menuID string, req *models.UpdateMenuRequest) (*models.MenuUpdateResponse, error) {
	family, err := s.getFamilyForUser(userID)
	if err != nil {
		return nil, err
	}

	// 获取现有菜单
	menu, err := s.menuRepo.GetMenuByID(menuID, family.ID)
	if err != nil {
		if errors.Is(err, repositories.ErrMenuNotFound) {
			return nil, ErrMenuNotFound
		}
		return nil, fmt.Errorf("failed to get menu: %w", err)
	}

	// 更新日期（如果提供）
	if req.Date != "" {
		date, err := parseDate(req.Date)
		if err != nil {
			return nil, ErrInvalidMenuDate
		}
		menu.Date = date
	}

	// 更新餐次（如果提供）
	if req.MealType != "" {
		if !isValidMealType(req.MealType) {
			return nil, ErrInvalidMealType
		}
		menu.MealType = req.MealType
	}

	// 更新菜式列表（如果提供）
	if len(req.DishIDs) > 0 {
		// 验证所有菜式都属于该家庭
		if err = s.validateDishesInFamily(family.ID, req.DishIDs); err != nil {
			return nil, err
		}

		if err = s.menuRepo.UpdateMenuWithDishes(menu, req.DishIDs); err != nil {
			if errors.Is(err, repositories.ErrMenuNotFound) {
				return nil, ErrMenuNotFound
			}
			return nil, fmt.Errorf("failed to update menu: %w", err)
		}
	} else {
		// 只更新菜单基本信息，保留原有菜式列表
		// 获取原有菜式ID列表
		existingDishIDs, err := s.menuRepo.GetMenuDishes(menu.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get existing dishes: %w", err)
		}

		if err = s.menuRepo.UpdateMenuWithDishes(menu, existingDishIDs); err != nil {
			if errors.Is(err, repositories.ErrMenuNotFound) {
				return nil, ErrMenuNotFound
			}
			return nil, fmt.Errorf("failed to update menu: %w", err)
		}
	}

	// 获取更新后的菜单详情
	detail, err := s.buildMenuDetail(menu)
	if err != nil {
		return nil, fmt.Errorf("failed to build menu detail: %w", err)
	}

	return (*models.MenuUpdateResponse)(detail), nil
}

// buildMenuDetail 构建菜单详情
func (s *MenuService) buildMenuDetail(menu *models.Menu) (*models.MenuDetail, error) {
	// 获取菜单的菜式ID列表
	dishIDs, err := s.menuRepo.GetMenuDishes(menu.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get menu dishes: %w", err)
	}

	// 获取菜式详情
	dishes, err := s.getDishSummaries(dishIDs, menu.FamilyID)
	if err != nil {
		return nil, fmt.Errorf("failed to get dish summaries: %w", err)
	}

	return &models.MenuDetail{
		MenuID:    menu.ID,
		FamilyID:  menu.FamilyID,
		Date:      formatDate(menu.Date),
		MealType:  menu.MealType,
		CreatedBy: menu.CreatedBy,
		Source:    menu.Source,
		Dishes:    dishes,
		CreatedAt: menu.CreatedAt,
		UpdatedAt: menu.UpdatedAt,
	}, nil
}

// validateDishesInFamily 验证所有菜式都属于该家庭
func (s *MenuService) validateDishesInFamily(familyID string, dishIDs []string) error {
	for _, dishID := range dishIDs {
		dish, err := s.dishRepo.GetDishByID(dishID, familyID)
		if err != nil {
			if errors.Is(err, repositories.ErrDishNotFound) {
				return ErrDishNotFound
			}
			return fmt.Errorf("failed to get dish: %w", err)
		}
		if dish.FamilyID != familyID {
			return ErrDishNotInFamily
		}
	}
	return nil
}

// getDishSummaries 获取菜式摘要列表
func (s *MenuService) getDishSummaries(dishIDs []string, familyID string) ([]*models.DishSummary, error) {
	if len(dishIDs) == 0 {
		return []*models.DishSummary{}, nil
	}

	dishes := make([]*models.DishSummary, 0, len(dishIDs))
	for _, dishID := range dishIDs {
		dish, err := s.dishRepo.GetDishByID(dishID, familyID)
		if err != nil {
			if errors.Is(err, repositories.ErrDishNotFound) {
				continue // 跳过不存在的菜式
			}
			return nil, fmt.Errorf("failed to get dish: %w", err)
		}

		dishes = append(dishes, &models.DishSummary{
			DishID:      dish.ID,
			Name:        dish.Name,
			Category:    dish.Category,
			Description: dish.Description,
			ImageURL:    dish.ImageURL,
			CreatedAt:   dish.CreatedAt,
			UpdatedAt:   dish.UpdatedAt,
		})
	}

	return dishes, nil
}

func (s *MenuService) getFamilyForUser(userID string) (*models.Family, error) {
	family, err := s.familyRepo.GetFamilyByUserID(userID)
	if err != nil {
		if errors.Is(err, repositories.ErrFamilyNotFound) {
			return nil, ErrFamilyNotFound
		}
		return nil, fmt.Errorf("failed to get family: %w", err)
	}
	return family, nil
}

// parseDate 解析日期字符串
func parseDate(dateStr string) (time.Time, error) {
	return time.Parse("2006-01-02", strings.TrimSpace(dateStr))
}

// formatDate 格式化日期
func formatDate(date time.Time) string {
	return date.Format("2006-01-02")
}

// isValidMealType 验证餐次类型
func isValidMealType(mealType string) bool {
	return mealType == models.MealTypeBreakfast ||
		mealType == models.MealTypeLunch ||
		mealType == models.MealTypeDinner
}
