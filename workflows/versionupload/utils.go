package versionupload

import (
	"context"
	"fmt"
	"io"

	"github.com/satisfactorymodding/smr-api/generated/ent"
	"github.com/satisfactorymodding/smr-api/storage"
)

func downloadMod(ctx context.Context, mod *ent.Mod, uploadID string) ([]byte, error) {
	modFile, err := storage.GetMod(ctx, mod.ID, mod.Name, uploadID)
	if err != nil {
		return nil, fmt.Errorf("failed getting mod: %w", err)
	}

	// TODO Optimize
	fileData, err := io.ReadAll(modFile)
	if err != nil {
		return nil, fmt.Errorf("failed reading mod file: %w", err)
	}

	return fileData, nil
}
