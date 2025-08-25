package app_service

import (
	internal_models "backend/app_service/models"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	linkmocks "backend/drivers/LinkDriver/mocks"
	"backend/helpers"
	"backend/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// newTestAnalyticsHandler reuses AppHandler with only linkDriver & tokenHelper relevant here.
func newTestAnalyticsHandler(m *linkmocks.LinkDriver) *AppHandler {
	return &AppHandler{
		linkDriver:  m,
		tokenHelper: helpers.NewTokenHelper("test_secret"),
		nowFunc: func() time.Time {
			return time.Unix(1700000000, 0)
		},
	}
}

func buildAnalyticsRouter(h *AppHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	api := r.Group("/app")
	h.Routes(api)
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

func authHeader(h *AppHandler, email, username string) map[string]string {
	tok, _ := h.tokenHelper.GenerateToken(email, username, h.nowFunc().Unix(), time.Hour)
	return map[string]string{"AuthToken": tok}
}

func resetLinkMock(m *linkmocks.LinkDriver) {
	m.ExpectedCalls = nil
	m.Calls = nil
}

func TestGetAnalytics_ByChartType(t *testing.T) {
	m := &linkmocks.LinkDriver{}
	h := newTestAnalyticsHandler(m)
	r := buildAnalyticsRouter(h)

	start := int64(1700000000)
	end := start + 2*24*3600
	testLinks := []*models.Link{
		{
			Title: "First",
			AccessEntries: []models.AccessEntry{
				{
					HourStart:  start,
					Accesses:   3,
					CountryMap: map[string]int{"US": 2, "DE": 1},
					OsMap:      map[string]int{"ios": 2, "android": 1},
				},
				{
					HourStart:  start + 24*3600,
					Accesses:   5,
					CountryMap: map[string]int{"US": 1, "FR": 4},
					OsMap:      map[string]int{"android": 5},
				},
			},
		},
		{
			Title: "Second",
			AccessEntries: []models.AccessEntry{
				{
					HourStart:  start + 6*3600,
					Accesses:   2,
					CountryMap: map[string]int{"US": 2},
					OsMap:      map[string]int{"web": 2},
				},
			},
		},
	}

	type tc struct {
		name       string
		chartCode  string
		mockSetup  func()
		wantStatus int
		check      func(t *testing.T, body []byte)
	}

	commonMock := func() {
		m.On("GetLinksForMember", mock.AnythingOfType("string")).Return(testLinks, nil)
	}

	cases := []tc{
		{
			name:       "sessions",
			chartCode:  "sessions",
			mockSetup:  commonMock,
			wantStatus: 200,
			check: func(t *testing.T, body []byte) {
				var obj map[string]any
				_ = json.Unmarshal(body, &obj)
				assert.Contains(t, obj, "sessions")
			},
		},
		{
			name:       "links",
			chartCode:  "links",
			mockSetup:  commonMock,
			wantStatus: 200,
			check: func(t *testing.T, body []byte) {
				var obj map[string]any
				_ = json.Unmarshal(body, &obj)
				assert.Contains(t, obj, "links")
			},
		},
		{
			name:       "platforms",
			chartCode:  "platforms",
			mockSetup:  commonMock,
			wantStatus: 200,
			check: func(t *testing.T, body []byte) {
				var obj map[string]any
				_ = json.Unmarshal(body, &obj)
				assert.Contains(t, obj, "values")
			},
		},
		{
			name:       "countries",
			chartCode:  "countries",
			mockSetup:  commonMock,
			wantStatus: 200,
			check: func(t *testing.T, body []byte) {
				var obj map[string]any
				_ = json.Unmarshal(body, &obj)
				assert.Contains(t, obj, "values")
			},
		},
		{
			name:       "invalid chart",
			chartCode:  "unknown",
			mockSetup:  func() {},
			wantStatus: 400,
		},
	}

	headers := authHeader(h, "user@example.com", "user")

	for _, cse := range cases {
		t.Run(cse.name, func(t *testing.T) {
			resetLinkMock(m)
			cse.mockSetup()
			url := "/app/analytics?chart_code=" + intToStrParam(cse.chartCode) +
				"&chart_code=" + cse.chartCode +
				"&start=" + intToStr(start) +
				"&end=" + intToStr(end)
			// (If duplicate chart_code param is unintended, simplify above line.)
			url = "/app/analytics?chart_code=" + cse.chartCode + "&start=" + intToStr(start) + "&end=" + intToStr(end)
			rr := doReq(r, "GET", url, nil, headers)
			assert.Equal(t, cse.wantStatus, rr.Code)
			m.AssertExpectations(t)
			if cse.check != nil && rr.Code == 200 {
				cse.check(t, rr.Body.Bytes())
			}
		})
	}
}

func buildDirectTestLinks(start int64) []*models.Link {
	day0 := start
	day1 := start + 24*3600
	return []*models.Link{
		{
			Title: "First",
			AccessEntries: []models.AccessEntry{
				{
					HourStart:  day0,
					Accesses:   3,
					CountryMap: map[string]int{"US": 2, "DE": 1},
					OsMap:      map[string]int{"ios": 2, "android": 1},
				},
				{
					HourStart:  day1,
					Accesses:   5,
					CountryMap: map[string]int{"US": 1, "FR": 4},
					OsMap:      map[string]int{"android": 5},
				},
			},
		},
		{
			Title: "Second",
			AccessEntries: []models.AccessEntry{
				{
					HourStart:  day0 + 6*3600,
					Accesses:   2,
					CountryMap: map[string]int{"US": 2},
					OsMap:      map[string]int{"web": 2},
				},
			},
		},
	}
}

func newDirectCtx(email string, start, end int64) (*gin.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req, _ := http.NewRequest("GET", "/x?start="+intToStr(start)+"&end="+intToStr(end), nil)
	c.Request = req
	c.Set("email", email)
	return c, w
}

func TestGetSessionAnalytics_Direct(t *testing.T) {
	m := &linkmocks.LinkDriver{}
	h := newTestAnalyticsHandler(m)
	start := int64(1700000000)
	end := start + 2*24*3600
	links := buildDirectTestLinks(start)

	m.On("GetLinksForMember", "user@example.com").Return(links, nil)

	c, _ := newDirectCtx("user@example.com", start, end)
	resp, err := h.GetSessionAnalytics(c)
	require.NoError(t, err)
	require.NotNil(t, resp)

	assert.IsType(t, &internal_models.SessionAnalytics{}, resp)
	require.Len(t, resp.SessionsPerDay, 2)
	assert.Equal(t, []int{5, 5}, resp.SessionsPerDay)
	assert.Equal(t, 10, resp.TotalSessions)

	m.AssertExpectations(t)
}

func TestGetLinksAnalytics_Direct(t *testing.T) {
	m := &linkmocks.LinkDriver{}
	h := newTestAnalyticsHandler(m)
	start := int64(1700000000)
	end := start + 2*24*3600
	links := buildDirectTestLinks(start)

	m.On("GetLinksForMember", "user@example.com").Return(links, nil)

	c, _ := newDirectCtx("user@example.com", start, end)
	resp, err := h.GetLinksAnalytics(c)
	require.NoError(t, err)
	require.NotNil(t, resp)

	require.Len(t, resp.Links, 2)

	l0 := resp.Links[0]
	assert.Equal(t, "First", l0.LinkTitle)
	assert.Equal(t, 0, l0.LinkId)
	assert.Equal(t, []int{3, 5}, l0.SessionsPerDay)
	assert.Equal(t, 8, l0.TotalSessions)
	assert.Equal(t, "android", l0.LinkPlatform)
	assert.Equal(t, "FR", l0.LinkCountry)

	l1 := resp.Links[1]
	assert.Equal(t, "Second", l1.LinkTitle)
	assert.Equal(t, 1, l1.LinkId)
	assert.Equal(t, []int{2, 0}, l1.SessionsPerDay)
	assert.Equal(t, 2, l1.TotalSessions)
	assert.Equal(t, "web", l1.LinkPlatform)
	assert.Equal(t, "US", l1.LinkCountry)

	m.AssertExpectations(t)
}

func TestGetPlatformAnalytics_Direct(t *testing.T) {
	m := &linkmocks.LinkDriver{}
	h := newTestAnalyticsHandler(m)
	start := int64(1700000000)
	end := start + 2*24*3600
	links := buildDirectTestLinks(start)

	m.On("GetLinksForMember", "user@example.com").Return(links, nil)

	c, _ := newDirectCtx("user@example.com", start, end)
	resp, err := h.GetPlatformAnalytics(c)
	require.NoError(t, err)
	require.NotNil(t, resp)

	assert.Equal(t, 10, resp.Total)
	got := map[string]int{}
	for _, p := range resp.Values {
		got[p.Name] = p.Value
	}
	assert.Equal(t, 6, got["android"])
	assert.Equal(t, 2, got["ios"])
	assert.Equal(t, 2, got["web"])

	m.AssertExpectations(t)
}

func TestGetCountryAnalytics_Direct(t *testing.T) {
	m := &linkmocks.LinkDriver{}
	h := newTestAnalyticsHandler(m)
	start := int64(1700000000)
	end := start + 2*24*3600
	links := buildDirectTestLinks(start)

	m.On("GetLinksForMember", "user@example.com").Return(links, nil)

	c, _ := newDirectCtx("user@example.com", start, end)
	resp, err := h.GetCountryAnalytics(c)
	require.NoError(t, err)
	require.NotNil(t, resp)

	assert.Equal(t, 10, resp.Total)
	got := map[string]int{}
	for _, p := range resp.Values {
		got[p.Name] = p.Value
	}
	assert.Equal(t, 5, got["US"])
	assert.Equal(t, 4, got["FR"])
	assert.Equal(t, 1, got["DE"])

	m.AssertExpectations(t)
}

// intToStr helper (kept manual to avoid adding strconv if undesired).
func intToStr(v int64) string {
	return int64ToString(v)
}

func int64ToString(v int64) string {
	if v == 0 {
		return "0"
	}
	buf := [20]byte{}
	i := len(buf)
	for v > 0 {
		i--
		buf[i] = byte('0' + v%10)
		v /= 10
	}
	return string(buf[i:])
}

// dummy to silence unused function if earlier build line retained
func intToStrParam(s string) string { return s }
