package versionupload

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"path"

	"github.com/Vilsol/slox"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"

	"github.com/satisfactorymodding/smr-api/db"
	"github.com/satisfactorymodding/smr-api/generated/ent"
	"github.com/satisfactorymodding/smr-api/storage"
	"github.com/satisfactorymodding/smr-api/validation"
)

type ScanModOnVirusTotalArgs struct {
	ModID     string `json:"mod_id"`
	VersionID string `json:"version_id"`
}

func (*A) ScanModOnVirusTotalActivity(ctx context.Context, args ScanModOnVirusTotalArgs) (bool, error) {
	ctx, span := otel.Tracer("ficsit-app").Start(ctx, "ScanModOnVirusTotal")
	defer span.End()

	slox.Info(ctx, "starting virus scan of mod", slog.String("mod", args.ModID), slog.String("version", args.VersionID))

	version, err := db.From(ctx).Version.Get(ctx, args.VersionID)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		return false, err
	}
	link := storage.GenerateDownloadLink(ctx, version.Key)

	response, err := http.Get(link)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		return false, fmt.Errorf("failed to download mod file: %w", err)
	}
	defer response.Body.Close()

	fileData, err := io.ReadAll(response.Body)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		return false, fmt.Errorf("failed to read mod file: %w", err)
	}

	archive, err := zip.NewReader(bytes.NewReader(fileData), int64(len(fileData)))
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		return false, fmt.Errorf("failed to unzip mod file: %w", err)
	}

	toScan := make([]io.Reader, 0)
	names := make([]string, 0)
	for _, file := range archive.File {
		if path.Ext(file.Name) == ".dll" || path.Ext(file.Name) == ".so" {
			open, err := file.Open()
			if err != nil {
				span.SetStatus(codes.Error, err.Error())
				span.RecordError(err)
				return false, fmt.Errorf("failed to open mod file: %w", err)
			}
			defer open.Close()

			toScan = append(toScan, open)
			names = append(names, path.Base(file.Name))
		}
	}

	success, err := scanAndSaveResults(ctx, toScan, names, args)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		return false, err
	}

	return success, nil
}

func scanAndSaveResults(ctx context.Context, toScan []io.Reader, names []string, args ScanModOnVirusTotalArgs) (bool, error) {
	ctx, span := otel.Tracer("ficsit-app").Start(ctx, "scanAndSaveResults")
	defer span.End()

	scanResults, scanErr := validation.ScanFiles(ctx, toScan, names)
	// Check error later, because we can have partial results to save, even in the case of an error

	if err := saveScanResults(ctx, scanResults, args); err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		return false, err
	}

	if scanErr != nil {
		span.SetStatus(codes.Error, scanErr.Error())
		span.RecordError(scanErr)
		return false, scanErr
	}

	if len(scanResults) != len(toScan) {
		return false, nil
	}

	for _, result := range scanResults {
		if !result.Safe {
			return false, nil
		}
	}

	return true, nil
}

func saveScanResults(ctx context.Context, scanResults []validation.ScanResult, args ScanModOnVirusTotalArgs) error {
	err := db.From(ctx).VirustotalResult.MapCreateBulk(scanResults,
		func(c *ent.VirustotalResultCreate, i int) {
			c.SetSafe(scanResults[i].Safe).
				SetVersionID(args.VersionID).
				SetFileName(scanResults[i].FileName).
				SetHash(scanResults[i].Hash)
		},
	).OnConflict().
		DoNothing().
		Exec(ctx)
	return err
}
