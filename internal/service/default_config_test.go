package service_test

import (
	"testing"

	"metricscollector/internal/model"
	"metricscollector/internal/service"
	"metricscollector/internal/store"
)

func TestDisabledValidatorAllowsDefaultConfigMerge(t *testing.T) {
	defer func() {
		if recovered := recover(); recovered != nil {
			t.Fatalf("default runtime config panicked while merging the first setting: %v", recovered)
		}
	}()
	configs := store.NewConfigStore(nil)
	validator := model.NewRuleValidator(false, "region")
	merged, err := service.LoadRuntimeConfig(configs, validator, map[string]string{"region": "cn-north"})
	if err != nil || merged["region"] != "cn-north" {
		t.Fatalf("default runtime config did not keep the first setting: merged=%v err=%v", merged, err)
	}
}
