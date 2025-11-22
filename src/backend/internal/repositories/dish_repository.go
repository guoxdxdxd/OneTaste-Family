package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"onetaste-family/backend/internal/models"
	"onetaste-family/backend/pkg/database"
)

var (
	// ErrDishNotFound 菜式不存在
	ErrDishNotFound = errors.New("dish not found")
)

// DishRepository 菜式数据访问层
type DishRepository struct {
	db *sql.DB
}

// NewDishRepository 创建菜式仓储
func NewDishRepository() *DishRepository {
	return &DishRepository{
		db: database.GetDB(),
	}
}

// CountByFamily 统计家庭菜式数量
func (r *DishRepository) CountByFamily(familyID string) (int, error) {
	query := `SELECT COUNT(*) FROM dishes WHERE family_id = $1 AND deleted_at IS NULL`

	var count int
	if err := r.db.QueryRow(query, familyID).Scan(&count); err != nil {
		return 0, fmt.Errorf("failed to count dishes: %w", err)
	}

	return count, nil
}

// ExistsByName 判断菜式名称是否存在（同一家庭内）
func (r *DishRepository) ExistsByName(familyID string, name string, excludeID string) (bool, error) {
	query := `SELECT EXISTS(
		SELECT 1 FROM dishes
		WHERE family_id = $1 AND deleted_at IS NULL AND LOWER(name) = LOWER($2)
	`
	args := []interface{}{familyID, name}

	if excludeID != "" {
		query += fmt.Sprintf(" AND id <> $%d", len(args)+1)
		args = append(args, excludeID)
	}

	query += ")"

	var exists bool
	if err := r.db.QueryRow(query, args...).Scan(&exists); err != nil {
		return false, fmt.Errorf("failed to check dish name: %w", err)
	}

	return exists, nil
}

// CreateDishWithDetails 创建菜式并保存详情
func (r *DishRepository) CreateDishWithDetails(dish *models.Dish, ingredients []*models.Ingredient, steps []*models.CookingStep) error {
	ctx := context.Background()
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction failed: %w", err)
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	insertDish := `
		INSERT INTO dishes (id, family_id, name, category, description, image_url, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING created_at, updated_at
	`

	err = tx.QueryRowContext(
		ctx,
		insertDish,
		dish.ID,
		dish.FamilyID,
		dish.Name,
		nullString(dish.Category),
		nullString(dish.Description),
		nullString(dish.ImageURL),
		dish.CreatedBy,
	).Scan(&dish.CreatedAt, &dish.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to insert dish: %w", err)
	}

	if err = r.insertIngredients(ctx, tx, dish.ID, ingredients); err != nil {
		return err
	}

	if err = r.insertCookingSteps(ctx, tx, dish.ID, steps); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction failed: %w", err)
	}

	return nil
}

// UpdateDishWithDetails 更新菜式及其详情
func (r *DishRepository) UpdateDishWithDetails(dish *models.Dish, ingredients []*models.Ingredient, steps []*models.CookingStep) error {
	ctx := context.Background()
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction failed: %w", err)
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	updateDish := `
		UPDATE dishes
		SET name = $1, category = $2, description = $3, image_url = $4, updated_at = NOW()
		WHERE id = $5 AND family_id = $6 AND deleted_at IS NULL
		RETURNING updated_at
	`

	err = tx.QueryRowContext(
		ctx,
		updateDish,
		dish.Name,
		nullString(dish.Category),
		nullString(dish.Description),
		nullString(dish.ImageURL),
		dish.ID,
		dish.FamilyID,
	).Scan(&dish.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrDishNotFound
		}
		return fmt.Errorf("failed to update dish: %w", err)
	}

	if _, err = tx.ExecContext(ctx, `DELETE FROM dish_ingredients WHERE dish_id = $1`, dish.ID); err != nil {
		return fmt.Errorf("failed to delete old ingredients: %w", err)
	}

	if _, err = tx.ExecContext(ctx, `DELETE FROM cooking_steps WHERE dish_id = $1`, dish.ID); err != nil {
		return fmt.Errorf("failed to delete old cooking steps: %w", err)
	}

	if err = r.insertIngredients(ctx, tx, dish.ID, ingredients); err != nil {
		return err
	}

	if err = r.insertCookingSteps(ctx, tx, dish.ID, steps); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction failed: %w", err)
	}

	return nil
}

