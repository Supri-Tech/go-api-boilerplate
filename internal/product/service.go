package product

import (
	"context"
	"errors"
)

type Service interface {
	GetProduct(ctx context.Context) ([]Product, error)
	GetProductByID(ctx context.Context, id int64) (*Product, error)
	AddProduct(ctx context.Context, product Product) (*Product, error)
	UpdateProduct(ctx context.Context, product Product) (*Product, error)
	DeleteProduct(ctx context.Context, id int64) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo}
}

func (svc *service) GetProduct(ctx context.Context) ([]Product, error) {
	products, err := svc.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (svc *service) GetProductByID(ctx context.Context, id int64) (*Product, error) {
	product, err := svc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, errors.New("product not found")
	}
	return product, nil
}

func (svc *service) AddProduct(ctx context.Context, product Product) (*Product, error) {
	if product.ProductName == "" {
		return nil, errors.New("product name is required")
	}
	if product.ProductPrice <= 0 {
		return nil, errors.New("product price must be greater than 0")
	}
	if product.ProductStock < 0 {
		return nil, errors.New("product stock cannot be negative")
	}
	return svc.repo.Create(ctx, product)
}

func (svc *service) UpdateProduct(ctx context.Context, product Product) (*Product, error) {
	return svc.repo.Update(ctx, product)
}

func (svc *service) DeleteProduct(ctx context.Context, id int64) error {
	return svc.repo.Delete(ctx, id)
}
