package workflows

import (
	"context"
	"log"
	"log/slog"
	"os"

	"github.com/Vilsol/slox"
	"github.com/spf13/viper"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/contrib/opentelemetry"
	"go.temporal.io/sdk/interceptor"
	"go.temporal.io/sdk/worker"
)

type workflowKey struct{}

const RepoTaskQueue = "REPO_TASK_QUEUE"

func InitializeWorkflows(ctx context.Context) (context.Context, func()) {
	tracingInterceptor, err := opentelemetry.NewTracingInterceptor(opentelemetry.TracerOptions{})
	if err != nil {
		log.Fatalln("unable to create tracing interceptor", err)
	}

	c, err := client.Dial(client.Options{
		HostPort:     viper.GetString("temporal.host"),
		Logger:       slox.From(ctx),
		Interceptors: []interceptor.ClientInterceptor{tracingInterceptor},
	})
	if err != nil {
		log.Fatalln("unable to create Temporal client", err)
	}

	Workflows.Statistics.InitializeStatisticsWorkflow(ctx, c, RepoTaskQueue)

	w := worker.New(c, RepoTaskQueue, worker.Options{
		BackgroundActivityContext: ctx,
	})

	w.RegisterWorkflow(Workflows.Statistics.StatisticsWorkflow)
	w.RegisterActivity(Workflows.Statistics.UpdateStatisticsActivity)

	w.RegisterWorkflow(Workflows.VersionUpload.FinalizeVersionUploadWorkflow)
	w.RegisterActivity(Workflows.VersionUpload.CompleteUploadMultipartModActivity)
	w.RegisterActivity(Workflows.VersionUpload.ExtractModInfoActivity)
	w.RegisterActivity(Workflows.VersionUpload.ExtractMetadataActivity)
	w.RegisterActivity(Workflows.VersionUpload.RenameVersionActivity)
	w.RegisterActivity(Workflows.VersionUpload.SeparateModTargetsActivity)
	w.RegisterActivity(Workflows.VersionUpload.CreateVersionInDatabaseActivity)
	w.RegisterActivity(Workflows.VersionUpload.ApproveAndPublishModActivity)
	w.RegisterActivity(Workflows.VersionUpload.StoreRedisStateActivity)
	w.RegisterActivity(Workflows.VersionUpload.ScanModOnVirusTotalActivity)
	w.RegisterActivity(Workflows.VersionUpload.RemoveModActivity)

	w.RegisterWorkflow(Workflows.UpdateModFromStorage.UpdateModDataFromStorageWorkflow)
	w.RegisterActivity(Workflows.UpdateModFromStorage.UpdateModDataFromStorageActivity)

	if err := w.Start(); err != nil {
		slox.Error(ctx, "unable to start worker", slog.Any("err", err))
		os.Exit(1)
	}

	return context.WithValue(ctx, workflowKey{}, c), func() {
		w.Stop()
		c.Close()
	}
}

func Client(ctx context.Context) client.Client {
	c := ctx.Value(workflowKey{})
	if c == nil {
		return nil
	}
	return c.(client.Client)
}

func TransferContext(source context.Context, target context.Context) context.Context {
	c := source.Value(workflowKey{})
	if c == nil {
		return target
	}
	return context.WithValue(target, workflowKey{}, c)
}
