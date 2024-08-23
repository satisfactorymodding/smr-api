package versionupload

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/Vilsol/slox"

	"github.com/satisfactorymodding/smr-api/db"
	"github.com/satisfactorymodding/smr-api/util"
	"github.com/satisfactorymodding/smr-api/validation"
)

type ExtractMetadataArgs struct {
	ModID    string             `json:"mod_id"`
	UploadID string             `json:"upload_id"`
	ModInfo  validation.ModInfo `json:"mod_info"`
}

func (*A) ExtractMetadataActivity(ctx context.Context, args ExtractMetadataArgs) (*string, error) {
	metadata, err := extractMetadata(ctx, args.ModID, args.UploadID, args.ModInfo)
	if err != nil {
		slox.Error(ctx, "failed to extract metadata", slog.Any("err", err), slog.String("mod_id", args.ModID), slog.String("upload_id", args.UploadID))
		return nil, nil
	}
	return metadata, nil
}

func extractMetadata(ctx context.Context, modID string, uploadID string, modInfo validation.ModInfo) (*string, error) {
	mod, err := db.From(ctx).Mod.Get(ctx, modID)
	if err != nil {
		return nil, err
	}

	ctx = slox.With(ctx, slog.String("mod_id", mod.ID), slog.String("upload_id", uploadID))

	fileData, err := downloadMod(ctx, mod, uploadID)
	if err != nil {
		return nil, err
	}

	metadata, err := validation.ExtractMetadata(ctx, fileData, modInfo.GameVersion, modInfo.ModReference)
	if err != nil {
		return nil, err
	}

	jsonData, err := json.Marshal(metadata)
	if err != nil {
		return nil, fmt.Errorf("failed json serialization: %w", err)
	}

	return util.Ptr(string(jsonData)), nil
}
