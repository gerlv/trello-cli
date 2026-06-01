package main

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"github.com/gerlv/trello-cli/internal/credentials"
	"github.com/gerlv/trello-cli/internal/trello"
)

func TestSearchCardsCommand(t *testing.T) {
	setupTestAuth(t)
	credStore.Set("default", credentials.Credentials{APIKey: "k", Token: "t", AuthMode: "manual"})
	apiClient = &mockAPI{
		searchCardsFn: func(ctx context.Context, query string) (trello.CardSearchResult, error) {
			if query != "alpha" {
				t.Fatalf("query = %q, want alpha", query)
			}
			return trello.CardSearchResult{
				Query: "alpha",
				Cards: []trello.Card{{ID: "c1", Name: "Alpha card"}},
			}, nil
		},
	}

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetArgs([]string{"search", "cards", "--query", "alpha"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("search cards failed: %v", err)
	}
	var envelope map[string]any
	if err := json.Unmarshal(buf.Bytes(), &envelope); err != nil {
		t.Fatalf("invalid JSON: %v\nraw: %s", err, buf.String())
	}
	data := envelope["data"].(map[string]any)
	if data["query"] != "alpha" {
		t.Fatalf("data = %+v", data)
	}
	cards := data["cards"].([]any)
	if len(cards) != 1 {
		t.Fatalf("cards = %+v", data["cards"])
	}
	card := cards[0].(map[string]any)
	if card["id"] != "c1" || card["name"] != "Alpha card" {
		t.Fatalf("card = %+v", card)
	}
}

func TestSearchBoardsCommand(t *testing.T) {
	setupTestAuth(t)
	credStore.Set("default", credentials.Credentials{APIKey: "k", Token: "t", AuthMode: "manual"})
	apiClient = &mockAPI{
		searchBoardsFn: func(ctx context.Context, query string) (trello.BoardSearchResult, error) {
			if query != "alpha" {
				t.Fatalf("query = %q, want alpha", query)
			}
			return trello.BoardSearchResult{
				Query:  "alpha",
				Boards: []trello.Board{{ID: "b1", Name: "Alpha board"}},
			}, nil
		},
	}

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetArgs([]string{"search", "boards", "--query", "alpha"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("search boards failed: %v", err)
	}
	var envelope map[string]any
	if err := json.Unmarshal(buf.Bytes(), &envelope); err != nil {
		t.Fatalf("invalid JSON: %v\nraw: %s", err, buf.String())
	}
	data := envelope["data"].(map[string]any)
	if data["query"] != "alpha" {
		t.Fatalf("data = %+v", data)
	}
	boards := data["boards"].([]any)
	if len(boards) != 1 {
		t.Fatalf("boards = %+v", data["boards"])
	}
	board := boards[0].(map[string]any)
	if board["id"] != "b1" || board["name"] != "Alpha board" {
		t.Fatalf("board = %+v", board)
	}
}
