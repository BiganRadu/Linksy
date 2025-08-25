package app_service

import (
	"backend/models"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	linkmocks "backend/drivers/LinkDriver/mocks"
	"backend/helpers"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// test handler constructor
func newTestAppHandler(mockLD *linkmocks.LinkDriver) *AppHandler {
	return &AppHandler{
		awsClient:   nil,
		s3Bucket:    "test-bucket",
		linkDriver:  mockLD,
		tokenHelper: helpers.NewTokenHelper("test_secret"),
		nowFunc: func() time.Time {
			return time.Unix(1700000000, 0)
		},
	}
}

func buildRouter(h *AppHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	api := r.Group("/app")
	h.Routes(api)
	return r
}

func doJSON(r http.Handler, method, path string, body any, headers map[string]string) *httptest.ResponseRecorder {
	var buf bytes.Buffer
	if body != nil {
		_ = json.NewEncoder(&buf).Encode(body)
	}
	req, _ := http.NewRequest(method, path, &buf)
	req.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	return rr
}

func authHeaders(h *AppHandler, email, username string) map[string]string {
	tok, _ := h.tokenHelper.GenerateToken(email, username, h.nowFunc().Unix(), time.Hour)
	return map[string]string{"AuthToken": tok}
}

func resetMock(m *linkmocks.LinkDriver) {
	m.ExpectedCalls = nil
	m.Calls = nil
}

// ---------- GET /app/member-info ----------
func TestGetMemberInfo(t *testing.T) {
	m := &linkmocks.LinkDriver{}
	h := newTestAppHandler(m)
	r := buildRouter(h)

	type tc struct {
		name       string
		headers    map[string]string
		mockSetup  func()
		wantStatus int
		check      func(*testing.T, *httptest.ResponseRecorder)
	}

	cases := []tc{
		{
			name:       "success",
			headers:    authHeaders(h, "user@example.com", "user"),
			mockSetup:  func() {},
			wantStatus: 200,
			check: func(t *testing.T, rr *httptest.ResponseRecorder) {
				var resp map[string]any
				require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &resp))
				assert.Equal(t, "user@example.com", resp["email"])
			},
		},
		{
			name:       "unauthorized",
			headers:    map[string]string{},
			mockSetup:  func() {},
			wantStatus: 401,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			resetMock(m)
			c.mockSetup()
			rr := doJSON(r, "GET", "/app/member-info", nil, c.headers)
			assert.Equal(t, c.wantStatus, rr.Code)
			m.AssertExpectations(t)
			if c.check != nil {
				c.check(t, rr)
			}
		})
	}
}

// ---------- GET /app/link?id= ----------
func TestGetLink(t *testing.T) {
	m := &linkmocks.LinkDriver{}
	h := newTestAppHandler(m)
	r := buildRouter(h)

	type tc struct {
		name       string
		id         string
		headers    map[string]string
		mockSetup  func()
		wantStatus int
		check      func(t *testing.T, body []byte)
	}

	cases := []tc{
		{
			name:    "success",
			id:      "link123",
			headers: authHeaders(h, "u@example.com", "u"),
			mockSetup: func() {
				m.On("GetLinkByID", "link123").Return(&models.Link{
					ID:             "link123",
					Title:          "Example",
					ReferencedLink: "https://example.com",
				}, nil)
			},
			wantStatus: 200,
			check: func(t *testing.T, body []byte) {
				var resp models.Link
				assert.NoError(t, json.Unmarshal(body, &resp))
				assert.Equal(t, "link123", resp.ID)
			},
		},
		{
			name:       "not found",
			id:         "missing",
			headers:    authHeaders(h, "u@example.com", "u"),
			mockSetup:  func() { m.On("GetLinkByID", "missing").Return(nil, nil) },
			wantStatus: 404,
		},
		{
			name:    "driver error",
			id:      "bad",
			headers: authHeaders(h, "u@example.com", "u"),
			mockSetup: func() {
				m.On("GetLinkByID", "bad").Return(nil, assert.AnError)
			},
			wantStatus: 500,
		},
		{
			name:    "empty id",
			id:      "",
			headers: authHeaders(h, "u@example.com", "u"),
			mockSetup: func() {
				m.On("GetLinkByID", "").Return(nil, nil)
			},
			wantStatus: 404,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			resetMock(m)
			c.mockSetup()
			q := "?link_id=" + c.id
			rr := doJSON(r, "GET", "/app/link"+q, nil, c.headers)
			assert.Equal(t, c.wantStatus, rr.Code)
			if c.check != nil && rr.Code == 200 {
				c.check(t, rr.Body.Bytes())
			}
			m.AssertExpectations(t)
		})
	}
}

