package main

import (
	"github.com/gerlv/trello-cli/internal/auth"
	"github.com/gerlv/trello-cli/internal/contract"
	"github.com/gerlv/trello-cli/internal/trello"
	"github.com/spf13/cobra"
)

var membersCmd = &cobra.Command{
	Use:   "members",
	Short: "Manage board members",
}

var membersListCmd = &cobra.Command{
	Use:   "list",
	Short: "List members on a board",
	RunE: func(cmd *cobra.Command, args []string) error {
		boardID, _ := cmd.Flags().GetString("board")
		if err := contract.RequireFlag("board", boardID); err != nil {
			return err
		}
		creds, err := auth.RequireAuth(getCredStore(), "default")
		if err != nil {
			return err
		}
		members, err := getAPIClient(creds).ListMembers(cmd.Context(), boardID)
		if err != nil {
			return err
		}
		return output(cmd.OutOrStdout(), members)
	},
}

var membersAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a member to a card",
	RunE: func(cmd *cobra.Command, args []string) error {
		cardID, _ := cmd.Flags().GetString("card")
		memberID, _ := cmd.Flags().GetString("member")
		if err := contract.RequireFlag("card", cardID); err != nil {
			return err
		}
		if err := contract.RequireFlag("member", memberID); err != nil {
			return err
		}
		creds, err := auth.RequireAuth(getCredStore(), "default")
		if err != nil {
			return err
		}
		if err := getAPIClient(creds).AddMemberToCard(cmd.Context(), cardID, memberID); err != nil {
			return err
		}
		return output(cmd.OutOrStdout(), trello.ActionResult{Success: true, ID: memberID})
	},
}

var membersRemoveCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove a member from a card",
	RunE: func(cmd *cobra.Command, args []string) error {
		cardID, _ := cmd.Flags().GetString("card")
		memberID, _ := cmd.Flags().GetString("member")
		if err := contract.RequireFlag("card", cardID); err != nil {
			return err
		}
		if err := contract.RequireFlag("member", memberID); err != nil {
			return err
		}
		creds, err := auth.RequireAuth(getCredStore(), "default")
		if err != nil {
			return err
		}
		if err := getAPIClient(creds).RemoveMemberFromCard(cmd.Context(), cardID, memberID); err != nil {
			return err
		}
		return output(cmd.OutOrStdout(), trello.ActionResult{Success: true, ID: memberID})
	},
}

func init() {
	membersListCmd.Flags().String("board", "", "Board ID")
	membersAddCmd.Flags().String("card", "", "Card ID")
	membersAddCmd.Flags().String("member", "", "Member ID")
	membersRemoveCmd.Flags().String("card", "", "Card ID")
	membersRemoveCmd.Flags().String("member", "", "Member ID")

	membersCmd.AddCommand(membersListCmd)
	membersCmd.AddCommand(membersAddCmd)
	membersCmd.AddCommand(membersRemoveCmd)
	rootCmd.AddCommand(membersCmd)
}
