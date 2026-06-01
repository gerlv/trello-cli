package main

import (
	"github.com/gerlv/trello-cli/internal/auth"
	"github.com/gerlv/trello-cli/internal/contract"
	"github.com/gerlv/trello-cli/internal/trello"
	"github.com/spf13/cobra"
)

var cardsCmd = &cobra.Command{
	Use:   "cards",
	Short: "Manage Trello cards",
}

var cardsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List cards by board or list",
	RunE: func(cmd *cobra.Command, args []string) error {
		boardID, _ := cmd.Flags().GetString("board")
		listID, _ := cmd.Flags().GetString("list")
		if err := contract.RequireExactlyOne(map[string]string{
			"board": boardID,
			"list":  listID,
		}); err != nil {
			return err
		}

		creds, err := auth.RequireAuth(getCredStore(), "default")
		if err != nil {
			return err
		}
		api := getAPIClient(creds)
		if boardID != "" {
			cards, err := api.ListCardsByBoard(cmd.Context(), boardID)
			if err != nil {
				return err
			}
			return output(cmd.OutOrStdout(), cards)
		}

		cards, err := api.ListCardsByList(cmd.Context(), listID)
		if err != nil {
			return err
		}
		return output(cmd.OutOrStdout(), cards)
	},
}

var cardsGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a card by ID",
	RunE: func(cmd *cobra.Command, args []string) error {
		cardID, _ := cmd.Flags().GetString("card")
		if err := contract.RequireFlag("card", cardID); err != nil {
			return err
		}

		creds, err := auth.RequireAuth(getCredStore(), "default")
		if err != nil {
			return err
		}
		card, err := getAPIClient(creds).GetCard(cmd.Context(), cardID)
		if err != nil {
			return err
		}
		return output(cmd.OutOrStdout(), card)
	},
}

var cardsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a card",
	RunE: func(cmd *cobra.Command, args []string) error {
		listID, _ := cmd.Flags().GetString("list")
		name, _ := cmd.Flags().GetString("name")
		desc, _ := cmd.Flags().GetString("desc")
		due, _ := cmd.Flags().GetString("due")
		labels, _ := cmd.Flags().GetString("labels")
		members, _ := cmd.Flags().GetString("members")

		if err := contract.RequireFlag("list", listID); err != nil {
			return err
		}
		if err := contract.RequireFlag("name", name); err != nil {
			return err
		}
		if err := contract.ValidateISO8601Optional(due); err != nil {
			return err
		}

		params := trello.CreateCardParams{IDList: listID, Name: name}
		if desc != "" {
			params.Desc = &desc
		}
		if due != "" {
			params.Due = &due
		}
		if labels != "" {
			params.Labels = &labels
		}
		if members != "" {
			params.Members = &members
		}

		creds, err := auth.RequireAuth(getCredStore(), "default")
		if err != nil {
			return err
		}
		card, err := getAPIClient(creds).CreateCard(cmd.Context(), params)
		if err != nil {
			return err
		}
		return output(cmd.OutOrStdout(), card)
	},
}

var cardsUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a card",
	RunE: func(cmd *cobra.Command, args []string) error {
		cardID, _ := cmd.Flags().GetString("card")
		name, _ := cmd.Flags().GetString("name")
		desc, _ := cmd.Flags().GetString("desc")
		due, _ := cmd.Flags().GetString("due")
		labels, _ := cmd.Flags().GetString("labels")
		members, _ := cmd.Flags().GetString("members")

		if err := contract.RequireFlag("card", cardID); err != nil {
			return err
		}
		if err := contract.RequireAtLeastOne(map[string]string{
			"name":    name,
			"desc":    desc,
			"due":     due,
			"labels":  labels,
			"members": members,
		}); err != nil {
			return err
		}
		if err := contract.ValidateISO8601Optional(due); err != nil {
			return err
		}

		params := trello.UpdateCardParams{}
		if name != "" {
			params.Name = &name
		}
		if desc != "" {
			params.Desc = &desc
		}
		if due != "" {
			params.Due = &due
		}
		if labels != "" {
			params.Labels = &labels
		}
		if members != "" {
			params.Members = &members
		}

		creds, err := auth.RequireAuth(getCredStore(), "default")
		if err != nil {
			return err
		}
		card, err := getAPIClient(creds).UpdateCard(cmd.Context(), cardID, params)
		if err != nil {
			return err
		}
		return output(cmd.OutOrStdout(), card)
	},
}

