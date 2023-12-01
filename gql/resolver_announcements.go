package gql

import (
	"context"
	"fmt"
	"log/slog"

	"gopkg.in/go-playground/validator.v9"

	"github.com/satisfactorymodding/smr-api/db"
	"github.com/satisfactorymodding/smr-api/generated"
	"github.com/satisfactorymodding/smr-api/generated/conv"
	"github.com/satisfactorymodding/smr-api/generated/ent/announcement"
	"github.com/satisfactorymodding/smr-api/util"
)

func (r *mutationResolver) CreateAnnouncement(ctx context.Context, announcement generated.NewAnnouncement) (*generated.Announcement, error) {
	wrapper, ctx := WrapMutationTrace(ctx, "createAnnouncement")
	defer wrapper.end()

	val := ctx.Value(util.ContextValidator{}).(*validator.Validate)
	if err := val.Struct(&announcement); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	result, err := db.From(ctx).Announcement.
		Create().
		SetMessage(announcement.Message).
		SetImportance(string(announcement.Importance)).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	return (*conv.AnnouncementImpl)(nil).Convert(result), nil
}

func (r *mutationResolver) DeleteAnnouncement(ctx context.Context, announcementID string) (bool, error) {
	wrapper, ctx := WrapMutationTrace(ctx, "deleteAnnouncement")
	defer wrapper.end()

	slog.Info("deleting announcement", slog.String("id", announcementID))

	if err := db.From(ctx).Announcement.DeleteOneID(announcementID).Exec(ctx); err != nil {
		return false, err
	}

	return true, nil
}

func (r *mutationResolver) UpdateAnnouncement(ctx context.Context, announcementID string, announcement generated.UpdateAnnouncement) (*generated.Announcement, error) {
	wrapper, ctx := WrapMutationTrace(ctx, "updateAnnouncement")
	defer wrapper.end()

	val := ctx.Value(util.ContextValidator{}).(*validator.Validate)
	if err := val.Struct(&announcement); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	update := db.From(ctx).Announcement.UpdateOneID(announcementID)
	SetINNOEF(announcement.Message, update.SetMessage)
	SetINNOEF((*string)(announcement.Importance), update.SetImportance)

	result, err := update.Save(ctx)
	if err != nil {
		return nil, err
	}

	return (*conv.AnnouncementImpl)(nil).Convert(result), nil
}

func (r *queryResolver) GetAnnouncement(ctx context.Context, announcementID string) (*generated.Announcement, error) {
	wrapper, ctx := WrapQueryTrace(ctx, "getAnnouncement")
	defer wrapper.end()

	result, err := db.From(ctx).Announcement.Get(ctx, announcementID)
	if err != nil {
		return nil, err
	}

	return (*conv.AnnouncementImpl)(nil).Convert(result), nil
}

func (r *queryResolver) GetAnnouncements(ctx context.Context) ([]*generated.Announcement, error) {
	wrapper, ctx := WrapQueryTrace(ctx, "getAnnouncements")
	defer wrapper.end()

	result, err := db.From(ctx).Debug().Announcement.Query().All(ctx)
	if err != nil {
		return nil, err
	}

	return (*conv.AnnouncementImpl)(nil).ConvertSlice(result), nil
}

func (r *queryResolver) GetAnnouncementsByImportance(ctx context.Context, importance generated.AnnouncementImportance) ([]*generated.Announcement, error) {
	wrapper, ctx := WrapQueryTrace(ctx, "getAnnouncementsByImportance")
	defer wrapper.end()

	result, err := db.From(ctx).Announcement.
		Query().
		Where(announcement.ImportanceEQ(string(importance))).
		All(ctx)
	if err != nil {
		return nil, err
	}

	return (*conv.AnnouncementImpl)(nil).ConvertSlice(result), nil
}
