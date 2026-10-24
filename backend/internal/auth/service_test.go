package auth

import (
	"context"
	"errors"
	"testing"
	"time"

	"market/internal/platform"
)

type authRepositoryFake struct {
	createUser func(context.Context, string, string, string) (User, error)
	getUser    func(context.Context, any) (User, error)
}

func (f *authRepositoryFake) CreateUser(ctx context.Context, username, passwordHash, accountType string) (User, error) {
	return f.createUser(ctx, username, passwordHash, accountType)
}

func (f *authRepositoryFake) GetUser(ctx context.Context, identifier any) (User, error) {
	return f.getUser(ctx, identifier)
}

type authTokenManagerFake struct {
	issue  func(int64) (platform.AccessToken, error)
	verify func(string) (int64, error)
}

func (f *authTokenManagerFake) IssueAccessToken(userID int64) (platform.AccessToken, error) {
	return f.issue(userID)
}

func (f *authTokenManagerFake) VerifyAccessToken(raw string) (int64, error) {
	return f.verify(raw)
}

func TestServiceRegister(t *testing.T) {
	t.Run("valid input", func(t *testing.T) {
		var gotUsername, gotPasswordHash, gotAccountType string
		repo := &authRepositoryFake{
			createUser: func(_ context.Context, username, passwordHash, accountType string) (User, error) {
				gotUsername = username
				gotPasswordHash = passwordHash
				gotAccountType = accountType
				return User{ID: 1, Username: username, PasswordHash: passwordHash, AccountType: accountType}, nil
			},
		}
		service := NewService(repo, &authTokenManagerFake{})

		got, err := service.Register(context.Background(), "  alice ", "password", "seller")

		if err != nil {
			t.Fatalf("Register() error = %v", err)
		}
		if got.Username != "alice" || gotUsername != "alice" || gotAccountType != "seller" {
			t.Fatalf("unexpected user: %+v, username=%q, accountType=%q", got, gotUsername, gotAccountType)
		}
		if gotPasswordHash == "" || ComparePassword(gotPasswordHash, "password") != nil {
			t.Fatal("password was not hashed correctly")
		}
	})

	tests := []struct {
		name     string
		username string
		password string
		wantErr  error
	}{
		{name: "empty username", username: " ", password: "pass", wantErr: ErrEmptyUsername},
		{name: "short password", username: "alice", password: "123", wantErr: ErrTooShortPassword},
		{name: "long password", username: "alice", password: string(make([]byte, 73)), wantErr: ErrTooLongPassword},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			called := false
			repo := &authRepositoryFake{
				createUser: func(context.Context, string, string, string) (User, error) {
					called = true
					return User{}, nil
				},
			}
			service := NewService(repo, &authTokenManagerFake{})

			_, err := service.Register(context.Background(), tt.username, tt.password, "user")

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Register() error = %v, want %v", err, tt.wantErr)
			}
			if called {
				t.Fatal("repository was called for invalid input")
			}
		})
	}
}

func TestServiceLogin(t *testing.T) {
	hash, err := HashPassword("password")
	if err != nil {
		t.Fatal(err)
	}
	expiresAt := time.Now().Add(time.Hour)
	var gotUserID int64
	repo := &authRepositoryFake{
		getUser: func(_ context.Context, identifier any) (User, error) {
			if identifier != "alice" {
				t.Fatalf("GetUser() identifier = %#v", identifier)
			}
			return User{ID: 42, Username: "alice", PasswordHash: hash}, nil
		},
	}
	tokens := &authTokenManagerFake{
		issue: func(userID int64) (platform.AccessToken, error) {
			gotUserID = userID
			return platform.AccessToken{Value: "access-token", ExpiresAt: expiresAt}, nil
		},
	}
	service := NewService(repo, tokens)

	got, err := service.Login(context.Background(), "alice", "password")

	if err != nil {
		t.Fatalf("Login() error = %v", err)
	}
	if got.AccessToken != "access-token" || !got.ExpiresAt.Equal(expiresAt) || gotUserID != 42 {
		t.Fatalf("unexpected login result: %+v, userID=%d", got, gotUserID)
	}
}

func TestServiceLoginPasswordMismatch(t *testing.T) {
	hash, err := HashPassword("password")
	if err != nil {
		t.Fatal(err)
	}
	service := NewService(
		&authRepositoryFake{
			getUser: func(context.Context, any) (User, error) {
				return User{ID: 1, PasswordHash: hash}, nil
			},
		},
		&authTokenManagerFake{},
	)

	_, err = service.Login(context.Background(), "alice", "wrong")

	if !errors.Is(err, ErrPasswordMismatch) {
		t.Fatalf("Login() error = %v, want %v", err, ErrPasswordMismatch)
	}
}

func TestServiceGetUserInfo(t *testing.T) {
	repo := &authRepositoryFake{
		getUser: func(_ context.Context, identifier any) (User, error) {
			if identifier != int64(7) {
				t.Fatalf("GetUser() identifier = %#v", identifier)
			}
			return User{ID: 7, Username: "alice"}, nil
		},
	}
	service := NewService(repo, &authTokenManagerFake{})

	got, err := service.GetUserInfo(context.Background(), 7)

	if err != nil {
		t.Fatalf("GetUserInfo() error = %v", err)
	}
	if got.ID != 7 || got.Username != "alice" {
		t.Fatalf("unexpected user: %+v", got)
	}
}
