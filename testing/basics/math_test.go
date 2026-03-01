package testing_basics

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)


func TestValidateUser(t *testing.T) {
	cases := []struct {
		name string
		user *User
		expErr error
	}{
		{
			name: 	"err_username_empty",
			user:	&User{
				ID: 1,
				Username: "",
				Email: "jonhdoe@mail.ru",
				Age: 25,
			},
			expErr: errors.New("username must be longer than 3 characters"),
		},
		{
			name: 	"err_username_lre_3",
			user:	&User{
				ID: 1,
				Username: "abc",
				Email: "jonhdoe@mail.ru",
				Age: 25,
			},
			expErr: errors.New("username must be longer than 3 characters"),
		},
		{			
			name: 	"err_invalid_email",
			user:	&User{
				ID: 1,
				Username: "Jonh Doe",
				Email: "jonhdoemail.ru",
				Age: 25,
			},
			expErr:  errors.New("invalid email format"),
		},
		{			
			name: 	"err_underage",
			user:	&User{
				ID: 1,
				Username: "Jonh Doe",
				Email: "jonhdoe@mail.ru",
				Age: 10,
			},
			expErr:  errors.New("user must be at least 18 years old"),
		},
		
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateUser(*tt.user)
			if err != nil{
				require.Error(t, err)
				require.EqualError(t, err, tt.expErr.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestCalculateDiscount(t *testing.T) {
	cases := []struct{
		name 	string
		amount 	float64
		expRes 	float64
		expErr 	error
	}{
		{
			name: "err_less_than_0",
			amount: -2,
			expRes: 0,
			expErr: errors.New("negative amount"),
		},

		{
			name: "err_more_than_5000",
			amount: 50000,
			expRes: 5000,
			expErr: nil,
		},

		{
			name: "err_between_1000_and_5000",
			amount: 1600,
			expRes: 80,
			expErr: nil,
		},

		{
			name: "err_less_than_0",
			amount: -4343,
			expRes: 0,
			expErr: errors.New("negative amount"),
		},
	}

	for _, tCase := range cases {
		t.Run(tCase.name, func(t *testing.T) {
			res, err := CalculateDiscount(tCase.amount)
			if tCase.expErr != nil {
				require.Error(t, err)
				require.EqualError(t, err , tCase.expErr.Error())
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tCase.expRes, res)
		})
	}
}