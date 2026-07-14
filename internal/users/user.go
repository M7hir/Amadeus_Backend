package user

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"time"

	"amadeus.m7hir.net/internal/jsonlog"
	"amadeus.m7hir.net/internal/validator"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type UserModel struct {
	DB     *sql.DB
	Logger *jsonlog.Logger
}

var (
	ErrDuplicateEmail = errors.New("duplicate email")
	ErrEditConflict   = errors.New("edit conflict")
	ErrRecordNotFound = errors.New("record not found")
)

var AnonymousUser = &User{}

type User struct {
	Id           string    `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	Email        string    `json:"email"`
	PasswordHash Password  `json:"password_hash"`
	Activated    bool      `json:"activated"`
	Version      int       `json:"version"`
}

func (u *User) IsAnonymousUser() bool {
	return u == AnonymousUser
}

type Password struct {
	plaintext *string
	hash      []byte
}

func (p *Password) Set(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return err
	}
	p.plaintext = &plaintextPassword
	p.hash = hash
	return nil
}

func (p *Password) Matches(plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plaintextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}
	return true, nil
}

func (p Password) Value() (driver.Value, error) {
	if len(p.hash) == 0 {
		return nil, errors.New("password hash is not set")
	}

	return string(p.hash), nil
}

func (p *Password) Scan(src any) error {
	switch value := src.(type) {
	case nil:
		p.hash = nil
		return nil
	case string:
		p.hash = []byte(value)
		return nil
	case []byte:
		p.hash = append(p.hash[:0], value...)
		return nil
	default:
		return fmt.Errorf("unsupported type %T for password hash", src)
	}
}

func ValidateEmail(v *validator.Validator, email string) {
	v.Check(email != "", "email", "must be provided")
	v.Check(validator.Matches(email, validator.EmailRX), "email", "must be a valid email address")
}

func Validatepassword(v *validator.Validator, password string) {
	v.Check(password != "", "password", "must be provided")
	v.Check(len(password) >= 8, "password", "must be at least 8 bytes")
	v.Check(len(password) <= 72, "password", "must not be longer than 72 bytes long")
	v.Check(validator.Matches(password, validator.UpperMatch), password, "must have atleast one upper case character")
	v.Check(validator.Matches(password, validator.LowerMatch), password, "must have atleast one lower case character")
	v.Check(validator.Matches(password, validator.SpecialMatch), password, "must have atleast one special character")
	v.Check(validator.Matches(password, validator.NumberMatch), password, "must have atleast one number character")
}

func ValidateUser(v *validator.Validator, user *User) {
	v.Check(user.FirstName != "", "First Name", "must be provided")
	v.Check(len(user.FirstName) <= 500, "First Name", "must not be more than 500 bytes")

	v.Check(user.LastName != "", "Last Name", "must be provided")
	v.Check(len(user.LastName) <= 500, "Last Name", "must not be more than 500 bytes")

	ValidateEmail(v, user.Email)

	if user.PasswordHash.plaintext != nil {
		Validatepassword(v, *user.PasswordHash.plaintext)
	}

	if user.PasswordHash.hash == nil {
		panic("missing password hash for user")
	}

}

func (m UserModel) InsertUser(user *User) error {
	query := `INSERT INTO users (first_name , last_name, email, password_hash, activated)
	VALUES ($1,$2,$3,$4,$5)
	RETURNING id, created_at , version`

	args := []interface{}{user.FirstName, user.LastName, user.Email, user.PasswordHash, user.Activated}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&user.Id, &user.CreatedAt, &user.Version)
	if err != nil {
		switch {
		case isDuplicateEmailError(err):
			return ErrDuplicateEmail
		default:
			return err
		}
	}
	return nil
}

func isDuplicateEmailError(err error) bool {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		return pqErr.Code == "23505"
	}

	return false
}

func (m UserModel) GetUser(email string) (*User, error) {
	query := `SELECT id, created_at, updated_at, first_name, last_name, email, password_hash, activated, version
    FROM users
    WHERE email = $1`

	var user User
	args := []interface{}{email}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(
		&user.Id,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.PasswordHash.hash,
		&user.Activated,
		&user.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &user, nil

}

func (m UserModel) UpdateUser(user *User) error {
	query := `UPDATE users 
	SET first_name=$1,last_name = $2,email = $3,password_hash = $4,activated = $5,
	version = version + 1
	WHERE id = $6 AND version = $7
	RETURNING version`

	args := []interface{}{user.FirstName, user.LastName, user.Email, user.PasswordHash, user.Activated, user.Id, user.Version}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&user.Version)
	if err != nil {
		switch {
		case isDuplicateEmailError(err):
			return ErrDuplicateEmail
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil
}

func (m UserModel) GetForToken(tokenScope, tokenPlaintext string) (*User, error) {
	tokenHash := sha256.Sum256([]byte(tokenPlaintext))

	query := `
		SELECT users.id, users.created_at, users.first_name,users.last_name, users.email,
users.password_hash, users.activated, users.version
FROM users
INNER JOIN tokens
ON users.id = tokens.user_id
WHERE tokens.hash = $1
AND tokens.scope = $2
AND tokens.expiry > $3`

	args := []interface{}{tokenHash[:], tokenScope, time.Now()}

	var user User

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(
		&user.Id,
		&user.CreatedAt,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.PasswordHash.hash,
		&user.Activated,
		&user.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}
