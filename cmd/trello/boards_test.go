package main

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/gerlv/trello-cli/internal/credentials"
	"github.com/gerlv/trello-cli/internal/trello"
)

func TestBoardsListCommand(t *testing.T) {
	setupTestAuth(t)
	credStore.Set("default", credentials.Credentials{APIKey: "k", Token: "t", AuthMode: "manual"})
	apiClient = &mockAPI{
		listBoardsFn: func(ctx context.Context) ([]trello.Board, error) {
			return []trello.Board{
				{ID: "b1", Name: "Board One"},
				{ID: "b2", Name: "Board Two"},
			}, nil
		},
	}

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetArgs([]string{"boards", "list"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("boards list failed: %v", err)
	}

	var envelope map[string]any
	if err := json.Unmarshal(buf.Bytes(), &envelope); err != nil {
		t.Fatalf("invalid JSON: %v\nraw: %s", err, buf.String())
	}

	if envelope["ok"] != true {
		t.Errorf("ok = %v, want true", envelope["ok"])
	}
	data, ok := envelope["data"].([]any)
	if !ok {
		t.Fatal("data should be an array")
	}
	if len(data) != 2 {
		t.Errorf("len(data) = %d, want 2", len(data))
	}
}

func TestBoardsListRequiresAuth(t *testing.T) {
	setupTestAuth(t)

	rootCmd.SetArgs([]string{"boards", "list"})
	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("boards list should fail without auth")
	}
}

func TestBoardsGetCommand(t *testing.T) {
	setupTestAuth(t)
	credStore.Set("default", credentials.Credentials{APIKey: "k", Token: "t", AuthMode: "manual"})
	apiClient = &mockAPI{
		getBoardFn: func(ctx context.Context, id string) (trello.Board, error) {
			if id != "b1" {
				t.Errorf("board ID = %q, want %q", id, "b1")
			}
			return trello.Board{ID: "b1", Name: "My Board"}, nil
		},
	}

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetArgs([]string{"boards", "get", "--board", "b1"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("boards get failed: %v", err)
	}

	var envelope map[string]any
	if err := json.Unmarshal(buf.Bytes(), &envelope); err != nil {
		t.Fatalf("invalid JSON: %v\nraw: %s", err, buf.String())
	}

	data := envelope["data"].(map[string]any)
	if data["id"] != "b1" {
		t.Errorf("data.id = %v, want b1", data["id"])
	}
}

func TestBoardsGetMissingFlag(t *testing.T) {
	setupTestAuth(t)
	credStore.Set("default", credentials.Credentials{APIKey: "k", Token: "t", AuthMode: "manual"})

	rootCmd.SetArgs([]string{"boards", "get"})
	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("boards get should fail without --board")
	}
}

func TestBoardsCreateCommand(t *testing.T) {
	setupTestAuth(t)
	credStore.Set("default", credentials.Credentials{APIKey: "k", Token: "t", AuthMode: "manual"})
	desc := "Board description"
	orgID := "org1"
	sourceBoardID := "template1"
	apiClient = &mockAPI{
		createBoardFn: func(ctx context.Context, params trello.CreateBoardParams) (trello.Board, error) {
			if params.Name != "Project Board" {
				t.Fatalf("params.Name = %q, want %q", params.Name, "Project Board")
			}
			if params.Desc == nil || *params.Desc != desc {
				t.Fatalf("params.Desc = %v, want %q", params.Desc, desc)
			}
			if params.DefaultLists == nil || !*params.DefaultLists {
				t.Fatalf("params.DefaultLists = %v, want true", params.DefaultLists)
			}
			if params.DefaultLabels == nil || !*params.DefaultLabels {
				t.Fatalf("params.DefaultLabels = %v, want true", params.DefaultLabels)
			}
			if params.IDOrganization == nil || *params.IDOrganization != orgID {
				t.Fatalf("params.IDOrganization = %v, want %q", params.IDOrganization, orgID)
			}
			if params.IDBoardSource == nil || *params.IDBoardSource != sourceBoardID {
				t.Fatalf("params.IDBoardSource = %v, want %q", params.IDBoardSource, sourceBoardID)
			}
			return trello.Board{ID: "b3", Name: params.Name, Desc: desc}, nil
		},
	}

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetArgs([]string{
		"boards", "create",
		"--name", "Project Board",
		"--desc", desc,
		"--default-lists",
		"--default-labels",
		"--organization", orgID,
		"--source-board", sourceBoardID,
	})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("boards create failed: %v", err)
	}

	var envelope map[string]any
	if err := json.Unmarshal(buf.Bytes(), &envelope); err != nil {
		t.Fatalf("invalid JSON: %v\nraw: %s", err, buf.String())
	}

	data := envelope["data"].(map[string]any)
	if data["id"] != "b3" {
		t.Errorf("data.id = %v, want b3", data["id"])
	}
}

func TestBoardsCreateMissingName(t *testing.T) {
	setupTestAuth(t)
	credStore.Set("default", credentials.Credentials{APIKey: "k", Token: "t", AuthMode: "manual"})
	assertContractCode(t, executeRootArgs("boards", "create"), "VALIDATION_ERROR")
}

func TestBoardsHelpIncludesCreate(t *testing.T) {
	setupTestAuth(t)

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)
	rootCmd.SetArgs([]string{"boards", "--help"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("boards help failed: %v", err)
	}

	help := buf.String()
	for _, want := range []string{"create", "Create a board"} {
		if !strings.Contains(help, want) {
			t.Fatalf("help missing %q\nhelp:\n%s", want, help)
		}
	}
}

func TestBoardsCreateHelpIncludesFlags(t *testing.T) {
	setupTestAuth(t)

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)
	rootCmd.SetArgs([]string{"boards", "create", "--help"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("boards create help failed: %v", err)
	}

	help := buf.String()
	for _, want := range []string{"--name", "--desc", "--default-lists", "--default-labels", "--organization", "--source-board"} {
		if !strings.Contains(help, want) {
			t.Fatalf("help missing %q\nhelp:\n%s", want, help)
		}
	}
}
