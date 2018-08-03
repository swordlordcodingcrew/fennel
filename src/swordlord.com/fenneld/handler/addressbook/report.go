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
	"net/http"
	"swordlord.com/fenneld/handler"
	"github.com/gorilla/mux"
	"github.com/beevik/etree"
	"fmt"
	"strings"
	"swordlord.com/fennelcore/db/tablemodule"
	"swordlord.com/fennelcore/db/model"
)

func Report(w http.ResponseWriter, req *http.Request){

	// REPORT https://dav.fruux.com/addressbooks/user/addressbookid/

	/*
	<?xml version="1.0" encoding="UTF-8"?>
<A:sync-collection xmlns:A="DAV:">
  <A:sync-token>http://sabre.io/ns/sync/18</A:sync-token>
  <A:sync-level>1</A:sync-level>
  <A:prop>
    <A:getetag/>
  </A:prop>
</A:sync-collection>

	<?xml version="1.0"?>
<d:multistatus xmlns:d="DAV:" xmlns:s="http://sabredav.org/ns" xmlns:fx="http://fruux.com/ns"
xmlns:cal="urn:ietf:params:xml:ns:caldav" xmlns:cs="http://calendarserver.org/ns/" xmlns:card="urn:ietf:params:xml:ns:carddav">
  <d:sync-token>http://sabre.io/ns/sync/18</d:sync-token>
</d:multistatus>
	*/

	vars := mux.Vars(req)
	sUser := vars["user"]
	sAB := vars["addressbook"]

	doc := etree.NewDocument()
	size, err := doc.ReadFrom(req.Body)
	if err != nil || size == 0 {

		fmt.Printf("Error reading XML Body. Error: %s Size: %v", err, size)

		handler.RespondWithMessage(w, http.StatusInternalServerError, "")
		return
	}

	root := doc.Root()
	name := root.Tag

	switch name {

		case "sync-collection":
			handleReportSyncCollection(w, req.RequestURI, root, sUser, sAB)

		case "addressbook-multiget":
			handleReportAddressbookMultiget(w, req.RequestURI, root, sUser, sAB)

		default:
			if name != "text" {
				fmt.Println("CARD-Report: root not handled: " + name)
				handler.RespondWithMessage(w, http.StatusInternalServerError, "")
			}
	}
}

func handleReportSyncCollection(w http.ResponseWriter, uri string, nodeQuery *etree.Element, sUser string, sAB string) {

	dRet, ms := handler.GetMultistatusDocWOResponseTag()

	println(ms)

	// TODO check filter:
	// <A:sync-token>http://sabre.io/ns/sync/18</A:sync-token>
	// <A:sync-level>1</A:sync-level>
	syncTokenQ := nodeQuery.FindElement("./sync-token/")
	if syncTokenQ != nil {

	}
	syncLevelQ := nodeQuery.FindElement("./sync-level/")
	if syncLevelQ != nil {

	}

	// TODO return etoken

	handler.SendETreeDocument(w, http.StatusMultiStatus, dRet)
}

func handleReportAddressbookMultiget(w http.ResponseWriter, uri string, nodeQuery *etree.Element, sUser string, sAB string) {

	dRet, ms := handler.GetMultistatusDocWOResponseTag()

	println(ms)

	// TODO check filter:
	// payload += "<A:prop xmlns:A=\"DAV:\">\n\r";
	// payload += "<A:getetag/>\n\r";
	// payload += "<D:address-data/>\n\r";
	// payload += "</A:prop>\n\r";
	eTagElement := nodeQuery.FindElement("./prop/getetag/")
	getETag := eTagElement != nil

	addressDataElement := nodeQuery.FindElement("./prop/address-data/")
	getAddressData := addressDataElement != nil

	reqDocs := nodeQuery.FindElements("./href/")

	arrVCards := make([]string, len(reqDocs))

	for i, reqDoc := range reqDocs {

		// get the last element, which should contain the filename
		arrPath := strings.Split(reqDoc.Text(), "/")
		pathCount := len(arrPath)

		if pathCount < 2 {
			// TODO continue is suboptimal, we should add a bool to the array and send a 404 for said file.
			continue
		}

		filename := arrPath[pathCount - 1]

		arrFile := strings.Split(filename, ".")
		if len(arrFile) < 2 {
			// TODO continue is suboptimal, we should add a bool to the array and send a 404 for said file.
			continue
		}

		arrVCards[i] = arrFile[0]
	}

	err, rows := tablemodule.FindVCardsFromAddressbook(sAB, arrVCards)
	if err != nil {

		handler.RespondWithMessage(w, http.StatusInternalServerError, "")
	}

	for _, row := range rows {

		if getAddressData {
			handleReportVCardReply(ms, row, uri, getETag, getAddressData)
		}
	}

	handler.SendETreeDocument(w, http.StatusMultiStatus, dRet)
}

func handleReportVCardReply(ms *etree.Element, vcard *model.VCARD, uri string, getETag bool, getAddressData bool){

	ps := handler.AddResponseToMultistat(ms, uri + vcard.Pkey + ".vcf")

	if getETag {

		prop := ps.CreateElement("getetag")
		prop.Space = "d"
		// TODO fill correct text
		prop.SetText("E1")
	}

	if getAddressData {

		prop := ps.CreateElement("address-data")
		prop.Space = "card"
		// TODO correct format

		// content = content.replace(/&/g,'&amp;');
		// content = content.replace(/\r\n|\r|\n/g,'&#13;\r\n');

		prop.SetText(vcard.Content)
	}

	handler.AddStatusToPropstat(http.StatusOK, ps)
}