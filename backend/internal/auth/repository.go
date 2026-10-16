package auth

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

// WithTx runs fn in a database transaction. The transaction is committed
// when fn returns nil and rolled back when fn returns an error.
func (r *Repository) WithTx(ctx context.Context, fn func(pgx.Tx) error) error {
	return pgx.BeginFunc(ctx, r.pool, fn)
}

func (r *Repository) CreateUser(ctx context.Context, username, passwordHash, accountType string) (User, error) {
	var user User

	row := r.pool.QueryRow(ctx, `INSERT INTO users (username, password_hash, account_type) VALUES ($1, $2, $3) RETURNING id, username, password_hash, balance, account_type`, username, passwordHash, accountType)
	err := row.Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Balance, &user.AccountType)

	var pgErr *pgconn.PgError

	if errors.As(err, &pgErr) && pgErr.Code == "23505" && pgErr.ConstraintName == "users_username_unique" {
		return User{}, ErrUserAlreadyExist
	}

	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (r *Repository) GetUser(ctx context.Context, identifier any) (User, error) {
	switch t := identifier.(type) {
	case string:
		return r.getUserByUsername(ctx, t)
	case int:
		return r.getUserByID(ctx, int64(t))
	case int64:
		return r.getUserByID(ctx, t)
	default:
		return User{}, ErrWrongIdentificator
	}
}

func (r *Repository) SetBalance(ctx context.Context, userID, balance int64) (User, error) {
	var user User

	row := r.pool.QueryRow(ctx, `UPDATE users SET balance = $1 WHERE id = $2 RETURNING id, username, password_hash, balance`, balance, userID)
	err := row.Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Balance, &user.AccountType)

	if err != nil {
		return User{}, err
	}

	return user, nil
}

// DebitBalanceTx atomically subtracts amount from a user's balance. The
// caller owns tx and decides when the whole transaction is committed.
func (r *Repository) DebitBalanceTx(ctx context.Context, tx pgx.Tx, userID, amount int64) (User, error) {
	if amount <= 0 {
		return User{}, ErrInvalidBalance
	}

	var user User
	err := tx.QueryRow(ctx, `
		UPDATE users
		SET balance = balance - $1
		WHERE id = $2 AND balance >= $1
		RETURNING id, username, password_hash, balance, account_type
	`, amount, userID).Scan(
		&user.ID,
		&user.Username,
		&user.PasswordHash,
		&user.Balance,
		&user.AccountType,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		var exists bool
		if existsErr := tx.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM users WHERE id = $1)`, userID).Scan(&exists); existsErr != nil {
			return User{}, existsErr
		}
		if !exists {
			return User{}, ErrUserNotFound
		}
		return User{}, ErrInsufficientBalance
	}
	if err != nil {
		return User{}, err
	}

	return user, nil
}

// CreditBalanceTx atomically adds amount to a user's balance. The caller owns
// tx and decides when the whole transaction is committed.
func (r *Repository) CreditBalanceTx(ctx context.Context, tx pgx.Tx, userID, amount int64) (User, error) {
	if amount <= 0 {
		return User{}, ErrInvalidBalance
	}

	var user User
	err := tx.QueryRow(ctx, `
		UPDATE users
		SET balance = balance + $1
		WHERE id = $2
		RETURNING id, username, password_hash, balance, account_type
	`, amount, userID).Scan(
		&user.ID,
		&user.Username,
		&user.PasswordHash,
		&user.Balance,
		&user.AccountType,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return User{}, ErrUserNotFound
	}
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (r *Repository) getUserByID(ctx context.Context, id int64) (User, error) {
	var user User

	row := r.pool.QueryRow(ctx, `SELECT id, username, password_hash, balance, account_type FROM users WHERE id = $1`, id)
	err := row.Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Balance, &user.AccountType)

	if errors.Is(err, pgx.ErrNoRows) {
		return User{}, ErrUserNotFound
	}

	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (r *Repository) getUserByUsername(ctx context.Context, name string) (User, error) {
	var user User

	row := r.pool.QueryRow(ctx, `SELECT id, username, password_hash, balance, account_type FROM users WHERE username = $1`, name)
	err := row.Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Balance, &user.AccountType)

	if errors.Is(err, pgx.ErrNoRows) {
		return User{}, ErrUserNotFound
	}

	if err != nil {
		return User{}, err
	}

	return user, nil
}
