package helpers

import (
	"backend/constants"
	"backend/models"
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/skip2/go-qrcode"
	"golang.org/x/net/html"
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
func GetTitleAndIcon(pageURL string) (string, string, error) {
	resp, err := http.Get(pageURL)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("failed to fetch page, status code: %d", resp.StatusCode)
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return "", "", err
	}

	var title, icon string
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "title" && n.FirstChild != nil {
			title = n.FirstChild.Data
		}

		if n.Type == html.ElementNode && n.Data == "link" {
			rel := ""
			href := ""
			for _, attr := range n.Attr {
				switch strings.ToLower(attr.Key) {
				case "rel":
					rel = strings.ToLower(attr.Val)
				case "href":
					href = attr.Val
				}
			}
			if (strings.Contains(rel, "icon") || strings.Contains(rel, "shortcut icon")) && href != "" && icon == "" {
				icon = resolveURL(pageURL, href)
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	if icon == "" {
		// fallback to /favicon.ico
		parsed, err := url.Parse(pageURL)
		if err == nil {
			icon = parsed.Scheme + "://" + parsed.Host + "/favicon.ico"
		}
	}

	if title == "" {
		title = "Unknown Title"
	}

	return title, icon, nil
}

func resolveURL(base, ref string) string {
	u, err := url.Parse(ref)
	if err != nil {
		return ref
	}
	if u.IsAbs() {
		return u.String()
	}
	baseURL, err := url.Parse(base)
	if err != nil {
		return ref
	}
	return baseURL.ResolveReference(u).String()
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
