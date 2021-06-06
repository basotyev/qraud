package postgresql

import (
	"armani_auth/pkg/models"
	"context"
	"errors"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"
)

func (m *UserModel) Insert(name, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err }
	stmt := `INSERT INTO users (name, email, hashed_password, created) VALUES($1, $2, $3, $4)`

	_, err = m.DB.Exec(context.Background(), stmt, name, email, string(hashedPassword), time.Now())

	if err != nil {
		if strings.Contains(err.Error(), "users_email_key"){
			return models.ErrDuplicateEmail
		}
	}
	return nil
}


type UserModel struct {
	DB *pgxpool.Pool
}


func (m *UserModel) Authenticate(email, password string) (models.User, error) {

	var id int
	var hashedPassword []byte
	var email1 string
	var name string
	stmt := "SELECT id, hashed_password, email, name FROM users WHERE email = $1 AND active = TRUE"
	row := m.DB.QueryRow(context.Background(),stmt,email)
	err := row.Scan(&id, &hashedPassword, &name, &email1)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.User{}, models.ErrInvalidCredentials
		} else {
			return models.User{}, err
		}
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) { return models.User{}, models.ErrInvalidCredentials
		} else {
			return models.User{}, err
		}
	}

	use := models.User{
		ID:             id,
		Name:           name,
		Email:          email1,
		HashedPassword: "",
		Created:        time.Time{},
		Active:         true,
	}

	return use, nil
}

func (m *UserModel) Get(id int) (*models.User, error) {
	return nil, nil
}