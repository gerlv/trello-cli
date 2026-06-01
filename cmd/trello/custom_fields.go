package main

import (
	"fmt"

	"github.com/gerlv/trello-cli/internal/auth"
	"github.com/gerlv/trello-cli/internal/contract"
	"github.com/gerlv/trello-cli/internal/trello"
	"github.com/spf13/cobra"
)

var allowedFieldTypes = map[string]bool{
	"text":     true,
	"number":   true,
	"date":     true,
	"checkbox": true,
	"list":     true,
}

// ── top-level command ──────────────────────────────────────────────

var customFieldsCmd = &cobra.Command{
	Use:   "custom-fields",
	Short: "Manage Trello custom fields",
}

// ── definition subcommands ─────────────────────────────────────────

var cfListCmd = &cobra.Command{
	Use:   "list",
	Short: "List custom fields on a board",
	RunE: func(cmd *cobra.Command, args []string) error {
		boardID, _ := cmd.Flags().GetString("board")
		if err := contract.RequireFlag("board", boardID); err != nil {
			return err
		}
		creds, err := auth.RequireAuth(getCredStore(), "default")
		if err != nil {
			return err
		}
		result, err := getAPIClient(creds).ListCustomFieldsByBoard(cmd.Context(), boardID)
		if err != nil {
			return err
		}
		return output(cmd.OutOrStdout(), result)
	},
}

var cfGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a custom field by ID",
	RunE: func(cmd *cobra.Command, args []string) error {
		fieldID, _ := cmd.Flags().GetString("field")
		if err := contract.RequireFlag("field", fieldID); err != nil {
			return err
		}
		creds, err := auth.RequireAuth(getCredStore(), "default")
		if err != nil {
			return err
		}
		result, err := getAPIClient(creds).GetCustomField(cmd.Context(), fieldID)
		if err != nil {
			return err
		}
		return output(cmd.OutOrStdout(), result)
	},
}

var cfCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a custom field on a board",
	RunE: func(cmd *cobra.Command, args []string) error {
		boardID, _ := cmd.Flags().GetString("board")
		name, _ := cmd.Flags().GetString("name")
		fieldType, _ := cmd.Flags().GetString("type")
		cardFront, _ := cmd.Flags().GetBool("card-front")
		options, _ := cmd.Flags().GetStringArray("option")

		if err := contract.RequireFlag("board", boardID); err != nil {
			return err
		}
		if err := contract.RequireFlag("name", name); err != nil {
			return err
		}
		if err := contract.RequireFlag("type", fieldType); err != nil {
			return err
		}
		if !allowedFieldTypes[fieldType] {
			return contract.NewError(contract.ValidationError, fmt.Sprintf("--type must be one of: text, number, date, checkbox, list; got %q", fieldType))
		}
		optionChanged := cmd.Flags().Changed("option")
		if optionChanged && fieldType != "list" {
			return contract.NewError(contract.ValidationError, "--option is only valid when --type is list")
		}

		params := trello.CreateCustomFieldParams{
			IDModel:   boardID,
			ModelType: "board",
			Name:      name,
			Type:      fieldType,
			Display:   trello.CustomFieldDisplay{CardFront: cardFront},
		}
		if optionChanged {
			for _, optText := range options {
				params.Options = append(params.Options, trello.CustomFieldOption{
					Value: trello.CustomFieldOptionValue{Text: optText},
				})
			}
		}

		creds, err := auth.RequireAuth(getCredStore(), "default")
		if err != nil {
			return err
		}
		result, err := getAPIClient(creds).CreateCustomField(cmd.Context(), params)
		if err != nil {
			return err
		}
		return output(cmd.OutOrStdout(), result)
	},
}

var cfUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a custom field",
	RunE: func(cmd *cobra.Command, args []string) error {
		fieldID, _ := cmd.Flags().GetString("field")
		name, _ := cmd.Flags().GetString("name")
		cardFront, _ := cmd.Flags().GetBool("card-front")

		if err := contract.RequireFlag("field", fieldID); err != nil {
			return err
		}

		cardFrontChanged := cmd.Flags().Changed("card-front")
		if name == "" && !cardFrontChanged {
			return contract.NewError(contract.ValidationError, "at least one of --name or --card-front must be provided")
		}

		params := trello.UpdateCustomFieldParams{}
		if name != "" {
			params.Name = &name
		}
		if cardFrontChanged {
			params.Display = &trello.CustomFieldDisplay{CardFront: cardFront}
		}

		creds, err := auth.RequireAuth(getCredStore(), "default")
		if err != nil {
			return err
		}
		result, err := getAPIClient(creds).UpdateCustomField(cmd.Context(), fieldID, params)
		if err != nil {
			return err
		}
		return output(cmd.OutOrStdout(), result)
	},
}

var cfDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a custom field",
	RunE: func(cmd *cobra.Command, args []string) error {
		fieldID, _ := cmd.Flags().GetString("field")
		if err := contract.RequireFlag("field", fieldID); err != nil {
			return err
		}
		creds, err := auth.RequireAuth(getCredStore(), "default")
		if err != nil {
			return err
		}
		if err := getAPIClient(creds).DeleteCustomField(cmd.Context(), fieldID); err != nil {
			return err
		}
		return output(cmd.OutOrStdout(), trello.DeleteResult{Deleted: true, ID: fieldID})
	},
}

// ── options subgroup ───────────────────────────────────────────────

var cfOptionsCmd = &cobra.Command{
	Use:   "options",
	Short: "Manage custom field options",
}

var cfOptionsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List options for a custom field",
	RunE: func(cmd *cobra.Command, args []string) error {
		fieldID, _ := cmd.Flags().GetString("field")
		if err := contract.RequireFlag("field", fieldID); err != nil {
			return err
		}
		creds, err := auth.RequireAuth(getCredStore(), "default")
		if err != nil {
			return err
		}
		result, err := getAPIClient(creds).ListCustomFieldOptions(cmd.Context(), fieldID)
		if err != nil {
			return err
		}
		return output(cmd.OutOrStdout(), result)
	},
}

var cfOptionsAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add an option to a custom field",
	RunE: func(cmd *cobra.Command, args []string) error {
		fieldID, _ := cmd.Flags().GetString("field")
		text, _ := cmd.Flags().GetString("text")
		color, _ := cmd.Flags().GetString("color")

		if err := contract.RequireFlag("field", fieldID); err != nil {
			return err
		}
		if err := contract.RequireFlag("text", text); err != nil {
			return err
		}

		params := trello.CreateCustomFieldOptionParams{
			Value: trello.CustomFieldOptionValue{Text: text},
			Color: color,
		}

		creds, err := auth.RequireAuth(getCredStore(), "default")
		if err != nil {
			return err
		}
		result, err := getAPIClient(creds).CreateCustomFieldOption(cmd.Context(), fieldID, params)
		if err != nil {
			return err
		}
		return output(cmd.OutOrStdout(), result)
	},
}

var cfOptionsUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a custom field option",
	RunE: func(cmd *cobra.Command, args []string) error {
		fieldID, _ := cmd.Flags().GetString("field")
		optionID, _ := cmd.Flags().GetString("option")
		text, _ := cmd.Flags().GetString("text")
		color, _ := cmd.Flags().GetString("color")

		if err := contract.RequireFlag("field", fieldID); err != nil {
			return err
		}
		if err := contract.RequireFlag("option", optionID); err != nil {
			return err
		}
		if err := contract.RequireAtLeastOne(map[string]string{"text": text, "color": color}); err != nil {
			return err
		}

		params := trello.UpdateCustomFieldOptionParams{}
		if text != "" {
			params.Value = &trello.CustomFieldOptionValue{Text: text}
		}
		if color != "" {
			params.Color = &color
		}

		creds, err := auth.RequireAuth(getCredStore(), "default")
		if err != nil {
			return err
		}
		result, err := getAPIClient(creds).UpdateCustomFieldOption(cmd.Context(), fieldID, optionID, params)
		if err != nil {
			return err
		}
		return output(cmd.OutOrStdout(), result)
	},
}

var cfOptionsDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a custom field option",
	RunE: func(cmd *cobra.Command, args []string) error {
		fieldID, _ := cmd.Flags().GetString("field")
		optionID, _ := cmd.Flags().GetString("option")

		if err := contract.RequireFlag("field", fieldID); err != nil {
			return err
		}
		if err := contract.RequireFlag("option", optionID); err != nil {
			return err
		}

		creds, err := auth.RequireAuth(getCredStore(), "default")
		if err != nil {
			return err
		}
		if err := getAPIClient(creds).DeleteCustomFieldOption(cmd.Context(), fieldID, optionID); err != nil {
			return err
		}
		return output(cmd.OutOrStdout(), trello.DeleteResult{Deleted: true, ID: optionID})
	},
}

// ── items subgroup ─────────────────────────────────────────────────

var cfItemsCmd = &cobra.Command{
	Use:   "items",
	Short: "Manage custom field items on cards",
}

var cfItemsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List custom field items on a card",
	RunE: func(cmd *cobra.Command, args []string) error {
		cardID, _ := cmd.Flags().GetString("card")
		if err := contract.RequireFlag("card", cardID); err != nil {
			return err
		}
		creds, err := auth.RequireAuth(getCredStore(), "default")
		if err != nil {
			return err
		}
		result, err := getAPIClient(creds).ListCardCustomFieldItems(cmd.Context(), cardID)
		if err != nil {
			return err
		}
		return output(cmd.OutOrStdout(), result)
	},
}

var cfItemsSetCmd = &cobra.Command{
	Use:   "set",
	Short: "Set a custom field value on a card",
	RunE: func(cmd *cobra.Command, args []string) error {
		cardID, _ := cmd.Flags().GetString("card")
		fieldID, _ := cmd.Flags().GetString("field")
		text, _ := cmd.Flags().GetString("text")
		number, _ := cmd.Flags().GetString("number")
		date, _ := cmd.Flags().GetString("date")
		checked, _ := cmd.Flags().GetString("checked")
		option, _ := cmd.Flags().GetString("option")

		if err := contract.RequireFlag("card", cardID); err != nil {
			return err
		}
		if err := contract.RequireFlag("field", fieldID); err != nil {
			return err
		}

		valueFlags := map[string]string{
			"text":    text,
			"number":  number,
			"date":    date,
			"checked": checked,
			"option":  option,
		}
		if err := contract.RequireExactlyOne(valueFlags); err != nil {
			return err
		}

		if date != "" {
			if err := contract.ValidateISO8601Optional(date); err != nil {
				return err
			}
		}

		params := trello.SetCardCustomFieldItemParams{}
		switch {
		case text != "":
			params.Value = trello.CardCustomFieldItemValue{Text: text}
		case number != "":
			params.Value = trello.CardCustomFieldItemValue{Number: number}
		case date != "":
			params.Value = trello.CardCustomFieldItemValue{Date: date}
		case checked != "":
			params.Value = trello.CardCustomFieldItemValue{Checked: checked}
		case option != "":
			params.IDValue = option
		}

		creds, err := auth.RequireAuth(getCredStore(), "default")
		if err != nil {
			return err
		}
		result, err := getAPIClient(creds).SetCardCustomFieldItem(cmd.Context(), cardID, fieldID, params)
		if err != nil {
			return err
		}
		return output(cmd.OutOrStdout(), result)
	},
}

var cfItemsClearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear a custom field value on a card",
	RunE: func(cmd *cobra.Command, args []string) error {
		cardID, _ := cmd.Flags().GetString("card")
		fieldID, _ := cmd.Flags().GetString("field")

		if err := contract.RequireFlag("card", cardID); err != nil {
			return err
		}
		if err := contract.RequireFlag("field", fieldID); err != nil {
			return err
		}

		creds, err := auth.RequireAuth(getCredStore(), "default")
		if err != nil {
			return err
		}
		if err := getAPIClient(creds).ClearCardCustomFieldItem(cmd.Context(), cardID, fieldID); err != nil {
			return err
		}
		return output(cmd.OutOrStdout(), trello.DeleteResult{Deleted: true, ID: fieldID})
	},
}

