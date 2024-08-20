package workflows

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"path"
	"time"

	"github.com/Vilsol/slox"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"

	"github.com/satisfactorymodding/smr-api/db"
	"github.com/satisfactorymodding/smr-api/integrations"
	"github.com/satisfactorymodding/smr-api/storage"
	"github.com/satisfactorymodding/smr-api/validation"
)

func scanModOnVirusTotalActivity(ctx context.Context, modID string, versionID string, approveAfter bool) error {
	ctx, span := otel.Tracer("ficsit-app").Start(ctx, "ScanModOnVirusTotal")
	defer span.End()

	slox.Info(ctx, "starting virus scan of mod", slog.String("mod", modID), slog.String("version", versionID))

	version, err := db.From(ctx).Version.Get(ctx, versionID)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		return err
	}
	link := storage.GenerateDownloadLink(version.Key)

	response, _ := http.Get(link)

	fileData, err := io.ReadAll(response.Body)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		return fmt.Errorf("failed to read mod file: %w", err)
	}

	archive, err := zip.NewReader(bytes.NewReader(fileData), int64(len(fileData)))
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		return fmt.Errorf("failed to unzip mod file: %w", err)
	}

	toScan := make([]io.Reader, 0)
	names := make([]string, 0)
	for _, file := range archive.File {
		if path.Ext(file.Name) == ".dll" || path.Ext(file.Name) == ".so" {
			open, err := file.Open()
			if err != nil {
				span.SetStatus(codes.Error, err.Error())
				span.RecordError(err)
				return fmt.Errorf("failed to open mod file: %w", err)
			}

			toScan = append(toScan, open)
			names = append(names, path.Base(file.Name))
		}
	}

	success, err := validation.ScanFiles(ctx, toScan, names)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		return err
	}

	if !success {
		slox.Warn(ctx, "mod failed to pass virus scan", slog.String("mod", modID), slog.String("version", versionID))
		return nil
	}

	if approveAfter {
		slox.Info(ctx, "approving mod after successful virus scan", slog.String("mod", modID), slog.String("version", versionID))

		if err := version.Update().SetApproved(true).Exec(ctx); err != nil {
			span.SetStatus(codes.Error, err.Error())
			span.RecordError(err)
			return err
		}

		if err := db.From(ctx).Mod.UpdateOneID(modID).SetLastVersionDate(time.Now()).Exec(ctx); err != nil {
			span.SetStatus(codes.Error, err.Error())
			span.RecordError(err)
			return err
		}

		go integrations.NewVersion(db.ReWrapCtx(ctx), version)
	}

	return nil
}
