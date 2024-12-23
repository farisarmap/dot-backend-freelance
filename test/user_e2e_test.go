package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
)

const baseURL = "http://localhost:8080"

type ResponseFormat struct {
	Status  string          `json:"status"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

func TestUserCRUD(t *testing.T) {
	var createdUserID int

	testCases := []struct {
		name       string
		method     string
		urlFn      func() string
		body       string
		wantStatus int
		checkFunc  func(t *testing.T, resp *http.Response)
	}{
		{
			name:   "CreateUser - Success",
			method: http.MethodPost,
			urlFn: func() string {
				return baseURL + "/users"
			},
			body:       `{"name":"E2E Test User","email":"e2e@testing.com"}`,
			wantStatus: http.StatusCreated,
			checkFunc: func(t *testing.T, resp *http.Response) {
				bodyBytes, _ := io.ReadAll(resp.Body)

				var response ResponseFormat
				if err := json.Unmarshal(bodyBytes, &response); err != nil {
					t.Fatalf("CreateUser unmarshal error: %v", err)
				}
				if response.Status != "success" {
					t.Fatalf("CreateUser status not success, got: %s", response.Status)
				}
				var userData map[string]interface{}
				if err := json.Unmarshal(response.Data, &userData); err != nil {
					t.Fatalf("CreateUser data parse error: %v", err)
				}
				if idVal, ok := userData["id"].(float64); ok {
					createdUserID = int(idVal)
				} else {
					t.Fatalf("CreateUser no valid ID returned")
				}
			},
		},
		{
			name:   "CreateUser - Fail (Invalid Email)",
			method: http.MethodPost,
			urlFn: func() string {
				return baseURL + "/users"
			},
			body:       `{"name":"Invalid Email","email":"not-an-email"}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:   "CreateUser - Fail (Missing name)",
			method: http.MethodPost,
			urlFn: func() string {
				return baseURL + "/users"
			},
			body:       `{"email":"no-name@testing.com"}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:   "GetAllUsers - Success",
			method: http.MethodGet,
			urlFn: func() string {
				return baseURL + "/users"
			},
			wantStatus: http.StatusOK,
		},
		{
			name:   "GetUserByID - Success",
			method: http.MethodGet,
			urlFn: func() string {
				return fmt.Sprintf("%s/users/%d", baseURL, createdUserID)
			},
			wantStatus: http.StatusOK,
		},
		{
			name:   "UpdateUser - Success",
			method: http.MethodPut,
			urlFn: func() string {
				return fmt.Sprintf("%s/users/%d", baseURL, createdUserID)
			},
			body:       `{"name":"E2E Updated Name","email":"e2e_updated@testing.com"}`,
			wantStatus: http.StatusOK,
		},
		{
			name:   "UpdateUser - Fail (Invalid email)",
			method: http.MethodPut,
			urlFn: func() string {
				return fmt.Sprintf("%s/users/%d", baseURL, createdUserID)
			},
			body:       `{"name":"Xyz","email":"not-email"}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:   "PartialUpdateUser - Success (Name only)",
			method: http.MethodPatch,
			urlFn: func() string {
				return fmt.Sprintf("%s/users/%d", baseURL, createdUserID)
			},
			body:       `{"name":"E2E Patched Name"}`,
			wantStatus: http.StatusOK,
		},
		{
			name:   "PartialUpdateUser - Fail (Invalid email)",
			method: http.MethodPatch,
			urlFn: func() string {
				return fmt.Sprintf("%s/users/%d", baseURL, createdUserID)
			},
			body:       `{"email":"this-is-not-email"}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:   "DeleteUser - Success",
			method: http.MethodDelete,
			urlFn: func() string {
				return fmt.Sprintf("%s/users/%d", baseURL, createdUserID)
			},
			wantStatus: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var req *http.Request
			var err error

			if tc.body != "" {
				req, err = http.NewRequest(tc.method, tc.urlFn(), bytes.NewBufferString(tc.body))
			} else {
				req, err = http.NewRequest(tc.method, tc.urlFn(), nil)
			}
			if err != nil {
				t.Fatalf("[%s] request creation error: %v", tc.name, err)
			}
			req.Header.Set("Content-Type", "application/json")

			// Kirim request
			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				t.Fatalf("[%s] request error: %v", tc.name, err)
			}
			defer resp.Body.Close()

			// Cek status code
			if resp.StatusCode != tc.wantStatus {
				t.Fatalf("[%s] expected %d, got %d", tc.name, tc.wantStatus, resp.StatusCode)
			}

			// Jika ada checkFunc, jalankan
			if tc.checkFunc != nil {
				tc.checkFunc(t, resp)
			}
		})
	}
}
