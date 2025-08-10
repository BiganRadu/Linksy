package models

// SessionAnalytics represents the analytics data for user sessions.
type SessionAnalytics struct {
	SessionsPerDay []int `json:"sessions"`
	TotalSessions  int   `json:"total"`
}

// LinkAnalytics represents the analytics data for a specific link.
type LinkAnalytics struct {
	LinkId         int    `json:"id"`
	LinkTitle      string `json:"title"`
	LinkPlatform   string `json:"platform"`
	LinkCountry    string `json:"country"`
	SessionsPerDay []int  `json:"sessions"`
	TotalSessions  int    `json:"total"`
}

// LinksAnalytics represents the analytics data for multiple links.
type LinksAnalytics struct {
	Links []LinkAnalytics `json:"links"`
}

// Pair represents a key-value pair used in analytics data.
type Pair struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

// PieAnalytics represents the analytics data for pie charts.
type PieAnalytics struct {
	Total  int    `json:"total"`
	Values []Pair `json:"values"`
}
