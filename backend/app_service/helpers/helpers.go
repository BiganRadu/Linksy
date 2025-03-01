package helpers

import (
	"backend/constants"
	"backend/models"
	"crypto/rand"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"regexp"
	"strings"
)

func RandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		result[i] = charset[num.Int64()]
	}
	return string(result)
}

func GetTitleAndIcon(url string) (string, string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("failed to fetch page, status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}
	html := string(body)

	// Extract the title using regex
	titleRe := regexp.MustCompile(`(?i)<title>(.*?)</title>`)
	titleMatches := titleRe.FindStringSubmatch(html)
	title := "Unknown Title"
	if len(titleMatches) > 1 {
		title = titleMatches[1]
	}

	// Extract the favicon URL using regex
	iconRe := regexp.MustCompile(`(?i)<link[^>]+rel=["']?icon["']?[^>]+href=["']([^"']+)["']`)
	iconMatches := iconRe.FindStringSubmatch(html)
	iconURL := ""
	if len(iconMatches) > 1 {
		iconURL = iconMatches[1]

		// Handle relative URLs
		if !strings.HasPrefix(iconURL, "http") {
			if strings.HasPrefix(iconURL, "//") {
				iconURL = "https:" + iconURL
			} else {
				iconURL = url + "/" + strings.TrimLeft(iconURL, "/")
			}
		}
	}

	return title, iconURL, nil
}

func IsIpAllowed(Ip uint32, Link *models.Link) bool {
	switch Link.AccessMode {
	case constants.LinkAccessModes.Default:
		return true
	case constants.LinkAccessModes.IpWhiteList:
		for _, ip := range Link.AllowedIps {
			if Ip == ip {
				return true
			}
		}
		return false
	case constants.LinkAccessModes.IpBlackList:
		for _, ip := range Link.AllowedIps {
			if Ip == ip {
				return false
			}
		}
		return true
	}
	return false
}
