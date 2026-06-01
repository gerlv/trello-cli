package main

import (
	"github.com/gerlv/trello-cli/internal/auth"
	"github.com/gerlv/trello-cli/internal/contract"
	"github.com/gerlv/trello-cli/internal/trello"
	"github.com/spf13/cobra"
)

var listsCmd = &cobra.Command{
	Use:   "lists",
	Short: "Manage Trello lists",
}

var listsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List lists on a board",
	RunE: func(cmd *cobra.Command, args []string) error {
		boardID, _ := cmd.Flags().GetString("board")
		if err := contract.RequireFlag("board", boardID); err != nil {
			return err
		}

		creds, err := auth.RequireAuth(getCredStore(), "default")
		if err != nil {
			return err
		}
		lists, err := getAPIClient(creds).ListLists(cmd.Context(), boardID)
		if err != nil {
			return err
		}
		return output(cmd.OutOrStdout(), lists)
	},
}

var listsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a list on a board",
	RunE: func(cmd *cobra.Command, args []string) error {
		boardID, _ := cmd.Flags().GetString("board")
		name, _ := cmd.Flags().GetString("name")
		if err := contract.RequireFlag("board", boardID); err != nil {
			return err
		}
		if err := contract.RequireFlag("name", name); err != nil {
			return err
		}

		creds, err := auth.RequireAuth(getCredStore(), "default")
		if err != nil {
			return err
		}
		list, err := getAPIClient(creds).CreateList(cmd.Context(), boardID, name)
		if err != nil {
			return err
		}
		return output(cmd.OutOrStdout(), list)
	},
}

var listsUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a list",
	RunE: func(cmd *cobra.Command, args []string) error {
		listID, _ := cmd.Flags().GetString("list")
		name, _ := cmd.Flags().GetString("name")
		pos, _ := cmd.Flags().GetFloat64("pos")
		posChanged := cmd.Flags().Changed("pos")

		if err := contract.RequireFlag("list", listID); err != nil {
			return err
		}

		mutationFlags := map[string]string{"name": name}
		if posChanged {
			mutationFlags["pos"] = "set"
		} else {
			mutationFlags["pos"] = ""
		}
		if err := contract.RequireAtLeastOne(mutationFlags); err != nil {
			return err
		}

		params := trello.UpdateListParams{}
		if name != "" {
			params.Name = &name
		}
		if posChanged {
			params.Pos = &pos
		}

		creds, err := auth.RequireAuth(getCredStore(), "default")
		if err != nil {
			return err
		}
		list, err := getAPIClient(creds).UpdateList(cmd.Context(), listID, params)
		if err != nil {
			return err
		}
		return output(cmd.OutOrStdout(), list)
	},
}

var listsArchiveCmd = &cobra.Command{
	Use:   "archive",
	Short: "Archive a list",
	RunE: func(cmd *cobra.Command, args []string) error {
		listID, _ := cmd.Flags().GetString("list")
		if err := contract.RequireFlag("list", listID); err != nil {
			return err
		}

		creds, err := auth.RequireAuth(getCredStore(), "default")
		if err != nil {
			return err
		}
		list, err := getAPIClient(creds).ArchiveList(cmd.Context(), listID)
		if err != nil {
			return err
		}
		return output(cmd.OutOrStdout(), list)
	},
}

var listsMoveCmd = &cobra.Command{
	Use:   "move",
	Short: "Move a list to another board",
	RunE: func(cmd *cobra.Command, args []string) error {
		listID, _ := cmd.Flags().GetString("list")
		boardID, _ := cmd.Flags().GetString("board")
		pos, _ := cmd.Flags().GetFloat64("pos")
		posChanged := cmd.Flags().Changed("pos")

		if err := contract.RequireFlag("list", listID); err != nil {
			return err
		}
		if err := contract.RequireFlag("board", boardID); err != nil {
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
		list, err := getAPIClient(creds).MoveList(cmd.Context(), listID, boardID, posPtr)
		if err != nil {
			return err
		}
		return output(cmd.OutOrStdout(), list)
	},
}

func init() {
	listsListCmd.Flags().String("board", "", "Board ID")

	listsCreateCmd.Flags().String("board", "", "Board ID")
	listsCreateCmd.Flags().String("name", "", "List name")

	listsUpdateCmd.Flags().String("list", "", "List ID")
	listsUpdateCmd.Flags().String("name", "", "New list name")
	listsUpdateCmd.Flags().Float64("pos", 0, "New list position")

	listsArchiveCmd.Flags().String("list", "", "List ID")

	listsMoveCmd.Flags().String("list", "", "List ID")
	listsMoveCmd.Flags().String("board", "", "Destination board ID")
	listsMoveCmd.Flags().Float64("pos", 0, "New list position")

	listsCmd.AddCommand(listsListCmd)
	listsCmd.AddCommand(listsCreateCmd)
	listsCmd.AddCommand(listsUpdateCmd)
	listsCmd.AddCommand(listsArchiveCmd)
	listsCmd.AddCommand(listsMoveCmd)
	rootCmd.AddCommand(listsCmd)
}
