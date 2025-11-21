package repositories

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/lib/pq"
	"onetaste-family/backend/internal/models"
	"onetaste-family/backend/pkg/database"
)

// IngredientRepository 基础食材仓储
type IngredientRepository struct {
	db *sql.DB
}

// NewIngredientRepository 创建仓储
func NewIngredientRepository() *IngredientRepository {
	return &IngredientRepository{
		db: database.GetDB(),
	}
}

// GetActiveByIDs 根据ID集合获取可用的基础食材
func (r *IngredientRepository) GetActiveByIDs(ids []string) (map[string]*models.BasicIngredient, error) {
	result := make(map[string]*models.BasicIngredient, len(ids))
	if len(ids) == 0 {
		return result, nil
	}

	query := `
		SELECT id, name, name_en, category, default_unit, storage_days, description
		FROM ingredients
		WHERE id = ANY($1) AND is_active = TRUE
	`

	rows, err := r.db.Query(query, pq.Array(ids))
	if err != nil {
		return nil, fmt.Errorf("failed to query ingredients: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		item := &models.BasicIngredient{}
		var nameEn sql.NullString
		var category sql.NullString
		var unit sql.NullString
		var storage sql.NullInt64
		var description sql.NullString
		if err := rows.Scan(
			&item.ID,
			&item.Name,
			&nameEn,
			&category,
			&unit,
			&storage,
			&description,
		); err != nil {
			return nil, fmt.Errorf("failed to scan ingredient: %w", err)
		}

		item.ID = strings.TrimSpace(item.ID)
		item.NameEN = nullableString(nameEn)
		item.Category = nullableString(category)
		item.DefaultUnit = nullableString(unit)
		item.Description = nullableString(description)
		if storage.Valid {
			value := int(storage.Int64)
			item.StorageDays = &value
		}

		result[item.ID] = item
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate ingredients: %w", err)
	}

	return result, nil
}
