package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"onetaste-family/backend/internal/models"
	"onetaste-family/backend/pkg/database"
)

var (
	// ErrFamilyNotFound 家庭不存在
	ErrFamilyNotFound = errors.New("family not found")
	// ErrFamilyMemberExists 成员已存在
	ErrFamilyMemberExists = errors.New("family member already exists")
)

// FamilyRepository 家庭数据访问层
type FamilyRepository struct {
	db *sql.DB
}

// NewFamilyRepository 创建家庭仓储
func NewFamilyRepository() *FamilyRepository {
	return &FamilyRepository{
		db: database.GetDB(),
	}
}

// BeginTx 开启事务
func (r *FamilyRepository) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return r.db.BeginTx(ctx, nil)
}

// IsUserInFamily 判断用户是否在家庭中
func (r *FamilyRepository) IsUserInFamily(userID int64) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM family_members WHERE user_id = $1 AND status = $2)`

	var exists bool
	if err := r.db.QueryRow(query, userID, models.FamilyMemberStatusActive).Scan(&exists); err != nil {
		return false, fmt.Errorf("failed to check user family membership: %w", err)
	}

	return exists, nil
}

// GetFamilyByUserID 根据用户ID查询家庭
func (r *FamilyRepository) GetFamilyByUserID(userID int64) (*models.Family, error) {
	query := `
		SELECT f.id, f.name, f.description, f.owner_id, f.max_dishes, f.status, f.created_at, f.updated_at
		FROM families f
		INNER JOIN family_members fm ON fm.family_id = f.id
		WHERE fm.user_id = $1 AND fm.status = $2 AND f.status = $3
		LIMIT 1
	`

	family := &models.Family{}
	err := r.db.QueryRow(
		query,
		userID,
		models.FamilyMemberStatusActive,
		models.FamilyStatusActive,
	).Scan(
		&family.ID,
		&family.Name,
		&family.Description,
		&family.OwnerID,
		&family.MaxDishes,
		&family.Status,
		&family.CreatedAt,
		&family.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrFamilyNotFound
		}
		return nil, fmt.Errorf("failed to get family by user id: %w", err)
	}

	return family, nil
}

// GetFamilyByID 根据家庭ID获取家庭信息
func (r *FamilyRepository) GetFamilyByID(familyID int64) (*models.Family, error) {
	query := `
		SELECT id, name, description, owner_id, max_dishes, status, created_at, updated_at
		FROM families
		WHERE id = $1 AND status = $2
	`

	family := &models.Family{}
	err := r.db.QueryRow(query, familyID, models.FamilyStatusActive).Scan(
		&family.ID,
		&family.Name,
		&family.Description,
		&family.OwnerID,
		&family.MaxDishes,
		&family.Status,
		&family.CreatedAt,
		&family.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrFamilyNotFound
		}
		return nil, fmt.Errorf("failed to get family by id: %w", err)
	}

	return family, nil
}

// IsUserInSpecificFamily 判断用户是否属于指定家庭
func (r *FamilyRepository) IsUserInSpecificFamily(userID, familyID int64) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM family_members
			WHERE user_id = $1 AND family_id = $2 AND status = $3
		)
	`

	var exists bool
	if err := r.db.QueryRow(query, userID, familyID, models.FamilyMemberStatusActive).Scan(&exists); err != nil {
		return false, fmt.Errorf("failed to check user in family: %w", err)
	}

	return exists, nil
}

// CountFamilyMembers 统计家庭成员数量
func (r *FamilyRepository) CountFamilyMembers(familyID int64) (int, error) {
	query := `SELECT COUNT(*) FROM family_members WHERE family_id = $1 AND status = $2`

	var count int
	if err := r.db.QueryRow(query, familyID, models.FamilyMemberStatusActive).Scan(&count); err != nil {
		return 0, fmt.Errorf("failed to count family members: %w", err)
	}

	return count, nil
}

