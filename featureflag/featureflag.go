package featureflag

import (
	"os"
	"strings"
)

type featureFlag struct {
	environment string
}

func New(environment string) *featureFlag {
	return &featureFlag{
		environment: environment,
	}
}

type FeatureFlag interface {
	CanBeSkipped() bool
	Get(flag string) bool
	GetExplicit(flag string) bool
}

func (f *featureFlag) CanBeSkipped() bool {
	return f.environment == "development" || f.environment == "local"
}

func (f *featureFlag) Get(flag string) bool {
	flag = strings.ToLower(flag)
	if f.CanBeSkipped() {
		return true
	}

	return f.GetExplicit(flag)
}

func (f *featureFlag) GetExplicit(flag string) bool {
	flag = strings.ToLower(flag)
	val := os.Getenv(flag)

	return val == "1" || val == "true"
}
