package repository

import (
	"fmt"

	"github.com/ezep02/microservicios/internal/auth/types"
	"gorm.io/gorm"
)

type AuthRepository struct {
	db *gorm.DB
}

func NewAuthRepository(DBC *gorm.DB) *AuthRepository {
	return &AuthRepository{
		db: DBC,
	}
}

func (AuthR *AuthRepository) UserSignUp(user *types.User) (*types.UserResponse, error) {
	if err := AuthR.db.Create(user).Error; err != nil {
		return nil, err
	}

	// Retornar el usuario creado
	return &types.UserResponse{
		Model:   user.Model,
		Name:    user.Name,
		Age:     user.Age,
		Email:   user.Email,
		Admin:   user.Admin,
		Is_user: user.Is_user,
	}, nil
} // Register

func (AuthR *AuthRepository) UserSignIn(user *types.UserRequest) (*types.UserSignInResponse, error) {

	//verificar si el usuario que
	var u types.User
	// result := r.db.Where("id = ?", id).First(&user)

	result := AuthR.db.Where("email = ?", user.Email).First(&u)

	if result.Error != nil {

		// Imprimir el error para depuración
		fmt.Printf("Error fetching user: %v\n", gorm.ErrRecordNotFound)
		return nil, gorm.ErrRecordNotFound
	}

	return &types.UserSignInResponse{
		Model:    u.Model,
		Name:     u.Name,
		Age:      u.Age,
		Password: u.Password,
		Email:    u.Email,
		Admin:    u.Admin,
		Is_user:  u.Is_user,
	}, nil
} // Login

func (AuthR *AuthRepository) CreateSession(ses *types.Session) (*types.Session, error) {

	if err := AuthR.db.Create(ses).Error; err != nil {
		return nil, err
	}

	return ses, nil
}

func (AuthR *AuthRepository) GetSession(id string) (*types.Session, error) {

	var u types.Session

	result := AuthR.db.Where("id = ?", id).First(&u)

	if result.Error != nil {

		// Imprimir el error para depuración
		fmt.Printf("Error fetching user: %v\n", gorm.ErrRecordNotFound)
		return nil, gorm.ErrRecordNotFound
	}

	return &u, nil
} // Obtiene la session del usuario

func (AuthR *AuthRepository) RevokeSession(id string) error {

	// Actualizar la columna is_revoked a true donde id coincide
	result := AuthR.db.Model(&types.Session{}).Where("id = ?", id).Updates(map[string]interface{}{"is_revoked": true})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil

} // Revoca la session del usuario

func (AuthR *AuthRepository) DeleteSession(id string) error {
	var s types.Session

	result := AuthR.db.Where("id = ?", id).Delete(&s)

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
} // Elimina la session del usuario
