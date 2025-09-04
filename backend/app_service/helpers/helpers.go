package helpers

import (
	"backend/constants"
	"backend/models"
	"context"
	"crypto/rand"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/skip2/go-qrcode"
)

// RandomString generates a random alphanumeric string of the specified length.
func RandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		result[i] = charset[num.Int64()]
	}
	return string(result)
}

// GetTitleAndIcon fetches the title and icon URL from a given webpage URL.
// It returns the title, icon URL, and any error encountered during the process.
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

// IsIpAllowed checks if a given IP address is allowed to access a link based on the link's access mode and allowed IPs.
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

// GenerateQRCode generates a QR code for the given text and saves it to the specified filename.
func GenerateQRCode(text, filename string, size int) error {
	err := qrcode.WriteFile(text, qrcode.Medium, size, "/app/pictures/"+filename)
	if err != nil {
		fmt.Println("Failed to generate QR code: %v", err)
		return err
	}
	return nil
}

// UploadToS3 uploads a file to an AWS S3 bucket and returns the URL of the uploaded object.
func UploadToS3(awsClient *s3.Client, bucketName, objectKey, filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	_, err = awsClient.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
		Body:   file,
	})
	if err != nil {
		return "", fmt.Errorf("upload failed: %w", err)
	}

	return fmt.Sprintf("https://%s.s3.amazonaws.com/%s", bucketName, objectKey), nil
}
