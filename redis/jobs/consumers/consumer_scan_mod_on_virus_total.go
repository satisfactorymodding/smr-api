package consumers

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"path"
	"time"

	"github.com/Vilsol/slox"
	"github.com/vmihailenco/taskq/v3"

	"github.com/satisfactorymodding/smr-api/db"
	"github.com/satisfactorymodding/smr-api/integrations"
	"github.com/satisfactorymodding/smr-api/redis/jobs/tasks"
	"github.com/satisfactorymodding/smr-api/storage"
	"github.com/satisfactorymodding/smr-api/validation"
)

func init() {
	tasks.ScanModOnVirusTotalTask = taskq.RegisterTask(&taskq.TaskOptions{
		Name:    "consumer_scan_mod_on_virus_total",
		Handler: ScanModOnVirusTotalConsumer,
	})
}

func ScanModOnVirusTotalConsumer(ctx context.Context, payload []byte) error {
	var task tasks.ScanModOnVirusTotalData
	if err := json.Unmarshal(payload, &task); err != nil {
		return fmt.Errorf("failed to unmarshal task data: %w", err)
	}

	slox.Info(ctx, "starting virus scan of mod", slog.String("mod", task.ModID), slog.String("version", task.VersionID))

	version, err := db.From(ctx).Version.Get(ctx, task.VersionID)
	if err != nil {
		return err
	}
	link := storage.GenerateDownloadLink(version.Key)

	response, _ := http.Get(link)

	fileData, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("failed to read mod file: %w", err)
	}

	archive, err := zip.NewReader(bytes.NewReader(fileData), int64(len(fileData)))
	if err != nil {
		return fmt.Errorf("failed to unzip mod file: %w", err)
	}

	toScan := make([]io.Reader, 0)
	names := make([]string, 0)
	for _, file := range archive.File {
		if path.Ext(file.Name) == ".dll" || path.Ext(file.Name) == ".so" {
			open, err := file.Open()
			if err != nil {
				return fmt.Errorf("failed to open mod file: %w", err)
			}

			toScan = append(toScan, open)
			names = append(names, path.Base(file.Name))
		}
	}

	success, err := validation.ScanFiles(ctx, toScan, names)
	if err != nil {
		return err
	}

	if !success {
		slox.Warn(ctx, "mod failed to pass virus scan", slog.String("mod", task.ModID), slog.String("version", task.VersionID))
		return nil
	}

	if task.ApproveAfter {
		slox.Info(ctx, "approving mod after successful virus scan", slog.String("mod", task.ModID), slog.String("version", task.VersionID))

		if err := version.Update().SetApproved(true).Exec(ctx); err != nil {
			return err
		}

		if err := db.From(ctx).Mod.UpdateOneID(task.ModID).SetLastVersionDate(time.Now()).Exec(ctx); err != nil {
			return err
		}

		go integrations.NewVersion(db.ReWrapCtx(ctx), version)
	}

	return nil
}
