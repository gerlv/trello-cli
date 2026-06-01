package main

import (
	"github.com/gerlv/trello-cli/internal/auth"
	"github.com/gerlv/trello-cli/internal/contract"
	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search Trello resources",
}

var searchCardsCmd = &cobra.Command{
	Use:   "cards",
	Short: "Search cards",
	RunE: func(cmd *cobra.Command, args []string) error {
		query, _ := cmd.Flags().GetString("query")
		if err := contract.RequireFlag("query", query); err != nil {
			return err
		}
		creds, err := auth.RequireAuth(getCredStore(), "default")
		if err != nil {
			return err
		}
		result, err := getAPIClient(creds).SearchCards(cmd.Context(), query)
		if err != nil {
			return err
		}
		return output(cmd.OutOrStdout(), result)
	},
}

var searchBoardsCmd = &cobra.Command{
	Use:   "boards",
	Short: "Search boards",
	RunE: func(cmd *cobra.Command, args []string) error {
		query, _ := cmd.Flags().GetString("query")
		if err := contract.RequireFlag("query", query); err != nil {
			return err
		}
		creds, err := auth.RequireAuth(getCredStore(), "default")
		if err != nil {
			return err
		}
		result, err := getAPIClient(creds).SearchBoards(cmd.Context(), query)
		if err != nil {
			return err
		}
		return output(cmd.OutOrStdout(), result)
	},
}

func init() {
	searchCardsCmd.Flags().String("query", "", "Search query")
	searchBoardsCmd.Flags().String("query", "", "Search query")

	searchCmd.AddCommand(searchCardsCmd)
	searchCmd.AddCommand(searchBoardsCmd)
	rootCmd.AddCommand(searchCmd)
}
