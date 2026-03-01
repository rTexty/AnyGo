package testingadvance

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

type MockWalletRepo struct {
	MockAddMoney func(userId int, amount float64) error
}

func (m *MockWalletRepo) AddMoney(userId int, amount float64) error {
	return m.MockAddMoney(userId, amount)
}

func TestWallerService_Deposti(t *testing.T) {
	cases := []struct {
		name    string     // Имя кейса
		req     DepositReq // Входные данные для метода Deposit
		mockErr error      // Поведение БД: что должен вернуть наш мок?
		expErr  string     // Ожидаемый результат: какую ошибку вернет сам Deposit?
	}{
		{
			name:    "Успешное пополнение",
			req:     DepositReq{UserID: 1, Amount: 500},
			mockErr: nil, // База данных отработала без ошибок
			expErr:  "",  // Значит и Deposit должен вернуть успешный результат (без ошибки)
		},
		{
			name:    "Ошибка базы данных",
			req:     DepositReq{UserID: 1, Amount: 500},
			mockErr: errors.New("timeout"), // Имитируем падение БД
			expErr:  "db error: timeout",   // Deposit должен обернуть ошибку и вернуть её нам
		},
	}
    for _, tt := range cases{
        t.Run(tt.name, func(t *testing.T) {
            mockRepo := &MockWalletRepo{
                func(userId int, amount float64) error {
                    return tt.mockErr
                },
            }
            service := NewWalletService(mockRepo)
            err := service.Deposit(tt.req)
             if tt.expErr != "" {
                 require.Error(t, err)
                 require.EqualError(t, err, tt.expErr)
             } else {
                 require.NoError(t, err)
             }
        })
    }
}
