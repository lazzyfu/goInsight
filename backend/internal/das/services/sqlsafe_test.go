package services

import "testing"

func TestValidateIdentifier(t *testing.T) {
	valid := []string{"db1", "table_name", "A1_$"}
	for _, v := range valid {
		if err := validateIdentifier(v, "field"); err != nil {
			t.Fatalf("expected %q valid, got %v", v, err)
		}
	}

	invalid := []string{"db-name", "db.name", "db name", "db`name", "db' OR 1=1 --"}
	for _, v := range invalid {
		if err := validateIdentifier(v, "field"); err == nil {
			t.Fatalf("expected %q invalid", v)
		}
	}
}
