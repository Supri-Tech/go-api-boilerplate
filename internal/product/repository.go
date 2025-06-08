package product

import (
	"context"
	"database/sql"
	"errors"
	"time"

	sq "github.com/Masterminds/squirrel"
)

type Repository interface {
	GetAll(ctx context.Context) ([]Product, error)
	GetByID(ctx context.Context, id int64) (*Product, error)
	Create(ctx context.Context, product Product) (*Product, error)
	Update(ctx context.Context, product Product) (*Product, error)
	Delete(ctx context.Context, id int64) error
}

type repository struct {
	db  *sql.DB
	psq sq.StatementBuilderType
}

func NewRepository(db *sql.DB) Repository {
	return &repository{
		db:  db,
		psq: sq.StatementBuilder.PlaceholderFormat(sq.Question),
	}
}

func (repo *repository) GetAll(ctx context.Context) ([]Product, error) {
	query := repo.psq.Select("id", "product_name", "product_price", "product_stock", "created_at").From("products")

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := repo.db.QueryContext(ctx, sqlStr, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []Product

	for rows.Next() {
		var p Product
		err := rows.Scan(&p.ID, &p.ProductName, &p.ProductPrice, &p.ProductStock, &p.CreatedAt)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

func (repo *repository) GetByID(ctx context.Context, id int64) (*Product, error) {
	query := repo.psq.Select("id", "product_name", "product_price", "product_stock", "created_at").From("products").Where(sq.Eq{"id": id}).Limit(1)

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	row := repo.db.QueryRowContext(ctx, sqlStr, args...)

	var product Product
	err = row.Scan(&product.ID, &product.ProductName, &product.ProductPrice, &product.ProductStock, &product.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &product, nil
}

func (repo *repository) Create(ctx context.Context, product Product) (*Product, error) {
	now := time.Now()
	product.CreatedAt = now
	product.UpdatedAt = now

	query := repo.psq.Insert("products").Columns("product_name", "product_price", "product_stock", "created_at", "updated_at").Values(product.ProductName, product.ProductPrice, product.ProductStock, product.CreatedAt, product.UpdatedAt)

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	result, err := repo.db.ExecContext(ctx, sqlStr, args...)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	product.ID = id
	return &product, nil
}

func (repo *repository) Update(ctx context.Context, product Product) (*Product, error) {
	query := repo.psq.Update("products").Set("product_name", product.ProductName).Set("product_price", product.ProductPrice).Set("product_stock", product.ProductStock).Where(sq.Eq{"id": product.ID}).Suffix("RETURNING id, product_name, product_price, product_stock")

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	row := repo.db.QueryRowContext(ctx, sqlStr, args...)

	var updated Product
	err = row.Scan(&updated.ID, &updated.ProductName, &updated.ProductPrice, &updated.ProductStock)
	if err != nil {
		return nil, err
	}
	return &updated, nil
}

func (repo *repository) Delete(ctx context.Context, id int64) error {
	query := repo.psq.Delete("products").Where(sq.Eq{"id": id})

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return err
	}
	_, err = repo.db.ExecContext(ctx, sqlStr, args...)
	return err
}
