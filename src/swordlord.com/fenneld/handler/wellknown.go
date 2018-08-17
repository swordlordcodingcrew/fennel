package handler
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
	"net/http"
	"github.com/gorilla/mux"
)


func OnWellKnown(w http.ResponseWriter, req *http.Request){

	vars := mux.Vars(req)
	sParam := vars["param"]

	// there are different redirect targets depending on the .well-known url
	// the user wants
	// Never seen most of these, though
	switch sParam {

		case "caldav":
			RespondWithRedirect(w, req, "/cal/")
		case "carddav":
			RespondWithRedirect(w, req, "/card/")
		default:
			RespondWithRedirect(w, req, "/p/")
	}
}

func OnWellKnownNoParam(w http.ResponseWriter, req *http.Request){

	RespondWithRedirect(w, req, "/p/")
}
