package main

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"github.com/gerlv/trello-cli/internal/credentials"
	"github.com/gerlv/trello-cli/internal/trello"
)

func TestChecklistsListCommand(t *testing.T) {
	setupTestAuth(t)
	credStore.Set("default", credentials.Credentials{APIKey: "k", Token: "t", AuthMode: "manual"})
	apiClient = &mockAPI{
		listChecklistsFn: func(ctx context.Context, cardID string) ([]trello.Checklist, error) {
			if cardID != "c1" {
				t.Errorf("card ID = %q, want c1", cardID)
			}
			return []trello.Checklist{{ID: "cl1", Name: "Checklist"}}, nil
		},
	}

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetArgs([]string{"checklists", "list", "--card", "c1"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("checklists list failed: %v", err)
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

func TestChecklistsCreateCommand(t *testing.T) {
	setupTestAuth(t)
	credStore.Set("default", credentials.Credentials{APIKey: "k", Token: "t", AuthMode: "manual"})
	apiClient = &mockAPI{
		createChecklistFn: func(ctx context.Context, cardID, name string) (trello.Checklist, error) {
			if cardID != "c1" || name != "Checklist" {
				t.Fatalf("card/name = %q/%q", cardID, name)
			}
			return trello.Checklist{ID: "cl1", Name: name}, nil
		},
	}

	if err := executeRootArgs("checklists", "create", "--card", "c1", "--name", "Checklist"); err != nil {
		t.Fatalf("checklists create failed: %v", err)
	}
}

func TestChecklistsDeleteCommand(t *testing.T) {
	setupTestAuth(t)
	credStore.Set("default", credentials.Credentials{APIKey: "k", Token: "t", AuthMode: "manual"})
	apiClient = &mockAPI{
		deleteChecklistFn: func(ctx context.Context, checklistID string) error {
			if checklistID != "cl1" {
				t.Errorf("checklist ID = %q, want cl1", checklistID)
			}
			return nil
		},
	}

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetArgs([]string{"checklists", "delete", "--checklist", "cl1"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("checklists delete failed: %v", err)
	}

	var envelope map[string]any
	if err := json.Unmarshal(buf.Bytes(), &envelope); err != nil {
		t.Fatalf("invalid JSON: %v\nraw: %s", err, buf.String())
	}
	data := envelope["data"].(map[string]any)
	if data["deleted"] != true || data["id"] != "cl1" {
		t.Fatalf("data = %+v", data)
	}
}

func TestChecklistsItemsAddCommand(t *testing.T) {
	setupTestAuth(t)
	credStore.Set("default", credentials.Credentials{APIKey: "k", Token: "t", AuthMode: "manual"})
	apiClient = &mockAPI{
		addCheckItemFn: func(ctx context.Context, checklistID, name string) (trello.CheckItem, error) {
			if checklistID != "cl1" || name != "Item" {
				t.Fatalf("checklist/name = %q/%q", checklistID, name)
			}
			return trello.CheckItem{ID: "i1", Name: name}, nil
		},
	}

	if err := executeRootArgs("checklists", "items", "add", "--checklist", "cl1", "--name", "Item"); err != nil {
		t.Fatalf("checklists items add failed: %v", err)
	}
}

func TestChecklistsItemsUpdateCommand(t *testing.T) {
	setupTestAuth(t)
	credStore.Set("default", credentials.Credentials{APIKey: "k", Token: "t", AuthMode: "manual"})
	apiClient = &mockAPI{
		updateCheckItemFn: func(ctx context.Context, cardID, itemID, state string) (trello.CheckItem, error) {
			if cardID != "c1" || itemID != "i1" || state != "complete" {
				t.Fatalf("card/item/state = %q/%q/%q", cardID, itemID, state)
			}
			return trello.CheckItem{ID: itemID, State: state}, nil
		},
	}

	if err := executeRootArgs("checklists", "items", "update", "--card", "c1", "--item", "i1", "--state", "complete"); err != nil {
		t.Fatalf("checklists items update failed: %v", err)
	}
}

func TestChecklistsItemsDeleteCommand(t *testing.T) {
	setupTestAuth(t)
	credStore.Set("default", credentials.Credentials{APIKey: "k", Token: "t", AuthMode: "manual"})
	apiClient = &mockAPI{
		deleteCheckItemFn: func(ctx context.Context, checklistID, itemID string) error {
			if checklistID != "cl1" || itemID != "i1" {
				t.Fatalf("checklist/item = %q/%q", checklistID, itemID)
			}
			return nil
		},
	}

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetArgs([]string{"checklists", "items", "delete", "--checklist", "cl1", "--item", "i1"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("checklists items delete failed: %v", err)
	}

	var envelope map[string]any
	if err := json.Unmarshal(buf.Bytes(), &envelope); err != nil {
		t.Fatalf("invalid JSON: %v\nraw: %s", err, buf.String())
	}
	data := envelope["data"].(map[string]any)
	if data["deleted"] != true || data["id"] != "i1" {
		t.Fatalf("data = %+v", data)
	}
}

func TestChecklistsValidation(t *testing.T) {
	setupTestAuth(t)
	credStore.Set("default", credentials.Credentials{APIKey: "k", Token: "t", AuthMode: "manual"})

	assertContractCode(t, executeRootArgs("checklists", "list"), "VALIDATION_ERROR")
	assertContractCode(t, executeRootArgs("checklists", "create", "--card", "c1"), "VALIDATION_ERROR")
	assertContractCode(t, executeRootArgs("checklists", "delete"), "VALIDATION_ERROR")
	assertContractCode(t, executeRootArgs("checklists", "items", "add", "--checklist", "cl1"), "VALIDATION_ERROR")
	assertContractCode(t, executeRootArgs("checklists", "items", "update", "--card", "c1", "--item", "i1", "--state", "done"), "VALIDATION_ERROR")
	assertContractCode(t, executeRootArgs("checklists", "items", "delete", "--checklist", "cl1"), "VALIDATION_ERROR")
}
