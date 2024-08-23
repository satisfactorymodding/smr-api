package versionupload

import (
	"context"
	"log/slog"

	"github.com/Vilsol/slox"
	"go.temporal.io/sdk/temporal"

	"github.com/satisfactorymodding/smr-api/db"
	"github.com/satisfactorymodding/smr-api/storage"
	"github.com/satisfactorymodding/smr-api/validation"
)

type SeparateModTargetsArgs struct {
	ModID   string             `json:"mod_id"`
	ModInfo validation.ModInfo `json:"mod_info"`
	FileKey string             `json:"file_key"`
}

type ModTargetData struct {
	TargetName string `json:"target_name"`
	Key        string `json:"key"`
	Hash       string `json:"hash"`
	Size       int64  `json:"size"`
}

func (*A) SeparateModTargetsActivity(ctx context.Context, args SeparateModTargetsArgs) ([]ModTargetData, error) {
	mod, err := db.From(ctx).Mod.Get(ctx, args.ModID)
	if err != nil {
		return nil, err
	}

	ctx = slox.With(ctx, slog.String("mod_id", mod.ID), slog.String("version", args.ModInfo.Version))

	fileData, err := downloadMod(ctx, mod, args.ModInfo.Version)
	if err != nil {
		return nil, err
	}

	if args.ModInfo.Type == validation.MultiTargetUEPlugin {
		targets := make([]ModTargetData, 0)

		for _, target := range args.ModInfo.Targets {
			slox.Info(ctx, "separating mod", slog.String("target", target), slog.String("mod", mod.Name), slog.String("version", args.ModInfo.Version))
			key, hash, size, err := storage.SeparateModTarget(ctx, fileData, mod.ID, mod.Name, args.ModInfo.Version, target)
			if err != nil {
				return nil, temporal.NewNonRetryableApplicationError("failed to separate mod", "fatal", err)
			}
			targets = append(targets, ModTargetData{
				TargetName: target,
				Key:        key,
				Hash:       hash,
				Size:       size,
			})
		}

		return targets, nil
	}

	// A single Windows target for legacy mod formats
	return []ModTargetData{{
		TargetName: "Windows",
		Key:        args.FileKey,
		Hash:       args.ModInfo.Hash,
		Size:       args.ModInfo.Size,
	}}, nil
}
