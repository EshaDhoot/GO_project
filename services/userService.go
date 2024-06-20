package services

import (
	"context"
	"go_project/models"
	"go_project/repository"
	"log"
)

type UserService struct {
	UserRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{
		UserRepo: userRepo,
	}
}

func (s *UserService) CreateUser(ctx context.Context, user *models.User) error {
	err := s.UserRepo.InsertData(ctx, user)
	if err != nil {
		log.Printf("UserService: error inserting user data: %v", err)
		return err
	}

	log.Println("UserService: user created successfully")
	return nil
}

func (s *UserService) FindUserByPhone(ctx context.Context, PhoneNumber string) (*models.User, error) {
	user, err := s.UserRepo.FindUserByPhone(ctx, PhoneNumber)
	if err != nil {
		log.Printf("UserService: cannot find user by phone number")
		return nil, err
	}

	log.Println("UserService: Found user by phone number")
	return user, nil
}

func (s *UserService) FindUserByEmail(ctx context.Context, EmailId string) (*models.User, error) {
	user, err := s.UserRepo.FindUserByEmail(ctx, EmailId)
	if err != nil {
		log.Printf("UserService: cannot find user by email")
		return nil, err
	}

	log.Println("UserService: Found user by email")
	return user, nil
}