// GetDishByID 根据ID获取菜式
func (r *DishRepository) GetDishByID(dishID, familyID string) (*models.Dish, error) {
	query := `
		SELECT id, family_id, name, category, description, image_url, created_by, created_at, updated_at
		FROM dishes
		WHERE id = $1 AND family_id = $2 AND deleted_at IS NULL
	`

	dish := &models.Dish{}
	var category, description, image sql.NullString
	if err := r.db.QueryRow(
		query,
		dishID,
		familyID,
	).Scan(
		&dish.ID,
		&dish.FamilyID,
		&dish.Name,
		&category,
		&description,
		&image,
		&dish.CreatedBy,
		&dish.CreatedAt,
		&dish.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrDishNotFound
		}
		return nil, fmt.Errorf("failed to get dish: %w", err)
	}

	dish.Category = nullableString(category)
	dish.Description = nullableString(description)
	dish.ImageURL = nullableString(image)

	return dish, nil
}

// GetIngredients 获取菜式食材列表
func (r *DishRepository) GetIngredients(dishID string) ([]*models.Ingredient, error) {
	query := `
		SELECT
			di.id,
			di.dish_id,
			di.ingredient_id,
			bi.name,
			bi.name_en,
			bi.category,
			bi.default_unit,
			bi.storage_days,
			di.amount,
			di.unit,
			di.notes,
			di.sort_order
		FROM dish_ingredients di
		JOIN ingredients bi ON di.ingredient_id = bi.id
		WHERE di.dish_id = $1
		ORDER BY di.sort_order ASC, di.id ASC
	`

	rows, err := r.db.Query(query, dishID)
	if err != nil {
		return nil, fmt.Errorf("failed to query ingredients: %w", err)
	}
	defer rows.Close()

	var ingredients []*models.Ingredient
	for rows.Next() {
		ingredient := &models.Ingredient{}
		var nameEn sql.NullString
		var category sql.NullString
		var defaultUnit sql.NullString
		var storage sql.NullInt64
		var notes sql.NullString
		if err := rows.Scan(
			&ingredient.ID,
			&ingredient.DishID,
			&ingredient.IngredientID,
			&ingredient.IngredientName,
			&nameEn,
			&category,
			&defaultUnit,
			&storage,
			&ingredient.Amount,
			&ingredient.Unit,
			&notes,
			&ingredient.SortOrder,
		); err != nil {
			return nil, fmt.Errorf("failed to scan ingredient: %w", err)
		}

		ingredient.IngredientID = strings.TrimSpace(ingredient.IngredientID)
		ingredient.IngredientNameEn = nullableString(nameEn)
		ingredient.Category = nullableString(category)
		ingredient.DefaultUnit = nullableString(defaultUnit)
		ingredient.Notes = nullableString(notes)
		if storage.Valid {
			value := int(storage.Int64)
			ingredient.StorageDays = &value
		}

		ingredients = append(ingredients, ingredient)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate ingredients: %w", err)
	}

	return ingredients, nil
}

// GetCookingSteps 获取烹饪步骤
func (r *DishRepository) GetCookingSteps(dishID string) ([]*models.CookingStep, error) {
	query := `
		SELECT id, dish_id, step_order, content, image_url
		FROM cooking_steps
		WHERE dish_id = $1
		ORDER BY step_order ASC, id ASC
	`

	rows, err := r.db.Query(query, dishID)
	if err != nil {
		return nil, fmt.Errorf("failed to query cooking steps: %w", err)
	}
	defer rows.Close()

	var steps []*models.CookingStep
	for rows.Next() {
		step := &models.CookingStep{}
		var image sql.NullString
		if err := rows.Scan(
			&step.ID,
			&step.DishID,
			&step.Order,
			&step.Content,
			&image,
		); err != nil {
			return nil, fmt.Errorf("failed to scan cooking step: %w", err)
		}

		step.ImageURL = nullableString(image)
		steps = append(steps, step)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate cooking steps: %w", err)
	}

	return steps, nil
}

