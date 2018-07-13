package cmd

import (
	"github.com/spf13/cobra"
	"swordlord.com/fennelcore/db/tablemodule"
	"fmt"
)

// domainCmd represents the domain command
var userCmd = &cobra.Command{
	Use:   "user",
	Short: "Add, change and manage your users.",
	Long: `Add, change and manage your users. Requires a subcommand.`,
	RunE: nil,
}

var userListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all users.",
	Long: `List all users.`,
	RunE: ListUser,
}

var userAddCmd = &cobra.Command{
	Use:   "add [username] [password]",
	Short: "Add new user to this instance of Wombag.",
	Long: `Add new user to this instance of Wombag.`,
	Args: cobra.ExactArgs(2),
	RunE: AddUser,
}

var userUpdateCmd = &cobra.Command{
	Use:   "update [userid] [password]",
	Short: "Update the password of the user.",
	Long: `Update the password of the user.`,
	Args: cobra.ExactArgs(2),
	RunE: UpdateUser,
}

var userDeleteCmd = &cobra.Command{
	Use:   "delete [userid]",
	Short: "Deletes a user and all of her devices.",
	Long: `Deletes a user and all of his or her devices.`,
	Args: cobra.ExactArgs(1),
	RunE: DeleteUser,
}

func ListUser(cmd *cobra.Command, args []string) error {

	tablemodule.ListUser()

	return nil
}

func AddUser(cmd *cobra.Command, args []string) error {

	if len(args) < 2 {
		return fmt.Errorf("command 'add' needs a user name and a password")
	}

	tablemodule.AddUser(args[0], args[1])

	return nil
}

func UpdateUser(cmd *cobra.Command, args []string) error {

	if len(args) < 2 {
		return fmt.Errorf("command 'update' needs a user identification and a new password")
	}

	tablemodule.UpdateUser(args[0], args[1])

	return nil
}

func DeleteUser(cmd *cobra.Command, args []string) error {

	if len(args) < 1 {
		return fmt.Errorf("command 'delete' needs a user identification")
	}

	tablemodule.DeleteUser(args[0])

	return nil
}

func init() {
	RootCmd.AddCommand(userCmd)

	userCmd.AddCommand(userListCmd)
	userCmd.AddCommand(userAddCmd)
	userCmd.AddCommand(userUpdateCmd)
	userCmd.AddCommand(userDeleteCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// domainCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// domainCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
