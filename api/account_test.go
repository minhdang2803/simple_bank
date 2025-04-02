package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	mockdb "github.com/minhdang2803/simple_bank/db/mock"
	db "github.com/minhdang2803/simple_bank/db/sqlc"
	"github.com/minhdang2803/simple_bank/utils"
	"github.com/stretchr/testify/require"
)

func getGetAccountAPITestCases() []struct {
	name          string
	accountID     int64
	buildStubs    func(store *mockdb.MockStore)
	checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
} {
	account := RandomAccount()
	testCases := []struct {
		name          string
		accountID     int64
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{{
		name:      "OK",
		accountID: account.ID,
		buildStubs: func(store *mockdb.MockStore) {
			store.EXPECT().
				GetAccount(gomock.Any(), gomock.Eq(account.ID)).
				Times(1).Return(account, nil)
		},
		checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
			require.Equal(t, http.StatusOK, recorder.Code)
			requireBodyMatchAccount(t, recorder, account)
		},
	},

		{
			name:      "Not found",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).Return(db.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:      "Internal Error",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:      "Internal Error",
			accountID: -1,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		}}
	return testCases
}

func TestGetAccountAPI(t *testing.T) {
	testCases := getGetAccountAPITestCases()
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			store := mockdb.NewMockStore(ctrl)
			///
			tc.buildStubs(store)

			server := NewServer(store)
			recoder := httptest.NewRecorder()
			url := fmt.Sprintf("/accounts/%d", tc.accountID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			//
			server.router.ServeHTTP(recoder, request)
			tc.checkResponse(t, recoder)
		})
	}
}

func TestListAccountAPI(t *testing.T) {
	listAccount := []db.Account{}
	for i := 0; i < 5; i++ {
		newAccount := RandomAccount()
		listAccount = append(listAccount, newAccount)
	}

	testCases := []struct {
		name          string
		params        db.ListAccountsParams
		buildStubs    func(store *mockdb.MockStore, param db.ListAccountsParams)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Happy case",
			params: db.ListAccountsParams{
				Limit:  5,
				Offset: 1,
			},
			buildStubs: func(store *mockdb.MockStore, param db.ListAccountsParams) {
				store.EXPECT().ListAccounts(gomock.Any(), gomock.Eq(param)).Times(1).Return(listAccount, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccounts(t, recorder.Body, listAccount)
			},
		},
		{
			name: "Internal server error",
			params: db.ListAccountsParams{
				Limit:  5,
				Offset: 1,
			},
			buildStubs: func(store *mockdb.MockStore, param db.ListAccountsParams) {
				store.EXPECT().ListAccounts(gomock.Any(), gomock.Eq(param)).Times(1).Return([]db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "Not found",
			params: db.ListAccountsParams{
				Limit:  5,
				Offset: 1,
			},
			buildStubs: func(store *mockdb.MockStore, param db.ListAccountsParams) {
				store.EXPECT().ListAccounts(gomock.Any(), gomock.Eq(param)).Times(1).Return([]db.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},

		{
			name: "Bad Request 1",
			params: db.ListAccountsParams{
				Limit:  5,
				Offset: -1,
			},
			buildStubs: func(store *mockdb.MockStore, param db.ListAccountsParams) {
				store.EXPECT().ListAccounts(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Bad Request 2",
			params: db.ListAccountsParams{
				Limit:  0,
				Offset: 1,
			},
			buildStubs: func(store *mockdb.MockStore, param db.ListAccountsParams) {
				store.EXPECT().ListAccounts(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for index := range testCases {
		tc := testCases[index]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			store := mockdb.NewMockStore(ctrl)

			// build stubs
			tc.buildStubs(store, tc.params)
			//
			server := NewServer(store)
			recorder := httptest.NewRecorder()
			url := fmt.Sprintf("/accounts?page_size=%d&page_id=%d", tc.params.Limit, tc.params.Offset)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			//serve the request
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}

}

func TestCreateAccountAPI(t *testing.T) {

	testCase := []struct {
		name          string
		body          func(account *db.Account) []byte
		buildStubs    func(store *mockdb.MockStore, account *db.Account)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder, account *db.Account)
	}{{
		name: "Happy case",
		body: func(account *db.Account) []byte {
			request, _ := json.Marshal(CreateAccountRequest{
				Owner:    account.Owner,
				Currency: account.Currency,
			})
			return request
		},
		buildStubs: func(store *mockdb.MockStore, account *db.Account) {
			store.EXPECT().CreateAccount(gomock.Any(), gomock.Any()).Times(1).Return(*account, nil)
		},
		checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, account *db.Account) {
			require.Equal(t, http.StatusOK, recorder.Code)
			requireBodyMatchAccount(t, recorder, *account)
		},
	},
		{
			name: "StatusInternalServerError case",
			body: func(account *db.Account) []byte {
				request, _ := json.Marshal(CreateAccountRequest{
					Owner:    account.Owner,
					Currency: account.Currency,
				})
				return request
			},
			buildStubs: func(store *mockdb.MockStore, account *db.Account) {
				store.EXPECT().CreateAccount(gomock.Any(), gomock.Any()).Times(1).Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, account *db.Account) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "BadRequest case",
			body: func(account *db.Account) []byte {
				request, _ := json.Marshal(CreateAccountRequest{
					Owner: account.Owner,
				})
				return request
			},
			buildStubs: func(store *mockdb.MockStore, account *db.Account) {
				store.EXPECT().CreateAccount(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, account *db.Account) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for index := range testCase {
		tc := testCase[index]
		account := RandomAccount()
		t.Run(
			tc.name,
			func(t *testing.T) {

				ctrl := gomock.NewController(t)
				defer ctrl.Finish()
				store := mockdb.NewMockStore(ctrl)
				tc.buildStubs(store, &account)

				server := NewServer(store)
				recorder := httptest.NewRecorder()
				url := "/accounts"
				request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(tc.body(&account)))
				require.NoError(t, err)

				// Add the appropriate headers
				request.Header.Set("Content-Type", "application/json")

				//Serve the request
				server.router.ServeHTTP(recorder, request)
				tc.checkResponse(t, recorder, &account)
			},
		)
	}

}

func requireBodyMatchAccounts(t *testing.T, body *bytes.Buffer, acccounts []db.Account) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotAccounts []db.Account
	err = json.Unmarshal(data, &gotAccounts)
	require.NoError(t, err)

	require.Equal(t, acccounts, gotAccounts)
}

func requireBodyMatchAccount(t *testing.T, body *httptest.ResponseRecorder, account db.Account) {
	data, err := io.ReadAll(body.Body)
	require.NoError(t, err)

	var gotAccount db.Account
	err = json.Unmarshal(data, &gotAccount)
	require.NoError(t, err)

	require.Equal(t, account, gotAccount)
}
func RandomAccount() db.Account {
	return db.Account{
		ID:       utils.RandomInt(1, 1000),
		Owner:    utils.RandomOwner(),
		Balance:  utils.RandomMoney(),
		Currency: utils.RandomCurrency2(),
	}
}
