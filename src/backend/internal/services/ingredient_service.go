package services

import (
	"fmt"
	"strings"

	"onetaste-family/backend/internal/models"
	"onetaste-family/backend/internal/repositories"
)

// IngredientService 食材服务
type IngredientService struct {
	ingredientRepo *repositories.IngredientRepository
}

// NewIngredientService 创建服务
func NewIngredientService() *IngredientService {
	return &IngredientService{
		ingredientRepo: repositories.NewIngredientRepository(),
	}
}

// SearchIngredients 模糊搜索
func (s *IngredientService) SearchIngredients(req *models.IngredientSearchQuery) ([]*models.IngredientSearchResult, error) {
	keyword := strings.TrimSpace(req.Keyword)
	if keyword == "" {
		return nil, fmt.Errorf("keyword required")
	}

	return s.ingredientRepo.SearchActiveIngredients(keyword, req.Limit)
}

// GetIngredientsByCategory 分类分页查询
func (s *IngredientService) GetIngredientsByCategory(req *models.IngredientCategoryQuery) (*models.IngredientCategoryListResponse, error) {
	category := strings.TrimSpace(req.Category)
	keyword := strings.TrimSpace(req.Keyword)

	items, total, err := s.ingredientRepo.GetActiveByCategory(category, keyword, req.Page, req.PageSize)
	if err != nil {
		return nil, err
	}

	return &models.IngredientCategoryListResponse{
		Items:    items,
		Page:     req.Page,
		PageSize: req.PageSize,
		Total:    total,
	}, nil
}
