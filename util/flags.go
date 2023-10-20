package util

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type FeatureFlag string

const (
	FeatureFlagAllowMultiTargetUpload = "allow_multi_target_upload"
)

func FlagEnabled(flag FeatureFlag) bool {
	return viper.GetBool("feature_flags." + string(flag))
}

func PrintFeatureFlags() {
	for _, flag := range []FeatureFlag{FeatureFlagAllowMultiTargetUpload} {
		log.Info().Str("flag", string(flag)).Bool("enabled", FlagEnabled(flag)).Msg("flag")
	}
}
