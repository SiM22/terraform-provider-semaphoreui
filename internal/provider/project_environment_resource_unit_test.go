package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-semaphoreui/semaphoreui/models"
)

func testProjectEnvironmentSecretsList(t *testing.T, secrets []ProjectEnvironmentSecretModel) types.List {
	t.Helper()

	list, diags := types.ListValueFrom(context.Background(), types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"id":    types.Int64Type,
			"type":  types.StringType,
			"name":  types.StringType,
			"value": types.StringType,
		},
	}, secrets)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics creating secrets list: %v", diags)
	}

	return list
}

func TestConvertProjectEnvironmentModelToEnvironmentRequestIncludesUpdatedSecretValue(t *testing.T) {
	ctx := context.Background()

	prev := &ProjectEnvironmentModel{
		ProjectID: types.Int64Value(1),
		Name:      types.StringValue("dockhand"),
		Secrets: testProjectEnvironmentSecretsList(t, []ProjectEnvironmentSecretModel{
			{
				ID:    types.Int64Value(42),
				Name:  types.StringValue("OP_SERVICE_ACCOUNT_TOKEN"),
				Type:  types.StringValue("env"),
				Value: types.StringValue("old-value"),
			},
		}),
	}
	env := ProjectEnvironmentModel{
		ID:        types.Int64Value(2),
		ProjectID: types.Int64Value(1),
		Name:      types.StringValue("dockhand"),
		Secrets: testProjectEnvironmentSecretsList(t, []ProjectEnvironmentSecretModel{
			{
				ID:    types.Int64Value(42),
				Name:  types.StringValue("OP_SERVICE_ACCOUNT_TOKEN"),
				Type:  types.StringValue("env"),
				Value: types.StringValue("new-value"),
			},
		}),
	}

	req := convertProjectEnvironmentModelToEnvironmentRequest(ctx, env, prev)
	if len(req.Secrets) != 1 {
		t.Fatalf("expected 1 secret request, got %d", len(req.Secrets))
	}

	secretReq := req.Secrets[0]
	if secretReq.Operation != "update" {
		t.Fatalf("expected update operation, got %q", secretReq.Operation)
	}
	if secretReq.Secret != "new-value" {
		t.Fatalf("expected updated secret value to be sent, got %q", secretReq.Secret)
	}
}

func TestConvertEnvironmentResponseToProjectEnvironmentModelDoesNotMaskMissingSecrets(t *testing.T) {
	ctx := context.Background()

	prev := &ProjectEnvironmentModel{
		ProjectID: types.Int64Value(1),
		Name:      types.StringValue("dockhand"),
		Secrets: testProjectEnvironmentSecretsList(t, []ProjectEnvironmentSecretModel{
			{
				ID:    types.Int64Value(42),
				Name:  types.StringValue("OP_SERVICE_ACCOUNT_TOKEN"),
				Type:  types.StringValue("env"),
				Value: types.StringValue("old-value"),
			},
		}),
	}
	environment := &models.Environment{
		ID:        2,
		ProjectID: 1,
		Name:      "dockhand",
		JSON:      "{}",
		Env:       "{}",
		Secrets:   []*models.EnvironmentSecret{},
	}

	model := convertEnvironmentResponseToProjectEnvironmentModel(ctx, environment, prev)
	if model.Secrets.IsNull() {
		t.Fatal("expected empty secrets list, got null")
	}

	var secrets []ProjectEnvironmentSecretModel
	diags := model.Secrets.ElementsAs(ctx, &secrets, false)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics reading secrets list: %v", diags)
	}
	if len(secrets) != 0 {
		t.Fatalf("expected no secrets when API returns none, got %d", len(secrets))
	}
}
