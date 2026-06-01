package main

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"github.com/gerlv/trello-cli/internal/credentials"
	"github.com/gerlv/trello-cli/internal/trello"
)

func TestCustomFieldsListCommand(t *testing.T) {
	setupTestAuth(t)
	credStore.Set("default", credentials.Credentials{APIKey: "k", Token: "t", AuthMode: "manual"})
	apiClient = &mockAPI{
		listCustomFieldsByBoardFn: func(ctx context.Context, boardID string) ([]trello.CustomField, error) {
			if boardID != "b1" {
				t.Fatalf("boardID = %q, want b1", boardID)
			}
			return []trello.CustomField{{ID: "cf1", Name: "Status", Type: "list"}}, nil
		},
	}

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetArgs([]string{"custom-fields", "list", "--board", "b1"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("custom-fields list failed: %v", err)
	}

	var envelope map[string]any
	if err := json.Unmarshal(buf.Bytes(), &envelope); err != nil {
		t.Fatalf("invalid JSON: %v\nraw: %s", err, buf.String())
	}
	data := envelope["data"].([]any)
	if len(data) != 1 {
		t.Fatalf("len(data) = %d, want 1", len(data))
	}
}

func TestCustomFieldsGetCommand(t *testing.T) {
	setupTestAuth(t)
	credStore.Set("default", credentials.Credentials{APIKey: "k", Token: "t", AuthMode: "manual"})
	apiClient = &mockAPI{
		getCustomFieldFn: func(ctx context.Context, fieldID string) (trello.CustomField, error) {
			if fieldID != "cf1" {
				t.Fatalf("fieldID = %q, want cf1", fieldID)
			}
			return trello.CustomField{ID: "cf1", Name: "Status"}, nil
		},
	}

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetArgs([]string{"custom-fields", "get", "--field", "cf1"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("custom-fields get failed: %v", err)
	}

	var envelope map[string]any
	if err := json.Unmarshal(buf.Bytes(), &envelope); err != nil {
		t.Fatalf("invalid JSON: %v\nraw: %s", err, buf.String())
	}
	data := envelope["data"].(map[string]any)
	if data["id"] != "cf1" {
		t.Fatalf("data.id = %v, want cf1", data["id"])
	}
}

func TestCustomFieldsCreateCommand(t *testing.T) {
	setupTestAuth(t)
	credStore.Set("default", credentials.Credentials{APIKey: "k", Token: "t", AuthMode: "manual"})
	apiClient = &mockAPI{
		createCustomFieldFn: func(ctx context.Context, params trello.CreateCustomFieldParams) (trello.CustomField, error) {
			if params.IDModel != "b1" || params.Name != "Status" || params.Type != "list" {
				t.Fatalf("params = %+v", params)
			}
			if !params.Display.CardFront {
				t.Fatalf("display.cardFront = false, want true")
			}
			if len(params.Options) != 2 {
				t.Fatalf("len(params.Options) = %d, want 2", len(params.Options))
			}
			if params.Options[0].Value.Text != "To Do" || params.Options[1].Value.Text != "Done" {
				t.Fatalf("options = %+v", params.Options)
			}
			return trello.CustomField{ID: "cf1", Name: params.Name, Type: params.Type}, nil
		},
	}

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetArgs([]string{
		"custom-fields", "create",
		"--board", "b1",
		"--name", "Status",
		"--type", "list",
		"--card-front",
		"--option", "To Do",
		"--option", "Done",
	})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("custom-fields create failed: %v", err)
	}
}

func TestCustomFieldsUpdateCommand(t *testing.T) {
	setupTestAuth(t)
	credStore.Set("default", credentials.Credentials{APIKey: "k", Token: "t", AuthMode: "manual"})
	apiClient = &mockAPI{
		updateCustomFieldFn: func(ctx context.Context, fieldID string, params trello.UpdateCustomFieldParams) (trello.CustomField, error) {
			if fieldID != "cf1" {
				t.Fatalf("fieldID = %q, want cf1", fieldID)
			}
			if params.Name == nil || *params.Name != "Status Updated" {
				t.Fatalf("params.Name = %v", params.Name)
			}
			if params.Display == nil || params.Display.CardFront {
				t.Fatalf("params.Display = %+v, want cardFront false", params.Display)
			}
			return trello.CustomField{ID: fieldID, Name: *params.Name}, nil
		},
	}

	if err := executeRootArgs("custom-fields", "update", "--field", "cf1", "--name", "Status Updated", "--card-front=false"); err != nil {
		t.Fatalf("custom-fields update failed: %v", err)
	}
}

func TestCustomFieldsDeleteCommand(t *testing.T) {
	setupTestAuth(t)
	credStore.Set("default", credentials.Credentials{APIKey: "k", Token: "t", AuthMode: "manual"})
	apiClient = &mockAPI{
		deleteCustomFieldFn: func(ctx context.Context, fieldID string) error {
			if fieldID != "cf1" {
				t.Fatalf("fieldID = %q, want cf1", fieldID)
			}
			return nil
		},
	}

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetArgs([]string{"custom-fields", "delete", "--field", "cf1"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("custom-fields delete failed: %v", err)
	}

	var envelope map[string]any
	if err := json.Unmarshal(buf.Bytes(), &envelope); err != nil {
		t.Fatalf("invalid JSON: %v\nraw: %s", err, buf.String())
	}
	data := envelope["data"].(map[string]any)
	if data["deleted"] != true || data["id"] != "cf1" {
		t.Fatalf("data = %+v", data)
	}
}

func TestCustomFieldsValidation(t *testing.T) {
	setupTestAuth(t)
	credStore.Set("default", credentials.Credentials{APIKey: "k", Token: "t", AuthMode: "manual"})

	assertContractCode(t, executeRootArgs("custom-fields", "list"), "VALIDATION_ERROR")
	assertContractCode(t, executeRootArgs("custom-fields", "get"), "VALIDATION_ERROR")
	assertContractCode(t, executeRootArgs("custom-fields", "create", "--board", "b1", "--name", "Status"), "VALIDATION_ERROR")
	assertContractCode(t, executeRootArgs("custom-fields", "create", "--board", "b1", "--name", "Status", "--type", "bogus"), "VALIDATION_ERROR")
	assertContractCode(t, executeRootArgs("custom-fields", "create", "--board", "b1", "--name", "Status", "--type", "text", "--option", "bad"), "VALIDATION_ERROR")
	assertContractCode(t, executeRootArgs("custom-fields", "update", "--field", "cf1"), "VALIDATION_ERROR")
	assertContractCode(t, executeRootArgs("custom-fields", "delete"), "VALIDATION_ERROR")
}

func TestCustomFieldOptionsCommands(t *testing.T) {
	setupTestAuth(t)
	credStore.Set("default", credentials.Credentials{APIKey: "k", Token: "t", AuthMode: "manual"})
	apiClient = &mockAPI{
		listCustomFieldOptionsFn: func(ctx context.Context, fieldID string) ([]trello.CustomFieldOption, error) {
			if fieldID != "cf1" {
				t.Fatalf("fieldID = %q, want cf1", fieldID)
			}
			return []trello.CustomFieldOption{{ID: "opt1", Value: trello.CustomFieldOptionValue{Text: "To Do"}}}, nil
		},
		createCustomFieldOptionFn: func(ctx context.Context, fieldID string, params trello.CreateCustomFieldOptionParams) (trello.CustomFieldOption, error) {
			if fieldID != "cf1" || params.Value.Text != "Done" || params.Color != "green" {
				t.Fatalf("fieldID/params = %q/%+v", fieldID, params)
			}
			return trello.CustomFieldOption{ID: "opt2", Value: params.Value, Color: params.Color}, nil
		},
		updateCustomFieldOptionFn: func(ctx context.Context, fieldID, optionID string, params trello.UpdateCustomFieldOptionParams) (trello.CustomFieldOption, error) {
			if fieldID != "cf1" || optionID != "opt1" {
				t.Fatalf("fieldID/optionID = %q/%q", fieldID, optionID)
			}
			if params.Value == nil || params.Value.Text != "Blocked" {
				t.Fatalf("params.Value = %+v", params.Value)
			}
			return trello.CustomFieldOption{ID: optionID, Value: *params.Value}, nil
		},
		deleteCustomFieldOptionFn: func(ctx context.Context, fieldID, optionID string) error {
			if fieldID != "cf1" || optionID != "opt1" {
				t.Fatalf("fieldID/optionID = %q/%q", fieldID, optionID)
			}
			return nil
		},
	}

	if err := executeRootArgs("custom-fields", "options", "list", "--field", "cf1"); err != nil {
		t.Fatalf("options list failed: %v", err)
	}
	if err := executeRootArgs("custom-fields", "options", "add", "--field", "cf1", "--text", "Done", "--color", "green"); err != nil {
		t.Fatalf("options add failed: %v", err)
	}
	if err := executeRootArgs("custom-fields", "options", "update", "--field", "cf1", "--option", "opt1", "--text", "Blocked"); err != nil {
		t.Fatalf("options update failed: %v", err)
	}
	if err := executeRootArgs("custom-fields", "options", "delete", "--field", "cf1", "--option", "opt1"); err != nil {
		t.Fatalf("options delete failed: %v", err)
	}
}

func TestCustomFieldOptionsValidation(t *testing.T) {
	setupTestAuth(t)
	credStore.Set("default", credentials.Credentials{APIKey: "k", Token: "t", AuthMode: "manual"})

	assertContractCode(t, executeRootArgs("custom-fields", "options", "list"), "VALIDATION_ERROR")
	assertContractCode(t, executeRootArgs("custom-fields", "options", "add", "--field", "cf1"), "VALIDATION_ERROR")
	assertContractCode(t, executeRootArgs("custom-fields", "options", "update", "--field", "cf1", "--option", "opt1"), "VALIDATION_ERROR")
	assertContractCode(t, executeRootArgs("custom-fields", "options", "delete", "--field", "cf1"), "VALIDATION_ERROR")
}

func TestCustomFieldItemsCommands(t *testing.T) {
	setupTestAuth(t)
	credStore.Set("default", credentials.Credentials{APIKey: "k", Token: "t", AuthMode: "manual"})
	apiClient = &mockAPI{
		listCardCustomFieldItemsFn: func(ctx context.Context, cardID string) ([]trello.CardCustomFieldItem, error) {
			if cardID != "c1" {
				t.Fatalf("cardID = %q, want c1", cardID)
			}
			return []trello.CardCustomFieldItem{{ID: "item1", IDCustomField: "cf1", IDValue: "opt1"}}, nil
		},
		setCardCustomFieldItemFn: func(ctx context.Context, cardID, fieldID string, params trello.SetCardCustomFieldItemParams) (trello.CardCustomFieldItem, error) {
			if cardID != "c1" || fieldID != "cf1" {
				t.Fatalf("cardID/fieldID = %q/%q", cardID, fieldID)
			}
			if params.IDValue != "opt1" {
				t.Fatalf("params = %+v", params)
			}
			return trello.CardCustomFieldItem{ID: "item1", IDCustomField: fieldID, IDValue: "opt1"}, nil
		},
		clearCardCustomFieldItemFn: func(ctx context.Context, cardID, fieldID string) error {
			if cardID != "c1" || fieldID != "cf1" {
				t.Fatalf("cardID/fieldID = %q/%q", cardID, fieldID)
			}
			return nil
		},
	}

	if err := executeRootArgs("custom-fields", "items", "list", "--card", "c1"); err != nil {
		t.Fatalf("items list failed: %v", err)
	}
	if err := executeRootArgs("custom-fields", "items", "set", "--card", "c1", "--field", "cf1", "--option", "opt1"); err != nil {
		t.Fatalf("items set failed: %v", err)
	}
	if err := executeRootArgs("custom-fields", "items", "clear", "--card", "c1", "--field", "cf1"); err != nil {
		t.Fatalf("items clear failed: %v", err)
	}
}

func TestCustomFieldItemsValidation(t *testing.T) {
	setupTestAuth(t)
	credStore.Set("default", credentials.Credentials{APIKey: "k", Token: "t", AuthMode: "manual"})

	assertContractCode(t, executeRootArgs("custom-fields", "items", "list"), "VALIDATION_ERROR")
	assertContractCode(t, executeRootArgs("custom-fields", "items", "set", "--card", "c1", "--field", "cf1"), "VALIDATION_ERROR")
	assertContractCode(t, executeRootArgs("custom-fields", "items", "set", "--card", "c1", "--field", "cf1", "--text", "hello", "--option", "opt1"), "VALIDATION_ERROR")
	assertContractCode(t, executeRootArgs("custom-fields", "items", "set", "--card", "c1", "--field", "cf1", "--date", "not-a-date"), "VALIDATION_ERROR")
	assertContractCode(t, executeRootArgs("custom-fields", "items", "clear", "--card", "c1"), "VALIDATION_ERROR")
}
