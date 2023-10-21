package gql

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"gopkg.in/go-playground/validator.v9"

	"github.com/satisfactorymodding/smr-api/db/postgres"
	"github.com/satisfactorymodding/smr-api/generated"
	"github.com/satisfactorymodding/smr-api/util"
)

func (r *mutationResolver) CreateAnnouncement(ctx context.Context, announcement generated.NewAnnouncement) (*generated.Announcement, error) {
	wrapper, newCtx := WrapMutationTrace(ctx, "createAnnouncement")
	defer wrapper.end()

	val := ctx.Value(util.ContextValidator{}).(*validator.Validate)
	if err := val.Struct(&announcement); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	dbAnnouncement := &postgres.Announcement{
		Message:    announcement.Message,
		Importance: string(announcement.Importance),
	}

	resultAnnouncement, err := postgres.CreateAnnouncement(newCtx, dbAnnouncement)
	if err != nil {
		return nil, err
	}
	return DBAnnouncementToGenerated(resultAnnouncement), nil
}

func (r *mutationResolver) DeleteAnnouncement(ctx context.Context, announcementID string) (bool, error) {
	wrapper, newCtx := WrapMutationTrace(ctx, "deleteAnnouncement")
	defer wrapper.end()

	dbAnnouncement := postgres.GetAnnouncementByID(newCtx, announcementID)

	if dbAnnouncement == nil {
		return false, errors.New("announcement not found")
	}

	postgres.Delete(newCtx, &dbAnnouncement)

	return true, nil
}

func (r *mutationResolver) UpdateAnnouncement(ctx context.Context, announcementID string, announcement generated.UpdateAnnouncement) (*generated.Announcement, error) {
	wrapper, newCtx := WrapMutationTrace(ctx, "updateAnnouncement")
	defer wrapper.end()

	val := ctx.Value(util.ContextValidator{}).(*validator.Validate)
	if err := val.Struct(&announcement); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	dbAnnouncement := postgres.GetAnnouncementByID(newCtx, announcementID)

	if dbAnnouncement == nil {
		return nil, errors.New("guide not found")
	}

	SetStringINNOE(announcement.Message, &dbAnnouncement.Message)
	SetStringINNOE((*string)(announcement.Importance), &dbAnnouncement.Importance)

	postgres.Save(newCtx, &dbAnnouncement)

	return DBAnnouncementToGenerated(dbAnnouncement), nil
}

func (r *queryResolver) GetAnnouncement(ctx context.Context, announcementID string) (*generated.Announcement, error) {
	wrapper, newCtx := WrapQueryTrace(ctx, "getAnnouncement")
	defer wrapper.end()

	announcement := postgres.GetAnnouncementByID(newCtx, announcementID)

	return DBAnnouncementToGenerated(announcement), nil
}

func (r *queryResolver) GetAnnouncements(ctx context.Context) ([]*generated.Announcement, error) {
	wrapper, newCtx := WrapQueryTrace(ctx, "getAnnouncements")
	defer wrapper.end()

	announcements := postgres.GetAnnouncements(newCtx)

	return DBAnnouncementsToGeneratedSlice(announcements), nil
}

func (r *queryResolver) GetAnnouncementsByImportance(ctx context.Context, importance generated.AnnouncementImportance) ([]*generated.Announcement, error) {
	wrapper, newCtx := WrapQueryTrace(ctx, "getAnnouncementsByImportance")
	defer wrapper.end()

	announcements := postgres.GetAnnouncementsByImportance(newCtx, string(importance))

	return DBAnnouncementsToGeneratedSlice(announcements), nil
}
