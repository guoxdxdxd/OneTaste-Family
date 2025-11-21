package repositories

import (
	"database/sql"
	"fmt"
	"strings"

	"onetaste-family/backend/internal/models"
)

// SearchActiveIngredients 模糊搜索启用食材
func (r *IngredientRepository) SearchActiveIngredients(keyword string, limit int) ([]*models.IngredientSearchResult, error) {
	keyword = strings.TrimSpace(keyword)
	if keyword == "" {
		return []*models.IngredientSearchResult{}, nil
	}

	query := `
        SELECT id, name, category, default_unit, storage_days
        FROM ingredients
        WHERE is_active = TRUE AND (name ILIKE $1 OR name_en ILIKE $1)
        ORDER BY name ASC
        LIMIT $2
    `

	rows, err := r.db.Query(query, "%"+keyword+"%", limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query ingredients: %w", err)
	}
	defer rows.Close()

	var results []*models.IngredientSearchResult
	for rows.Next() {
		item := &models.IngredientSearchResult{}
		var category sql.NullString
		var unit sql.NullString
		var storage sql.NullInt64
		if err := rows.Scan(&item.IngredientID, &item.Name, &category, &unit, &storage); err != nil {
			return nil, fmt.Errorf("failed to scan ingredient: %w", err)
		}

		item.IngredientID = strings.TrimSpace(item.IngredientID)
		item.Category = nullableString(category)
		item.DefaultUnit = nullableString(unit)
		if storage.Valid {
			value := int(storage.Int64)
			item.StorageDays = &value
		}

		results = append(results, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate ingredients: %w", err)
	}

	return results, nil
}

// GetActiveByCategory 分类分页查询
func (r *IngredientRepository) GetActiveByCategory(category, keyword string, page, pageSize int) ([]*models.IngredientSearchResult, int64, error) {
	args := []interface{}{category}
	where := "WHERE is_active = TRUE AND category = $1"
	placeholder := 2

	if strings.TrimSpace(keyword) != "" {
		where += fmt.Sprintf(" AND (name ILIKE $%d OR name_en ILIKE $%d)", placeholder, placeholder)
		args = append(args, "%"+keyword+"%")
		placeholder++
	}

	countQuery := "SELECT COUNT(*) FROM ingredients " + where
	var total int64
	if err := r.db.QueryRow(countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count ingredients: %w", err)
	}

	limitPlaceholder := placeholder
	offsetPlaceholder := placeholder + 1

	listQuery := fmt.Sprintf(`
        SELECT id, name, category, default_unit, storage_days
        FROM ingredients
        %s
        ORDER BY name ASC
        LIMIT $%d OFFSET $%d
    `, where, limitPlaceholder, offsetPlaceholder)

	offset := (page - 1) * pageSize
	dataArgs := append(append([]interface{}{}, args...), pageSize, offset)

	rows, err := r.db.Query(listQuery, dataArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query ingredients: %w", err)
	}
	defer rows.Close()

	var items []*models.IngredientSearchResult
	for rows.Next() {
		item := &models.IngredientSearchResult{}
		var category sql.NullString
		var unit sql.NullString
		var storage sql.NullInt64
		if err := rows.Scan(&item.IngredientID, &item.Name, &category, &unit, &storage); err != nil {
			return nil, 0, fmt.Errorf("failed to scan ingredient: %w", err)
		}

		item.IngredientID = strings.TrimSpace(item.IngredientID)
		item.Category = nullableString(category)
		item.DefaultUnit = nullableString(unit)
		if storage.Valid {
			value := int(storage.Int64)
			item.StorageDays = &value
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("failed to iterate ingredients: %w", err)
	}

	return items, total, nil
}
