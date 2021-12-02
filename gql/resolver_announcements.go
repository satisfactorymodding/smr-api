package gql

import (
	"context"
	"github.com/satisfactorymodding/smr-api/db/postgres"
	"github.com/satisfactorymodding/smr-api/generated"
)

func (r *mutationResolver) CreateAnnouncement(ctx context.Context, announcement generated.NewAnnouncement) (*generated.Announcement, error) {
	return nil, nil
}

func (r *mutationResolver) DeleteAnnouncement(ctx context.Context, announcementID string) (bool, error) {
	return false, nil
}

func (r *mutationResolver) UpdateAnnouncement(ctx context.Context, announcementID string, announcement generated.UpdateAnnouncement) (*generated.Announcement, error) {
	return nil, nil
}

func (r *queryResolver) GetAnnouncement(ctx context.Context, announcementID string) (*generated.Announcement, error) {
	wrapper, newCtx := WrapQueryTrace(ctx, "getAnnouncement")
	defer wrapper.end()

	announcement := postgres.GetAnnouncementByID(announcementID, &newCtx)

	return DBAnnouncementToGenerated(announcement), nil
}

func (r *queryResolver) GetAnnouncements(ctx context.Context) (*generated.GetAnnouncements, error) {
	return nil, nil
	//wrapper, newCtx := WrapQueryTrace(ctx, "getAnnouncements")
	//defer wrapper.end()
	//
	//announcement := postgres.GetAnnouncements(&newCtx)
	//
	//return DBAnnouncementToGenerated(announcement), nil
}

func (r *queryResolver) GetAnnouncementsByImportance(ctx context.Context, importance string) (*generated.GetAnnouncements, error) {
	return nil, nil
}

//type getAnnouncementsResolver struct{ *Resolver }