// ---------- GET /app/member-links ----------
func TestGetMemberLinks(t *testing.T) {
	m := &linkmocks.LinkDriver{}
	h := newTestAppHandler(m)
	r := buildRouter(h)

	type tc struct {
		name       string
		headers    map[string]string
		mockSetup  func()
		wantStatus int
		wantCount  int
	}

	cases := []tc{
		{
			name:    "empty",
			headers: authHeaders(h, "u@example.com", "u"),
			mockSetup: func() {
				m.On("GetLinksForMember", "u@example.com").Return([]*models.Link{}, nil)
			},
			wantStatus: 200,
			wantCount:  0,
		},
		{
			name:    "multiple",
			headers: authHeaders(h, "u@example.com", "u"),
			mockSetup: func() {
				m.On("GetLinksForMember", "u@example.com").Return([]*models.Link{
					{ID: "1", Title: "One", ReferencedLink: "https://one.test", CreatedAt: 111},
					{ID: "2", Title: "Two", ReferencedLink: "https://two.test", CreatedAt: 222},
				}, nil)
			},
			wantStatus: 200,
			wantCount:  2,
		},
		{
			name:       "unauthorized",
			headers:    map[string]string{},
			mockSetup:  func() {},
			wantStatus: 401,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			resetMock(m)
			c.mockSetup()
			rr := doJSON(r, "GET", "/app/member-links", nil, c.headers)
			assert.Equal(t, c.wantStatus, rr.Code)

			if c.wantStatus == 200 {
				var resp struct {
					Links []struct {
						ID string `json:"id"`
					} `json:"Links"`
				}
				require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &resp))
				assert.Equal(t, c.wantCount, len(resp.Links))
			}

			m.AssertExpectations(t)
		})
	}
}

// ---------- GET /app/member-qrs ----------
func TestGetMemberQRs(t *testing.T) {
	m := &linkmocks.LinkDriver{}
	h := newTestAppHandler(m)
	r := buildRouter(h)

	resetMock(m)
	m.On("GetQRsForMember", "u@example.com").Return([]*models.Link{
		{ID: "qr1", Title: "QR One", ReferencedLink: "https://one.test", CreatedAt: 111, HasQR: true, QRLink: "https://cdn/qr1.png"},
	}, nil)

	rr := doJSON(r, "GET", "/app/member-qrs", nil, authHeaders(h, "u@example.com", "u"))
	assert.Equal(t, 200, rr.Code)
	var resp struct {
		Links []struct {
			ID        string `json:"id"`
			QRPicture string `json:"qr_picture"`
		} `json:"links"`
	}
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &resp))
	assert.Len(t, resp.Links, 1)
	assert.Equal(t, "qr1", resp.Links[0].ID)
	m.AssertExpectations(t)
}

// ---------- POST /app/create-link ----------
func TestCreateLink(t *testing.T) {
	m := &linkmocks.LinkDriver{}
	h := newTestAppHandler(m)
	r := buildRouter(h)

	// Local test server to provide a valid referenced link without external network dependency
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		_, _ = w.Write([]byte(`<html><head><title>Test Page</title><link rel="icon" href="/favicon.ico"></head><body>OK</body></html>`))
	}))
	defer ts.Close()

	type tc struct {
		name       string
		body       any
		headers    map[string]string
		mockSetup  func()
		wantStatus int
	}

	cases := []tc{
		{
			name: "success",
			body: map[string]any{
				"title":           "My Link",
				"referenced_link": ts.URL,
			},
			headers: authHeaders(h, "u@example.com", "u"),
			mockSetup: func() {
				m.On("GetLinkByID", mock.AnythingOfType("string")).Return(nil, assert.AnError)
				m.On("UpsertLink", mock.MatchedBy(func(l *models.Link) bool { return l.MemberEmail == "u@example.com" && l.ReferencedLink == ts.URL })).Return(nil)
			},
			wantStatus: 200,
		},
		{
			name: "validation error (missing referenced_link)",
			body: map[string]any{
				"title": "No URL",
			},
			headers:    authHeaders(h, "u@example.com", "u"),
			mockSetup:  func() {},
			wantStatus: 400,
		},
		{
			name: "driver error",
			body: map[string]any{
				"title":           "Err",
				"referenced_link": ts.URL,
			},
			headers: authHeaders(h, "u@example.com", "u"),
			mockSetup: func() {
				m.On("GetLinkByID", mock.AnythingOfType("string")).Return(nil, assert.AnError)
				m.On("UpsertLink", mock.AnythingOfType("*models.Link")).Return(assert.AnError)
			},
			wantStatus: 500,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			resetMock(m)
			c.mockSetup()
			rr := doJSON(r, "POST", "/app/create-link", c.body, c.headers)
			assert.Equal(t, c.wantStatus, rr.Code)
			m.AssertExpectations(t)
		})
	}
}

