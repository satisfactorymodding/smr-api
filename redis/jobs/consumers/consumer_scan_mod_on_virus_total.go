package consumers

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"path"
	"time"

	"github.com/satisfactorymodding/smr-api/db/postgres"
	"github.com/satisfactorymodding/smr-api/integrations"
	"github.com/satisfactorymodding/smr-api/redis/jobs/tasks"
	"github.com/satisfactorymodding/smr-api/storage"
	"github.com/satisfactorymodding/smr-api/util"
	"github.com/satisfactorymodding/smr-api/validation"

	"github.com/pkg/errors"

	"github.com/rs/zerolog/log"
	"github.com/vmihailenco/taskq/v3"
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
		return errors.Wrap(err, "failed to unmarshal task data")
	}

	log.Ctx(ctx).Info().Msgf("starting virus scan of mod %s version %s", task.ModID, task.VersionID)

	version := postgres.GetVersion(ctx, task.VersionID)
	link := storage.GenerateDownloadLink(version.Key)

	response, _ := http.Get(link)

	fileData, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return errors.Wrap(err, "failed to read mod file")
	}

	archive, err := zip.NewReader(bytes.NewReader(fileData), int64(len(fileData)))

	if err != nil {
		return errors.Wrap(err, "failed to unzip mod file")
	}

	toScan := make([]io.Reader, 0)
	names := make([]string, 0)
	for _, file := range archive.File {
		if path.Ext(file.Name) == ".dll" || path.Ext(file.Name) == ".so" {
			open, err := file.Open()

			if err != nil {
				return errors.Wrap(err, "failed to open mod file")
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
		log.Ctx(ctx).Warn().Msgf("mod %s version %s failed to pass virus scan", task.ModID, task.VersionID)
		return nil
	}

	if task.ApproveAfter {
		log.Ctx(ctx).Info().Msgf("approving mod %s version %s after successful virus scan", task.ModID, task.VersionID)
		version.Approved = true
		postgres.Save(ctx, &version)

		mod := postgres.GetModByID(ctx, task.ModID)
		now := time.Now()
		mod.LastVersionDate = &now
		postgres.Save(ctx, &mod)

		go integrations.NewVersion(util.ReWrapCtx(ctx), version)
	}

	go storage.DeleteCombinedMod(ctx, task.ModID, mod.Name, task.VersionID)

	return nil
}
