package member_service

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"backend/helpers"
	"backend/models"

	membermocks "backend/drivers/member_driver/mocks"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func newTestHandler(md *membermocks.MemberDriver) *MemberHandler {
	return &MemberHandler{
		memberDriver: md,
		tokenHelper:  helpers.NewTokenHelper("test_secret"),
		nowFunc:      func() time.Time { return time.Unix(1700000000, 0) },
	}
}

func buildRouter(h *MemberHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/register", h.Register)
	r.POST("/login", h.Login)
	r.POST("/logout", h.Logout)
	r.POST("/change-password", h.ChangePassword)
	r.POST("/change-name", h.ChangeName)
	return r
}

func doReq(r http.Handler, method, path string, body any, headers map[string]string) *httptest.ResponseRecorder {
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

func resetMock(m *membermocks.MemberDriver) {
	m.ExpectedCalls = nil
	m.Calls = nil
}

func TestRegister(t *testing.T) {
	mockDrv := &membermocks.MemberDriver{}
	h := newTestHandler(mockDrv)
	r := buildRouter(h)

	type tc struct {
		name       string
		body       any
		mockSetup  func()
		wantStatus int
		check      func(*testing.T)
	}

	var captured *models.Member

	cases := []tc{
		{
			name: "success",
			body: map[string]any{
				"email":    "user@example.com",
				"username": "tester",
				"password": "plainpass",
			},
			mockSetup: func() {
				mockDrv.On("CountMembersWithEmail", "user@example.com").Return(0, nil)
				mockDrv.
					On("UpsertMember", mock.MatchedBy(func(m *models.Member) bool {
						captured = m
						return m.Email == "user@example.com" && m.Username == "tester" && m.Token != "" && !m.ID.IsZero()
					})).
					Return(nil)
			},
			wantStatus: 200,
			check: func(t *testing.T) {
				require.NotNil(t, captured)
				assert.NoError(t, bcrypt.CompareHashAndPassword([]byte(captured.Password), []byte("plainpass")))
			},
		},
		{
			name: "duplicate email",
			body: map[string]any{
				"email":    "dup@example.com",
				"username": "dup",
				"password": "x",
			},
			mockSetup: func() {
				mockDrv.On("CountMembersWithEmail", "dup@example.com").Return(1, nil)
			},
			wantStatus: 400,
		},
	}

	for _, cse := range cases {
		t.Run(cse.name, func(t *testing.T) {
			resetMock(mockDrv)
			captured = nil
			cse.mockSetup()
			rr := doReq(r, "POST", "/register", cse.body, nil)
			assert.Equal(t, cse.wantStatus, rr.Code)
			if cse.check != nil {
				cse.check(t)
			}
			mockDrv.AssertExpectations(t)
		})
	}
}

func TestLogin(t *testing.T) {
	mockDrv := &membermocks.MemberDriver{}
	h := newTestHandler(mockDrv)
	r := buildRouter(h)

	hashedOK, _ := bcrypt.GenerateFromPassword([]byte("secret"), bcrypt.DefaultCost)
	hashedWrong, _ := bcrypt.GenerateFromPassword([]byte("right"), bcrypt.DefaultCost)

	type tc struct {
		name        string
		body        any
		mockSetup   func()
		wantStatus  int
		expectToken bool
	}

	cases := []tc{
		{
			name: "success",
			body: map[string]any{
				"email":    "login@example.com",
				"password": "secret",
			},
			mockSetup: func() {
				member := &models.Member{
					Email:     "login@example.com",
					Username:  "u1",
					Password:  string(hashedOK),
					CreatedAt: h.nowFunc().Unix(),
				}
				mockDrv.On("GetMemberByEmail", "login@example.com").Return(member, nil)
				mockDrv.On("UpsertMember", mock.AnythingOfType("*models.Member")).Return(nil)
			},
			wantStatus:  200,
			expectToken: true,
		},
		{
			name: "invalid password",
			body: map[string]any{
				"email":    "lp@example.com",
				"password": "wrong",
			},
			mockSetup: func() {
				member := &models.Member{
					Email:    "lp@example.com",
					Username: "u",
					Password: string(hashedWrong),
				}
				mockDrv.On("GetMemberByEmail", "lp@example.com").Return(member, nil)
			},
			wantStatus: 400,
		},
	}

	for _, cse := range cases {
		t.Run(cse.name, func(t *testing.T) {
			resetMock(mockDrv)
			cse.mockSetup()
			rr := doReq(r, "POST", "/login", cse.body, nil)
			assert.Equal(t, cse.wantStatus, rr.Code)
			if cse.expectToken && rr.Code == 200 {
				var resp map[string]string
				require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &resp))
				assert.NotEmpty(t, resp["token"])
			}
			mockDrv.AssertExpectations(t)
		})
	}
}

