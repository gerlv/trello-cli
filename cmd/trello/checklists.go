package main

import (
	"github.com/gerlv/trello-cli/internal/auth"
	"github.com/gerlv/trello-cli/internal/contract"
	"github.com/gerlv/trello-cli/internal/trello"
	"github.com/spf13/cobra"
)

func validateState(state string) error {
	if state != "complete" && state != "incomplete" {
		return contract.NewError(contract.ValidationError, "--state must be 'complete' or 'incomplete'")
	}
	return nil
}

var checklistsCmd = &cobra.Command{
	Use:   "checklists",
	Short: "Manage checklists",
}

var checklistsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List checklists on a card",
	RunE: func(cmd *cobra.Command, args []string) error {
		cardID, _ := cmd.Flags().GetString("card")
		if err := contract.RequireFlag("card", cardID); err != nil {
			return err
		}
		creds, err := auth.RequireAuth(getCredStore(), "default")
		if err != nil {
			return err
		}
		checklists, err := getAPIClient(creds).ListChecklists(cmd.Context(), cardID)
		if err != nil {
			return err
		}
		return output(cmd.OutOrStdout(), checklists)
	},
}

var checklistsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a checklist on a card",
	RunE: func(cmd *cobra.Command, args []string) error {
		cardID, _ := cmd.Flags().GetString("card")
		name, _ := cmd.Flags().GetString("name")
		if err := contract.RequireFlag("card", cardID); err != nil {
			return err
		}
		if err := contract.RequireFlag("name", name); err != nil {
			return err
		}
		creds, err := auth.RequireAuth(getCredStore(), "default")
		if err != nil {
			return err
		}
		checklist, err := getAPIClient(creds).CreateChecklist(cmd.Context(), cardID, name)
		if err != nil {
			return err
		}
		return output(cmd.OutOrStdout(), checklist)
	},
}

var checklistsDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a checklist",
	RunE: func(cmd *cobra.Command, args []string) error {
		checklistID, _ := cmd.Flags().GetString("checklist")
		if err := contract.RequireFlag("checklist", checklistID); err != nil {
			return err
		}
		creds, err := auth.RequireAuth(getCredStore(), "default")
		if err != nil {
			return err
		}
		if err := getAPIClient(creds).DeleteChecklist(cmd.Context(), checklistID); err != nil {
			return err
		}
		return output(cmd.OutOrStdout(), trello.DeleteResult{Deleted: true, ID: checklistID})
	},
}

var checklistsItemsCmd = &cobra.Command{
	Use:   "items",
	Short: "Manage checklist items",
}

var checklistsItemsAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add an item to a checklist",
	RunE: func(cmd *cobra.Command, args []string) error {
		checklistID, _ := cmd.Flags().GetString("checklist")
		name, _ := cmd.Flags().GetString("name")
		if err := contract.RequireFlag("checklist", checklistID); err != nil {
			return err
		}
		if err := contract.RequireFlag("name", name); err != nil {
			return err
		}
		creds, err := auth.RequireAuth(getCredStore(), "default")
		if err != nil {
			return err
		}
		item, err := getAPIClient(creds).AddCheckItem(cmd.Context(), checklistID, name)
		if err != nil {
			return err
		}
		return output(cmd.OutOrStdout(), item)
	},
}

var checklistsItemsUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a checklist item state",
	RunE: func(cmd *cobra.Command, args []string) error {
		cardID, _ := cmd.Flags().GetString("card")
		itemID, _ := cmd.Flags().GetString("item")
		state, _ := cmd.Flags().GetString("state")
		if err := contract.RequireFlag("card", cardID); err != nil {
			return err
		}
		if err := contract.RequireFlag("item", itemID); err != nil {
			return err
		}
		if err := contract.RequireFlag("state", state); err != nil {
			return err
		}
		if err := validateState(state); err != nil {
			return err
		}
		creds, err := auth.RequireAuth(getCredStore(), "default")
		if err != nil {
			return err
		}
		item, err := getAPIClient(creds).UpdateCheckItem(cmd.Context(), cardID, itemID, state)
		if err != nil {
			return err
		}
		return output(cmd.OutOrStdout(), item)
	},
}

var checklistsItemsDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a checklist item",
	RunE: func(cmd *cobra.Command, args []string) error {
		checklistID, _ := cmd.Flags().GetString("checklist")
		itemID, _ := cmd.Flags().GetString("item")
		if err := contract.RequireFlag("checklist", checklistID); err != nil {
			return err
		}
		if err := contract.RequireFlag("item", itemID); err != nil {
			return err
		}
		creds, err := auth.RequireAuth(getCredStore(), "default")
		if err != nil {
			return err
		}
		if err := getAPIClient(creds).DeleteCheckItem(cmd.Context(), checklistID, itemID); err != nil {
			return err
		}
		return output(cmd.OutOrStdout(), trello.DeleteResult{Deleted: true, ID: itemID})
	},
}

func init() {
	checklistsListCmd.Flags().String("card", "", "Card ID")
	checklistsCreateCmd.Flags().String("card", "", "Card ID")
	checklistsCreateCmd.Flags().String("name", "", "Checklist name")
	checklistsDeleteCmd.Flags().String("checklist", "", "Checklist ID")

	checklistsItemsAddCmd.Flags().String("checklist", "", "Checklist ID")
	checklistsItemsAddCmd.Flags().String("name", "", "Item name")
	checklistsItemsUpdateCmd.Flags().String("card", "", "Card ID")
	checklistsItemsUpdateCmd.Flags().String("item", "", "Checklist item ID")
	checklistsItemsUpdateCmd.Flags().String("state", "", "Item state (complete|incomplete)")
	checklistsItemsDeleteCmd.Flags().String("checklist", "", "Checklist ID")
	checklistsItemsDeleteCmd.Flags().String("item", "", "Checklist item ID")

	checklistsCmd.AddCommand(checklistsListCmd)
	checklistsCmd.AddCommand(checklistsCreateCmd)
	checklistsCmd.AddCommand(checklistsDeleteCmd)
	checklistsItemsCmd.AddCommand(checklistsItemsAddCmd)
	checklistsItemsCmd.AddCommand(checklistsItemsUpdateCmd)
	checklistsItemsCmd.AddCommand(checklistsItemsDeleteCmd)
	checklistsCmd.AddCommand(checklistsItemsCmd)
	rootCmd.AddCommand(checklistsCmd)
}
