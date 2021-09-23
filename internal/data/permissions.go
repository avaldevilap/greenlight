package data

import (
	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Permission struct {
	ID   uuid.UUID      `json:"id" gorm:"default:gen_random_uuid()"`
	Code pq.StringArray `json:"code" gorm:"type:text[]"`
}

func (p Permission) Include(code string) bool {
	for i := range p.Code {
		if code == p.Code[i] {
			return true
		}
	}
	return false
}

type PermissionModel struct {
	DB *gorm.DB
}

func (m *PermissionModel) GetAllForUser(userID uint) (*Permission, error) {
	var permissions *Permission
	if err := m.DB.Joins(
		"INNER JOIN users_permissions ON users_permissions.permission_id = permissions.id",
	).Joins(
		"INNER JOIN users ON users_permissions.user_id = users.id",
	).Where(
		"users.id = ?", userID,
	).Find(&permissions).Error; err != nil {
		return nil, err
	}

	return permissions, nil
}

func (m *PermissionModel) AddForUser(userID uint, codes ...string) error {
	var user User
	if err := m.DB.Find(&user, userID).Association("Permissions").Replace(pq.Array(codes)); err != nil {
		return err
	}

	return m.DB.Save(user).Error
}
