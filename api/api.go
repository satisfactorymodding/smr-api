package api

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/Vilsol/slox"
	"github.com/felixge/fgprof"
	"github.com/labstack/echo-contrib/pprof"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"
	echoSwagger "github.com/swaggo/echo-swagger"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"go.opentelemetry.io/otel/trace"
	"gopkg.in/go-playground/validator.v9"

	"github.com/satisfactorymodding/smr-api/auth"
	"github.com/satisfactorymodding/smr-api/config"
	"github.com/satisfactorymodding/smr-api/dataloader"
	"github.com/satisfactorymodding/smr-api/db"
	"github.com/satisfactorymodding/smr-api/db/postgres"
	"github.com/satisfactorymodding/smr-api/generated"
	"github.com/satisfactorymodding/smr-api/gql"
	"github.com/satisfactorymodding/smr-api/migrations"
	"github.com/satisfactorymodding/smr-api/nodes"
	"github.com/satisfactorymodding/smr-api/oauth"
	"github.com/satisfactorymodding/smr-api/redis"
	"github.com/satisfactorymodding/smr-api/redis/jobs"
	"github.com/satisfactorymodding/smr-api/storage"
	"github.com/satisfactorymodding/smr-api/util"
	"github.com/satisfactorymodding/smr-api/validation"

	// Load REST docs
	_ "github.com/satisfactorymodding/smr-api/docs"
	// Load redis consumers
	_ "github.com/satisfactorymodding/smr-api/redis/jobs/consumers"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return fmt.Errorf("validation error: %w", err)
	}
	return nil
}

func Initialize(baseCtx context.Context) context.Context {
	ctx := config.InitializeConfig(baseCtx)

	if os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT") != "" {
		cleanup := installExportPipeline(ctx)
		defer cleanup()
	}

	redis.InitializeRedis(ctx)
	postgres.InitializePostgres(ctx)
	ctx, err := db.WithDB(ctx)
	if err != nil {
		panic(err)
	}

	storage.InitializeStorage(ctx)
	oauth.InitializeOAuth()
	util.InitializeSecurity()
	validation.InitializeValidator()
	auth.InitializeAuth()
	jobs.InitializeJobs(ctx)
	validation.InitializeVirusTotal()
	util.PrintFeatureFlags(ctx)

	return ctx
}

func Migrate(ctx context.Context) {
	migrations.RunMigrations(ctx)
}

var e *echo.Echo

