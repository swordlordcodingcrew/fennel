package principal
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
	"fmt"
	"github.com/beevik/etree"
)

func Propfind(w http.ResponseWriter, req *http.Request){

	//vars := mux.Vars(req)
	//sUser := vars["user"]

	dRet, propstat := handler.GetMultistatusDoc(req.RequestURI)

	doc := etree.NewDocument()
	size, err := doc.ReadFrom(req.Body)
	if err != nil || size == 0 {

		fmt.Printf("Error reading XML Body. Error: %s Size: %v", err, size)

		handler.SendMultiStatus(w, http.StatusNotFound, dRet, propstat)
		return
	}

	// find query parameters and store in props
	// could probably be faster with compiled path...
	// propfindPath := etree.MustCompilePath("/propfind/prop/*")
	propsQuery := doc.FindElements("/propfind/prop/*")

	// create new element to store response in
	prop := propstat.CreateElement("prop")
	prop.Space = "d"

	// let helper function fill prop element with requested props
	fillPropfindResponse(prop, propsQuery)

	// add status based on query
	status := propstat.CreateElement("status")
	status.Space = "d"

	if len(prop.ChildElements()) > 0 {

		status.SetText("HTTP/1.1 200 OK")

	} else {

		status.SetText("HTTP/1.1 404 Not Found")
	}

	// send response to client
	handler.SendETreeDocument(w, http.StatusMultiStatus, dRet)
}

// TODO: handle as expected, this is a cheap workaround
func Proppatch(w http.ResponseWriter, req *http.Request){

	dRet, propstat := handler.GetMultistatusDoc(req.RequestURI)

	// create new element to store response in
	prop := propstat.CreateElement("prop")
	prop.Space = "d"

	davd := prop.CreateElement("default-alarm-vevent-date")
	davd.Space = "cal"

	// add status
	status := propstat.CreateElement("status")
	status.Space = "d"
	status.SetText("HTTP/1.1 403 Forbidden")

	handler.SendETreeDocument(w, http.StatusMultiStatus, dRet)
}

func Report(w http.ResponseWriter, req *http.Request){

	handler.RespondWithMessage(w, http.StatusOK, "Report not implemented yet")

}


func Options(w http.ResponseWriter, req *http.Request){

	handler.RespondWithStandardOptions(w, http.StatusOK, "")
}

func fillPropfindResponse(node *etree.Element, props []*etree.Element) {

	// TODO
	token := ""

	for _, e := range props {

		//fmt.Println(e.Tag)
		name := e.Tag
		switch(name) {

		case "checksum-versions":
			//response += "";

		case "sync-token":
			prop := node.CreateElement("sync-token")
			prop.Space = "d"
			prop.SetText("https://swordlord.com/ns/sync/" + token)

		case "supported-report-set":
			//response += getSupportedReportSet(comm);

		case "principal-URL":
			//response += "<d:principal-URL><d:href>/p/" + comm.getUser().getUserName() + "/</d:href></d:principal-URL>\r\n";

		case "displayname":
			//response += "<d:displayname>" + comm.getUser().getUserName() + "</d:displayname>";

		case "principal-collection-set":
			//response += "<d:principal-collection-set><d:href>/p/</d:href></d:principal-collection-set>";

		case "current-user-principal":
			//response += "<d:current-user-principal><d:href>/p/" + comm.getUser().getUserName() + "/</d:href></d:current-user-principal>";

		case "calendar-home-set":
			//response += "<cal:calendar-home-set><d:href>/cal/" + comm.getUser().getUserName() + "</d:href></cal:calendar-home-set>";

		case "schedule-outbox-URL":
			//response += "<cal:schedule-outbox-URL><d:href>/cal/" + comm.getUser().getUserName() + "/outbox</d:href></cal:schedule-outbox-URL>";

		case "calendar-user-address-set":
			//response += getCalendarUserAddressSet(comm);

		case "notification-URL":
			//response += "<cs:notification-URL><d:href>/cal/" + comm.getUser().getUserName() + "/notifications/</d:href></cs:notification-URL>";

		case "getcontenttype":
			//response += "";

		case "addressbook-home-set":
			//response += "<card:addressbook-home-set><d:href>/card/" + comm.getUser().getUserName() + "/</d:href></card:addressbook-home-set>";

		case "directory-gateway":
			//response += "";

		case "email-address-set":
			//response += "<cs:email-address-set><cs:email-address>lord test at swordlord.com</cs:email-address></cs:email-address-set>";

		case "resource-id":

		default:
			if name != "text" {
				fmt.Println("CAL-PF: not handled: " + name)
			}
		}
	}
}