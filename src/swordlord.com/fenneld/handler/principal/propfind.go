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

	sUser, ok := req.Context().Value("auth_user").(string)
	if !ok {
		// TODO fail when there is no user, since this can't really happen!
		sUser = ""
	}

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
	fillPropfindResponse(prop, propsQuery, sUser)

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

func fillPropfindResponse(node *etree.Element, props []*etree.Element, sUser string) {

	// TODO
	token := ""

	for _, e := range props {

		//fmt.Println(e.Tag)
		name := e.Tag
		switch(name) {

		case "checksum-versions":
			//";

		case "sync-token":
			prop := node.CreateElement("sync-token")
			prop.Space = "d"
			prop.SetText("https://swordlord.com/ns/sync/" + token)

		case "supported-report-set":
			fillSupportedReportSet(node)

		case "principal-URL":
			//<d:principal-URL><d:href>/p/" + comm.getUser().getUserName() + "/</d:href></d:principal-URL>\r\n";
			prop := node.CreateElement("principal-URL")
			prop.Space = "d"

			href := prop.CreateElement("href")
			href.Space = "d"

			href.SetText("/p/" + sUser + "/")

		case "displayname":
			//<d:displayname>" + comm.getUser().getUserName() + "</d:displayname>";
			prop := node.CreateElement("displayname")
			prop.Space = "d"
			prop.SetText(sUser)

		case "principal-collection-set":
			//<d:principal-collection-set><d:href>/p/</d:href></d:principal-collection-set>";
			prop := node.CreateElement("principal-collection-set")
			prop.Space = "d"

			href := prop.CreateElement("href")
			href.Space = "d"

			href.SetText("/p/")

		case "current-user-principal":
			//<d:current-user-principal><d:href>/p/" + comm.getUser().getUserName() + "/</d:href></d:current-user-principal>";
			prop := node.CreateElement("current-user-principal")
			prop.Space = "d"

			href := prop.CreateElement("href")
			href.Space = "d"

			href.SetText("/p/" + sUser + "/")

		case "calendar-home-set":
			//<cal:calendar-home-set><d:href>/cal/" + comm.getUser().getUserName() + "</d:href></cal:calendar-home-set>";
			prop := node.CreateElement("calendar-home-set")
			prop.Space = "cal"

			href := prop.CreateElement("href")
			href.Space = "d"

			href.SetText("/cal/" + sUser + "/")

		case "schedule-outbox-URL":
			//<cal:schedule-outbox-URL><d:href>/cal/" + comm.getUser().getUserName() + "/outbox</d:href></cal:schedule-outbox-URL>";
			prop := node.CreateElement("schedule-outbox-URL")
			prop.Space = "cal"

			href := prop.CreateElement("href")
			href.Space = "d"
			href.SetText("/cal/" + sUser + "/outbox/")

		case "calendar-user-address-set":
			prop := node.CreateElement("calendar-user-address-set")
			prop.Space = "cal"

			href := prop.CreateElement("href")
			href.Space = "d"
			href.SetText("mailto:lord test at swordlord.com")

			href2 := prop.CreateElement("href")
			href2.Space = "d"
			href2.SetText("/p/" + sUser + "/")

		case "notification-URL":
			//<cs:notification-URL><d:href>/cal/" + comm.getUser().getUserName() + "/notifications/</d:href></cs:notification-URL>";
			prop := node.CreateElement("notification-URL")
			prop.Space = "cs"

			href := prop.CreateElement("href")
			href.Space = "d"

			href.SetText("/cal/" + sUser + "/notifications/")

		case "getcontenttype":
			//";

		case "addressbook-home-set":
			//<card:addressbook-home-set><d:href>/card/" + comm.getUser().getUserName() + "/</d:href></card:addressbook-home-set>";
			prop := node.CreateElement("addressbook-home-set")
			prop.Space = "card"

			href := prop.CreateElement("href")
			href.Space = "d"

			href.SetText("/card/" + sUser + "/")

		case "directory-gateway":
			//";

		case "email-address-set":
			//<cs:email-address-set><cs:email-address>lord test at swordlord.com</cs:email-address></cs:email-address-set>";
			prop := node.CreateElement("email-address-set")
			prop.Space = "cs"

			ea := prop.CreateElement("email-address")
			ea.Space = "cs"

			// todo load user email from db
			ea.SetText("lord test at swordlord.com")

		case "resource-id":

		default:
			if name != "text" {
				fmt.Println("CAL-PF: not handled: " + name)
			}
		}
	}
}


func fillSupportedReportSet(node *etree.Element) {

	/*
        <d:supported-report-set>\r\n";
        	<d:supported-report>\r\n";
        		<d:report>\r\n";
        			<d:expand-property/>\r\n";
        		</d:report>\r\n";
        	</d:supported-report>\r\n";
        	<d:supported-report>\r\n";
        		<d:report>\r\n";
        			<d:principal-property-search/>\r\n";
        		</d:report>\r\n";
        	</d:supported-report>\r\n";
        	<d:supported-report>\r\n";
        		<d:report>\r\n";
        			<d:principal-search-property-set/>\r\n";
        		</d:report>\r\n";
        	</d:supported-report>\r\n";
        </d:supported-report-set>\r\n";
*/
	srs := node.CreateElement("supported-report-set")
	srs.Space = "d"

	// ---
	sr1 := srs.CreateElement("supported-report")
	sr1.Space = "d"

	r1 := sr1.CreateElement("report")
	r1.Space = "d"

	ep := r1.CreateElement("expand-property")
	ep.Space = "d"

	// ---
	sr2 := srs.CreateElement("supported-report")
	sr2.Space = "d"

	r2 := sr2.CreateElement("report")
	r2.Space = "d"

	pps := r2.CreateElement("expand-property")
	pps.Space = "d"

	// ---
	sr3 := srs.CreateElement("supported-report")
	sr3.Space = "d"

	r3 := sr3.CreateElement("report")
	r3.Space = "d"

	psps := r3.CreateElement("principal-search-property-set")
	psps.Space = "d"
	
}