func Setup(ctx context.Context) {
	if viper.GetBool("profiler") {
		go func() {
			debugServer := echo.New()
			pprof.Register(debugServer)
			debugServer.Any("/debug/fgprof", echo.WrapHandler(fgprof.Handler()))
			debugServer.HideBanner = true
			debugServer.HidePort = true
			address := ":6060"
			slog.Info("starting profiler", slog.String("address", address))
			debugServer.Logger.Fatal(debugServer.Start(address))
		}()
	}

	db.RunAsyncStatisticLoop(ctx)

	dataValidator := validator.New()

	e = echo.New()
	e.HideBanner = true
	e.Validator = &CustomValidator{validator: dataValidator}

	e.Pre(middleware.RemoveTrailingSlash())

	e.Static("/static", "static")

	v1 := e.Group("/v1")

	v1.Use(func(handlerFunc echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			newCtx := slox.With(ctx.Request().Context(), slog.String("facade", "REST"))
			ctx.SetRequest(ctx.Request().WithContext(newCtx))
			return handlerFunc(ctx)
		}
	})

	nodes.RegisterOAuthRoutes(v1.Group("/oauth"))
	nodes.RegisterUserRoutes(v1.Group("/user"))
	nodes.RegisterUsersRoutes(v1.Group("/users"))
	nodes.RegisterModRoutes(v1.Group("/mod"))
	nodes.RegisterModsRoutes(v1.Group("/mods"))
	nodes.RegisterVersionRoutes(v1.Group("/version"))
	nodes.RegisterSMLRoutes(v1.Group("/sml"))

	v2 := e.Group("/v2")

	v2.Use(func(handlerFunc echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			newCtx := slox.With(ctx.Request().Context(), slog.String("facade", "GQL"))
			newCtx = context.WithValue(newCtx, util.ContextHeader{}, ctx.Request().Header)
			newCtx = context.WithValue(newCtx, util.ContextRequest{}, ctx.Request())
			newCtx = context.WithValue(newCtx, util.ContextResponse{}, ctx.Response().Writer)
			newCtx = context.WithValue(newCtx, util.ContextValidator{}, dataValidator)
			ctx.SetRequest(ctx.Request().WithContext(newCtx))
			return handlerFunc(ctx)
		}
	})

	v2.Any("", echo.WrapHandler(playground.Handler("GraphQL Playground", "/v2/query")))

	schema := generated.NewExecutableSchema(generated.Config{
		Resolvers:  &gql.Resolver{},
		Directives: gql.MakeDirective(),
	})

	v2Query := v2.Group("/query")

	v2Query.Use(func(handlerFunc echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			if ctx.Request().Method == "GET" &&
				ctx.Request().Header.Get("Authorization") == "" {
				ctx.Response().Header().Add("Cache-Control", "public, max-age=60, s-maxage=60")
			}

			return handlerFunc(ctx)
		}
	})

	v2Query.Use(dataloader.Middleware())

	gqlHandler := handler.New(schema)

	gqlHandler.AddTransport(transport.Options{})
	gqlHandler.AddTransport(transport.GET{})
	gqlHandler.AddTransport(transport.POST{})
	gqlHandler.AddTransport(transport.MultipartForm{
		MaxUploadSize: 100 << 20,
		MaxMemory:     100 << 20,
	})

	gqlHandler.SetQueryCache(lru.New(5000))

	gqlHandler.Use(extension.Introspection{})
	gqlHandler.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New(5000),
	})

	v2Query.Any("", echo.WrapHandler(gqlHandler))

	e.Any("/analytics*", func(ctx echo.Context) error {
		util.HandleRequestAndRedirect(ctx.Response(), ctx.Request())
		return nil
	})

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			defer func() {
				if r := recover(); r != nil {
					err, ok := r.(error)
					if !ok {
						err = fmt.Errorf("%v", r)
					}
					stack := make([]byte, 4<<10)
					length := runtime.Stack(stack, true)
					c.Logger().Printf("[PANIC RECOVER] %v %s\n", err, stack[:length])
					c.Error(err)
				}
			}()
			return next(c)
		}
	})

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		Skipper:          middleware.DefaultSkipper,
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
		AllowCredentials: true,
	}))

	if os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT") != "" {
		e.Use(otelecho.Middleware("ficsit-api"))
	}

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			res := c.Response()
			start := time.Now()

			if err := next(c); err != nil {
				c.Error(err)
			}
			stop := time.Now()

			p := req.URL.Path
			if p == "" {
				p = "/"
			}

			spanContext := trace.SpanContextFromContext(req.Context()) //nolint:contextcheck

			bytesIn := req.Header.Get(echo.HeaderContentLength)
			if bytesIn == "" {
				bytesIn = "0"
			}

			slog.Info(
				"Handled request",
				slog.String("time_rfc3339", time.Now().Format(time.RFC3339)),
				slog.String("remote_ip", c.RealIP()),
				slog.String("host", req.Host),
				slog.String("uri", req.RequestURI),
				slog.String("method", req.Method),
				slog.String("path", p),
				slog.String("referer", req.Referer()),
				slog.String("user_agent", req.UserAgent()),
				slog.Int("status", res.Status),
				slog.Int64("latency", stop.Sub(start).Nanoseconds()/1000),
				slog.String("latency_human", stop.Sub(start).String()),
				slog.String("bytes_in", bytesIn),
				slog.Int64("bytes_out", res.Size),
				slog.String("trace_id", spanContext.TraceID().String()),
			)

			return nil
		}
	})

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signals
		_ = e.Close()
	}()
}

func Serve() {
	address := fmt.Sprintf(":%d", viper.GetInt("port"))
	slog.Info("starting server", slog.String("address", address))

	e.HidePort = true
	e.Logger.Error(e.Start(address))
}

func installExportPipeline(ctx context.Context) func() {
	client := otlptracehttp.NewClient()
	exporter, err := otlptrace.New(ctx, client)
	if err != nil {
		slog.Error("creating OTLP trace exporter", slog.Any("err", err))
		os.Exit(1)
	}

	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(newResource()),
	)
	otel.SetTracerProvider(tracerProvider)

	return func() {
		if err := tracerProvider.Shutdown(ctx); err != nil {
			slog.Error("stopping tracer provider", slog.Any("err", err))
			os.Exit(1)
		}
	}
}

func newResource() *resource.Resource {
	r, _ := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("ficsit-app-api"),
			semconv.ServiceVersionKey.String("0.0.1"),
		),
	)
	return r
}

func Start() {
	ctx := Initialize(context.Background())
	Migrate(ctx)
	Setup(ctx)
	Serve()
}

func Stop() error {
	if err := e.Close(); err != nil {
		return fmt.Errorf("failed to stop http server: %w", err)
	}
	return nil
}
