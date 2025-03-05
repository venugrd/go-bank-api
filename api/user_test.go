package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	mockdb "github.com/gurukanth/simplebank/db/mock"
	db "github.com/gurukanth/simplebank/db/sqlc"
	"github.com/gurukanth/simplebank/util"
	"github.com/stretchr/testify/require"
)

func TestCreateUserAPI(t *testing.T) {
	user := randomUser()

	testCases := []struct {
		name          string
		req           createUserRequest
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			req: createUserRequest{
				Username: user.Username,
				FullName: user.FullName,
				Email:    user.Email,
				Password: user.HashedPassword,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Eq(db.CreateUserParams{
						Username:       user.Username,
						FullName:       user.FullName,
						Email:          user.Email,
						HashedPassword: user.HashedPassword,
					})).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchUser(t, recorder.Body, user)
			},
		},
		{
			name: "InvalidEmail",
			req: createUserRequest{
				Username: user.Username,
				FullName: user.FullName,
				Email:    "abc.email.com",
				Password: user.HashedPassword,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidUsername",
			req: createUserRequest{
				Username: "Hello#1",
				FullName: user.FullName,
				Email:    user.Email,
				Password: user.HashedPassword,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InternalError",
			req: createUserRequest{
				Username: user.Username,
				FullName: user.FullName,
				Email:    user.Email,
				Password: user.HashedPassword,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Eq(db.CreateUserParams{
						Username:       user.Username,
						FullName:       user.FullName,
						Email:          user.Email,
						HashedPassword: user.HashedPassword,
					})).
					Times(1).
					Return(db.User{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {

			var buf bytes.Buffer
			err := json.NewEncoder(&buf).Encode(tc.req)
			require.NoError(t, err)

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)

			//build stubs
			tc.buildStubs(store)
			//start test server and send request
			server := NewServer(store)
			recorder := httptest.NewRecorder()

			request, err := http.NewRequest(http.MethodPost, "/users", &buf)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)

			//check response
			tc.checkResponse(t, recorder)
		})
	}
}

func randomUser() db.User {
	user := db.User{
		Username:       util.RandomUsername(),
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
		HashedPassword: util.RandomString(10),
	}
	log.Println("Created Randome User:", user)
	return user
}

func requireBodyMatchUser(t *testing.T, body *bytes.Buffer, user db.User) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotUser userResponse
	err = json.Unmarshal(data, &gotUser)
	require.NoError(t, err)

	log.Println("Go User:", gotUser)
	require.Equal(t, gotUser.Username, user.Username)
	require.Equal(t, gotUser.Email, user.Email)
	require.Equal(t, gotUser.FullName, user.FullName)
}
