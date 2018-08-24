package calendar
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
	"swordlord.com/fennelcore/db/tablemodule"
	"github.com/beevik/etree"
	"fmt"
	"swordlord.com/fennelcore/db/model"
)

// TODO check if on root, if yes, answer differently
func PropfindRoot(w http.ResponseWriter, req *http.Request) {

	sUser, ok := req.Context().Value("auth_user").(string)
	if !ok {
		// TODO fail when there is no user, since this can't really happen!
		sUser = ""
	}

	dRet, propstat := handler.GetMultistatusDoc("/cal/" + sUser + "/")

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
	// TODO fix empty CAL (nil preferred)
	fillPropfindResponse(prop, sUser, model.CAL{}, propsQuery)

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

func PropfindUser(w http.ResponseWriter, req *http.Request){

	// TODO
	handler.RespondWithMessage(w, http.StatusMultiStatus, "Not implemented yet")
}

func PropfindInbox(w http.ResponseWriter, req *http.Request){

	dRet, propstat := handler.GetMultistatusDoc(req.RequestURI)

	handler.SendMultiStatus(w, http.StatusOK, dRet, propstat)
}

func PropfindOutbox(w http.ResponseWriter, req *http.Request){

	// TODO
	handler.RespondWithMessage(w, http.StatusMultiStatus, "Not implemented yet")
}

func PropfindNotification(w http.ResponseWriter, req *http.Request){

	dRet, propstat := handler.GetMultistatusDoc(req.RequestURI)

	handler.SendMultiStatus(w, http.StatusOK, dRet, propstat)
}

func PropfindCalendar(w http.ResponseWriter, req *http.Request){

	vars := mux.Vars(req)
	sCal := vars["calendar"]

	sUser, ok := req.Context().Value("auth_user").(string)
	if !ok {
		// TODO fail when there is no user, since this can't really happen!
		sUser = ""
	}

	dRet, propstat := handler.GetMultistatusDoc("/cal/" + sUser + "/" + sCal + "/")

	cal, err := tablemodule.GetCal(sCal)
	if err != nil {

		fmt.Println(err)

		handler.SendMultiStatus(w, http.StatusNotFound, dRet, propstat)
		return
	}

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
	fillPropfindResponse(prop, sUser, cal, propsQuery)

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

func fillPropfindResponse(node *etree.Element, user string, cal model.CAL, props []*etree.Element) {

	// TODO
	response := ""
	token := ""

	for _, e := range props {

		//fmt.Println(e.Tag)
		name := e.Tag
		switch(name) {

		//case "add-member":

		case "allowed-sharing-modes":
			response += "<cs:allowed-sharing-modes><cs:can-be-shared/><cs:can-be-published/></cs:allowed-sharing-modes>";

			//case "autoprovisioned":
			//case "bulk-requests":

		case "calendar-color":
			response += "<xical:calendar-color xmlns:xical=\"http://apple.com/ns/ical/\">" + cal.Colour + "</xical:calendar-color>";

			//case "calendar-description":
			//case "calendar-free-busy-set":
			//response += "<d:response><d:href>/</d:href></d:response>";

		case "calendar-order":
			response += "" // "<xical:calendar-order xmlns:xical=\"http://apple.com/ns/ical/\">" + cal.Order + "</xical:calendar-order>";

		case "calendar-timezone":
			var timezone = cal.Timezone;
			//timezone = timezone.replace(/\r\n|\r|\n/g,"&#13;\r\n");
			response += "<cal:calendar-timezone>" + timezone + "</cal:calendar-timezone>";

		case "current-user-privilege-set":
			response += "" //getCurrentUserPrivilegeSet();

		case "current-user-principal":
			response += ""
			// <d:current-user-principal><d:href>/p/" + username + "/</d:href></d:current-user-principal>

			//case "default-alarm-vevent-date":
			//case "default-alarm-vevent-datetime":

		case "displayname":
			response += "<d:displayname>" + cal.Displayname + "</d:displayname>"

			//case "language-code":
			//case "location-code":

		case "owner":
			response += "<d:owner><d:href>/p/" + user +"/</d:href></d:owner>"

		case "principal-collection-set":
			//"<d:principal-collection-set><d:href>/p/</d:href></d:principal-collection-set>"

			// TODO Fix URL
		case "pre-publish-url":
			response += "<cs:pre-publish-url><d:href>https://127.0.0.1/cal/" + user + "/" + cal.Pkey + "</d:href></cs:pre-publish-url>";

			//case "publish-url":
			//case "push-transports":
			//case "pushkey":
			//case "quota-available-bytes":
			//case "quota-used-bytes":
			//case "refreshrate":
			//case "resource-id":

		case "resourcetype":
			response += "<d:resourcetype><d:collection/><cal:calendar/></d:resourcetype>";

		case "schedule-calendar-transp":
			response += "<cal:schedule-calendar-transp><cal:opaque/></cal:schedule-calendar-transp>";

			//case "schedule-default-calendar-URL":
			//case "source":
			//case "subscribed-strip-alarms":
			//case "subscribed-strip-attachments":
			//case "subscribed-strip-todos":
			//case "supported-calendar-component-set":

		case "supported-calendar-component-sets":
			response += "<cal:supported-calendar-component-set><cal:comp name=\"VEVENT\"/></cal:supported-calendar-component-set>";

		case "supported-report-set":
			response += "" //getSupportedReportSet(false);

		case "getctag":
			prop := node.CreateElement("getctag")
			prop.Space = "cs"
			prop.SetText("https://swordlord.com/ns/sync/" + token)

			//case "getetag":
			// no response?

			//case "checksum-versions":
			// no response?

		case "sync-token":
			prop := node.CreateElement("sync-token")
			prop.Space = "d"
			prop.SetText("https://swordlord.com/ns/sync/" + token)

		case "acl":
			response += "" // getACL(comm)

			//case "getcontenttype":
			//response += "<d:getcontenttype>text/calendar;charset=utf-8</d:getcontenttype>";

		default:
			if name != "text" {
				fmt.Println("CAL-PF: not handled: " + name)
			}
		}
	}
}

/*
Version     string   `xml:"version,attr"`

"<B:mkcalendar xmlns:B=\"urn:ietf:params:xml:ns:caldav\">\n\r";
payload += "<A:set xmlns:A=\"DAV:\">\n\r";
payload += "<A:prop>\n\r";
payload += "<B:supported-calendar-component-set>\n\r";
payload += "    <B:comp name=\"VEVENT\"/>\n\r";
payload += "</B:supported-calendar-component-set>\n\r";
payload += "<A:displayname>Three</A:displayname>\n\r";
payload += "<D:calendar-order xmlns:D=\"http://apple.com/ns/ical/\">4</D:calendar-order>\n\r";
payload += "<B:schedule-calendar-transp>\n\r";
payload += "    <B:transparent/>\n\r";
payload += "</B:schedule-calendar-transp>\n\r";
payload += "<B:calendar-timezone>BEGIN:VCALENDAR&#13;\n\r";
payload += "VERSION:2.0&#13;\n\r";
payload += "CALSCALE:GREGORIAN&#13;\n\r";
payload += "BEGIN:VTIMEZONE&#13;\n\r";
payload += "TZID:Europe/Zurich&#13;\n\r";
payload += "BEGIN:DAYLIGHT&#13;\n\r";
payload += "TZOFFSETFROM:+0100&#13;\n\r";
payload += "RRULE:FREQ=YEARLY;BYMONTH=3;BYDAY=-1SU&#13;\n\r";
payload += "DTSTART:19810329T020000&#13;\n\r";
payload += "TZNAME:GMT+2&#13;\n\r";
payload += "TZOFFSETTO:+0200&#13;\n\r";
payload += "END:DAYLIGHT&#13;\n\r";
payload += "BEGIN:STANDARD&#13;\n\r";
payload += "TZOFFSETFROM:+0200&#13;\n\r";
payload += "RRULE:FREQ=YEARLY;BYMONTH=10;BYDAY=-1SU&#13;\n\r";
payload += "DTSTART:19961027T030000&#13;\n\r";
payload += "TZNAME:GMT+1&#13;\n\r";
payload += "TZOFFSETTO:+0100&#13;\n\r";
payload += "END:STANDARD&#13;\n\r";
payload += "END:VTIMEZONE&#13;\n\r";
payload += "END:VCALENDAR&#13;\n\r";
payload += "</B:calendar-timezone>\n\r";
payload += "<D:calendar-color xmlns:D=\"http://apple.com/ns/ical/\"\n\r";
payload += "symbolic-color=\"yellow\">#FFCC00</D:calendar-color>\n\r";
payload += "</A:prop>\n\r";
payload += "</A:set>\n\r";
payload += "</B:mkcalendar>\n\r";
*/
