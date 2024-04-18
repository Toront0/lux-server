package services

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"context"
	"guthub.com/Toront0/lux-server/internal/types"

	"fmt"
)

type AuthStorer interface {
	CreateUser(firstName, lastName, email, password string) (*types.AuthUser, error)
	GetUserBy(columnName string, value interface{}) (*types.LoginUser, error)
	
}

type authStore struct {
	conn *pgxpool.Pool
}

func NewAuthStore(conn *pgxpool.Pool) *authStore {
	return &authStore{
		conn: conn,
	}
}

func (s *authStore) CreateUser(firstName, lastName, email, password string) (*types.AuthUser, error) {
	acc := &types.AuthUser{}

	err := s.conn.QueryRow(context.Background(), `insert into users (first_name, last_name, email, password) values($1, $2, $3, $4) returning id, first_name, last_name, profile_img`, firstName, lastName, email, password).Scan(&acc.ID, &acc.FirstName, &acc.LastName, &acc.ProfileImg)

	// defaultTime := time.Date(1970, time.January, 1, 23, 0, 0, 0, time.UTC)

	// acc.VipFinishedAt = &defaultTime

	if err != nil {
		fmt.Printf("could not create the user %s", err)
		return &types.AuthUser{}, err
	}


	return acc, nil
}

func (s *authStore) GetUserBy(columnName string, value interface{}) (*types.LoginUser, error) {
	acc := &types.LoginUser{}

	query := fmt.Sprintf("select id, first_name, last_name, password, profile_img from users t1 where %s = $1", columnName)

	err := s.conn.QueryRow(context.Background(), query, value).Scan(&acc.ID, &acc.FirstName, &acc.LastName, &acc.Password, &acc.ProfileImg)

	if err != nil {
		fmt.Printf("could not get an user %s", err)
		return &types.LoginUser{}, err
	}

	return acc, nil
}

func (s *authStore) InsertEmailCode(email string, code int) error {

	_, err := s.conn.Exec(context.Background(), `insert into email_codes (email, code) values ($1, $2)`, email, code)

	return err
}

func (s *authStore) DeleteCodeIfExist(email string) error {

	_, err := s.conn.Exec(context.Background(), `delete from email_codes where email = $1`, email)

	return err
}

func (s *authStore) VerifyCode(email string, code string) (bool, error) {
	var res int
	
	err := s.conn.QueryRow(context.Background(), `select id from email_codes where email = $1 and code = $2`, email, code).Scan(&res)
	
	if err != nil {
		return false, err
	}

	if res == 0 {
		return false, nil
	} else {
		return true, nil
	}
}

func (s *authStore) ChangePassword(password string) error {

	_, err := s.conn.Exec(context.Background(), `update users set password = $1`, password)

	return err
}