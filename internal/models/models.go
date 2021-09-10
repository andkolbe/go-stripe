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
	Image          string `json:"image"`
	IsRecurring    bool   `json:"is_recurring"`
	PlanID         string `json:"plan_id"`
}

// the type for all orders
type Order struct {
	ID            int `json:"id"`
	WidgetID      int `json:"widget_id"`
	TransactionID int `json:"transaction_id"`
	CustomerID    int `json:"customer_id"`
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
	ExpiryMonth         int    `json:"expiry_month"`
	ExpiryYear          int    `json:"expiry_year"`
	BankReturnCode      string `json:"bank_return_code"`
	TransactionStatusID int    `json:"transaction_status_id"`
	PaymentIntent       string `json:"payment_intent"`
	PaymentMethod       string `json:"payment_method"`
}

// the type for users
type User struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

// the type for customers
type Customer struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

func (m *DBModel) GetWidget(id int) (Widget, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var widget Widget

	row := m.DB.QueryRowContext(ctx, `
		SELECT id, name, description, inventory_level, price, COALESCE(image, ''), is_recurring, plan_id
		FROM widgets 
		WHERE id = ?
		`, id)
	err := row.Scan(
		&widget.ID,
		&widget.Name,
		&widget.Description,
		&widget.InventoryLevel,
		&widget.Price,
		&widget.Image,
		&widget.IsRecurring,
		&widget.PlanID,
	)
	if err != nil {
		return widget, err
	}

	return widget, nil
}

// inserts a new transaction and returns its id
func (m *DBModel) InsertTransaction(txn Transaction) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `
		INSERT INTO transactions (amount, currency, last_four, bank_return_code, transaction_status_id, expiry_month, expiry_year, payment_intent, payment_method)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := m.DB.ExecContext(ctx, stmt,
		txn.Amount,
		txn.Currency,
		txn.LastFour,
		txn.BankReturnCode,
		txn.TransactionStatusID,
		txn.ExpiryMonth,
		txn.ExpiryYear,
		txn.PaymentIntent,
		txn.PaymentMethod,
	)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// inserts a new order and returns its id
func (m *DBModel) InsertOrder(order Order) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `
		INSERT INTO orders (widget_id, transaction_id, status_id, customer_id, quantity, amount)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	result, err := m.DB.ExecContext(ctx, stmt,
		order.WidgetID,
		order.TransactionID,
		order.StatusID,
		order.CustomerID,
		order.Quantity,
		order.Amount,
	)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// inserts a new customer and returns its id
func (m *DBModel) InsertCustomer(c Customer) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `
		INSERT INTO customers (first_name, last_name, email)
		VALUES (?, ?, ?)
	`

	result, err := m.DB.ExecContext(ctx, stmt,
		c.FirstName,
		c.LastName,
		c.Email,
	)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}
