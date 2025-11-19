package models

import "time"

const (
	// FamilyStatusActive 家庭正常状态
	FamilyStatusActive = 1
	// FamilyStatusDisabled 家庭已解散/禁用
	FamilyStatusDisabled = 0
)

const (
	// FamilyMemberStatusActive 成员有效状态
	FamilyMemberStatusActive = 1
	// FamilyMemberStatusInactive 成员离开家庭
	FamilyMemberStatusInactive = 0
)

const (
	// FamilyRoleOwner 家庭管理员角色
	FamilyRoleOwner = "owner"
	// FamilyRoleMember 家庭普通成员角色
	FamilyRoleMember = "member"
)

// Family 家庭模型
type Family struct {
	ID          int64     `json:"family_id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description,omitempty" db:"description"`
	OwnerID     int64     `json:"owner_id" db:"owner_id"`
	MaxDishes   int       `json:"max_dishes" db:"max_dishes"`
	Status      int       `json:"status" db:"status"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// FamilyMember 家庭成员模型
type FamilyMember struct {
	ID        int64     `json:"id" db:"id"`
	FamilyID  int64     `json:"family_id" db:"family_id"`
	UserID    int64     `json:"user_id" db:"user_id"`
	Role      string    `json:"role" db:"role"`
	Status    int       `json:"status" db:"status"`
	JoinedAt  time.Time `json:"joined_at" db:"joined_at"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// CreateFamilyRequest 创建家庭请求
type CreateFamilyRequest struct {
	Name        string `json:"name" binding:"required,max=100" example:"张家的厨房"` // 家庭名称
	Description string `json:"description" binding:"max=500" example:"温馨的家庭"`   // 家庭描述，可选
}

// FamilyResponse 家庭创建响应
type FamilyResponse struct {
	FamilyID    int64  `json:"family_id" example:"1"`
	Name        string `json:"name" example:"张家的厨房"`
	Description string `json:"description,omitempty" example:"温馨的家庭"`
	MemberCount int    `json:"member_count" example:"1"`
	MaxDishes   int    `json:"max_dishes" example:"30"`
}

// FamilyInfoResponse 获取家庭信息响应
type FamilyInfoResponse struct {
	FamilyID    int64  `json:"family_id" example:"1"`
	Name        string `json:"name" example:"张家的厨房"`
	Description string `json:"description,omitempty" example:"温馨的家庭"`
	OwnerID     int64  `json:"owner_id" example:"1001"`
	MemberCount int    `json:"member_count" example:"3"`
	DishCount   int    `json:"dish_count" example:"15"`
	MaxDishes   int    `json:"max_dishes" example:"30"`
}

// FamilyInviteRequest 扫码加入家庭请求
type FamilyInviteRequest struct {
	FamilyID        int64  `json:"family_id" binding:"required" example:"1"`
	FamilyName      string `json:"family_name" binding:"required" example:"张家的厨房"`
	InviterID       int64  `json:"inviter_id" binding:"required" example:"1001"`
	InviterNickname string `json:"inviter_nickname" binding:"required" example:"张三"`
	Action          string `json:"action" binding:"required" example:"accept"`
}

// FamilyJoinResponse 扫码加入家庭响应
type FamilyJoinResponse struct {
	FamilyID   int64     `json:"family_id" example:"1"`
	MemberRole string    `json:"member_role" example:"member"`
	JoinedAt   time.Time `json:"joined_at" example:"2024-01-01T00:00:00Z"`
}

// FamilyMemberInfo 家庭成员列表项
type FamilyMemberInfo struct {
	UserID   int64     `json:"user_id" example:"1002"`
	Nickname string    `json:"nickname" example:"李四"`
	Avatar   string    `json:"avatar,omitempty" example:"https://.../avatar.jpg"`
	Role     string    `json:"role" example:"member"`
	JoinedAt time.Time `json:"joined_at" example:"2024-01-01T00:00:00Z"`
}
