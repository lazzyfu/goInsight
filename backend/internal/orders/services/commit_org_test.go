package services

import (
	"reflect"
	"testing"

	"gorm.io/datatypes"
)

func TestCollectAccessibleOrganizationKeys(t *testing.T) {
	t.Run("collect and dedupe parent and current keys", func(t *testing.T) {
		bindings := []userOrganizationBinding{
			{
				Key:  "0-1-2",
				Path: datatypes.JSON([]byte(`["0-1"]`)),
			},
			{
				Key:  "0-5-6",
				Path: datatypes.JSON([]byte(`["0-5"]`)),
			},
			{
				Key:  "0-1-2",
				Path: datatypes.JSON([]byte(`["0-1"]`)),
			},
		}

		keys, err := collectAccessibleOrganizationKeys(bindings)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		want := []string{"0-1", "0-1-2", "0-5", "0-5-6"}
		if !reflect.DeepEqual(keys, want) {
			t.Fatalf("unexpected keys, got=%v, want=%v", keys, want)
		}
	})

	t.Run("invalid path json should fail", func(t *testing.T) {
		bindings := []userOrganizationBinding{
			{
				Key:  "0-1-2",
				Path: datatypes.JSON([]byte(`{}`)),
			},
		}

		if _, err := collectAccessibleOrganizationKeys(bindings); err == nil {
			t.Fatal("expected error for invalid organization path json")
		}
	})
}