func TestLogout(t *testing.T) {
	mockDrv := &membermocks.MemberDriver{}
	h := newTestHandler(mockDrv)
	r := buildRouter(h)

	token, _ := h.tokenHelper.GenerateToken("lg@example.com", "u", h.nowFunc().Unix(), time.Hour)

	type tc struct {
		name       string
		headers    map[string]string
		mockSetup  func()
		wantStatus int
	}

	cases := []tc{
		{
			name: "success",
			headers: map[string]string{
				"AuthToken": token,
			},
			mockSetup: func() {
				mockDrv.On("SetTokenForMember", "lg@example.com", "").Return(nil)
			},
			wantStatus: 200,
		},
		{
			name:       "missing token",
			headers:    map[string]string{},
			mockSetup:  func() {},
			wantStatus: 401,
		},
	}

	for _, cse := range cases {
		t.Run(cse.name, func(t *testing.T) {
			resetMock(mockDrv)
			cse.mockSetup()
			rr := doReq(r, "POST", "/logout", nil, cse.headers)
			assert.Equal(t, cse.wantStatus, rr.Code)
			mockDrv.AssertExpectations(t)
		})
	}
}

func TestChangePassword(t *testing.T) {
	mockDrv := &membermocks.MemberDriver{}
	h := newTestHandler(mockDrv)
	r := buildRouter(h)

	oldHashed, _ := bcrypt.GenerateFromPassword([]byte("oldpass"), bcrypt.DefaultCost)
	member := &models.Member{
		Username: "cp",
		Email:    "cp@example.com",
		Password: string(oldHashed),
	}
	token, _ := h.tokenHelper.GenerateToken(member.Email, "u", h.nowFunc().Unix(), time.Hour)

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
			body: map[string]string{
				"old_password": "oldpass",
				"new_password": "newpass123",
			},
			headers: map[string]string{"AuthToken": token},
			mockSetup: func() {
				mockDrv.On("GetMemberByEmail", "cp@example.com").Return(member, nil)
				mockDrv.On("UpsertMember", mock.MatchedBy(func(m *models.Member) bool {
					// Email must match; new password must NOT match old hash; must match new plain password when hashed.
					if m.Email != "cp@example.com" {
						return false
					}
					if bcrypt.CompareHashAndPassword([]byte(m.Password), []byte("newpass123")) != nil {
						return false
					}
					return true
				})).Return(nil)
			},
			wantStatus: 200,
		},
		{
			name: "old password wrong",
			body: map[string]string{
				"oldPassword": "WRONG",
				"newPassword": "newpass123",
			},
			headers: map[string]string{"AuthToken": token},
			mockSetup: func() {
				mockDrv.On("GetMemberByEmail", "cp@example.com").Return(member, nil)
			},
			wantStatus: 400,
		},
		{
			name: "missing token",
			body: map[string]string{
				"oldPassword": "oldpass",
				"newPassword": "x",
			},
			headers:    map[string]string{},
			mockSetup:  func() {},
			wantStatus: 401,
		},
	}

	for _, cse := range cases {
		t.Run(cse.name, func(t *testing.T) {
			resetMock(mockDrv)
			cse.mockSetup()
			rr := doReq(r, "POST", "/change-password", cse.body, cse.headers)
			assert.Equal(t, cse.wantStatus, rr.Code)
			mockDrv.AssertExpectations(t)
		})
	}
}

func TestChangeName(t *testing.T) {
	mockDrv := &membermocks.MemberDriver{}
	h := newTestHandler(mockDrv)
	r := buildRouter(h)

	member := &models.Member{
		Email:    "cn@example.com",
		Username: "before",
	}
	token, _ := h.tokenHelper.GenerateToken(member.Email, member.Username, h.nowFunc().Unix(), time.Hour)

	type tc struct {
		name       string
		body       any
		headers    map[string]string
		mockSetup  func()
		wantStatus int
		check      func(*testing.T)
	}

	var updated *models.Member

	cases := []tc{
		{
			name: "success",
			body: map[string]string{
				"new_name": "after",
			},
			headers: map[string]string{"AuthToken": token},
			mockSetup: func() {
				mockDrv.On("GetMemberByEmail", "cn@example.com").Return(member, nil)
				mockDrv.On("UpsertMember", mock.MatchedBy(func(m *models.Member) bool {
					updated = m
					return m.Username == "after"
				})).Return(nil)
			},
			wantStatus: 200,
			check: func(t *testing.T) {
				require.NotNil(t, updated)
				assert.Equal(t, "after", updated.Username)
			},
		},
		{
			name: "missing token",
			body: map[string]string{
				"new_name": "x",
			},
			headers:    map[string]string{},
			mockSetup:  func() {},
			wantStatus: 401,
		},
	}

	for _, cse := range cases {
		t.Run(cse.name, func(t *testing.T) {
			resetMock(mockDrv)
			updated = nil
			cse.mockSetup()
			rr := doReq(r, "POST", "/change-name", cse.body, cse.headers)
			assert.Equal(t, cse.wantStatus, rr.Code)
			if cse.check != nil {
				cse.check(t)
			}
			mockDrv.AssertExpectations(t)
		})
	}
}
