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

// calCmd represents the domain command
var calCmd = &cobra.Command{
	Use:   "cal",
	Short: "Add, change and manage calendars.",
	Long: `Add, change and manage calendars. Requires a subcommand.`,
	RunE: nil,
}

var calListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all calendars.",
	Long: `List all calendars.`,
	RunE: ListCal,
}

var calAddCmd = &cobra.Command{
	Use:   "add [username] [calendar]",
	Short: "Add new  calendar for given user.",
	Long: `Add new user to this instance of Wombag.`,
	Args: cobra.ExactArgs(2),
	RunE: AddCal,
}

var calDeleteCmd = &cobra.Command{
	Use:   "delete [calendar]",
	Short: "Deletes a calendar for given user.",
	Long: `Deletes a calendar for given user.`,
	Args: cobra.ExactArgs(1),
	RunE: DeleteCal,
}

func ListCal(cmd *cobra.Command, args []string) error {

	tablemodule.ListCal()

	return nil
}

func AddCal(cmd *cobra.Command, args []string) error {

	if len(args) < 2 {
		return fmt.Errorf("command 'add' needs a user name and a password")
	}

	// TODO
	//tablemodule.AddCal(args[0], args[1])

	return nil
}

func DeleteCal(cmd *cobra.Command, args []string) error {

	if len(args) < 1 {
		return fmt.Errorf("command 'delete' needs a user identification")
	}

	tablemodule.DeleteCal(args[0])

	return nil
}

func init() {

	// TODO reactivate once its running
	//RootCmd.AddCommand(calCmd)

	calCmd.AddCommand(calListCmd)
	calCmd.AddCommand(calAddCmd)
	calCmd.AddCommand(calDeleteCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// domainCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// domainCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
