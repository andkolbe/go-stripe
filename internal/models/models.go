package models

import (
	"context"
	"database/sql"
	"time"
)

// the type for database connection values
type DBModel struct {
	DB *sql.DB
}

// the wrapper for all models
type Models struct {
	DB DBModel
}

// returns a model type with database connection pool
func NewModels(db *sql.DB) Models {
	return Models{
		DB: DBModel{DB: db},
	}
}

// the type for all widgets
type Widget struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	InventoryLevel int    `json:"inventory_level"`
	Price          int    `json:"price"`
	Image          string    `json:"image"`
}

// the type for all orders
type Order struct {
	ID            int `json:"id"`
	WidgetID      int `json:"widget_id"`
	TransactionID int `json:"transaction_id"`
	StatusID      int `json:"status_id"`
	Quantity      int `json:"quantity"`
	Amount        int `json:"amount"`
}

// the type for order statuses
type Status struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// the type for transaction statuses
type TransactionStatus struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// the type for transactions
type Transaction struct {
	ID                  int    `json:"id"`
	Amount              int    `json:"amount"`
	Currency            string `json:"currency"`
	LastFour            string `json:"last_four"`
	BankReturnCode      string `json:"bank_return_code"`
	TransactionStatusId string `json:"transaction_status_id"`
}

// the type for users
type User struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

func (m *DBModel) GetWidget(id int) (Widget, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var widget Widget

	row := m.DB.QueryRowContext(ctx, "SELECT id, name FROM widgets WHERE id = ?", id)
	err := row.Scan(&widget.ID, &widget.Name)
	if err != nil {
		return widget, err
	}

	return widget, nil
}
