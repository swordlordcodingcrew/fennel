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
			"swordlord.com/fennelcore/db/tablemodule"
	"fmt"
		"github.com/beevik/etree"
	"swordlord.com/fennelcore/db/model"
	"strconv"
)

// TODO check if on root, if yes, answer differently
func PropfindRoot(w http.ResponseWriter, req *http.Request) {

	Propfind(w, req)
}

func Propfind(w http.ResponseWriter, req *http.Request){

	sUser, ok := req.Context().Value("auth_user").(string)
	if !ok {
		// TODO fail when there is no user, since this can't really happen!
		sUser = ""
	}

	dRet, ms := handler.GetMultistatusDocWOResponseTag()

	// TODO check if user exists?

	doc := etree.NewDocument()
	size, err := doc.ReadFrom(req.Body)
	if err != nil || size == 0 {

		fmt.Printf("Error reading XML Body. Error: %s Size: %v", err, size)

		handler.RespondWithMessage(w, http.StatusInternalServerError, "")
		return
	}

	// find query parameters and store in props
	// could probably be faster with compiled path...
	// propfindPath := etree.MustCompilePath("/propfind/prop/*")
	propsQuery := doc.FindElements("/propfind/prop/*")

	// get the propstat for the root response
	psRoot := handler.AddResponseWStatusToMultistat(ms, "/cal/" + sUser + "/", http.StatusOK)

	// let helper function fill prop element with requested props from the root
	fillPropfindResponseOnAddressbookRoot(psRoot, sUser, propsQuery)

	err, rowsADB := tablemodule.GetAddressbooksFromUser(sUser)
	if err != nil {

		fmt.Printf("Error getting ADB from User: %s", err)

		handler.RespondWithMessage(w, http.StatusInternalServerError, "")
		return
	}

	// let helper function fill prop elements for every single sub element
	fillPropfindResponseOnEachAddressbook(ms, sUser, rowsADB, propsQuery)

	handler.SendETreeDocument(w, http.StatusMultiStatus, dRet)
}

func fillPropfindResponseOnAddressbookRoot(psRoot *etree.Element, sUser string, propsQuery []*etree.Element){

	token := ""

	for _, e := range propsQuery {

		//fmt.Println(e.Tag)
		name := e.Tag
		switch(name) {

			case "current-user-privilege-set":
				fillCurrentUserPrivilegeSetRoot(psRoot, sUser)

			case "owner":
				// 		<d:owner>
				//          <d:href>/p/a3298271331/</d:href>
				//        </d:owner>
				prop := psRoot.CreateElement("owner")
				prop.Space = "d"

				href := prop.CreateElement("href")
				href.Space = "d"

				href.SetText("/p/" + sUser + "/")

			case "resourcetype":
				//         <d:resourcetype>
				//          <d:collection/>
				//        </d:resourcetype>
				prop := psRoot.CreateElement("resourcetype")
				prop.Space = "d"

				col := prop.CreateElement("collection")
				col.Space = "d"

			case "supported-report-set":
				fillSupportedReportSetRoot(psRoot)

			case "sync-token":
			prop := psRoot.CreateElement("sync-token")
			prop.Space = "d"
			prop.SetText("https://swordlord.com/ns/sync/" + token)

		default:
			if name != "text" {
				fmt.Println("CARD_AddressbookRoot-PF: not handled: " + name)
			}
		}
	}
}

func fillPropfindResponseOnEachAddressbook(ms *etree.Element, sUser string, rowsADB []*model.ADB, propsQuery []*etree.Element) {

	for _, row := range rowsADB{

		psRoot := handler.AddResponseWStatusToMultistat(ms, "/cal/" + sUser + "/" + row.Pkey + "/", http.StatusOK)

		for _, e := range propsQuery {

			//fmt.Println(e.Tag)
			name := e.Tag
			switch (name) {

				case "current-user-privilege-set":
					fillCurrentUserPrivilegeSetADB(psRoot, sUser, row)

				case "displayname": //>Addressbook</d:displayname>
					//         <d:displayname>Addressbook</d:displayname>
					prop := psRoot.CreateElement("displayname")
					prop.Space = "d"

					prop.SetText(row.Name)

				case "max-resource-size": //>1048576</card:max-resource-size>
					//        <card:max-resource-size>1048576</card:max-resource-size>
					prop := psRoot.CreateElement("max-resource-size")
					prop.Space = "card"

					prop.SetText("1048576")

				case "owner":
					// 		<d:owner>
					//          <d:href>/p/a3298271331/</d:href>
					//        </d:owner>
					prop := psRoot.CreateElement("owner")
					prop.Space = "d"

					href := prop.CreateElement("href")
					href.Space = "d"

					href.SetText("/p/" + sUser + "/")

				case "resourcetype":
					//         <d:resourcetype>
					//          <d:collection/>
					// 			<card:addressbook/>
					//        </d:resourcetype>
					prop := psRoot.CreateElement("resourcetype")
					prop.Space = "d"

					col := prop.CreateElement("collection")
					col.Space = "d"

					addrb := prop.CreateElement("addressbook")
					addrb.Space = "card"

				case "supported-report-set":
					fillSupportedReportSetADB(psRoot)

				case "sync-token":
					prop := psRoot.CreateElement("sync-token")
					prop.Space = "d"
					prop.SetText("https://swordlord.com/ns/sync/" + strconv.Itoa(row.Synctoken))

			default:
				if name != "text" {
					fmt.Println("CARD-Addressbook-PF: not handled: " + name)
				}
			}
		}
	}
}

