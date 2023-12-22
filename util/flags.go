package util

import (
	"context"
	"log/slog"

	"github.com/Vilsol/slox"
	"github.com/spf13/viper"
)

type FeatureFlag string

const (
	FeatureFlagAllowMultiTargetUpload = "allow_multi_target_upload"
)

func FlagEnabled(flag FeatureFlag) bool {
	return viper.GetBool("feature_flags." + string(flag))
}

func PrintFeatureFlags(ctx context.Context) {
	for _, flag := range []FeatureFlag{FeatureFlagAllowMultiTargetUpload} {
		slox.Info(ctx, "flag", slog.String("flag", string(flag)), slog.Bool("enabled", FlagEnabled(flag)))
	}
}
