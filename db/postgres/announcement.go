package postgres

import (
	"context"

	"github.com/patrickmn/go-cache"
	"github.com/satisfactorymodding/smr-api/util"
)

func CreateAnnouncement(announcement *Announcement, ctx *context.Context) (*Announcement, error) {
	announcement.ID = util.GenerateUniqueID()
	DBCtx(ctx).Create(&announcement)
	return announcement, nil
}

func GetAnnouncementByID(announcementID string, ctx *context.Context) *Announcement {
	cacheKey := "GetAnnouncementByID_" + announcementID

	if announcement, ok := dbCache.Get(cacheKey); ok {
		return announcement.(*Announcement)
	}

	var announcement Announcement
	DBCtx(ctx).Find(&announcement, "id = ?", announcementID)

	if announcement.ID == "" {
		return nil
	}

	dbCache.Set(cacheKey, announcement, cache.DefaultExpiration)

	return &announcement
}

func GetAnnouncements(ctx *context.Context) []Announcement {
	cacheKey := "GetAnnouncements"

	if announcements, ok := dbCache.Get(cacheKey); ok {
		return announcements.([]Announcement)
	}

	var announcements []Announcement
	DBCtx(ctx).Find(&announcements)

	dbCache.Set(cacheKey, announcements, cache.DefaultExpiration)

	return announcements
}

func GetAnnouncementsByImportance(importance string, ctx *context.Context) []Announcement {
	cacheKey := "GetAnnouncementsByImportance"

	if announcements, ok := dbCache.Get(cacheKey); ok {
		return announcements.([]Announcement)
	}

	var announcements []Announcement
	DBCtx(ctx).Find(&announcements, "importance = ?", importance)

	dbCache.Set(cacheKey, announcements, cache.DefaultExpiration)

	return announcements
}