// ---------- POST /app/update-link ----------
func TestUpdateLink(t *testing.T) {
	m := &linkmocks.LinkDriver{}
	h := newTestAppHandler(m)
	r := buildRouter(h)

	type tc struct {
		name       string
		body       any
		headers    map[string]string
		mockSetup  func()
		wantStatus int
	}

	cases := []tc{
		{
			name: "success",
			body: map[string]any{
				"id":              "link1",
				"title":           "Updated",
				"referenced_link": "https://example.com",
			},
			headers: authHeaders(h, "u@example.com", "u"),
			mockSetup: func() {
				m.On("GetLinkByID", "link1").Return(&models.Link{ID: "link1", MemberEmail: "u@example.com", CreatedAt: 123}, nil)
				m.On("UpsertLink", mock.AnythingOfType("*models.Link")).Return(nil)
			},
			wantStatus: 200,
		},
		{
			name: "invalid id",
			body: map[string]any{
				"title": "No ID",
			},
			headers:    authHeaders(h, "u@example.com", "u"),
			mockSetup:  func() { m.On("GetLinkByID", "").Return(nil, assert.AnError) },
			wantStatus: 500,
		},
		{
			name: "driver error",
			body: map[string]any{
				"id":              "bad",
				"title":           "Err",
				"referenced_link": "https://example.com",
			},
			headers: authHeaders(h, "u@example.com", "u"),
			mockSetup: func() {
				m.On("GetLinkByID", "bad").Return(&models.Link{ID: "bad", MemberEmail: "u@example.com", CreatedAt: 123}, nil)
				m.On("UpsertLink", mock.AnythingOfType("*models.Link")).Return(assert.AnError)
			},
			wantStatus: 500,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			resetMock(m)
			c.mockSetup()
			rr := doJSON(r, "POST", "/app/update-link", c.body, c.headers)
			assert.Equal(t, c.wantStatus, rr.Code)
			m.AssertExpectations(t)
		})
	}
}

// ---------- DELETE /app/delete-link ----------
func TestDeleteLink(t *testing.T) {
	m := &linkmocks.LinkDriver{}
	h := newTestAppHandler(m)
	r := buildRouter(h)

	resetMock(m)
	m.On("DeleteLink", "link1").Return(nil)

	rr := doJSON(r, "DELETE", "/app/delete-link?link_id=link1", nil, authHeaders(h, "u@example.com", "u"))
	assert.Equal(t, 200, rr.Code)
	m.AssertExpectations(t)
}

// ---------- DELETE /app/delete-qr ----------
func TestDeleteQr(t *testing.T) {
	m := &linkmocks.LinkDriver{}
	h := newTestAppHandler(m)
	r := buildRouter(h)

	resetMock(m)
	m.On("GetLinkByID", "qr1").Return(&models.Link{ID: "qr1", MemberEmail: "u@example.com", HasQR: true, QRLink: "some"}, nil)
	m.On("UpsertLink", mock.AnythingOfType("*models.Link")).Return(nil)

	rr := doJSON(r, "DELETE", "/app/delete-qr?link_id=qr1", nil, authHeaders(h, "u@example.com", "u"))
	assert.Equal(t, 200, rr.Code)
	m.AssertExpectations(t)
}
