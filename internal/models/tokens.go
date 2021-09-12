package models

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"log"
	"time"
)

const (
	ScopeAuthentication = "authentication"
)

type Token struct {
	PlainText string    `json:"token"`
	UserID    int64     `json:"-"`
	Hash      []byte    `json:"-"`
	Expiry    time.Time `json:"expiry"`
	Scope     string    `json:"-"`
}

// generates a token that lasts for ttl, and returns it
func GenerateToken(userID int, ttl time.Duration, scope string) (*Token, error) {
	token := &Token {
		UserID: int64(userID),
		Expiry: time.Now().Add(ttl),
		Scope: scope,
	}

	// make sure the token is secure
	randomBytes := make([]byte, 16)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}

	// the token being sent back to the end user
	token.PlainText = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)
	hash := sha256.Sum256([]byte(token.PlainText))
	token.Hash = hash[:]
	return token, nil
}

func (m *DBModel) InsertToken(t *Token, u User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// delete a preexisting token for the user if it exists
	stmt := `DELETE FROM tokens WHERE user_id = ?`
	_, err := m.DB.ExecContext(ctx, stmt, u.ID)
	if err != nil {
		return err
	}

	stmt = `
		INSERT INTO tokens (user_id, name, email, token_hash, expiry)
		VALUES (?, ?, ?, ?, ?)
	`
	_, err = m.DB.ExecContext(ctx, stmt,
		u.ID,
		u.LastName,
		u.Email,
		t.Hash,
		t.Expiry,
	)

	if err != nil {
		return err
	}
	return nil
}

func (m *DBModel) GetUserForToken(token string)  (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// convert the token we are receiving from the request into a hash so it matches what is in the database
	tokenHash := sha256.Sum256([]byte(token))
	var user User

	query := `
		SELECT u.id, u.first_name, u.last_name, u.email 
		FROM users u
		INNER JOIN tokens t
		ON (u.id = t.user_id) 
		WHERE t.token_hash = ?
		AND t.expiry > ?
	`
	err := m.DB.QueryRowContext(ctx, query, tokenHash[:], time.Now()).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
	)

	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &user, nil
}