package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"onetaste-family/backend/internal/models"
	"onetaste-family/backend/internal/utils"
	"onetaste-family/backend/pkg/database"
)

var (
	// ErrMenuNotFound 菜单不存在
	ErrMenuNotFound = errors.New("menu not found")
)

// MenuRepository 菜单数据访问层
type MenuRepository struct {
	db *sql.DB
}

// NewMenuRepository 创建菜单仓储
func NewMenuRepository() *MenuRepository {
	return &MenuRepository{
		db: database.GetDB(),
	}
}

// GetMenuByDateAndMealType 根据日期和餐次获取菜单
func (r *MenuRepository) GetMenuByDateAndMealType(familyID string, date time.Time, mealType string) (*models.Menu, error) {
	query := `
		SELECT id, family_id, date, meal_type, created_by, source, created_at, updated_at
		FROM menus
		WHERE family_id = $1 AND date = $2 AND meal_type = $3
	`

	menu := &models.Menu{}
	var source sql.NullString
	if err := r.db.QueryRow(query, familyID, date, mealType).Scan(
		&menu.ID,
		&menu.FamilyID,
		&menu.Date,
		&menu.MealType,
		&menu.CreatedBy,
		&source,
		&menu.CreatedAt,
		&menu.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrMenuNotFound
		}
		return nil, fmt.Errorf("failed to get menu: %w", err)
	}

	menu.Source = nullableString(source)
	if menu.Source == "" {
		menu.Source = models.MenuSourceManual
	}

	return menu, nil
}

// CreateMenuWithDishes 创建菜单并关联菜式
func (r *MenuRepository) CreateMenuWithDishes(menu *models.Menu, dishIDs []string) error {
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

	// 插入菜单
	insertMenu := `
		INSERT INTO menus (id, family_id, date, meal_type, created_by, source)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING created_at, updated_at
	`

	err = tx.QueryRowContext(
		ctx,
		insertMenu,
		menu.ID,
		menu.FamilyID,
		menu.Date,
		menu.MealType,
		menu.CreatedBy,
		nullString(menu.Source),
	).Scan(&menu.CreatedAt, &menu.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to insert menu: %w", err)
	}

	// 插入菜单菜式关联
	if err = r.insertMenuDishes(ctx, tx, menu.ID, dishIDs); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction failed: %w", err)
	}

	return nil
}

// UpdateMenuWithDishes 更新菜单并关联菜式
func (r *MenuRepository) UpdateMenuWithDishes(menu *models.Menu, dishIDs []string) error {
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

	// 更新菜单
	updateMenu := `
		UPDATE menus
		SET date = $1, meal_type = $2, source = $3, updated_at = NOW()
		WHERE id = $4 AND family_id = $5
		RETURNING updated_at
	`

	err = tx.QueryRowContext(
		ctx,
		updateMenu,
		menu.Date,
		menu.MealType,
		nullString(menu.Source),
		menu.ID,
		menu.FamilyID,
	).Scan(&menu.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrMenuNotFound
		}
		return fmt.Errorf("failed to update menu: %w", err)
	}

	// 删除旧的菜单菜式关联
	if _, err = tx.ExecContext(ctx, `DELETE FROM menu_dishes WHERE menu_id = $1`, menu.ID); err != nil {
		return fmt.Errorf("failed to delete old menu dishes: %w", err)
	}

	// 插入新的菜单菜式关联
	if err = r.insertMenuDishes(ctx, tx, menu.ID, dishIDs); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction failed: %w", err)
	}

	return nil
}

// GetMenuByID 根据ID获取菜单
func (r *MenuRepository) GetMenuByID(menuID, familyID string) (*models.Menu, error) {
	query := `
		SELECT id, family_id, date, meal_type, created_by, source, created_at, updated_at
		FROM menus
		WHERE id = $1 AND family_id = $2
	`

	menu := &models.Menu{}
	var source sql.NullString
	if err := r.db.QueryRow(query, menuID, familyID).Scan(
		&menu.ID,
		&menu.FamilyID,
		&menu.Date,
		&menu.MealType,
		&menu.CreatedBy,
		&source,
		&menu.CreatedAt,
		&menu.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrMenuNotFound
		}
		return nil, fmt.Errorf("failed to get menu: %w", err)
	}

	menu.Source = nullableString(source)
	if menu.Source == "" {
		menu.Source = models.MenuSourceManual
	}

	return menu, nil
}

// GetMenuDishes 获取菜单的菜式列表
func (r *MenuRepository) GetMenuDishes(menuID string) ([]string, error) {
	query := `
		SELECT dish_id
		FROM menu_dishes
		WHERE menu_id = $1
		ORDER BY created_at ASC
	`

	rows, err := r.db.Query(query, menuID)
	if err != nil {
		return nil, fmt.Errorf("failed to query menu dishes: %w", err)
	}
	defer rows.Close()

	var dishIDs []string
	for rows.Next() {
		var dishID string
		if err := rows.Scan(&dishID); err != nil {
			return nil, fmt.Errorf("failed to scan dish id: %w", err)
		}
		dishIDs = append(dishIDs, dishID)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate menu dishes: %w", err)
	}

	return dishIDs, nil
}

// GetMenusByDateRange 根据日期范围获取菜单列表
func (r *MenuRepository) GetMenusByDateRange(familyID string, startDate, endDate time.Time) ([]*models.Menu, error) {
	query := `
		SELECT id, family_id, date, meal_type, created_by, source, created_at, updated_at
		FROM menus
		WHERE family_id = $1 AND date >= $2 AND date <= $3
		ORDER BY date ASC, meal_type ASC
	`

	rows, err := r.db.Query(query, familyID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to query menus: %w", err)
	}
	defer rows.Close()

	var menus []*models.Menu
	for rows.Next() {
		menu := &models.Menu{}
		var source sql.NullString
		if err := rows.Scan(
			&menu.ID,
			&menu.FamilyID,
			&menu.Date,
			&menu.MealType,
			&menu.CreatedBy,
			&source,
			&menu.CreatedAt,
			&menu.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan menu: %w", err)
		}

		menu.Source = nullableString(source)
		if menu.Source == "" {
			menu.Source = models.MenuSourceManual
		}

		menus = append(menus, menu)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate menus: %w", err)
	}

	return menus, nil
}

// DeleteMenu 删除菜单
func (r *MenuRepository) DeleteMenu(menuID, familyID string) error {
	query := `DELETE FROM menus WHERE id = $1 AND family_id = $2`

	result, err := r.db.Exec(query, menuID, familyID)
	if err != nil {
		return fmt.Errorf("failed to delete menu: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return ErrMenuNotFound
	}

	return nil
}

// insertMenuDishes 插入菜单菜式关联
func (r *MenuRepository) insertMenuDishes(ctx context.Context, tx *sql.Tx, menuID string, dishIDs []string) error {
	if len(dishIDs) == 0 {
		return nil
	}

	query := `
		INSERT INTO menu_dishes (id, menu_id, dish_id)
		VALUES ($1, $2, $3)
		ON CONFLICT (menu_id, dish_id) DO NOTHING
	`

	for _, dishID := range dishIDs {
		menuDishID := utils.GenerateULID()
		if _, err := tx.ExecContext(ctx, query, menuDishID, menuID, dishID); err != nil {
			return fmt.Errorf("failed to insert menu dish: %w", err)
		}
	}

	return nil
}
