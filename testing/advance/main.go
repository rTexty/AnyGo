package testingadvance

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// --- Модели ---
type DepositReq struct {
	UserID int     `json:"user_id"`
	Amount float64 `json:"amount"`
}

// --- ЖЕСТКАЯ ЗАВИСИМОСТЬ (БД) ---
// Представим, что это реальная структура, которая ходит в Postgres
type PostgresDB struct{}

func (db *PostgresDB) AddMoney(userID int, amount float64) error {
	// Представим, что тут SQL-запрос: UPDATE wallets SET balance = balance + $1...
	fmt.Println("Делаем реальный запрос в БД...")
	return nil
}

type WalletRepository interface {
	AddMoney(userId int, amount float64) error
}

// --- СЕРВИС (Бизнес-логика + Валидация) ---
type WalletService struct {
	// 	db *PostgresDB // ЖЕСТКАЯ ПРИВЯЗКА! Без реальной БД сервис не создать.
	repo WalletRepository
}

func NewWalletService(repo WalletRepository) *WalletService {
	return &WalletService{
		repo: repo,
	}
}

// Пополнение баланса
func (s *WalletService) Deposit(req DepositReq) error {
	// 2. Поход в базу
	err := s.repo.AddMoney(req.UserID, req.Amount)
	if err != nil {
		return fmt.Errorf("db error: %w", err)
	}

	return nil
}

func (req *DepositReq) Vaildate() error {
	if req.UserID <= 0 {
		return errors.New("invalid user id")
	}
	if req.Amount <= 0 {
		return errors.New("amount must be greater than zero")
	}
	if req.Amount > 10000 {
		return errors.New("limit exceeded")
	}
	return nil
}

type WalletDeposit struct {
    ws WalletService
}

func NewWalletDeposit(ws *WalletService) *WalletDeposit{
    return &WalletDeposit{
        ws: *ws,
    }
}

// --- ХЕНДЛЕР (API Слой) ---
func (wd *WalletDeposit) DepositHandler(w http.ResponseWriter, r *http.Request) {
	var req DepositReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad json", http.StatusBadRequest)
		return
	}

	if err := wd.ws.Deposit(req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"success"}`))
}
