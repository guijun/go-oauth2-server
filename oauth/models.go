package oauth

import (
	"database/sql"
	"time"

	"github.com/ant0ine/go-json-rest/rest"
	"github.com/pborman/uuid"
)

// Client ...
type Client struct {
	ID          uint           `gorm:"primary_key"`
	ClientID    string         `sql:"type:varchar(254);unique;not null"`
	Secret      string         `sql:"type:varchar(60);not null"`
	RedirectURI sql.NullString `sql:"type:varchar(200)"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Scope ...
type Scope struct {
	ID          uint   `gorm:"primary_key"`
	Scope       string `sql:"type:varchar(200);unique;not null"`
	Description sql.NullString
	IsDefault   bool `sql:"default:false"`
}

// User ...
type User struct {
	ID        uint   `gorm:"primary_key"`
	Username  string `sql:"type:varchar(254);unique;not null"`
	Password  string `sql:"type:varchar(60);not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// RefreshToken ...
type RefreshToken struct {
	ID        uint          `gorm:"primary_key"`
	Token     string        `sql:"type:varchar(40);unique;not null"`
	ExpiresAt time.Time     `sql:"not null"`
	Scope     string        `sql:"type:varchar(200);not null"`
	ClientID  sql.NullInt64 `sql:"index;not null"`
	UserID    sql.NullInt64 `sql:"index"`
	Client    Client
	User      User
	CreatedAt time.Time
	UpdatedAt time.Time
}

// AccessToken ...
type AccessToken struct {
	ID        uint          `gorm:"primary_key"`
	Token     string        `sql:"type:varchar(40);unique;not null"`
	ExpiresAt time.Time     `sql:"not null"`
	Scope     string        `sql:"type:varchar(200);not null"`
	ClientID  sql.NullInt64 `sql:"index;not null"`
	UserID    sql.NullInt64 `sql:"index"`
	Client    Client
	User      User
	CreatedAt time.Time
	UpdatedAt time.Time
}

// AuthCode ...
type AuthCode struct {
	ID          uint           `gorm:"primary_key"`
	Code        string         `sql:"type:varchar(40);unique;not null"`
	RedirectURI sql.NullString `sql:"type:varchar(200)"`
	ExpiresAt   time.Time      `sql:"not null"`
	Scope       string         `sql:"type:varchar(200);not null"`
	ClientID    sql.NullInt64  `sql:"index;not null"`
	UserID      sql.NullInt64  `sql:"index"`
	Client      Client
	User        User
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func newAccessToken(accessTokenLifetime int, client *Client, user *User, scope string) *AccessToken {
	accessToken := &AccessToken{
		Token:     uuid.New(),
		ExpiresAt: time.Now().Add(time.Duration(accessTokenLifetime) * time.Second),
		Scope:     scope,
		Client:    *client,
	}
	if user != nil {
		accessToken.User = *user
	}
	return accessToken
}

func newRefreshToken(refreshTokenLifetime int, client *Client, user *User, scope string) *RefreshToken {
	refreshToken := &RefreshToken{
		Token:     uuid.New(),
		ExpiresAt: time.Now().Add(time.Duration(refreshTokenLifetime) * time.Second),
		Scope:     scope,
		Client:    *client,
	}
	if user != nil {
		refreshToken.User = *user
	}
	return refreshToken
}

func writeAccessToken(w rest.ResponseWriter, accessTokenLifetime int, accessToken *AccessToken, refreshToken *RefreshToken) {
	// Content-Type header must set charset in response
	// See https://github.com/ant0ine/go-json-rest/issues/156
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// Write access token to response
	w.WriteJson(map[string]interface{}{
		"id":            accessToken.ID,
		"access_token":  accessToken.Token,
		"expires_in":    accessTokenLifetime,
		"token_type":    "Bearer",
		"scope":         accessToken.Scope,
		"refresh_token": refreshToken.Token,
	})
}
