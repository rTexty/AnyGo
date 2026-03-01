package testingadvance

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

// test
func TestIntegrate(t *testing.T) {
    var resp map[string]any

	obj := DepositReq{
		UserID: 1,
		Amount: 250,
	}

	jsonData, err := json.Marshal(obj)
	require.NoError(t, err)

	req := httptest.NewRequest("POST", "/deposit", bytes.NewReader(jsonData))
	w := httptest.NewRecorder()
	repo := MockWalletRepo{
		func(userId int, amount float64) error {
			return nil
		},
	}
	ws := NewWalletService(&repo)

	dep := NewWalletDeposit(ws)
	dep.DepositHandler(w, req)
	res := w.Result()
	defer res.Body.Close()
	bytes, err := io.ReadAll(res.Body)
	require.NoError(t, err)
    err = json.Unmarshal(bytes, &resp)
    require.NoError(t, err)
    require.Equal(t, http.StatusOK, res.StatusCode)
    require.Equal(t, "success", resp["status"].(string))
}
