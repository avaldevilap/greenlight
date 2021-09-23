package data

import (
	"crypto/sha256"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var AnonymousUser = &User{}

type User struct {
	gorm.Model
	UUID        uuid.UUID    `json:"uuid" gorm:"type:uuid"`
	Name        string       `json:"name" validate:"required"`
	Email       string       `json:"email" gorm:"unique" validate:"required,email"`
	Password    string       `json:"-" validate:"required,min=8,max=72"`
	Activated   bool         `json:"activated"`
	Permissions []Permission `json:"permissions" gorm:"many2many:users_permissions"`
}

func (u *User) Validate() error {
	validate := validator.New()
	return validate.Struct(u)
}

func (u *User) IsAnonymous() bool {
	return u == AnonymousUser
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.UUID = uuid.New()
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), 12)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	return nil
}

func (u *User) CheckPassword() bool {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), 12)
	if err != nil {
		return false
	}
	return bcrypt.CompareHashAndPassword(hash, []byte(u.Password)) == nil
}

type UserModel struct {
	DB *gorm.DB
}

func (m *UserModel) Insert(user *User) error {
	return m.DB.Create(user).Error
}

func (m *UserModel) GetByEmail(email string) (*User, error) {
	var user User
	err := m.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (m *UserModel) Update(user *User) error {
	return m.DB.Save(user).Error
}

func (m *UserModel) Delete(id int64) error {
	return m.DB.Delete(&User{}, id).Error
}

func (m *UserModel) GetAll() ([]User, error) {
	var users []User
	err := m.DB.Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (m *UserModel) GetForToken(tokenScope, tokenPlaintext string) (*User, error) {
	tokenHash := sha256.Sum256([]byte(tokenPlaintext))
	var user User
	err := m.DB.Joins(
		"INNER JOIN tokens ON users.id = tokens.user_id",
	).Where(
		"tokens.hash = ?", tokenHash[:],
	).Where(
		"tokens.scope = ?", tokenScope,
	).Where(
		"tokens.expiry > ?", time.Now(),
	).Find(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
