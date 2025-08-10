package app_service

import (
	"backend/app_service/helpers"
	internal_models "backend/app_service/models"
	"backend/models"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

// GetRedirect handles the redirection request for a link.
func (a *AppHandler) GetRedirect(c *gin.Context) {
	// Get the link ID from the query string
	linkID := c.Query("link_id")

	var linkRequest internal_models.RedirectLinkRequest
	err := c.BindJSON(&linkRequest)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	link, err := a.linkDriver.GetLinkByID(linkID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if link == nil {
		c.JSON(404, gin.H{"error": "Link not found"})
		return
	}

	// Check if the link is private and if the IP is allowed
	if !helpers.IsIpAllowed(linkRequest.Ip, link) {
		c.JSON(403, gin.H{"error": "This ip is not allowed to access this link"})
		return
	}

	a.AddAccessEntry(&linkRequest, link)

	err = a.linkDriver.UpsertLink(link)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	response := internal_models.RedirectLinkResponse{
		RedirectLink: link.ReferencedLink,
	}

	c.JSON(200, response)
}

// GetMemberInfo retrieves the member's information from the context and returns it as a JSON response.
func (a *AppHandler) GetMemberInfo(c *gin.Context) {
	var memberInfo internal_models.MemberInfoResponse

	memberInfo.Email = c.GetString("email")
	memberInfo.Username = c.GetString("username")
	memberInfo.CreatedAt = c.GetInt64("created_at")

	c.JSON(200, memberInfo)
}

// GetLink retrieves a link by its ID and returns it as a JSON response.
func (a *AppHandler) GetLink(c *gin.Context) {
	linkID := c.Query("link_id")

	link, err := a.linkDriver.GetLinkByID(linkID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if link == nil {
		c.JSON(404, gin.H{"error": "Link not found"})
		return
	}

	c.JSON(200, link)
}

// GetMemberLinks retrieves all links for a member and returns them as a JSON response.
func (a *AppHandler) GetMemberLinks(c *gin.Context) {
	email := c.GetString("email")

	links, err := a.linkDriver.GetLinksForMember(email)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	var memberLinks []*internal_models.LinkDetails
	for _, link := range links {
		memberLinks = append(memberLinks, &internal_models.LinkDetails{
			ID:             link.ID,
			Title:          link.Title,
			Icon:           link.Icon,
			ReferencedLink: link.ReferencedLink,
			CreatedAt:      link.CreatedAt,
		})
	}

	c.JSON(200, internal_models.MemberLinksResponse{
		Links: memberLinks,
	})
}

// GetMemberQRs retrieves all QR codes for a member and returns them as a JSON response.
func (a *AppHandler) GetMemberQRs(c *gin.Context) {
	email := c.GetString("email")

	links, err := a.linkDriver.GetQRsForMember(email)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// Construct the response with link details
	var memberLinks []*internal_models.LinkDetails
	for _, link := range links {
		memberLinks = append(memberLinks, &internal_models.LinkDetails{
			ID:             link.ID,
			Title:          link.Title,
			Icon:           link.Icon,
			ReferencedLink: link.ReferencedLink,
			CreatedAt:      link.CreatedAt,
			QRPicture:      link.QRLink,
		})
	}

	c.JSON(200, internal_models.MemberLinksResponse{
		Links: memberLinks,
	})
}

// CreateLink creates a new link for a member and returns a success response.
func (a *AppHandler) CreateLink(c *gin.Context) {
	var link models.Link
	if err := c.BindJSON(&link); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}

	link.MemberEmail = c.GetString("email")
	link.ID = helpers.RandomString(6)
	link.CreatedAt = a.nowFunc().Unix()

	// Retrieve the title and icon for the referenced link from the URL
	title, icon, err := helpers.GetTitleAndIcon(link.ReferencedLink)
	if err != nil {
		c.JSON(400, gin.H{"error": "Failed to fetch details about the referenced link"})
		return
	}
	if link.Title == "" {
		link.Title = title
	}
	link.Icon = icon

	// Ensure the link ID is unique
	_, err = a.linkDriver.GetLinkByID(link.ID)
	for err == nil {
		link.ID = helpers.RandomString(6)
		_, err = a.linkDriver.GetLinkByID(link.ID)
	}

	// If the link should have a QR code, create it and save the image in S3
	if link.HasQR == true {
		qrLink, err := a.CreateQr(link.ID)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to create QR code"})
			return
		}
		link.QRLink = qrLink
	}

	err = a.linkDriver.UpsertLink(&link)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, nil)
}

