package product

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

)

type Service struct {
	db *sql.DB
}

func NewService(db *sql.DB) *Service {
	return &Service{db: db}
}

type Product struct {
	ID          int64  		`json:"id"`
	NameProduct string 		`json:"name"`
	Image       string 		`json:"image"`
	File  		string		`json:"file"`
	Information string 		`json:"information"`
	Count       int64  		`json:"count"`
	Account_ID  int64		`json:"accountid"`
	Created     string 		`json:"created"`
}

func (s *Service) ByID(ctx context.Context, id int64) (*Product, error) {
	item := &Product{}

	err := s.db.QueryRowContext(ctx, `
		SELECT id, name, image, file, information, count, account_id, created FROM product WHERE id = $1
	`, id).Scan(&item.ID, &item.NameProduct, &item.Image,  &item.File, &item.Information, &item.Count, &item.Account_ID, &item.Created)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("item not found")
	}

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return item, nil
}

func (s *Service) SaveProduct(ctx context.Context, id int64, name string, image string, file string, information string, count int64, account_id int64) (*Product, error) {
	if id != 0 {
		times := time.Now()
		time := strings.Split(times.String(), ".")
		_, err := s.db.ExecContext(ctx, `
			UPDATE product 
			SET name = $2, image = $3, file = $4, information = $5, count = $6, account_id = $7, created = $8 
			WHERE id = $1
			RETURNING id
		`, id, name, image, file, information, count, account_id, time[0])
		if err != nil {
			log.Println(err)
			return nil, err
		}
	}
	if id == 0 {
		times := time.Now()
		time := strings.Split(times.String(), ".")
		err := s.db.QueryRowContext(ctx, `
		INSERT INTO product (name, image, file, information, count, account_id, created) VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
		`, name, image, file, information, count, account_id, time[0]).Scan(
			&id,
		)
		if err != nil {
			log.Println(err)
			return nil, err
		}
	}
	return s.ByID(ctx, id)
}

func (s *Service) All(ctx context.Context) ([]*Product, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, name, image, file, information, count, account_id, created FROM product ORDER BY created
	`)
	if err != nil {
		log.Println("GetAll s.pool.Query error:", err)
		return nil, err
	}
	defer rows.Close()
	items := make([]*Product, 0)
	for rows.Next() {
		item := &Product{}
		err = rows.Scan(&item.ID, &item.NameProduct, &item.Image, &item.File, &item.Information, &item.Count, &item.Account_ID, &item.Created)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if rows.Err() != nil {
		log.Println("GetAll rows.Err error:", err)
		return nil, err
	}
	return items, nil
}

func (s *Service) RemoveByID(ctx context.Context, id int64) (*Product, error) {
	item, _ := s.ByID(ctx, id)
	err := s.db.QueryRowContext(ctx, `
	DELETE FROM product WHERE id = $1
	`, id)
	if err != nil {
		log.Println(err)
	}
	return item, nil
}

func (s *Service) BuyProduct(ctx context.Context, id int64, count int64) (*Product, error) {
		ProductID, err := s.ByID(ctx, id)
			if err != nil {
				log.Print(err)
			}
		if ((ProductID.Count == 1) || (ProductID.Count <= count)) {
			ProductID.Count -= ProductID.Count
			s.RemoveByID(ctx, id)
			return ProductID, nil
		}
		if (ProductID.Count > 0) {
			ProductID.Count -= count
			s.SaveProduct(ctx, ProductID.ID, ProductID.NameProduct, ProductID.Image, ProductID.File, ProductID.Information, ProductID.Count, ProductID.Account_ID)
			return ProductID, nil
		}
		return nil, errors.New(http.StatusText(http.StatusInternalServerError))
}
