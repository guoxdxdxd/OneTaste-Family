package services

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"

	"onetaste-family/backend/internal/models"
	"onetaste-family/backend/internal/repositories"
)

var (
	// ErrUserAlreadyInFamily 用户已属于某个家庭
	ErrUserAlreadyInFamily = errors.New("user already in a family")
	// ErrInvalidFamilyName 家庭名称非法
	ErrInvalidFamilyName = errors.New("invalid family name")
	// ErrInvalidFamilyDescription 家庭描述过长
	ErrInvalidFamilyDescription = errors.New("invalid family description")
	// ErrFamilyNotFound 家庭不存在
	ErrFamilyNotFound = errors.New("family not found")
	// ErrInvalidInviteAction 扫码邀请动作非法
	ErrInvalidInviteAction = errors.New("invalid invite action")
	// ErrInviterNotInFamily 邀请人不在家庭中
	ErrInviterNotInFamily = errors.New("inviter not in family")
	// ErrFamilyMemberLimitReached 家庭成员数量已达上限
	ErrFamilyMemberLimitReached = errors.New("family member limit reached")
	// ErrFamilyNameMismatch 二维码中的家庭名称不匹配
	ErrFamilyNameMismatch = errors.New("family name mismatch")
)

const (
	defaultMaxDishes           = 30
	familyNameMaxLength        = 100
	familyDescriptionMaxLength = 500
	maxFamilyMembers           = 10
)

// FamilyService 家庭业务逻辑层
type FamilyService struct {
	familyRepo *repositories.FamilyRepository
}

// NewFamilyService 创建FamilyService
func NewFamilyService() *FamilyService {
	return &FamilyService{
		familyRepo: repositories.NewFamilyRepository(),
	}
}

// CreateFamily 创建家庭
func (s *FamilyService) CreateFamily(req *models.CreateFamilyRequest, userID int64) (*models.FamilyResponse, error) {
	name := strings.TrimSpace(req.Name)
	if name == "" || utf8.RuneCountInString(name) > familyNameMaxLength {
		return nil, ErrInvalidFamilyName
	}

	description := strings.TrimSpace(req.Description)
	if utf8.RuneCountInString(description) > familyDescriptionMaxLength {
		return nil, ErrInvalidFamilyDescription
	}

	// 一个用户只能属于一个家庭
	inFamily, err := s.familyRepo.IsUserInFamily(userID)
	if err != nil {
		return nil, fmt.Errorf("check user family membership failed: %w", err)
	}
	if inFamily {
		return nil, ErrUserAlreadyInFamily
	}

	family := &models.Family{
		Name:        name,
		Description: description,
		OwnerID:     userID,
		MaxDishes:   defaultMaxDishes,
		Status:      models.FamilyStatusActive,
	}

	ctx := context.Background()
	tx, err := s.familyRepo.BeginTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin transaction failed: %w", err)
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	if err = s.familyRepo.CreateFamilyTx(ctx, tx, family); err != nil {
		return nil, fmt.Errorf("create family failed: %w", err)
	}

	member := &models.FamilyMember{
		FamilyID: family.ID,
		UserID:   userID,
		Role:     models.FamilyRoleOwner,
		Status:   models.FamilyMemberStatusActive,
	}

	if err = s.familyRepo.AddFamilyMemberTx(ctx, tx, member); err != nil {
		return nil, fmt.Errorf("add owner to family failed: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit transaction failed: %w", err)
	}

	return &models.FamilyResponse{
		FamilyID:    family.ID,
		Name:        family.Name,
		Description: family.Description,
		MemberCount: 1,
		MaxDishes:   family.MaxDishes,
	}, nil
}

// GetFamilyInfo 获取当前用户所属家庭信息
func (s *FamilyService) GetFamilyInfo(userID int64) (*models.FamilyInfoResponse, error) {
	family, err := s.familyRepo.GetFamilyByUserID(userID)
	if err != nil {
		if errors.Is(err, repositories.ErrFamilyNotFound) {
			return nil, ErrFamilyNotFound
		}
		return nil, fmt.Errorf("failed to get family: %w", err)
	}

	memberCount, dishCount, err := s.familyRepo.GetFamilyStats(family.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get family stats: %w", err)
	}

	return &models.FamilyInfoResponse{
		FamilyID:    family.ID,
		Name:        family.Name,
		Description: family.Description,
		OwnerID:     family.OwnerID,
		MemberCount: memberCount,
		DishCount:   dishCount,
		MaxDishes:   family.MaxDishes,
	}, nil
}

// JoinFamilyViaInvite 扫码加入家庭
func (s *FamilyService) JoinFamilyViaInvite(req *models.FamilyInviteRequest, userID int64) (*models.FamilyJoinResponse, error) {
	if strings.ToLower(req.Action) != "accept" {
		return nil, ErrInvalidInviteAction
	}

	inFamily, err := s.familyRepo.IsUserInFamily(userID)
	if err != nil {
		return nil, fmt.Errorf("check user family membership failed: %w", err)
	}
	if inFamily {
		return nil, ErrUserAlreadyInFamily
	}

	family, err := s.familyRepo.GetFamilyByID(req.FamilyID)
	if err != nil {
		if errors.Is(err, repositories.ErrFamilyNotFound) {
			return nil, ErrFamilyNotFound
		}
		return nil, fmt.Errorf("failed to get family: %w", err)
	}

	if req.FamilyName != "" && req.FamilyName != family.Name {
		return nil, ErrFamilyNameMismatch
	}

	inviterValid, err := s.familyRepo.IsUserInSpecificFamily(req.InviterID, req.FamilyID)
	if err != nil {
		return nil, fmt.Errorf("failed to validate inviter: %w", err)
	}
	if !inviterValid {
		return nil, ErrInviterNotInFamily
	}

	memberCount, err := s.familyRepo.CountFamilyMembers(req.FamilyID)
	if err != nil {
		return nil, fmt.Errorf("failed to count family members: %w", err)
	}
	if memberCount >= maxFamilyMembers {
		return nil, ErrFamilyMemberLimitReached
	}

	member := &models.FamilyMember{
		FamilyID: req.FamilyID,
		UserID:   userID,
		Role:     models.FamilyRoleMember,
		Status:   models.FamilyMemberStatusActive,
	}

	if err := s.familyRepo.AddFamilyMember(member); err != nil {
		if errors.Is(err, repositories.ErrFamilyMemberExists) {
			return nil, ErrUserAlreadyInFamily
		}
		return nil, fmt.Errorf("failed to add family member: %w", err)
	}

	return &models.FamilyJoinResponse{
		FamilyID:   member.FamilyID,
		MemberRole: member.Role,
		JoinedAt:   member.JoinedAt,
	}, nil
}

// GetFamilyMembers 获取家庭成员列表
func (s *FamilyService) GetFamilyMembers(userID int64) ([]*models.FamilyMemberInfo, error) {
	family, err := s.familyRepo.GetFamilyByUserID(userID)
	if err != nil {
		if errors.Is(err, repositories.ErrFamilyNotFound) {
			return nil, ErrFamilyNotFound
		}
		return nil, fmt.Errorf("failed to get family: %w", err)
	}

	members, err := s.familyRepo.GetFamilyMembers(family.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get family members: %w", err)
	}

	return members, nil
}
