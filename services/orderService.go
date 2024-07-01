package services

import (
	"context"
	"go_project/models"
	"go_project/repository"
	"log"

	
)

type OrderService struct {
	OrderRepo *repository.OrderRepository
}

func NewOrderService(orderRepo *repository.OrderRepository) *OrderService {
	return &OrderService{
		OrderRepo: orderRepo,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, order *models.Order) error {
	err := s.OrderRepo.InsertOrderData(ctx, order)
	if err != nil {
		log.Printf("OrderService: error inserting order data: %v", err)
		return err
	}

	log.Println("OrderService: order created successfully")
	return nil
}

func (s *OrderService) FindOrderByUserId(ctx context.Context, ID string) ([]models.Order, error) {
	order, err := s.OrderRepo.FindById(ctx, ID)
	if err != nil {
		log.Printf("OrderService: unable to fetch Order")
		return nil, err
	}

	log.Println("OrderService: successfully fetched Order")
	return order, nil
}
