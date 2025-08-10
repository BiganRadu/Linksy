package app_service

import (
	internal_models "backend/app_service/models"
	"sort"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetAnalytics handles the request to retrieve analytics data based on the chart type specified in the query parameter.
// It checks the chart type and calls the appropriate function to get the analytics data.
func (a *AppHandler) GetAnalytics(c *gin.Context) {
	chartType := c.Query("chart_code")
	if chartType == "sessions" {
		resp, err := a.GetSessionAnalytics(c)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, resp)
	} else if chartType == "links" {
		resp, err := a.GetLinksAnalytics(c)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, resp)
	} else if chartType == "platforms" {
		resp, err := a.GetPlatformAnalytics(c)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, resp)
	} else if chartType == "countries" {
		resp, err := a.GetCountryAnalytics(c)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, resp)
	} else {
		c.JSON(400, gin.H{"error": "Invalid chart type"})
	}
}

// GetSessionAnalytics retrieves session analytics for a member based on the provided start and end timestamps.
func (a *AppHandler) GetSessionAnalytics(c *gin.Context) (*internal_models.SessionAnalytics, error) {
	email := c.GetString("email")

	// Parse start and end timestamps from query parameters
	startStr := c.Query("start")
	endStr := c.Query("end")
	start, _ := strconv.ParseInt(startStr, 10, 64)
	end, _ := strconv.ParseInt(endStr, 10, 64)

	links, err := a.linkDriver.GetLinksForMember(email)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return nil, err
	}

	// Calculate the number of days between start and end timestamps
	days := (end - start) / 3600 / 24

	// Iterate through the links and their access entries to populate the response
	response := internal_models.SessionAnalytics{}
	response.SessionsPerDay = make([]int, days)
	for _, link := range links {
		for _, entry := range link.AccessEntries {
			if entry.HourStart >= start && entry.HourStart < end {
				dayIndex := (entry.HourStart - start) / 3600 / 24
				if dayIndex >= 0 && dayIndex < days {
					response.SessionsPerDay[dayIndex] += entry.Accesses
					response.TotalSessions += entry.Accesses
				}
			}
		}
	}

	return &response, nil
}

// GetLinksAnalytics retrieves analytics for each link of a member based on the provided start and end timestamps.
func (a *AppHandler) GetLinksAnalytics(c *gin.Context) (*internal_models.LinksAnalytics, error) {
	email := c.GetString("email")

	// Parse start and end timestamps from query parameters
	startStr := c.Query("start")
	endStr := c.Query("end")
	start, _ := strconv.ParseInt(startStr, 10, 64)
	end, _ := strconv.ParseInt(endStr, 10, 64)

	links, err := a.linkDriver.GetLinksForMember(email)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return nil, err
	}

	response := internal_models.LinksAnalytics{}
	index := 0
	for _, link := range links {
		platformMap := make(map[string]int)
		countryMap := make(map[string]int)
		linkAnalytics := internal_models.LinkAnalytics{
			LinkTitle:      link.Title,
			LinkId:         index,
			SessionsPerDay: make([]int, (end-start)/3600/24),
		}

		for _, entry := range link.AccessEntries {
			if entry.HourStart >= start && entry.HourStart < end {
				dayIndex := (entry.HourStart - start) / 3600 / 24
				if dayIndex >= 0 {
					linkAnalytics.SessionsPerDay[dayIndex] += entry.Accesses
					linkAnalytics.TotalSessions += entry.Accesses
					for country, count := range entry.CountryMap {
						countryMap[country] += count
					}
					for platform, count := range entry.OsMap {
						platformMap[platform] += count
					}
				}
			}
		}

		// Determine the most frequent platform and country for the link
		bestPlatform := ""
		bestPlatformCount := 0
		for platform, count := range platformMap {
			if count > bestPlatformCount {
				bestPlatformCount = count
				bestPlatform = platform
			}
		}
		linkAnalytics.LinkPlatform = bestPlatform
		bestCountry := ""
		bestCountryCount := 0
		for country, count := range countryMap {
			if count > bestCountryCount {
				bestCountryCount = count
				bestCountry = country
			}
		}
		linkAnalytics.LinkCountry = bestCountry

		response.Links = append(response.Links, linkAnalytics)
		index++
	}

	return &response, nil
}

// GetPlatformAnalytics retrieves platform analytics for a member based on the provided start and end timestamps.
func (a *AppHandler) GetPlatformAnalytics(c *gin.Context) (*internal_models.PieAnalytics, error) {
	email := c.GetString("email")

	// Parse start and end timestamps from query parameters
	startStr := c.Query("start")
	endStr := c.Query("end")
	start, _ := strconv.ParseInt(startStr, 10, 64)
	end, _ := strconv.ParseInt(endStr, 10, 64)

	links, err := a.linkDriver.GetLinksForMember(email)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return nil, err
	}

	// Create a map to hold platform counts
	platformMap := make(map[string]int)
	for _, link := range links {
		for _, entry := range link.AccessEntries {
			if entry.HourStart >= start && entry.HourStart < end {
				for platform, count := range entry.OsMap {
					platformMap[platform] += count
				}
			}
		}
	}

	response := internal_models.PieAnalytics{
		Total:  0,
		Values: make([]internal_models.Pair, 0),
	}

	for platform, count := range platformMap {
		response.Total += count
		response.Values = append(response.Values, internal_models.Pair{Name: platform, Value: count})
	}

	// Sort the values by count in descending order
	sort.Slice(response.Values, func(i, j int) bool {
		return response.Values[i].Value > response.Values[j].Value
	})

	return &response, nil
}

// GetCountryAnalytics retrieves country analytics for a member based on the provided start and end timestamps.
func (a *AppHandler) GetCountryAnalytics(c *gin.Context) (*internal_models.PieAnalytics, error) {
	email := c.GetString("email")

	// Parse start and end timestamps from query parameters
	startStr := c.Query("start")
	endStr := c.Query("end")
	start, _ := strconv.ParseInt(startStr, 10, 64)
	end, _ := strconv.ParseInt(endStr, 10, 64)

	links, err := a.linkDriver.GetLinksForMember(email)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return nil, err
	}

	// Create a map to hold country counts
	countryMap := make(map[string]int)
	for _, link := range links {
		for _, entry := range link.AccessEntries {
			if entry.HourStart >= start && entry.HourStart < end {
				for country, count := range entry.CountryMap {
					countryMap[country] += count
				}
			}
		}
	}

	response := internal_models.PieAnalytics{
		Total:  0,
		Values: make([]internal_models.Pair, 0),
	}

	for country, count := range countryMap {
		response.Total += count
		response.Values = append(response.Values, internal_models.Pair{Name: country, Value: count})
	}

	// Sort the values by count in descending order
	sort.Slice(response.Values, func(i, j int) bool {
		return response.Values[i].Value > response.Values[j].Value
	})

	return &response, nil
}
