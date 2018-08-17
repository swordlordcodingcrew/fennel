package cmd
/*-----------------------------------------------------------------------------
 **
 ** - Fennel -
 **
 ** your lightweight CalDAV and CardDAV server
 **
 ** Copyright 2018 by SwordLord - the coding crew - http://www.swordlord.com
 ** and contributing authors
 **
 ** This program is free software; you can redistribute it and/or modify it
 ** under the terms of the GNU Affero General Public License as published by the
 ** Free Software Foundation, either version 3 of the License, or (at your option)
 ** any later version.
 **
 ** This program is distributed in the hope that it will be useful, but WITHOUT
 ** ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or
 ** FITNESS FOR A PARTICULAR PURPOSE.  See the GNU Affero General Public License
 ** for more details.
 **
 ** You should have received a copy of the GNU Affero General Public License
 ** along with this program. If not, see <http://www.gnu.org/licenses/>.
 **
 **-----------------------------------------------------------------------------
 **
 ** Original Authors:
 ** LordEidi@swordlord.com
 ** LordCelery@swordlord.com
 **
-----------------------------------------------------------------------------*/

import (
	"github.com/spf13/cobra"
	"swordlord.com/fennelcore/db/tablemodule"
	"fmt"
)

// userCmd represents the domain command
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
	Short: "Add new user to this instance of Fennel.",
	Long: `Add new user to this instance of Fennel.`,
	Args: cobra.ExactArgs(2),
	RunE: AddUser,
}

var userUpdateCmd = &cobra.Command{
	Use:   "update [userid] [password] [comment]",
	Short: "Update the password and comment of the user.",
	Long: `Update the password of the user. Comment field can be left empty`,
	Args: cobra.MinimumNArgs(2),
	RunE: UpdateUser,
}

var userVerifyCmd = &cobra.Command{
	Use:   "verify [userid] [pwd]",
	Short: "Verifies the password of given user.",
	Long: `Verifies the password of the given user. Can be used to check if memorised password is correct.`,
	Args: cobra.ExactArgs(2),
	RunE: VerifyUser,
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

func VerifyUser(cmd *cobra.Command, args []string) error {

	if len(args) < 2 {
		return fmt.Errorf("command 'verify' needs a user name and a password")
	}

	isValid, _ := tablemodule.ValidateUserInDB(args[0], args[1])

	if isValid {

		fmt.Println("User was verified")
	} else {

		fmt.Println("User could not be verified")
	}

	return nil
}

func UpdateUser(cmd *cobra.Command, args []string) error {

	argCount := len(args)

	if argCount < 2 {
		return fmt.Errorf("command 'update' needs a user identification, a new password and optionally a comment")
	}

	comment := ""

	if argCount > 2 {

		comment = args[2]
	}

	tablemodule.UpdateUser(args[0], args[1], comment)

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
	userCmd.AddCommand(userVerifyCmd)
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
