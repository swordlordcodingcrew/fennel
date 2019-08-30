package main
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
	"fmt"
	"os"
	"swordlord.com/fennelcore"
	"swordlord.com/fennelcli/cmd"
)

func main() {

	// Initialise env and params
	fennelcore.InitConfig()
	fennelcore.InitLog()

	// Initialise database
	// if there is an error, this function will quit the app
	fennelcore.InitDatabase()
	defer fennelcore.CloseDB()

	// initialise the command structure
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)

		// yes, we deferred closing of the db, but that only works when ending with dignity
		fennelcore.CloseDB()
		os.Exit(1)
	}
}