func fillCurrentUserPrivilegeSetRoot(ps *etree.Element, user string) {

	fillCurrentUserPrivilegeStandardSet(ps)
}

func fillCurrentUserPrivilegeSetADB(ps *etree.Element, user string, row *model.ADB) {

	fillCurrentUserPrivilegeStandardSet(ps)
}

func fillCurrentUserPrivilegeStandardSet(ps *etree.Element) {

	// <d:current-user-privilege-set>
	//          <d:privilege>
	//            <d:write/>
	//          </d:privilege>
	//          <d:privilege>
	//            <d:write-acl/>
	//          </d:privilege>
	//          <d:privilege>
	//            <d:write-properties/>
	//          </d:privilege>
	//          <d:privilege>
	//            <d:write-content/>
	//          </d:privilege>
	//          <d:privilege>
	//            <d:bind/>
	//          </d:privilege>
	//          <d:privilege>
	//            <d:unbind/>
	//          </d:privilege>
	//          <d:privilege>
	//            <d:unlock/>
	//          </d:privilege>
	//          <d:privilege>
	//            <d:read/>
	//          </d:privilege>
	//          <d:privilege>
	//               <d:read-acl/>
	//          </d:privilege>
	//          <d:privilege>
	//            <d:read-current-user-privilege-set/>
	//          </d:privilege>
	//        </d:current-user-privilege-set>

	cups := ps.CreateElement("current-user-privilege-set")
	cups.Space = "d"

	addPrivilege(cups, "d", "write")
	addPrivilege(cups, "d", "write-acl")
	addPrivilege(cups, "d", "write-properties")
	addPrivilege(cups, "d", "write-content")
	addPrivilege(cups, "d", "bind")
	addPrivilege(cups, "d", "unbind")
	addPrivilege(cups, "d", "unlock")
	addPrivilege(cups, "d", "read")
	addPrivilege(cups, "d", "read-acl")
	addPrivilege(cups, "d", "read-current-user-privilege-set")
}

func addPrivilege(cups *etree.Element, namespace string, elementName string){

	p := cups.CreateElement("privilege")
	p.Space = "d"

	e := p.CreateElement(elementName)
	e.Space = namespace
}

func fillSupportedReportSetRoot(ps *etree.Element) {

	//         <d:supported-report-set>
	//          <d:supported-report>
	//            <d:report>
	//              <d:expand-property/>
	//            </d:report>
	//          </d:supported-report>
	//          <d:supported-report>
	//            <d:report>
	//              <d:principal-property-search/>
	//            </d:report>
	//          </d:supported-report>
	//          <d:supported-report>
	//            <d:report>
	//              <d:principal-search-property-set/>
	//            </d:report>
	//          </d:supported-report>
	//        </d:supported-report-set>
	srs := ps.CreateElement("supported-report-set")
	srs.Space = "d"

	addSupportedReportElement(srs, "d", "expand-property")
	addSupportedReportElement(srs, "d", "principal-property-search")
	addSupportedReportElement(srs, "d", "principal-search-property-set")
}

func fillSupportedReportSetADB(ps *etree.Element) {
	//         <d:supported-report-set>
	//          <d:supported-report>
	//            <d:report>
	//              <d:expand-property/>
	//            </d:report>
	//          </d:supported-report>
	//          <d:supported-report>
	//            <d:report>
	//              <d:principal-property-search/>
	//            </d:report>
	//          </d:supported-report>
	//          <d:supported-report>
	//            <d:report>
	//              <d:principal-search-property-set/>
	//            </d:report>
	//          </d:supported-report>
	//          <d:supported-report>
	//            <d:report>
	//              <card:addressbook-multiget/>
	//            </d:report>
	//          </d:supported-report>
	//          <d:supported-report>
	//            <d:report>
	//              <card:addressbook-query/>
	//            </d:report>
	//          </d:supported-report>
	//          <d:supported-report>
	//            <d:report>
	//              <d:sync-collection/>
	//            </d:report>
	//          </d:supported-report>
	//        </d:supported-report-set>

	srs := ps.CreateElement("supported-report-set")
	srs.Space = "d"

	addSupportedReportElement(srs, "d", "expand-property")
	addSupportedReportElement(srs, "d", "principal-property-search")
	addSupportedReportElement(srs, "d", "principal-search-property-set")
	addSupportedReportElement(srs, "card", "addressbook-multiget")
	addSupportedReportElement(srs, "card", "addressbook-query")
	addSupportedReportElement(srs, "d", "sync-collection")
}

func addSupportedReportElement(srs *etree.Element, namespace string, elementName string)  {

	sr := srs.CreateElement("supported-report")
	sr.Space = "d"

	r := sr.CreateElement("report")
	r.Space = "d"

	e := r.CreateElement(elementName)
	e.Space = namespace
}