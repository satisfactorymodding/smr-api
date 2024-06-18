package profiling

import (
	"runtime"

	otelpyroscope "github.com/grafana/otel-profiling-go"
	"github.com/grafana/pyroscope-go"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel"
)

func SetupProfiling() {
	if !viper.GetBool("production") && viper.IsSet("pyroscope_endpoint") {
		otel.SetTracerProvider(otelpyroscope.NewTracerProvider(otel.GetTracerProvider()))

		runtime.SetMutexProfileFraction(5)
		runtime.SetBlockProfileRate(5)

		_, err := pyroscope.Start(pyroscope.Config{
			ApplicationName: "ficsit-app-api",
			ServerAddress:   viper.GetString("pyroscope_endpoint"),
			ProfileTypes: []pyroscope.ProfileType{
				pyroscope.ProfileCPU,
				pyroscope.ProfileAllocObjects,
				pyroscope.ProfileAllocSpace,
				pyroscope.ProfileInuseObjects,
				pyroscope.ProfileInuseSpace,
				pyroscope.ProfileGoroutines,
				pyroscope.ProfileMutexCount,
				pyroscope.ProfileMutexDuration,
				pyroscope.ProfileBlockCount,
				pyroscope.ProfileBlockDuration,
			},
		})
		if err != nil {
			panic(err)
		}
	}
}
