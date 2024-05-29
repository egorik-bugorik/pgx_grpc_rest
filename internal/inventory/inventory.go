package inventory

import (
	"context"
	"time"
)

type Product struct {
	ID          string
	Name        string
	Description string
	Price       int
	CreatedAt   time.Time
	ModifiedAt  time.Time
}

type CreateProductParams struct {
	ID          string
	Name        string
	Description string
	Price       int
}
type UpdateProductParams struct {
	ID          string
	Name        *string
	Description *string
	Price       *int
}
type SearchProductsParams struct {
	QueryString string
	MinPrice    int
	MAxPice     int
	Pagination  Pagination
}
type SearchProductResponse struct {
	items []*Product
	Total int
}
type Pagination struct {
	Offset int
	Limit  int
}

// VALIDATE REQUEST PARAMS ....
func (p *CreateProductParams) validate() error {
	//TODO
	return nil

}

func (p *UpdateProductParams) validate() error {
	//TODO
	return nil

}

func (p *SearchProductsParams) validate() error {
	//TODO
	return nil

}

// SERVICE IMPLEMENTATION

func (s *Service) CreateProducts(ctx context.Context, p CreateProductParams) error {
	if err := p.validate(); err != nil {
		return err
	}

	return s.DB.CreateProducts(ctx, p)
}
func (s *Service) UpdateProducts(ctx context.Context, p UpdateProductParams) error {
	if err := p.validate(); err != nil {
		return err
	}

	return s.DB.UpdateProducts(ctx, p)
}
func (s *Service) DeleteProducts(ctx context.Context, id string) error {
	if id == "" {
		return ValidationError{"Empty id to delete"}
	}

	return s.DB.DeleteProducts(ctx, id)
}
func (s *Service) GetProducts(ctx context.Context, id string) (*Product, error) {
	if id == "" {
		return nil, ValidationError{"Empty id to delete"}
	}

	return s.DB.GetProducts(ctx, id)
}
func (s *Service) SearchProducts(ctx context.Context, p SearchProductsParams) (*SearchProductResponse, error) {
	if err := p.validate(); err != nil {
		return nil, err
	}

	return s.DB.SearchProducts(ctx, p)
}

type ValidationError struct {
	msg string
}

func (v ValidationError) Error() string {
	return v.msg
}
