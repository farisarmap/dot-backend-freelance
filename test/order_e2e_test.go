package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/farisarmap/dot-backend-freelance/internal/entity"
)

func TestOrderCRUD(t *testing.T) {
	var (
		createdUserID  int
		createdOrderID int
	)

	t.Run("Setup - CreateUser", func(t *testing.T) {
		userPayload := `{"name":"E2E Test User","email":"e2e_order@testing.com"}`
		req, err := http.NewRequest(http.MethodPost, baseURL+"/users", bytes.NewBufferString(userPayload))
		if err != nil {
			t.Fatalf("Setup - CreateUser request creation error: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Setup - CreateUser request error: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Setup - CreateUser expected status %d, got %d, body: %s", http.StatusCreated, resp.StatusCode, string(bodyBytes))
		}

		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Setup - CreateUser read body error: %v", err)
		}

		var response ResponseFormat
		if err := json.Unmarshal(bodyBytes, &response); err != nil {
			t.Fatalf("Setup - CreateUser unmarshal error: %v", err)
		}

		var userData entity.User
		if err := json.Unmarshal(response.Data, &userData); err != nil {
			t.Fatalf("Setup - CreateUser data parse error: %v", err)
		}

		if userData.ID == 0 {
			t.Fatalf("Setup - CreateUser no valid ID returned")
		}
		createdUserID = int(userData.ID)
	})

	defer func() {
		// Cleanup: Delete the User after tests
		t.Run("Cleanup - DeleteUser", func(t *testing.T) {
			req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/users/%d", baseURL, createdUserID), nil)
			if err != nil {
				t.Fatalf("Cleanup - DeleteUser request creation error: %v", err)
			}

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				t.Fatalf("Cleanup - DeleteUser request error: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				bodyBytes, _ := io.ReadAll(resp.Body)
				t.Fatalf("Cleanup - DeleteUser expected status %d, got %d, body: %s", http.StatusOK, resp.StatusCode, string(bodyBytes))
			}
		})
	}()

	testCases := []struct {
		name       string
		method     string
		urlFn      func() string
		body       string
		wantStatus int
		checkFunc  func(t *testing.T, resp *http.Response)
	}{
		{
			name:   "CreateOrder - Success",
			method: http.MethodPost,
			urlFn: func() string {
				return baseURL + "/orders"
			},
			body:       fmt.Sprintf(`{"order_name":"E2E Test Order","user_id":%d}`, createdUserID),
			wantStatus: http.StatusCreated,
			checkFunc: func(t *testing.T, resp *http.Response) {
				bodyBytes, _ := io.ReadAll(resp.Body)

				var response ResponseFormat
				if err := json.Unmarshal(bodyBytes, &response); err != nil {
					t.Fatalf("CreateOrder unmarshal error: %v", err)
				}
				if response.Status != "success" {
					t.Fatalf("CreateOrder status not success, got: %s", response.Status)
				}
				var orderData entity.Order
				if err := json.Unmarshal(response.Data, &orderData); err != nil {
					t.Fatalf("CreateOrder data parse error: %v", err)
				}
				if orderData.ID == 0 {
					t.Fatalf("CreateOrder no valid ID returned")
				}
				createdOrderID = int(orderData.ID)
			},
		},
		{
			name:   "CreateOrder - Fail (Invalid UserID)",
			method: http.MethodPost,
			urlFn: func() string {
				return baseURL + "/orders"
			},
			body:       `{"order_name":"Invalid User Order","user_id":999999}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:   "CreateOrder - Fail (Missing order_name)",
			method: http.MethodPost,
			urlFn: func() string {
				return baseURL + "/orders"
			},
			body:       fmt.Sprintf(`{"user_id":%d}`, createdUserID),
			wantStatus: http.StatusBadRequest,
		},
		{
			name:   "GetAllOrders - Success",
			method: http.MethodGet,
			urlFn: func() string {
				return baseURL + "/orders"
			},
			wantStatus: http.StatusOK,
		},
		{
			name:   "GetOrderByID - Success",
			method: http.MethodGet,
			urlFn: func() string {
				return fmt.Sprintf("%s/orders/%d", baseURL, createdOrderID)
			},
			wantStatus: http.StatusOK,
		},
		{
			name:   "GetOrderByID - Fail (Not Found)",
			method: http.MethodGet,
			urlFn: func() string {
				return fmt.Sprintf("%s/orders/%d", baseURL, 999999)
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:   "UpdateOrder - Success",
			method: http.MethodPut,
			urlFn: func() string {
				return fmt.Sprintf("%s/orders/%d", baseURL, createdOrderID)
			},
			body:       fmt.Sprintf(`{"order_name":"E2E Updated Order","user_id":%d}`, createdUserID),
			wantStatus: http.StatusOK,
			checkFunc: func(t *testing.T, resp *http.Response) {
				bodyBytes, _ := io.ReadAll(resp.Body)

				var response ResponseFormat
				if err := json.Unmarshal(bodyBytes, &response); err != nil {
					t.Fatalf("UpdateOrder unmarshal error: %v", err)
				}
				if response.Status != "success" {
					t.Fatalf("UpdateOrder status not success, got: %s", response.Status)
				}
				var orderData entity.Order
				if err := json.Unmarshal(response.Data, &orderData); err != nil {
					t.Fatalf("UpdateOrder data parse error: %v", err)
				}
				if orderData.OrderName != "E2E Updated Order" {
					t.Fatalf("UpdateOrder name not updated, got: %s", orderData.OrderName)
				}
			},
		},
		{
			name:   "UpdateOrder - Fail (Invalid order_name)",
			method: http.MethodPut,
			urlFn: func() string {
				return fmt.Sprintf("%s/orders/%d", baseURL, createdOrderID)
			},
			body:       fmt.Sprintf(`{"order_name":"","user_id":%d}`, createdUserID),
			wantStatus: http.StatusBadRequest,
		},
		{
			name:   "UpdateOrder - Fail (Invalid UserID)",
			method: http.MethodPut,
			urlFn: func() string {
				return fmt.Sprintf("%s/orders/%d", baseURL, createdOrderID)
			},
			body:       `{"order_name":"E2E Updated Order","user_id":999999}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:   "PartialUpdateOrder - Success (OrderName only)",
			method: http.MethodPatch,
			urlFn: func() string {
				return fmt.Sprintf("%s/orders/%d", baseURL, createdOrderID)
			},
			body:       `{"order_name":"E2E Patched Order"}`,
			wantStatus: http.StatusOK,
			checkFunc: func(t *testing.T, resp *http.Response) {
				bodyBytes, _ := io.ReadAll(resp.Body)

				var response ResponseFormat
				if err := json.Unmarshal(bodyBytes, &response); err != nil {
					t.Fatalf("PartialUpdateOrder unmarshal error: %v", err)
				}
				if response.Status != "success" {
					t.Fatalf("PartialUpdateOrder status not success, got: %s", response.Status)
				}
				var orderData entity.Order
				if err := json.Unmarshal(response.Data, &orderData); err != nil {
					t.Fatalf("PartialUpdateOrder data parse error: %v", err)
				}
				if orderData.OrderName != "E2E Patched Order" {
					t.Fatalf("PartialUpdateOrder name not updated, got: %s", orderData.OrderName)
				}
			},
		},
		{
			name:   "PartialUpdateOrder - Fail (Invalid order_name)",
			method: http.MethodPatch,
			urlFn: func() string {
				return fmt.Sprintf("%s/orders/%d", baseURL, createdOrderID)
			},
			body:       `{"order_name":""}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:   "DeleteOrder - Success",
			method: http.MethodDelete,
			urlFn: func() string {
				return fmt.Sprintf("%s/orders/%d", baseURL, createdOrderID)
			},
			wantStatus: http.StatusOK,
		},
		{
			name:   "DeleteOrder - Fail (Not Found)",
			method: http.MethodDelete,
			urlFn: func() string {
				// Order ini sudah dihapus
				return fmt.Sprintf("%s/orders/%d", baseURL, createdOrderID)
			},
			wantStatus: http.StatusBadRequest,
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
				bodyBytes, _ := io.ReadAll(resp.Body)
				t.Fatalf("[%s] expected %d, got %d, body: %s", tc.name, tc.wantStatus, resp.StatusCode, string(bodyBytes))
			}

			// Jika ada checkFunc, jalankan
			if tc.checkFunc != nil {
				tc.checkFunc(t, resp)
			}

			// Debugging: Print createdUserID dan createdOrderID
			// fmt.Println("UserID:", createdUserID, "OrderID:", createdOrderID)
		})
	}
}
