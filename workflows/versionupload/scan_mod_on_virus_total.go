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

	response, _ := http.Get(link)

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

			toScan = append(toScan, open)
			names = append(names, path.Base(file.Name))
		}
	}

	success, err := validation.ScanFiles(ctx, toScan, names)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		return false, err
	}

	if !success {
		slox.Warn(ctx, "mod failed to pass virus scan", slog.String("mod", args.ModID), slog.String("version", args.VersionID))
		return false, nil
	}

	return true, nil
}
