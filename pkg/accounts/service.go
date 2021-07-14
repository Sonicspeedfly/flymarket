package accounts

import (
	"context"
	"database/sql"
	"encoding/hex"
	"errors"
	"log"
	"math/rand"
	"strings"

	"github.com/jackc/pgx/v4"
	"golang.org/x/crypto/bcrypt"
)

var ErrNoSuchUser = errors.New("no such user")
var ErrInternal = errors.New("internal error")

type Service struct {
	db *sql.DB
}

func NewService(db *sql.DB) *Service {
	return &Service{db: db}
}

type Account struct {
	ID      	int64	`json:"id"`
	Phone   	string	`json:"phone"`
	Email   	string	`json:"email"`
	UserName	string	`json:"username"`
	Password	string	`json:"password"`
	Active 		bool	`json:"active"`
	Created		string	`json:"creted"`
}

type Registration struct {
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Email   	string	`json:"email"`
	Password string `json:"password"`
}

type Token struct {
	ID		int64 `json:"id"`
	Token   string `json:"token"`
}

func (s *Service) TonkenFound(ctx context.Context, token string) (*Account, error) {
	item := &Account{}
	infToken := &Token{}
	err := s.db.QueryRowContext(ctx, `
	SELECT account_id FROM account_tokens WHERE token = $1
	`, token).Scan(&infToken.ID)
	if err != nil {
		return nil, err
	}	
	item, _ = s.ByID(ctx, infToken.ID)
	return item, nil
}

func (s *Service) LoginAccount(ctx context.Context, phone string, pass string) (*Token, error) {
	item := &Account{}
	err := s.db.QueryRowContext(ctx, `
		SELECT id, phone, email, name, pass, active, created FROM account WHERE phone = $1
	`, phone).Scan(&item.ID, &item.Phone, &item.Email, &item.UserName, &item.Password, &item.Active, &item.Created)
	if err != nil {
		log.Println("login is invalid")
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(item.Password), []byte(pass))
	if err != nil {
		log.Print("password is invalid")
		return nil, err
	}
	buffer := make([]byte, 256)
	n, err := rand.Read(buffer)
	if n != len(buffer) || err != nil {
		return nil, errors.New("internal error")
	}

	token := hex.EncodeToString(buffer)
	infToken := &Token{
		ID: item.ID,
		Token: token,
	}
	_, err = s.db.QueryContext(ctx, `INSERT INTO account_tokens(token, account_id) VALUES($1, $2)`, token, item.ID)
	return infToken, nil
}  

func (s *Service) Autification(ctx context.Context, token string) (int64, error) {
	var id int64
	err := s.db.QueryRowContext(ctx, `SELECT account_id FROM account_tokens WHERE token = $1`, token).Scan(&id)
	
	if err == pgx.ErrNoRows {
		return 0, errors.New("account not found")
	}
	if err != nil {
		return 0, errors.New("internal error")
	}
	return id, nil
}

func (s *Service) ByID(ctx context.Context, id int64) (*Account, error) {
	item := &Account{}
	err := s.db.QueryRowContext(ctx, `
		SELECT id, phone, email, name, pass, active, created FROM account WHERE id = $1
	`, id).Scan(&item.ID, &item.Phone, &item.Email, &item.UserName, &item.Password, &item.Active, &item.Created)



	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("item not found")
	}

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return item, nil
}

func (s *Service) SaveAccount(ctx context.Context, phone string, email string, name string, pass string) (*Account, error) {
	var err error
	item := &Account{}
	registration := Registration{
		Name: name,
		Phone: phone,
		Email: email,
		Password: pass,
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(registration.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("customers Register bcrypt.GenerateFromPassword ERROR:", err)
		return nil, ErrInternal
	}

	err = s.db.QueryRowContext(ctx, `
		INSERT INTO account (name, phone, email, pass)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (phone) DO NOTHING
		RETURNING id, name, email, phone, pass, active, created
	`, registration.Name, registration.Phone, registration.Email, hash).Scan(
		&item.ID, &item.UserName, &item.Email, &item.Phone, &item.Password, &item.Active, &item.Created,
	)
	if err == pgx.ErrNoRows {
		return nil, ErrNoSuchUser
	}
	if err != nil {
		log.Println("customers Register s.pool.QueryRow ERROR:", err)
		return nil, ErrInternal
	}

	return item, nil
}

func (s *Service) RemoveByID(ctx context.Context, id int64) (*Account, error) {
	item, _ := s.ByID(ctx, id)
	err := s.db.QueryRowContext(ctx, `
	DELETE FROM account WHERE id = $1
	`, id)
	if err != nil {
		log.Println(err)
	}
	return item, nil
}

func (s *Service) All(ctx context.Context) ([]*Account, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, phone, email, name, pass, active, created FROM account ORDER BY created
	`)
	if err != nil {
		log.Println("GetAll s.pool.Query error:", err)
		return nil, err
	}
	defer rows.Close()
	items := make([]*Account, 0)
	for rows.Next() {
		item := &Account{}
		err = rows.Scan(&item.ID, &item.Phone, &item.Email, &item.UserName, &item.Password, &item.Created)
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

func (s *Service) IDByToken(ctx context.Context, token string) (int64, error) {
	var id int64
	err := s.db.QueryRowContext(ctx, `
		SELECT account_id FROM account_tokens WHERE token = $1
	`, token).Scan(&id)

	if err == pgx.ErrNoRows {
		return 0, ErrNoSuchUser
	}
	if err != nil {
		log.Println("account IDByToken s.db.QueryRowContext ERROR:", err)
		return 0, ErrInternal
	}

	return id, nil
}

func (s *Service) HasAnyRole(ctx context.Context, id int64, inRoles ...string) bool {
	var dbRoles []string
	err := s.db.QueryRowContext(ctx, `
		SELECT roles FROM account WHERE id = $1
	`, id).Scan(&dbRoles)
	if err != nil {
		log.Println("account HasAnyRole s.db.QueryRow ERROR:", err)
		return false
	}

	for _, inRole := range inRoles {
		for _, dbRole := range dbRoles {
			if dbRole == strings.ToUpper(inRole) {
				return true
			}
		}
	}

	return false
}
