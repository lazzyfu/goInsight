package services

import (
	"testing"

	"gorm.io/datatypes"
)

func TestNormalizeAndValidateApprovalDefinition(t *testing.T) {
	t.Run("valid definition should normalize and sort", func(t *testing.T) {
		raw := datatypes.JSON([]byte(`[
			{"type":" and ","stage":2,"stage_name":" 二级审批 ","approvers":["bob","alice","bob"]},
			{"type":"or","stage":1,"stage_name":"一级审批","approvers":[" tom ","jerry"]}
		]`))

		normalized, users, err := normalizeAndValidateApprovalDefinition(raw)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		expect := `[{"type":"OR","stage":1,"stage_name":"一级审批","approvers":["tom","jerry"]},{"type":"AND","stage":2,"stage_name":"二级审批","approvers":["bob","alice"]}]`
		if string(normalized) != expect {
			t.Fatalf("unexpected normalized definition: got=%s want=%s", string(normalized), expect)
		}

		if len(users) != 4 {
			t.Fatalf("unexpected users length: %d", len(users))
		}
	})

	t.Run("empty definition should fail", func(t *testing.T) {
		if _, _, err := normalizeAndValidateApprovalDefinition(nil); err == nil {
			t.Fatal("expected error for empty definition")
		}
	})

	t.Run("non continuous stage should fail", func(t *testing.T) {
		raw := datatypes.JSON([]byte(`[
			{"type":"AND","stage":2,"stage_name":"一级审批","approvers":["alice"]}
		]`))
		if _, _, err := normalizeAndValidateApprovalDefinition(raw); err == nil {
			t.Fatal("expected error for non continuous stages")
		}
	})

	t.Run("duplicate stage should fail", func(t *testing.T) {
		raw := datatypes.JSON([]byte(`[
			{"type":"AND","stage":1,"stage_name":"一级审批","approvers":["alice"]},
			{"type":"OR","stage":1,"stage_name":"二级审批","approvers":["bob"]}
		]`))
		if _, _, err := normalizeAndValidateApprovalDefinition(raw); err == nil {
			t.Fatal("expected error for duplicate stage")
		}
	})

	t.Run("invalid type should fail", func(t *testing.T) {
		raw := datatypes.JSON([]byte(`[
			{"type":"XOR","stage":1,"stage_name":"一级审批","approvers":["alice"]}
		]`))
		if _, _, err := normalizeAndValidateApprovalDefinition(raw); err == nil {
			t.Fatal("expected error for invalid approval type")
		}
	})

	t.Run("empty approver should fail", func(t *testing.T) {
		raw := datatypes.JSON([]byte(`[
			{"type":"AND","stage":1,"stage_name":"一级审批","approvers":["   "]}
		]`))
		if _, _, err := normalizeAndValidateApprovalDefinition(raw); err == nil {
			t.Fatal("expected error for empty approver")
		}
	})
}
