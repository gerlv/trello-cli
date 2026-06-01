package main

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"github.com/gerlv/trello-cli/internal/credentials"
	"github.com/gerlv/trello-cli/internal/trello"
)

func TestCardsListByBoardCommand(t *testing.T) {
	setupTestAuth(t)
	credStore.Set("default", credentials.Credentials{APIKey: "k", Token: "t", AuthMode: "manual"})
	apiClient = &mockAPI{
		listCardsByBoardFn: func(ctx context.Context, boardID string) ([]trello.Card, error) {
			if boardID != "b1" {
				t.Errorf("board ID = %q, want %q", boardID, "b1")
			}
			return []trello.Card{{ID: "c1", Name: "Card One"}}, nil
		},
	}

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetArgs([]string{"cards", "list", "--board", "b1"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("cards list failed: %v", err)
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

func TestCardsListByListCommand(t *testing.T) {
	setupTestAuth(t)
	credStore.Set("default", credentials.Credentials{APIKey: "k", Token: "t", AuthMode: "manual"})
	apiClient = &mockAPI{
		listCardsByListFn: func(ctx context.Context, listID string) ([]trello.Card, error) {
			if listID != "l1" {
				t.Errorf("list ID = %q, want %q", listID, "l1")
			}
			return []trello.Card{{ID: "c1", Name: "Card One"}}, nil
		},
	}

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetArgs([]string{"cards", "list", "--list", "l1"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("cards list failed: %v", err)
	}
}

func TestCardsListRequiresExactlyOne(t *testing.T) {
	setupTestAuth(t)
	credStore.Set("default", credentials.Credentials{APIKey: "k", Token: "t", AuthMode: "manual"})

	err := executeRootArgs("cards", "list")
	assertContractCode(t, err, "VALIDATION_ERROR")

	err = executeRootArgs("cards", "list", "--board", "b1", "--list", "l1")
	assertContractCode(t, err, "VALIDATION_ERROR")
}

func TestCardsGetCommand(t *testing.T) {
	setupTestAuth(t)
	credStore.Set("default", credentials.Credentials{APIKey: "k", Token: "t", AuthMode: "manual"})
	apiClient = &mockAPI{
		getCardFn: func(ctx context.Context, cardID string) (trello.Card, error) {
			if cardID != "c1" {
				t.Errorf("card ID = %q, want %q", cardID, "c1")
			}
			return trello.Card{ID: "c1", Name: "Card One"}, nil
		},
	}

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetArgs([]string{"cards", "get", "--card", "c1"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("cards get failed: %v", err)
	}
}

func TestCardsGetMissingCard(t *testing.T) {
	setupTestAuth(t)
	credStore.Set("default", credentials.Credentials{APIKey: "k", Token: "t", AuthMode: "manual"})
	assertContractCode(t, executeRootArgs("cards", "get"), "VALIDATION_ERROR")
}

func TestCardsCreateCommand(t *testing.T) {
	setupTestAuth(t)
	credStore.Set("default", credentials.Credentials{APIKey: "k", Token: "t", AuthMode: "manual"})
	desc := "Notes"
	due := "2026-03-13T12:00:00Z"
	labels := "lab1,lab2"
	members := "mem1,mem2"
	apiClient = &mockAPI{
		createCardFn: func(ctx context.Context, params trello.CreateCardParams) (trello.Card, error) {
			if params.IDList != "l1" || params.Name != "Card One" {
				t.Fatalf("params = %+v", params)
			}
			if params.Desc == nil || *params.Desc != desc {
				t.Fatalf("desc = %v", params.Desc)
			}
			if params.Due == nil || *params.Due != due {
				t.Fatalf("due = %v", params.Due)
			}
			if params.Labels == nil || *params.Labels != labels {
				t.Fatalf("labels = %v, want %q", params.Labels, labels)
			}
			if params.Members == nil || *params.Members != members {
				t.Fatalf("members = %v, want %q", params.Members, members)
			}
			return trello.Card{ID: "c1", Name: params.Name, IDList: params.IDList, Desc: desc, Due: &due}, nil
		},
	}

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetArgs([]string{"cards", "create", "--list", "l1", "--name", "Card One", "--desc", desc, "--due", due, "--labels", labels, "--members", members})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("cards create failed: %v", err)
	}
}

func TestCardsCreateMissingList(t *testing.T) {
	setupTestAuth(t)
	credStore.Set("default", credentials.Credentials{APIKey: "k", Token: "t", AuthMode: "manual"})
	assertContractCode(t, executeRootArgs("cards", "create", "--name", "Card One"), "VALIDATION_ERROR")
}

func TestCardsCreateMissingName(t *testing.T) {
	setupTestAuth(t)
	credStore.Set("default", credentials.Credentials{APIKey: "k", Token: "t", AuthMode: "manual"})
	assertContractCode(t, executeRootArgs("cards", "create", "--list", "l1"), "VALIDATION_ERROR")
}

func TestCardsCreateInvalidDue(t *testing.T) {
	setupTestAuth(t)
	credStore.Set("default", credentials.Credentials{APIKey: "k", Token: "t", AuthMode: "manual"})
	assertContractCode(t, executeRootArgs("cards", "create", "--list", "l1", "--name", "Card One", "--due", "tomorrow-ish"), "VALIDATION_ERROR")
}

func TestCardsCreateValidDue(t *testing.T) {
	setupTestAuth(t)
	credStore.Set("default", credentials.Credentials{APIKey: "k", Token: "t", AuthMode: "manual"})
	due := "2026-03-13T12:00:00Z"
	apiClient = &mockAPI{
		createCardFn: func(ctx context.Context, params trello.CreateCardParams) (trello.Card, error) {
			return trello.Card{ID: "c1", Name: params.Name, IDList: params.IDList, Due: params.Due}, nil
		},
	}

	if err := executeRootArgs("cards", "create", "--list", "l1", "--name", "Card One", "--due", due); err != nil {
		t.Fatalf("cards create with valid due failed: %v", err)
	}
}

func TestCardsUpdateCommand(t *testing.T) {
	setupTestAuth(t)
	credStore.Set("default", credentials.Credentials{APIKey: "k", Token: "t", AuthMode: "manual"})
	name := "Updated"
	apiClient = &mockAPI{
		updateCardFn: func(ctx context.Context, cardID string, params trello.UpdateCardParams) (trello.Card, error) {
			if cardID != "c1" {
				t.Errorf("card ID = %q, want c1", cardID)
			}
			if params.Name == nil || *params.Name != name {
				t.Fatalf("params.Name = %v", params.Name)
			}
			return trello.Card{ID: "c1", Name: name}, nil
		},
	}

	if err := executeRootArgs("cards", "update", "--card", "c1", "--name", name); err != nil {
		t.Fatalf("cards update failed: %v", err)
	}
}

func TestCardsUpdateMissingCard(t *testing.T) {
	setupTestAuth(t)
	credStore.Set("default", credentials.Credentials{APIKey: "k", Token: "t", AuthMode: "manual"})
	assertContractCode(t, executeRootArgs("cards", "update", "--name", "Updated"), "VALIDATION_ERROR")
}

func TestCardsUpdateRequiresMutation(t *testing.T) {
	setupTestAuth(t)
	credStore.Set("default", credentials.Credentials{APIKey: "k", Token: "t", AuthMode: "manual"})
	assertContractCode(t, executeRootArgs("cards", "update", "--card", "c1"), "VALIDATION_ERROR")
}

func TestCardsUpdateInvalidDue(t *testing.T) {
	setupTestAuth(t)
	credStore.Set("default", credentials.Credentials{APIKey: "k", Token: "t", AuthMode: "manual"})
	assertContractCode(t, executeRootArgs("cards", "update", "--card", "c1", "--due", "not-a-date"), "VALIDATION_ERROR")
}

func TestCardsMoveCommand(t *testing.T) {
	setupTestAuth(t)
	credStore.Set("default", credentials.Credentials{APIKey: "k", Token: "t", AuthMode: "manual"})
	apiClient = &mockAPI{
		moveCardFn: func(ctx context.Context, cardID, listID string, pos *float64) (trello.Card, error) {
			if cardID != "c1" || listID != "l2" {
				t.Fatalf("card/list IDs = %q/%q", cardID, listID)
			}
			if pos == nil || *pos != 2.5 {
				t.Fatalf("pos = %v", pos)
			}
			return trello.Card{ID: "c1", IDList: "l2"}, nil
		},
	}

	if err := executeRootArgs("cards", "move", "--card", "c1", "--list", "l2", "--pos", "2.5"); err != nil {
		t.Fatalf("cards move failed: %v", err)
	}
}

func TestCardsMoveMissingFlags(t *testing.T) {
	setupTestAuth(t)
	credStore.Set("default", credentials.Credentials{APIKey: "k", Token: "t", AuthMode: "manual"})
	assertContractCode(t, executeRootArgs("cards", "move", "--list", "l2"), "VALIDATION_ERROR")

	setupTestAuth(t)
	credStore.Set("default", credentials.Credentials{APIKey: "k", Token: "t", AuthMode: "manual"})
	assertContractCode(t, executeRootArgs("cards", "move", "--card", "c1"), "VALIDATION_ERROR")
}

func TestCardsArchiveCommand(t *testing.T) {
	setupTestAuth(t)
	credStore.Set("default", credentials.Credentials{APIKey: "k", Token: "t", AuthMode: "manual"})
	apiClient = &mockAPI{
		archiveCardFn: func(ctx context.Context, cardID string) (trello.Card, error) {
			return trello.Card{ID: cardID, Closed: true}, nil
		},
	}

	if err := executeRootArgs("cards", "archive", "--card", "c1"); err != nil {
		t.Fatalf("cards archive failed: %v", err)
	}
}

func TestCardsDeleteCommand(t *testing.T) {
	setupTestAuth(t)
	credStore.Set("default", credentials.Credentials{APIKey: "k", Token: "t", AuthMode: "manual"})
	apiClient = &mockAPI{
		deleteCardFn: func(ctx context.Context, cardID string) error {
			if cardID != "c1" {
				t.Errorf("card ID = %q, want c1", cardID)
			}
			return nil
		},
	}

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetArgs([]string{"cards", "delete", "--card", "c1"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("cards delete failed: %v", err)
	}

	var envelope map[string]any
	if err := json.Unmarshal(buf.Bytes(), &envelope); err != nil {
		t.Fatalf("invalid JSON: %v\nraw: %s", err, buf.String())
	}
	data := envelope["data"].(map[string]any)
	if data["deleted"] != true || data["id"] != "c1" {
		t.Fatalf("data = %+v", data)
	}
}

func TestCardsDeleteMissingCard(t *testing.T) {
	setupTestAuth(t)
	credStore.Set("default", credentials.Credentials{APIKey: "k", Token: "t", AuthMode: "manual"})
	assertContractCode(t, executeRootArgs("cards", "delete"), "VALIDATION_ERROR")
}
