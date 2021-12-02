package gql

import (
	"context"

	"github.com/pkg/errors"
	"github.com/satisfactorymodding/smr-api/util"
	"gopkg.in/go-playground/validator.v9"

	"github.com/satisfactorymodding/smr-api/db/postgres"
	"github.com/satisfactorymodding/smr-api/generated"
)

func (r *mutationResolver) CreateAnnouncement(ctx context.Context, announcement generated.NewAnnouncement) (*generated.Announcement, error) {
	wrapper, newCtx := WrapMutationTrace(ctx, "createAnnouncement")
	defer wrapper.end()

	val := ctx.Value(util.ContextValidator{}).(*validator.Validate)
	if err := val.Struct(&announcement); err != nil {
		return nil, errors.Wrap(err, "validation failed")
	}

	dbAnnouncement := &postgres.Announcement{
		Message:    announcement.Message,
		Importance: announcement.Importance,
	}

	resultAnnouncement, err := postgres.CreateAnnouncement(dbAnnouncement, &newCtx)
	if err != nil {
		return nil, err
	}
	return DBAnnouncementToGenerated(resultAnnouncement), nil
}

func (r *mutationResolver) DeleteAnnouncement(ctx context.Context, announcementID string) (bool, error) {
	wrapper, newCtx := WrapMutationTrace(ctx, "deleteAnnouncement")
	defer wrapper.end()

	dbAnnouncement := postgres.GetAnnouncementByID(announcementID, &newCtx)

	if dbAnnouncement == nil {
		return false, errors.New("announcement not found")
	}

	postgres.Delete(&dbAnnouncement, &newCtx)

	return true, nil
}

func (r *mutationResolver) UpdateAnnouncement(ctx context.Context, announcementID string, announcement generated.UpdateAnnouncement) (*generated.Announcement, error) {
	wrapper, newCtx := WrapMutationTrace(ctx, "updateAnnouncement")
	defer wrapper.end()

	val := ctx.Value(util.ContextValidator{}).(*validator.Validate)
	if err := val.Struct(&announcement); err != nil {
		return nil, errors.Wrap(err, "validation failed")
	}

	dbAnnouncement := postgres.GetAnnouncementByID(announcementID, &newCtx)

	if dbAnnouncement == nil {
		return nil, errors.New("guide not found")
	}

	SetStringINNOE(announcement.Message, &dbAnnouncement.Message)
	SetStringINNOE(announcement.Importance, &dbAnnouncement.Importance)

	postgres.Save(&dbAnnouncement, &newCtx)

	return DBAnnouncementToGenerated(dbAnnouncement), nil
}

func (r *queryResolver) GetAnnouncement(ctx context.Context, announcementID string) (*generated.Announcement, error) {
	wrapper, newCtx := WrapQueryTrace(ctx, "getAnnouncement")
	defer wrapper.end()

	announcement := postgres.GetAnnouncementByID(announcementID, &newCtx)

	return DBAnnouncementToGenerated(announcement), nil
}

func (r *queryResolver) GetAnnouncements(ctx context.Context) ([]*generated.Announcement, error) {
	wrapper, newCtx := WrapQueryTrace(ctx, "getAnnouncements")
	defer wrapper.end()

	announcements := postgres.GetAnnouncements(&newCtx)

	return DBAnnouncementsToGeneratedSlice(announcements), nil
}

func (r *queryResolver) GetAnnouncementsByImportance(ctx context.Context, importance string) ([]*generated.Announcement, error) {
	wrapper, newCtx := WrapQueryTrace(ctx, "getAnnouncementsByImportance")
	defer wrapper.end()

	announcements := postgres.GetAnnouncementsByImportance(importance, &newCtx)

	return DBAnnouncementsToGeneratedSlice(announcements), nil
}
