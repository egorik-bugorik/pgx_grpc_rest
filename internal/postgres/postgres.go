package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/henvic/pgtools"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"rest_api_postgres_clean/internal/database"
	"rest_api_postgres_clean/internal/inventory"
	"time"
)

type DB struct {
	pool *pgxpool.Pool
}

func (db DB) conn(ctx context.Context) database.PGXQuerier {
	if tx, ok := ctx.Value(txCtx{}).(pgx.Tx); ok && tx != nil {

		return tx
	}

	if res, ok := ctx.Value(connCtx{}).(*pgxpool.Conn); ok && res != nil {

		return res
	}
	return db.pool
}

type txCtx struct{}

type connCtx struct{}

func (D *DB) CreateProducts(ctx context.Context, p inventory.CreateProductParams) error {
	q := `insert into product(id,name,price,description) values ($1,$2,$3,$4)`
	switch _, err := D.conn(ctx).Exec(ctx, q, p.ID, p.Name, p.Price, p.Description); {

	case errors.As(err, context.Canceled), errors.As(err, context.DeadlineExceeded):
		{
			return err
		}
	case err != nil:
		if pgErr := D.pgError(err); pgErr != nil {
			return pgErr
		}
		log.Println("Coudln't create product!!! :::", err)
		return errors.New("Coudln't create product!!!")

	}
	return nil

}
func (db DB) pgError(err error) error {

	var pgErr pgconn.PgError
	if !errors.As(err, &pgErr) {
		return nil
	}
	if pgErr.Code == pgerrcode.UniqueViolation {
		return errors.New("Product  with id  this id already exists")
	}
	if pgErr.Code == pgerrcode.CheckViolation {
		switch pgErr.ConstraintName {
		case "product_id_check":
			return errors.New("product id is valid")
		case "product_name_check":
			return errors.New("product name is valid")
		case "product_price_check":
			return errors.New("product price is valid")
		}
	}
	return nil
}

func (D *DB) UpdateProducts(ctx context.Context, p inventory.UpdateProductParams) error {
	q := `update product set 
                   name = COALESCE($1,"name"),
                   description = COALESCE($1,"description"),
                   price = COALESCE($1,"price"),
                   modified_at = now() 
where id =$4`

	exec, err := D.conn(ctx).Exec(ctx, q, p.Name, p.Description, p.Price, p.ID)
	if errors.As(err, context.Canceled) || errors.As(err, context.DeadlineExceeded) {
		return err
	}
	if pgErr := D.pgError(err); pgErr != nil {
		return pgErr
	}
	if err != nil {
		log.Println("Coudln't update product :::", err)
		return errors.New("Coudln't update product ")
	}
	if exec.RowsAffected() == 0 {
		return ErrProductNotFound
	}

	return nil
}

var ErrProductNotFound = errors.New("Product not found!!!")

func (D *DB) DeleteProducts(ctx context.Context, id string) error {
	//TODO implement me
	panic("implement me")
}

type product struct {
	ID          string
	Name        string
	Description string
	Price       int
	CreatedAt   time.Time
	ModifiedAt  time.Time
}

func (p *product) dto() *inventory.Product {
	return &inventory.Product{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		CreatedAt:   p.CreatedAt,
		ModifiedAt:  p.ModifiedAt,
	}

}
func (D *DB) GetProducts(ctx context.Context, id string) (*inventory.Product, error) {
	var p product

	q := fmt.Sprintf(`SELECT %s FROM "product" WHERE id = $1`, pgtools.Wildcard(p))

	query, err := D.conn(ctx).Query(ctx, q, id)
	if err != nil {
		return nil, err
	}
	if errors.As(err, context.Canceled) || errors.As(err, context.DeadlineExceeded) {
		return nil, err
	}
	if err == nil {
		p, err = pgx.CollectOneRow(query, pgx.RowToStructByPos[product])
	}
	if err != nil {
		log.Printf("cannot get product from database: %v\n", err)
		return nil, errors.New("cannot get product from database")
	}
	return p.dto(), nil

}

func (D *DB) SearchProducts(ctx context.Context, p inventory.SearchProductsParams) (*inventory.SearchProductResponse, error) {
	//TODO implement me
	panic("implement me")
}


}

func NewDb(pool *pgxpool.Pool) DB {
	return DB{
		pool: pool,
	}
}

var _ inventory.DB = *DB{}