// GetDishList 获取菜式列表
func (r *DishRepository) GetDishList(familyID string, page, pageSize int, category, keyword string) ([]*models.DishSummary, int64, error) {
	var whereBuilder strings.Builder
	whereBuilder.WriteString("WHERE family_id = $1 AND deleted_at IS NULL")

	args := []interface{}{familyID}
	placeholder := 2

	if category != "" {
		whereBuilder.WriteString(fmt.Sprintf(" AND category = $%d", placeholder))
		args = append(args, category)
		placeholder++
	}

	if keyword != "" {
		whereBuilder.WriteString(fmt.Sprintf(" AND name ILIKE $%d", placeholder))
		args = append(args, "%"+keyword+"%")
		placeholder++
	}

	whereClause := whereBuilder.String()

	countQuery := "SELECT COUNT(*) FROM dishes " + whereClause
	var total int64
	if err := r.db.QueryRow(countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count dishes: %w", err)
	}

	limitPlaceholder := placeholder
	offsetPlaceholder := placeholder + 1

	listQuery := fmt.Sprintf(`
		SELECT id, name, category, description, image_url, created_at, updated_at
		FROM dishes
		%s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, limitPlaceholder, offsetPlaceholder)

	offset := (page - 1) * pageSize
	dataArgs := append(append([]interface{}{}, args...), pageSize, offset)

	rows, err := r.db.Query(listQuery, dataArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query dishes: %w", err)
	}
	defer rows.Close()

	var dishes []*models.DishSummary
	for rows.Next() {
		item := &models.DishSummary{}
		var category sql.NullString
		var description sql.NullString
		var image sql.NullString
		if err := rows.Scan(
			&item.DishID,
			&item.Name,
			&category,
			&description,
			&image,
			&item.CreatedAt,
			&item.UpdatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("failed to scan dish: %w", err)
		}

		item.Category = nullableString(category)
		item.Description = nullableString(description)
		item.ImageURL = nullableString(image)
		dishes = append(dishes, item)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("failed to iterate dishes: %w", err)
	}

	return dishes, total, nil
}

// SoftDeleteDish 软删除菜式
func (r *DishRepository) SoftDeleteDish(dishID, familyID string) error {
	ctx := context.Background()
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction failed: %w", err)
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	res, execErr := tx.ExecContext(ctx, `UPDATE dishes SET deleted_at = NOW() WHERE id = $1 AND family_id = $2 AND deleted_at IS NULL`, dishID, familyID)
	if execErr != nil {
		err = fmt.Errorf("failed to delete dish: %w", execErr)
		return err
	}

	rowsAffected, execErr := res.RowsAffected()
	if execErr != nil {
		err = fmt.Errorf("failed to fetch affected rows: %w", execErr)
		return err
	}
	if rowsAffected == 0 {
		err = ErrDishNotFound
		return err
	}

	if _, err = tx.ExecContext(ctx, `DELETE FROM dish_ingredients WHERE dish_id = $1`, dishID); err != nil {
		return fmt.Errorf("failed to delete ingredients: %w", err)
	}

	if _, err = tx.ExecContext(ctx, `DELETE FROM cooking_steps WHERE dish_id = $1`, dishID); err != nil {
		return fmt.Errorf("failed to delete cooking steps: %w", err)
	}

	if _, err = tx.ExecContext(ctx, `DELETE FROM menu_dishes WHERE dish_id = $1`, dishID); err != nil {
		return fmt.Errorf("failed to delete menu relation: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction failed: %w", err)
	}

	return nil
}

func (r *DishRepository) insertIngredients(ctx context.Context, tx *sql.Tx, dishID string, ingredients []*models.Ingredient) error {
	if len(ingredients) == 0 {
		return nil
	}

	query := `
		INSERT INTO dish_ingredients (id, dish_id, ingredient_id, amount, unit, notes, sort_order)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	for _, ingredient := range ingredients {
		ingredient.DishID = dishID

		if _, err := tx.ExecContext(
			ctx,
			query,
			ingredient.ID,
			ingredient.DishID,
			ingredient.IngredientID,
			ingredient.Amount,
			ingredient.Unit,
			nullString(ingredient.Notes),
			ingredient.SortOrder,
		); err != nil {
			return fmt.Errorf("failed to insert ingredient: %w", err)
		}
	}

	return nil
}

func (r *DishRepository) insertCookingSteps(ctx context.Context, tx *sql.Tx, dishID string, steps []*models.CookingStep) error {
	if len(steps) == 0 {
		return nil
	}

	query := `
		INSERT INTO cooking_steps (id, dish_id, step_order, content, image_url)
		VALUES ($1, $2, $3, $4, $5)
	`

	for _, step := range steps {
		step.DishID = dishID
		if _, err := tx.ExecContext(
			ctx,
			query,
			step.ID,
			step.DishID,
			step.Order,
			step.Content,
			nullString(step.ImageURL),
		); err != nil {
			return fmt.Errorf("failed to insert cooking step: %w", err)
		}
	}

	return nil
}

func nullableString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}

func nullString(value string) interface{} {
	if strings.TrimSpace(value) == "" {
		return nil
	}
	return value
}
