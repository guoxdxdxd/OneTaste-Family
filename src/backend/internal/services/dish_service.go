package services

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"onetaste-family/backend/internal/models"
	"onetaste-family/backend/internal/repositories"
	"onetaste-family/backend/internal/utils"
)

const (
	maxDishIngredients = 50
	maxDishSteps       = 50
)

var (
	// ErrDishLimitReached 超过家庭菜式数量限制
	ErrDishLimitReached = errors.New("dish limit reached")
	// ErrInvalidDishName 菜式名称非法
	ErrInvalidDishName = errors.New("invalid dish name")
	// ErrDishNameExists 菜式名称重复
	ErrDishNameExists = errors.New("dish name exists")
	// ErrDishPermissionDenied 无权限操作菜式
	ErrDishPermissionDenied = errors.New("no permission to modify dish")
	// ErrDishNotFound 菜式不存在
	ErrDishNotFound = errors.New("dish not found")
	// ErrInvalidDishIngredients 食材列表非法
	ErrInvalidDishIngredients = errors.New("invalid ingredients")
	// ErrInvalidDishSteps 烹饪步骤非法
	ErrInvalidDishSteps = errors.New("invalid cooking steps")
)

// DishService 菜式业务逻辑层
type DishService struct {
	dishRepo   *repositories.DishRepository
	familyRepo *repositories.FamilyRepository
}

// NewDishService 创建DishService
func NewDishService() *DishService {
	return &DishService{
		dishRepo:   repositories.NewDishRepository(),
		familyRepo: repositories.NewFamilyRepository(),
	}
}

// CreateDish 创建菜式
func (s *DishService) CreateDish(userID string, req *models.CreateDishRequest) (*models.DishCreateResponse, error) {
	family, err := s.getFamilyForUser(userID)
	if err != nil {
		return nil, err
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, ErrInvalidDishName
	}

	if len(req.Ingredients) == 0 || len(req.Ingredients) > maxDishIngredients {
		return nil, ErrInvalidDishIngredients
	}
	if len(req.Steps) == 0 || len(req.Steps) > maxDishSteps {
		return nil, ErrInvalidDishSteps
	}

	count, err := s.dishRepo.CountByFamily(family.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to count dishes: %w", err)
	}
	if count >= family.MaxDishes {
		return nil, ErrDishLimitReached
	}

	exists, err := s.dishRepo.ExistsByName(family.ID, name, "")
	if err != nil {
		return nil, fmt.Errorf("failed to check dish name: %w", err)
	}
	if exists {
		return nil, ErrDishNameExists
	}

	ingredients, err := convertIngredients(req.Ingredients)
	if err != nil {
		return nil, err
	}

	steps, err := convertCookingSteps(req.Steps)
	if err != nil {
		return nil, err
	}

	dish := &models.Dish{
		ID:          utils.GenerateULID(),
		FamilyID:    family.ID,
		Name:        name,
		Category:    strings.TrimSpace(req.Category),
		Description: strings.TrimSpace(req.Description),
		ImageURL:    strings.TrimSpace(req.ImageURL),
		CreatedBy:   userID,
	}

	if err := s.dishRepo.CreateDishWithDetails(dish, ingredients, steps); err != nil {
		return nil, fmt.Errorf("failed to create dish: %w", err)
	}

	return &models.DishCreateResponse{
		DishID: dish.ID,
		Name:   dish.Name,
	}, nil
}

