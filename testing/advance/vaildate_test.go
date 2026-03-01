package testingadvance

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidate(t *testing.T) {
    cases := []struct{
        name string
        req *DepositReq
        expErr string
    }{
        {
            name: "Negative UserID",
            req: &DepositReq{
                UserID: -3,
                Amount: 1234,
            },
            expErr: "invalid user id",
        },
        {
            name: "Negative Amount",
            req:   &DepositReq{
                UserID: 3,
                Amount: -1234,
            },
            expErr: "amount must be greater than zero",
        },
        {
            name: "Too Big Amount",
            req: &DepositReq{
                UserID: 1234,
                Amount: 12345,
            },
            expErr: "limit exceeded",
        },
        {
            name: "Happy Path",
            req: &DepositReq{
                UserID: 1234,
                Amount: 1234,
            },
            expErr: "",
        },
    }

    for _, tt := range cases {
        t.Run(tt.name, func(t *testing.T) {
            err := tt.req.Vaildate()
            if tt.expErr != ""{
                require.Error(t, err)
                require.EqualError(t, err, tt.expErr)
            } else {
                require.NoError(t, err)
            }
        })
    }
}
