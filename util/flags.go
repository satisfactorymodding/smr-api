package util

import "github.com/spf13/viper"

type FeatureFlag string

const (
	FeatureFlagAllowMultiTargetUpload = "allow_multi_target_upload"
)

func FlagEnabled(flag FeatureFlag) bool {
	return viper.GetBool("feature_flags." + string(flag))
}
