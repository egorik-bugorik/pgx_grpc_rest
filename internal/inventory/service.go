package inventory

import "context"

type Service struct {
	DB DB
}

func NewService(db DB) *Service {
	return &Service{DB: db}
}

type DB interface {
	CreateProducts(ctx context.Context, p CreateProductParams) error
	UpdateProducts(ctx context.Context, p UpdateProductParams) error
	DeleteProducts(ctx context.Context, id string) error
	GetProducts(ctx context.Context, id string) (*Product, error)
	SearchProducts(ctx context.Context, p SearchProductsParams) (*SearchProductResponse, error)
}
