package app_service

import (
	"backend/app_service/helpers"
	internal_models "backend/app_service/models"
	"backend/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

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

func (a *AppHandler) GetMemberInfo(c *gin.Context) {
	var memberInfo internal_models.MemberInfoResponse

	memberInfo.Email = c.GetString("email")
	memberInfo.Username = c.GetString("username")
	memberInfo.CreatedAt = c.GetInt64("created_at")

	c.JSON(200, memberInfo)
}

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

	title, icon, err := helpers.GetTitleAndIcon(link.ReferencedLink)
	if err != nil {
		c.JSON(400, gin.H{"error": "Failed to fetch details about the referenced link"})
		return
	}
	if link.Title == "" {
		link.Title = title
	}
	link.Icon = icon

	_, err = a.linkDriver.GetLinkByID(link.ID)
	for err == nil {
		link.ID = helpers.RandomString(6)
		_, err = a.linkDriver.GetLinkByID(link.ID)
	}

	err = a.linkDriver.UpsertLink(&link)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, nil)
}

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

	link.CreatedAt = currentSavedLink.CreatedAt
	link.MemberEmail = currentSavedLink.MemberEmail

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

func (a *AppHandler) DeleteLink(c *gin.Context) {
	linkID := c.Query("link_id")

	err := a.linkDriver.DeleteLink(linkID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, nil)
}

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
