package testing_basics

import (
	"errors"
	"strings"
)

// User представляет пользователя системы
type User struct {
	ID       int
	Username string
	Email    string
	Age      int
}

// ValidateUser проверяет, корректны ли данные пользователя.
// Правила:
// 1. Username не должен быть пустым и должен быть длиннее 3 символов.
// 2. Email должен содержать символ '@'.
// 3. Age должен быть больше или равен 18.
func ValidateUser(u User) error {
	if len(strings.TrimSpace(u.Username)) <= 3 {
		return errors.New("username must be longer than 3 characters")
	}
	if !strings.Contains(u.Email, "@") {
		return errors.New("invalid email format")
	}
	if u.Age < 18 {
		return errors.New("user must be at least 18 years old")
	}
	return nil
}

// CalculateDiscount рассчитывает скидку в зависимости от суммы покупок.
// Правила:
// - Сумма < 1000: скидка 0%
// - Сумма от 1000 до 5000: скидка 5%
// - Сумма > 5000: скидка 10%
// - Если сумма отрицательная, возвращает 0 и ошибку "negative amount"
func CalculateDiscount(amount float64) (float64, error) {
	if amount < 0 {
		return 0, errors.New("negative amount")
	}
	if amount >= 5000 {
		return amount * 0.10, nil
	}
	if amount >= 1000 {
		return amount * 0.05, nil
	}
	return 0, nil
}