// ── init: wire flags and command tree ──────────────────────────────

func init() {
	// definition subcommands flags
	cfListCmd.Flags().String("board", "", "Board ID")
	cfGetCmd.Flags().String("field", "", "Custom field ID")

	cfCreateCmd.Flags().String("board", "", "Board ID")
	cfCreateCmd.Flags().String("name", "", "Field name")
	cfCreateCmd.Flags().String("type", "", "Field type (text, number, date, checkbox, list)")
	cfCreateCmd.Flags().Bool("card-front", false, "Show on card front")
	cfCreateCmd.Flags().StringArray("option", nil, "Option value (repeatable, for list type)")

	cfUpdateCmd.Flags().String("field", "", "Custom field ID")
	cfUpdateCmd.Flags().String("name", "", "New field name")
	cfUpdateCmd.Flags().Bool("card-front", false, "Show on card front")

	cfDeleteCmd.Flags().String("field", "", "Custom field ID")

	// options subcommands flags
	cfOptionsListCmd.Flags().String("field", "", "Custom field ID")

	cfOptionsAddCmd.Flags().String("field", "", "Custom field ID")
	cfOptionsAddCmd.Flags().String("text", "", "Option text")
	cfOptionsAddCmd.Flags().String("color", "", "Option color")

	cfOptionsUpdateCmd.Flags().String("field", "", "Custom field ID")
	cfOptionsUpdateCmd.Flags().String("option", "", "Option ID")
	cfOptionsUpdateCmd.Flags().String("text", "", "New option text")
	cfOptionsUpdateCmd.Flags().String("color", "", "New option color")

	cfOptionsDeleteCmd.Flags().String("field", "", "Custom field ID")
	cfOptionsDeleteCmd.Flags().String("option", "", "Option ID")

	// items subcommands flags
	cfItemsListCmd.Flags().String("card", "", "Card ID")

	cfItemsSetCmd.Flags().String("card", "", "Card ID")
	cfItemsSetCmd.Flags().String("field", "", "Custom field ID")
	cfItemsSetCmd.Flags().String("text", "", "Text value")
	cfItemsSetCmd.Flags().String("number", "", "Number value")
	cfItemsSetCmd.Flags().String("date", "", "Date value (ISO 8601)")
	cfItemsSetCmd.Flags().String("checked", "", "Checked value (true/false)")
	cfItemsSetCmd.Flags().String("option", "", "Option ID")

	cfItemsClearCmd.Flags().String("card", "", "Card ID")
	cfItemsClearCmd.Flags().String("field", "", "Custom field ID")

	// wire command tree
	customFieldsCmd.AddCommand(cfListCmd)
	customFieldsCmd.AddCommand(cfGetCmd)
	customFieldsCmd.AddCommand(cfCreateCmd)
	customFieldsCmd.AddCommand(cfUpdateCmd)
	customFieldsCmd.AddCommand(cfDeleteCmd)

	cfOptionsCmd.AddCommand(cfOptionsListCmd)
	cfOptionsCmd.AddCommand(cfOptionsAddCmd)
	cfOptionsCmd.AddCommand(cfOptionsUpdateCmd)
	cfOptionsCmd.AddCommand(cfOptionsDeleteCmd)
	customFieldsCmd.AddCommand(cfOptionsCmd)

	cfItemsCmd.AddCommand(cfItemsListCmd)
	cfItemsCmd.AddCommand(cfItemsSetCmd)
	cfItemsCmd.AddCommand(cfItemsClearCmd)
	customFieldsCmd.AddCommand(cfItemsCmd)

	rootCmd.AddCommand(customFieldsCmd)
}
