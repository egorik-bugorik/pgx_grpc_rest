package api

import (
	"context"
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"rest_api_postgres_clean/internal/inventory"
)

type InventoryGRPC struct {
	UnimplementedInventoryServer
	inventory *inventory.Service
}

func (g *InventoryGRPC) SearchProducts(ctx context.Context, p *SearchProductsRequest) (*SearchProductsResponse, error) {
	params := inventory.SearchProductsParams{
		QueryString: p.QueryString,
	}
	if p.MinPrice != nil {
		params.MinPrice = int(*p.MinPrice)
	}
	if p.MaxPrice != nil {
		params.MAxPice = int(*p.MaxPrice)
	}
	page, pp := 1, 50
	if p.Page != nil {
		page = int(*p.Page)
	}
	params.Pagination = inventory.Pagination{
		Limit:  pp * page,
		Offset: pp * (page - 1),
	}
	products, err := g.inventory.SearchProducts(ctx, params)
	if err != nil {
		return nil, grpcAPIError(err)
	}

	items := []*Product{}
	for _, p := range products.Items {
		items = append(items, &Product{
			Id:          p.ID,
			Price:       int64(p.Price),
			Name:        p.Name,
			Description: p.Description,
		})
	}
	return &SearchProductsResponse{
		Total: int32(products.Total),
		Items: items,
	}, nil
}
func (g *InventoryGRPC) CreateProduct(ctx context.Context, p *CreateProductRequest) (*CreateProductResponse, error) {

	if err := g.inventory.CreateProducts(ctx, inventory.CreateProductParams{

		ID:          p.Id,
		Name:        p.Name,
		Description: p.Description,
		Price:       int(p.Price),
	}); err != nil {
		return nil, grpcAPIError(err)
	}
	return &CreateProductResponse{}, nil

	return nil, status.Errorf(codes.Unimplemented, "method CreateProduct not implemented")
}
func (g *InventoryGRPC) UpdateProduct(ctx context.Context, p *UpdateProductRequest) (*UpdateProductResponse, error) {

	params := inventory.UpdateProductParams{
		ID:          p.Id,
		Name:        p.Name,
		Description: p.Description,
	}
	if p.Price != nil {
		price := int(*p.Price)
		params.Price = &price
	}
	if err := g.inventory.UpdateProducts(ctx, params); err != nil {
		return nil, grpcAPIError(err)
	}
	return &UpdateProductResponse{}, nil
}
func (g *InventoryGRPC) DeleteProduct(ctx context.Context, p *DeleteProductRequest) (*DeleteProductResponse, error) {
	if err := g.inventory.DeleteProducts(ctx, p.Id); err != nil {
		return nil, grpcAPIError(err)
	}
	return &DeleteProductResponse{}, nil
}
func (g *InventoryGRPC) GetProduct(ctx context.Context, p *GetProductRequest) (*GetProductResponse, error) {
	product, err := g.inventory.GetProducts(ctx, p.Id)
	if err != nil {
		return nil, grpcAPIError(err)
	}
	if product == nil {
		return nil, status.Error(codes.NotFound, "product not found")
	}
	return &GetProductResponse{
		Id:          product.ID,
		Price:       int64(product.Price),
		Name:        product.Name,
		Description: product.Description,
		CreatedAt:   product.CreatedAt.String(),
		ModifiedAt:  product.ModifiedAt.String(),
	}, nil
}

func grpcAPIError(err error) error {
	switch {
	case err == context.DeadlineExceeded:
		return status.Error(codes.DeadlineExceeded, err.Error())
	case err == context.Canceled:
		return status.Error(codes.Canceled, err.Error())
	case errors.As(err, &inventory.ValidationError{}):
		return status.Errorf(codes.InvalidArgument, err.Error())
	default:
		return err
	}
}