// GetFamilyStats 获取家庭统计数据：成员数量、菜式数量
func (r *FamilyRepository) GetFamilyStats(familyID int64) (memberCount int, dishCount int, err error) {
	memberQuery := `SELECT COUNT(*) FROM family_members WHERE family_id = $1 AND status = $2`
	if err = r.db.QueryRow(memberQuery, familyID, models.FamilyMemberStatusActive).Scan(&memberCount); err != nil {
		return 0, 0, fmt.Errorf("failed to count family members: %w", err)
	}

	dishQuery := `SELECT COUNT(*) FROM dishes WHERE family_id = $1 AND deleted_at IS NULL`
	if err = r.db.QueryRow(dishQuery, familyID).Scan(&dishCount); err != nil {
		return 0, 0, fmt.Errorf("failed to count dishes: %w", err)
	}

	return memberCount, dishCount, nil
}

// CreateFamilyTx 在事务内创建家庭
func (r *FamilyRepository) CreateFamilyTx(ctx context.Context, tx *sql.Tx, family *models.Family) error {
	query := `
		INSERT INTO families (name, description, owner_id, max_dishes, status)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`

	return tx.QueryRowContext(
		ctx,
		query,
		family.Name,
		family.Description,
		family.OwnerID,
		family.MaxDishes,
		family.Status,
	).Scan(&family.ID, &family.CreatedAt, &family.UpdatedAt)
}

// AddFamilyMemberTx 在事务内添加成员
func (r *FamilyRepository) AddFamilyMemberTx(ctx context.Context, tx *sql.Tx, member *models.FamilyMember) error {
	query := `
		INSERT INTO family_members (family_id, user_id, role, status)
		VALUES ($1, $2, $3, $4)
		RETURNING id, joined_at, created_at, updated_at
	`

	return tx.QueryRowContext(
		ctx,
		query,
		member.FamilyID,
		member.UserID,
		member.Role,
		member.Status,
	).Scan(&member.ID, &member.JoinedAt, &member.CreatedAt, &member.UpdatedAt)
}

// AddFamilyMember 添加家庭成员（非事务）
func (r *FamilyRepository) AddFamilyMember(member *models.FamilyMember) error {
	query := `
		INSERT INTO family_members (family_id, user_id, role, status)
		VALUES ($1, $2, $3, $4)
		RETURNING id, joined_at, created_at, updated_at
	`

	err := r.db.QueryRow(
		query,
		member.FamilyID,
		member.UserID,
		member.Role,
		member.Status,
	).Scan(&member.ID, &member.JoinedAt, &member.CreatedAt, &member.UpdatedAt)

	if err != nil {
		if err.Error() == "pq: duplicate key value violates unique constraint \"family_members_family_id_user_id_key\"" {
			return ErrFamilyMemberExists
		}
		return fmt.Errorf("failed to add family member: %w", err)
	}

	return nil
}

// GetFamilyMembers 获取家庭成员列表
func (r *FamilyRepository) GetFamilyMembers(familyID int64) ([]*models.FamilyMemberInfo, error) {
	query := `
		SELECT fm.user_id, COALESCE(u.nickname, ''), COALESCE(u.avatar, ''), fm.role, fm.joined_at
		FROM family_members fm
		INNER JOIN users u ON u.id = fm.user_id
		WHERE fm.family_id = $1 AND fm.status = $2
		ORDER BY CASE WHEN fm.role = $3 THEN 0 ELSE 1 END, fm.joined_at ASC
	`

	rows, err := r.db.Query(query, familyID, models.FamilyMemberStatusActive, models.FamilyRoleOwner)
	if err != nil {
		return nil, fmt.Errorf("failed to query family members: %w", err)
	}
	defer rows.Close()

	var members []*models.FamilyMemberInfo
	for rows.Next() {
		member := &models.FamilyMemberInfo{}
		if err := rows.Scan(&member.UserID, &member.Nickname, &member.Avatar, &member.Role, &member.JoinedAt); err != nil {
			return nil, fmt.Errorf("failed to scan family member: %w", err)
		}
		members = append(members, member)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate family members: %w", err)
	}

	return members, nil
}
