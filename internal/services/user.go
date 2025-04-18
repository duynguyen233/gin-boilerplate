package services

import (
	"context"
	"time"
	"wine-be/internal/model"
	"wine-be/internal/repositories"
	"wine-be/pkg/utils/errs"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetUserByID(ctx context.Context, id int) (model.User, error) {
	return s.repo.GetUserByID(ctx, id)
}

func (s *UserService) GetUsers(ctx context.Context) ([]model.User, error) {
	return s.repo.GetUsers(ctx)
}

func (s *UserService) UserLogin(ctx context.Context, username, password string) (model.User, error) {
	user, err := s.repo.GetUserByUsername(ctx, username)
	if err != nil {
		return model.User{}, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return model.User{}, errs.BadRequestError{Message: "Invalid username or password"}
	}
	return user, nil
}

func (s *UserService) GenerateToken(userID int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})
	return token.SignedString([]byte("secret"))
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func (s *UserService) CreateUser(ctx context.Context, user model.User) (model.User, error) {
	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		return model.User{}, err
	}
	user.Password = hashedPassword
	return s.repo.CreateUser(ctx, user)
}

func (s *UserService) UpdateUser(ctx context.Context, user model.User) (model.User, error) {
	if user.Password != "" {
		hashedPassword, err := hashPassword(user.Password)
		if err != nil {
			return model.User{}, err
		}
		user.Password = hashedPassword
	}
	return s.repo.UpdateUser(ctx, user)
}

func (s *UserService) DeleteUser(ctx context.Context, id string) error {
	return s.repo.DeleteUser(ctx, id)
}

func (s *UserService) GetUserByToken(ctx context.Context, token string) (model.User, error) {
	claims, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	if err != nil {
		return model.User{}, err
	}
	userID := claims.Claims.(jwt.MapClaims)["user_id"].(float64)
	return s.repo.GetUserByID(ctx, int(userID))
}