// GetDishList 获取菜式列表
func (s *DishService) GetDishList(userID string, req *models.DishListRequest) (*models.DishListResponse, error) {
	family, err := s.getFamilyForUser(userID)
	if err != nil {
		return nil, err
	}

	category := strings.TrimSpace(req.Category)
	keyword := strings.TrimSpace(req.Keyword)

	dishes, total, err := s.dishRepo.GetDishList(family.ID, req.Page, req.PageSize, category, keyword)
	if err != nil {
		return nil, fmt.Errorf("failed to query dishes: %w", err)
	}

	return &models.DishListResponse{
		Dishes:   dishes,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

// GetDishDetail 获取菜式详情
func (s *DishService) GetDishDetail(userID, dishID string) (*models.DishDetailResponse, error) {
	family, err := s.getFamilyForUser(userID)
	if err != nil {
		return nil, err
	}

	dish, err := s.dishRepo.GetDishByID(dishID, family.ID)
	if err != nil {
		if errors.Is(err, repositories.ErrDishNotFound) {
			return nil, ErrDishNotFound
		}
		return nil, fmt.Errorf("failed to get dish: %w", err)
	}

	ingredients, err := s.dishRepo.GetIngredients(dish.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get ingredients: %w", err)
	}

	steps, err := s.dishRepo.GetCookingSteps(dish.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cooking steps: %w", err)
	}

	return buildDishDetailResponse(dish, ingredients, steps), nil
}

// UpdateDish 更新菜式
func (s *DishService) UpdateDish(userID, dishID string, req *models.UpdateDishRequest) (*models.DishDetailResponse, error) {
	family, err := s.getFamilyForUser(userID)
	if err != nil {
		return nil, err
	}

	dish, err := s.dishRepo.GetDishByID(dishID, family.ID)
	if err != nil {
		if errors.Is(err, repositories.ErrDishNotFound) {
			return nil, ErrDishNotFound
		}
		return nil, fmt.Errorf("failed to get dish: %w", err)
	}

	if dish.CreatedBy != userID && family.OwnerID != userID {
		return nil, ErrDishPermissionDenied
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, ErrInvalidDishName
	}

	if len(req.Ingredients) == 0 || len(req.Ingredients) > maxDishIngredients {
		return nil, ErrInvalidDishIngredients
	}
	if len(req.Steps) == 0 || len(req.Steps) > maxDishSteps {
		return nil, ErrInvalidDishSteps
	}

	if !strings.EqualFold(name, dish.Name) {
		exists, err := s.dishRepo.ExistsByName(family.ID, name, dish.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to check dish name: %w", err)
		}
		if exists {
			return nil, ErrDishNameExists
		}
	}

	ingredients, err := convertIngredients(req.Ingredients)
	if err != nil {
		return nil, err
	}

	steps, err := convertCookingSteps(req.Steps)
	if err != nil {
		return nil, err
	}

	dish.Name = name
	dish.Category = strings.TrimSpace(req.Category)
	dish.Description = strings.TrimSpace(req.Description)
	dish.ImageURL = strings.TrimSpace(req.ImageURL)

	if err := s.dishRepo.UpdateDishWithDetails(dish, ingredients, steps); err != nil {
		if errors.Is(err, repositories.ErrDishNotFound) {
			return nil, ErrDishNotFound
		}
		return nil, fmt.Errorf("failed to update dish: %w", err)
	}

	return buildDishDetailResponse(dish, ingredients, steps), nil
}

// DeleteDish 删除菜式
func (s *DishService) DeleteDish(userID, dishID string) error {
	family, err := s.getFamilyForUser(userID)
	if err != nil {
		return err
	}

	dish, err := s.dishRepo.GetDishByID(dishID, family.ID)
	if err != nil {
		if errors.Is(err, repositories.ErrDishNotFound) {
			return ErrDishNotFound
		}
		return fmt.Errorf("failed to get dish: %w", err)
	}

	if dish.CreatedBy != userID && family.OwnerID != userID {
		return ErrDishPermissionDenied
	}

	if err := s.dishRepo.SoftDeleteDish(dish.ID, family.ID); err != nil {
		if errors.Is(err, repositories.ErrDishNotFound) {
			return ErrDishNotFound
		}
		return fmt.Errorf("failed to delete dish: %w", err)
	}

	return nil
}

func (s *DishService) getFamilyForUser(userID string) (*models.Family, error) {
	family, err := s.familyRepo.GetFamilyByUserID(userID)
	if err != nil {
		if errors.Is(err, repositories.ErrFamilyNotFound) {
			return nil, ErrFamilyNotFound
		}
		return nil, fmt.Errorf("failed to get family: %w", err)
	}
	return family, nil
}

func convertIngredients(inputs []models.IngredientInput) ([]*models.Ingredient, error) {
	if len(inputs) == 0 {
		return nil, ErrInvalidDishIngredients
	}

	ingredients := make([]*models.Ingredient, 0, len(inputs))
	for index, input := range inputs {
		name := strings.TrimSpace(input.Name)
		unit := strings.TrimSpace(input.Unit)
		if name == "" || unit == "" {
			return nil, ErrInvalidDishIngredients
		}

		sortOrder := input.SortOrder
		if sortOrder <= 0 {
			sortOrder = index + 1
		}

		category := strings.TrimSpace(input.Category)
		ingredient := &models.Ingredient{
			ID:        utils.GenerateULID(),
			Name:      name,
			Amount:    input.Amount,
			Unit:      unit,
			Category:  category,
			SortOrder: sortOrder,
		}

		if input.StorageDays != nil && *input.StorageDays >= 0 {
			value := *input.StorageDays
			ingredient.StorageDays = &value
		}

		ingredients = append(ingredients, ingredient)
	}

	sort.SliceStable(ingredients, func(i, j int) bool {
		if ingredients[i].SortOrder == ingredients[j].SortOrder {
			return i < j
		}
		return ingredients[i].SortOrder < ingredients[j].SortOrder
	})

	return ingredients, nil
}

func convertCookingSteps(inputs []models.CookingStepInput) ([]*models.CookingStep, error) {
	if len(inputs) == 0 {
		return nil, ErrInvalidDishSteps
	}

	steps := make([]*models.CookingStep, 0, len(inputs))
	for idx, input := range inputs {
		content := strings.TrimSpace(input.Content)
		if content == "" {
			return nil, ErrInvalidDishSteps
		}

		order := input.Order
		if order <= 0 {
			order = idx + 1
		}

		step := &models.CookingStep{
			ID:       utils.GenerateULID(),
			Order:    order,
			Content:  content,
			ImageURL: strings.TrimSpace(input.ImageURL),
		}
		steps = append(steps, step)
	}

	sort.SliceStable(steps, func(i, j int) bool {
		if steps[i].Order == steps[j].Order {
			return i < j
		}
		return steps[i].Order < steps[j].Order
	})

	return steps, nil
}

func buildDishDetailResponse(dish *models.Dish, ingredients []*models.Ingredient, steps []*models.CookingStep) *models.DishDetailResponse {
	return &models.DishDetailResponse{
		DishID:      dish.ID,
		Name:        dish.Name,
		Category:    dish.Category,
		Description: dish.Description,
		ImageURL:    dish.ImageURL,
		Ingredients: ingredients,
		Steps:       steps,
		CreatedAt:   dish.CreatedAt,
		UpdatedAt:   dish.UpdatedAt,
	}
}
