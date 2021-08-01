package addressbook

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
	"github.com/gorilla/mux"
	"github.com/swordlordcodingcrew/fennel/fennelcore/db/tablemodule"
	"github.com/swordlordcodingcrew/fennel/fenneld/handler"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func Proppatch(w http.ResponseWriter, req *http.Request) {

	handler.RespondWithMessage(w, http.StatusOK, "Not implemented yet")
}

func Options(w http.ResponseWriter, req *http.Request) {

	handler.RespondWithStandardOptions(w, http.StatusOK, "")
}

func Put(w http.ResponseWriter, req *http.Request) {

	vars := mux.Vars(req)
	sAB := vars["addressbook"]
	sCard := vars["card"]

	sUser, ok := req.Context().Value("auth_user").(string)
	if !ok {

		fmt.Println("err: could not find auth_user in context")
		handler.RespondWithMessage(w, http.StatusPreconditionFailed, "No user supplied")
		return
	}

	bodyBuffer, _ := ioutil.ReadAll(req.Body)

	isGroup := strings.Contains(string(bodyBuffer), "X-ADDRESSBOOKSERVER-KIND:group")

	vcard, err := tablemodule.AddVCard(sCard, sUser, sAB, isGroup, string(bodyBuffer))
	if err != nil {

		fmt.Println("err: could not add vcard " + sCard)
		handler.RespondWithMessage(w, http.StatusPreconditionFailed, err.Error())
		return
	}

	handler.RespondWithMessage(w, http.StatusCreated, "VCARD added: "+vcard.Pkey)
}

func Get(w http.ResponseWriter, req *http.Request) {

	vars := mux.Vars(req)
	sCard := vars["card"]

	vcard, err := tablemodule.GetVCard(sCard)

	if err != nil {

		fmt.Println("err: could not find vcard " + sCard)

		handler.RespondWithMessage(w, http.StatusInternalServerError, err.Error())
		return
	}

	handler.RespondWithVCARD(w, http.StatusOK, vcard.Content)
}

func Delete(w http.ResponseWriter, req *http.Request) {

	vars := mux.Vars(req)
	sCard := vars["card"]

	err := tablemodule.DeleteVCard(sCard)

	if err != nil {
		log.Printf("Error with deleting VCard %q: %s\n", sCard)

		handler.RespondWithMessage(w, http.StatusInternalServerError, err.Error())

		return
	}

	handler.RespondWithMessage(w, http.StatusOK, "Deleted")
}

func Move(w http.ResponseWriter, req *http.Request) {

	handler.RespondWithMessage(w, http.StatusOK, "Not implemented yet")
}