// UpdateLink updates an existing link for a member and returns a success response.
func (a *AppHandler) UpdateLink(c *gin.Context) {
	var link models.Link
	if err := c.BindJSON(&link); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	currentSavedLink, err := a.linkDriver.GetLinkByID(link.ID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Invalid link"})
		return
	}

	// Copy the fields that should not be updated
	link.CreatedAt = currentSavedLink.CreatedAt
	link.MemberEmail = currentSavedLink.MemberEmail
	link.HasQR = currentSavedLink.HasQR
	link.QRLink = currentSavedLink.QRLink
	link.AccessEntries = currentSavedLink.AccessEntries

	// Retrieve the title and icon for the referenced link from the URL
	title, icon, err := helpers.GetTitleAndIcon(link.ReferencedLink)
	if err != nil {
		c.JSON(400, gin.H{"error": "Failed to fetch details about the referenced link"})
		return
	}
	if link.Title == "" {
		link.Title = title
	}
	link.Icon = icon

	err = a.linkDriver.UpsertLink(&link)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, nil)
}

// DeleteLink deletes a link by its ID and returns a success response.
func (a *AppHandler) DeleteLink(c *gin.Context) {
	linkID := c.Query("link_id")

	err := a.linkDriver.DeleteLink(linkID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, nil)
}

// DeleteQr deletes the QR code for a link by its ID and returns a success response.
func (a *AppHandler) DeleteQr(c *gin.Context) {
	linkID := c.Query("link_id")

	link, err := a.linkDriver.GetLinkByID(linkID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	link.QRLink = ""
	link.HasQR = false
	err = a.linkDriver.UpsertLink(link)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, nil)
}

// CreateQr generates a QR code for a link ID and uploads it to S3, returning the URL of the uploaded QR code.
func (a *AppHandler) CreateQr(ID string) (string, error) {
	err := helpers.GenerateQRCode("bit.ly/"+ID, ID+".png", 256)
	if err != nil {
		return "", err
	}

	return helpers.UploadToS3(a.awsClient, a.s3Bucket, "qrcodes/"+ID+".png", "/home/raduzew/CS2023-2027/GO/Linksy/backend/pictures/"+ID+".png")
}

// AddAccessEntry adds an access entry for a link based on the request data.
func (a *AppHandler) AddAccessEntry(req *internal_models.RedirectLinkRequest, Link *models.Link) {
	left, right := 0, len(Link.AccessEntries)-1
	startOfHour := time.Date(a.nowFunc().Year(), a.nowFunc().Month(), a.nowFunc().Day(), a.nowFunc().Hour(), 0, 0, 0, time.UTC).Unix()
	foundEntryIndex := -1
	for left <= right {
		mid := left + (right-left)/2
		if Link.AccessEntries[mid].HourStart == startOfHour {
			foundEntryIndex = mid
			break
		} else if Link.AccessEntries[mid].HourStart < startOfHour {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}

	// If no entry found for the current hour, create a new one
	if foundEntryIndex == -1 {
		Link.AccessEntries = append(Link.AccessEntries, models.AccessEntry{
			HourStart:  startOfHour,
			CountryMap: map[string]int{},
			DeviceMap:  map[string]int{},
			OsMap:      map[string]int{},
			Accesses:   0,
		})
		foundEntryIndex = len(Link.AccessEntries) - 1
	}

	Link.AccessEntries[foundEntryIndex].Accesses++
	Link.AccessEntries[foundEntryIndex].CountryMap[req.Country]++
	Link.AccessEntries[foundEntryIndex].DeviceMap[req.Device]++
	Link.AccessEntries[foundEntryIndex].OsMap[req.Os]++
}
