package user_service

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService interface {
	RegisterUser(ctx context.Context, user *User) error
	GetUser(ctx context.Context, id uuid.UUID) (*User, error)
	LoginUser(ctx context.Context, email, password string) (*string, error)
}

type userService struct {
	repo UserRepository
}

type CustomClaims struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
	jwt.StandardClaims
}

func NewUserService(repo UserRepository) UserService {
	return &userService{repo: repo}
}

// GetUser implements UserService.
func (s *userService) GetUser(ctx context.Context, id uuid.UUID) (*User, error) {
	return s.repo.GetUserByID(ctx, id)
}

// LoginUser implements UserService.
func (s *userService) LoginUser(ctx context.Context, email string, password string) (*string, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("User Not Found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("Invalid Password")
	}

	token, err := generateJWT(user.UserID, user.Email)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	return &token, nil
}

// RegisterUser implements UserService.
func (s *userService) RegisterUser(ctx context.Context, user *User) error {
	userDB, err := s.repo.GetUserByEmail(ctx, user.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("database error")
	}

	if userDB != nil && userDB.Email == user.Email {
		return errors.New("Email already exist")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to hash password")
	}

	user.Password = string(hashedPassword)
	user.UserID = uuid.New()

	return s.repo.CreateUser(ctx, user)
}

// utils for generate JsonWebToken
func generateJWT(userID uuid.UUID, email string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &CustomClaims{
		UserID: userID,
		Email:  email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte("your_secret_key"))
}
