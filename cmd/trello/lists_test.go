package main

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"github.com/gerlv/trello-cli/internal/contract"
	"github.com/gerlv/trello-cli/internal/credentials"
	"github.com/gerlv/trello-cli/internal/trello"
)

func TestListsListCommand(t *testing.T) {
	setupTestAuth(t)
	credStore.Set("default", credentials.Credentials{APIKey: "k", Token: "t", AuthMode: "manual"})
	apiClient = &mockAPI{
		listListsFn: func(ctx context.Context, boardID string) ([]trello.List, error) {
			if boardID != "b1" {
				t.Errorf("board ID = %q, want %q", boardID, "b1")
			}
			return []trello.List{
				{ID: "l1", Name: "Todo"},
				{ID: "l2", Name: "Done"},
			}, nil
		},
	}

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetArgs([]string{"lists", "list", "--board", "b1"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("lists list failed: %v", err)
	}

	var envelope map[string]any
	if err := json.Unmarshal(buf.Bytes(), &envelope); err != nil {
		t.Fatalf("invalid JSON: %v\nraw: %s", err, buf.String())
	}
	data, ok := envelope["data"].([]any)
	if !ok {
		t.Fatal("data should be an array")
	}
	if len(data) != 2 {
		t.Fatalf("len(data) = %d, want 2", len(data))
	}
}

func TestListsListMissingBoard(t *testing.T) {
	setupTestAuth(t)
	credStore.Set("default", credentials.Credentials{APIKey: "k", Token: "t", AuthMode: "manual"})

	rootCmd.SetArgs([]string{"lists", "list"})
	err := rootCmd.Execute()
	assertContractCode(t, err, contract.ValidationError)
}

func TestListsListRequiresAuth(t *testing.T) {
	setupTestAuth(t)

	rootCmd.SetArgs([]string{"lists", "list", "--board", "b1"})
	err := rootCmd.Execute()
	assertContractCode(t, err, contract.AuthRequired)
}

func TestListsCreateCommand(t *testing.T) {
	setupTestAuth(t)
	credStore.Set("default", credentials.Credentials{APIKey: "k", Token: "t", AuthMode: "manual"})
	apiClient = &mockAPI{
		createListFn: func(ctx context.Context, boardID, name string) (trello.List, error) {
			if boardID != "b1" {
				t.Errorf("board ID = %q, want %q", boardID, "b1")
			}
			if name != "New List" {
				t.Errorf("name = %q, want %q", name, "New List")
			}
			return trello.List{ID: "l1", Name: name, IDBoard: boardID}, nil
		},
	}

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetArgs([]string{"lists", "create", "--board", "b1", "--name", "New List"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("lists create failed: %v", err)
	}

	var envelope map[string]any
	if err := json.Unmarshal(buf.Bytes(), &envelope); err != nil {
		t.Fatalf("invalid JSON: %v\nraw: %s", err, buf.String())
	}
	data := envelope["data"].(map[string]any)
	if data["id"] != "l1" {
		t.Fatalf("data.id = %v, want l1", data["id"])
	}
}

func TestListsCreateMissingBoard(t *testing.T) {
	setupTestAuth(t)
	credStore.Set("default", credentials.Credentials{APIKey: "k", Token: "t", AuthMode: "manual"})

	rootCmd.SetArgs([]string{"lists", "create", "--name", "New List"})
	err := rootCmd.Execute()
	assertContractCode(t, err, contract.ValidationError)
}

func TestListsCreateMissingName(t *testing.T) {
	setupTestAuth(t)
	credStore.Set("default", credentials.Credentials{APIKey: "k", Token: "t", AuthMode: "manual"})

	rootCmd.SetArgs([]string{"lists", "create", "--board", "b1"})
	err := rootCmd.Execute()
	assertContractCode(t, err, contract.ValidationError)
}

func TestListsUpdateCommand(t *testing.T) {
	setupTestAuth(t)
	credStore.Set("default", credentials.Credentials{APIKey: "k", Token: "t", AuthMode: "manual"})
	apiClient = &mockAPI{
		updateListFn: func(ctx context.Context, listID string, params trello.UpdateListParams) (trello.List, error) {
			if listID != "l1" {
				t.Errorf("list ID = %q, want %q", listID, "l1")
			}
			if params.Name == nil || *params.Name != "Doing" {
				t.Fatalf("params.Name = %v, want Doing", params.Name)
			}
			return trello.List{ID: "l1", Name: "Doing"}, nil
		},
	}

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetArgs([]string{"lists", "update", "--list", "l1", "--name", "Doing"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("lists update failed: %v", err)
	}

	var envelope map[string]any
	if err := json.Unmarshal(buf.Bytes(), &envelope); err != nil {
		t.Fatalf("invalid JSON: %v\nraw: %s", err, buf.String())
	}
	data := envelope["data"].(map[string]any)
	if data["name"] != "Doing" {
		t.Fatalf("data.name = %v, want Doing", data["name"])
	}
}

func TestListsUpdateMissingList(t *testing.T) {
	setupTestAuth(t)
	credStore.Set("default", credentials.Credentials{APIKey: "k", Token: "t", AuthMode: "manual"})

	rootCmd.SetArgs([]string{"lists", "update", "--name", "Doing"})
	err := rootCmd.Execute()
	assertContractCode(t, err, contract.ValidationError)
}

func TestListsUpdateRequiresMutation(t *testing.T) {
	setupTestAuth(t)
	credStore.Set("default", credentials.Credentials{APIKey: "k", Token: "t", AuthMode: "manual"})

	rootCmd.SetArgs([]string{"lists", "update", "--list", "l1"})
	err := rootCmd.Execute()
	assertContractCode(t, err, contract.ValidationError)
}

func TestListsArchiveCommand(t *testing.T) {
	setupTestAuth(t)
	credStore.Set("default", credentials.Credentials{APIKey: "k", Token: "t", AuthMode: "manual"})
	apiClient = &mockAPI{
		archiveListFn: func(ctx context.Context, listID string) (trello.List, error) {
			if listID != "l1" {
				t.Errorf("list ID = %q, want %q", listID, "l1")
			}
			return trello.List{ID: "l1", Closed: true}, nil
		},
	}

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetArgs([]string{"lists", "archive", "--list", "l1"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("lists archive failed: %v", err)
	}

	var envelope map[string]any
	if err := json.Unmarshal(buf.Bytes(), &envelope); err != nil {
		t.Fatalf("invalid JSON: %v\nraw: %s", err, buf.String())
	}
	data := envelope["data"].(map[string]any)
	if data["closed"] != true {
		t.Fatalf("data.closed = %v, want true", data["closed"])
	}
}

func TestListsArchiveMissingList(t *testing.T) {
	setupTestAuth(t)
	credStore.Set("default", credentials.Credentials{APIKey: "k", Token: "t", AuthMode: "manual"})

	rootCmd.SetArgs([]string{"lists", "archive"})
	err := rootCmd.Execute()
	assertContractCode(t, err, contract.ValidationError)
}

func TestListsMoveCommand(t *testing.T) {
	setupTestAuth(t)
	credStore.Set("default", credentials.Credentials{APIKey: "k", Token: "t", AuthMode: "manual"})
	apiClient = &mockAPI{
		moveListFn: func(ctx context.Context, listID, boardID string, pos *float64) (trello.List, error) {
			if listID != "l1" {
				t.Errorf("list ID = %q, want %q", listID, "l1")
			}
			if boardID != "b2" {
				t.Errorf("board ID = %q, want %q", boardID, "b2")
			}
			if pos == nil || *pos != 2.5 {
				t.Fatalf("pos = %v, want 2.5", pos)
			}
			return trello.List{ID: "l1", IDBoard: "b2", Pos: 2.5}, nil
		},
	}

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetArgs([]string{"lists", "move", "--list", "l1", "--board", "b2", "--pos", "2.5"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("lists move failed: %v", err)
	}

	var envelope map[string]any
	if err := json.Unmarshal(buf.Bytes(), &envelope); err != nil {
		t.Fatalf("invalid JSON: %v\nraw: %s", err, buf.String())
	}
	data := envelope["data"].(map[string]any)
	if data["idBoard"] != "b2" {
		t.Fatalf("data.idBoard = %v, want b2", data["idBoard"])
	}
}

func TestListsMoveMissingList(t *testing.T) {
	setupTestAuth(t)
	credStore.Set("default", credentials.Credentials{APIKey: "k", Token: "t", AuthMode: "manual"})

	rootCmd.SetArgs([]string{"lists", "move", "--board", "b2"})
	err := rootCmd.Execute()
	assertContractCode(t, err, contract.ValidationError)
}

func TestListsMoveMissingBoard(t *testing.T) {
	setupTestAuth(t)
	credStore.Set("default", credentials.Credentials{APIKey: "k", Token: "t", AuthMode: "manual"})

	rootCmd.SetArgs([]string{"lists", "move", "--list", "l1"})
	err := rootCmd.Execute()
	assertContractCode(t, err, contract.ValidationError)
}