var cardsMoveCmd = &cobra.Command{
	Use:   "move",
	Short: "Move a card to another list",
	RunE: func(cmd *cobra.Command, args []string) error {
		cardID, _ := cmd.Flags().GetString("card")
		listID, _ := cmd.Flags().GetString("list")
		pos, _ := cmd.Flags().GetFloat64("pos")
		posChanged := cmd.Flags().Changed("pos")

		if err := contract.RequireFlag("card", cardID); err != nil {
			return err
		}
		if err := contract.RequireFlag("list", listID); err != nil {
			return err
		}

		var posPtr *float64
		if posChanged {
			posPtr = &pos
		}

		creds, err := auth.RequireAuth(getCredStore(), "default")
		if err != nil {
			return err
		}
		card, err := getAPIClient(creds).MoveCard(cmd.Context(), cardID, listID, posPtr)
		if err != nil {
			return err
		}
		return output(cmd.OutOrStdout(), card)
	},
}

var cardsArchiveCmd = &cobra.Command{
	Use:   "archive",
	Short: "Archive a card",
	RunE: func(cmd *cobra.Command, args []string) error {
		cardID, _ := cmd.Flags().GetString("card")
		if err := contract.RequireFlag("card", cardID); err != nil {
			return err
		}

		creds, err := auth.RequireAuth(getCredStore(), "default")
		if err != nil {
			return err
		}
		card, err := getAPIClient(creds).ArchiveCard(cmd.Context(), cardID)
		if err != nil {
			return err
		}
		return output(cmd.OutOrStdout(), card)
	},
}

var cardsDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a card",
	RunE: func(cmd *cobra.Command, args []string) error {
		cardID, _ := cmd.Flags().GetString("card")
		if err := contract.RequireFlag("card", cardID); err != nil {
			return err
		}

		creds, err := auth.RequireAuth(getCredStore(), "default")
		if err != nil {
			return err
		}
		if err := getAPIClient(creds).DeleteCard(cmd.Context(), cardID); err != nil {
			return err
		}
		return output(cmd.OutOrStdout(), trello.DeleteResult{Deleted: true, ID: cardID})
	},
}

func init() {
	cardsListCmd.Flags().String("board", "", "Board ID")
	cardsListCmd.Flags().String("list", "", "List ID")

	cardsGetCmd.Flags().String("card", "", "Card ID")

	cardsCreateCmd.Flags().String("list", "", "List ID")
	cardsCreateCmd.Flags().String("name", "", "Card name")
	cardsCreateCmd.Flags().String("desc", "", "Card description")
	cardsCreateCmd.Flags().String("due", "", "Due date (ISO-8601)")
	cardsCreateCmd.Flags().String("labels", "", "Comma-separated label IDs")
	cardsCreateCmd.Flags().String("members", "", "Comma-separated member IDs")

	cardsUpdateCmd.Flags().String("card", "", "Card ID")
	cardsUpdateCmd.Flags().String("name", "", "Card name")
	cardsUpdateCmd.Flags().String("desc", "", "Card description")
	cardsUpdateCmd.Flags().String("due", "", "Due date (ISO-8601)")
	cardsUpdateCmd.Flags().String("labels", "", "Comma-separated label IDs")
	cardsUpdateCmd.Flags().String("members", "", "Comma-separated member IDs")

	cardsMoveCmd.Flags().String("card", "", "Card ID")
	cardsMoveCmd.Flags().String("list", "", "Destination list ID")
	cardsMoveCmd.Flags().Float64("pos", 0, "New card position")

	cardsArchiveCmd.Flags().String("card", "", "Card ID")
	cardsDeleteCmd.Flags().String("card", "", "Card ID")

	cardsCmd.AddCommand(cardsListCmd)
	cardsCmd.AddCommand(cardsGetCmd)
	cardsCmd.AddCommand(cardsCreateCmd)
	cardsCmd.AddCommand(cardsUpdateCmd)
	cardsCmd.AddCommand(cardsMoveCmd)
	cardsCmd.AddCommand(cardsArchiveCmd)
	cardsCmd.AddCommand(cardsDeleteCmd)
	rootCmd.AddCommand(cardsCmd)
}
