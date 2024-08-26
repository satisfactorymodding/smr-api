package versionupload

import (
	"context"
	"log/slog"

	"github.com/Vilsol/slox"

	"github.com/satisfactorymodding/smr-api/generated"
	"github.com/satisfactorymodding/smr-api/redis"
)

type StoreRedisStateArgs struct {
	UploadID string                           `json:"upload_id"`
	Data     *generated.CreateVersionResponse `json:"data"`
	Err      string                           `json:"err"`
}

func (*A) StoreRedisStateActivity(ctx context.Context, args StoreRedisStateArgs) error {
	ctx = slox.With(ctx, slog.String("upload_id", args.UploadID))

	if err2 := redis.StoreVersionUploadState(args.UploadID, args.Data, args.Err); err2 != nil {
		slox.Error(ctx, "error storing redis state", slog.Any("err", err2))
		return err2
	}

	return nil
}
