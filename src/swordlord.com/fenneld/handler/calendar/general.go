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
	"encoding/xml"
	"github.com/gorilla/mux"
	"swordlord.com/fennelcore/db/tablemodule"
	"time"
	"github.com/beevik/etree"
	"fmt"
	"swordlord.com/fennelcore/db/model"
	"io/ioutil"
	"log"
)

type Xmlmakecalendar struct {
	XMLName xml.Name
	Set       Xmlset   `xml:"set"`
}

type Xmlset struct {
	XMLName xml.Name
	Prop	Xmlprop   `xml:"prop"`
}

type Xmlprop struct {
	XMLName 				xml.Name
	Displayname				string 	`xml:"displayname"`
	CalendarOrder			uint	`xml:"calendar-order"`
	CalendarTimezone		string	`xml:"calendar-timezone"`
	CalendarColour			string	`xml:"calendar-color"`
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
	sUser := vars["user"]
	sCal := vars["calendar"]

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

func Proppatch(w http.ResponseWriter, req *http.Request){

	handler.RespondWithMessage(w, http.StatusOK, "Not implemented yet")

}

func Report(w http.ResponseWriter, req *http.Request) {

	vars := mux.Vars(req)
	//sUser := vars["user"]
	sCalId := vars["calendar"]

	dRet, propstat := handler.GetMultistatusDoc(req.RequestURI)

	doc := etree.NewDocument()
	size, err := doc.ReadFrom(req.Body)
	if err != nil || size == 0 {

		fmt.Printf("Error reading XML Body. Error: %s Size: %v", err, size)

		handler.SendMultiStatus(w, http.StatusNotFound, dRet, propstat)
		return
	}

	root := doc.Root()
	name := root.Tag

	switch (name) {

		case "sync-collection":
			handleReportSyncCollection(w, req.RequestURI, root, sCalId)

		case "calendar-multiget":
			//handleReportCalendarMultiget(comm);

		case "calendar-query":
			handleReportCalendarQuery(w, req.RequestURI, root, sCalId)

		default:
			if name != "text" {
				fmt.Println("CAL-Report: not handled: " + name)
			}
	}
}

func Options(w http.ResponseWriter, req *http.Request){

	handler.RespondWithStandardOptions(w, http.StatusOK, "")
}

func handleReportCalendarQuery(w http.ResponseWriter, uri string, nodeCalendarQuery *etree.Element, sCalId string) {

	dRet, ms := handler.GetMultistatusDocWOResponseTag()

	cal, err := tablemodule.GetCal(sCalId)
	if err != nil {

		fmt.Println(err)

		propstat := handler.AddResponseToMultistat(ms, uri)

		handler.SendMultiStatus(w, http.StatusNotFound, dRet, propstat)
		return
	}

	retProps := nodeCalendarQuery.FindElements("./prop/*")

	// TODO: check filter:
	// <B:comp-filter name=\"VCALENDAR\">\n\r";
	//    <B:comp-filter name=\"VEVENT\">\n\r";
	//    <B:time-range start=\"" + now.subtract(1, "h").format("YMMDD[T]HH0000[Z]") + "\"/>\n\r";
	//    </B:comp-filter>\n\r";
	//</B:comp-filter>\n\r
	//
	// BEGIN:VEVENT.
	// DTSTART;TZID=Europe/Zurich:20161014T120000Z.
	// DTEND;TZID=Europe/Zurich:20161014T130000Z
	// parse when storing
	timerange := nodeCalendarQuery.FindElement("./filter/comp-filter[name='VCALENDAR']/comp-filter[name='VEVENT']/time-range")

	fmt.Println(timerange)

	rows, err := tablemodule.FindIcsByTimeslot(sCalId, time.Time{}, time.Time{})
	for _, row := range rows {

		propstat := handler.AddResponseToMultistat(ms, uri + "/" + row.Pkey + ".ics")

		// values to return: /B:calendar-query/A:prop
		for _, prop := range retProps {

			propName := prop.Tag
			switch(propName) {
				case "getetag":
					getETag := propstat.CreateElement("getetag")
					getETag.Space = "d"
					getETag.SetText("etag")
				//response += "<d:getetag>\"" + Number(date) + "\"</d:getetag>";

				case "getcontenttype":
					getCT := propstat.CreateElement("getcontenttype")
					getCT.Space = "d"
					getCT.SetText("text/calendar; charset=utf-8; component=" + cal.SupportedCalComponent)
				//response += "<d:getcontenttype>text/calendar; charset=utf-8; component=" + cal.supported_cal_component + "</d:getcontenttype>";

				case "calendar-data":
					getCD := propstat.CreateElement("calendar-data")
					getCD.Space = "cal"
					getCD.SetText(row.Content)
				//response += "<cal:calendar-data>" + ics.content + "</cal:calendar-data>"; // has to be cal: since a few lines below the namespace is cal: not c:

				default:
					if propName != "text" {
						fmt.Println("CAL-Query: not handled: " + propName)
					}
			}

		}

		handler.AddStatusToPropstat(http.StatusOK, propstat)

	}

	handler.SendETreeDocument(w, http.StatusMultiStatus, dRet)
}

func handleReportSyncCollection(w http.ResponseWriter, uri string, nodeSyncCollection *etree.Element, sCalId string) {

	dRet, propstat := handler.GetMultistatusDoc(uri)

	cal, err := tablemodule.GetCal(sCalId)
	if err != nil {

		fmt.Println(err)

		handler.SendMultiStatus(w, http.StatusNotFound, dRet, propstat)
		return
	}

	fmt.Println(cal.Pkey)

	rows, err := tablemodule.FindIcsByCalendar(sCalId)
	if err != nil {

		fmt.Println(err)

		handler.SendMultiStatus(w, http.StatusNotFound, dRet, propstat)
		return
	}

	for _, ics := range rows {

		for _, el := range nodeSyncCollection.ChildElements() {

			//fmt.Println(e.Tag)
			name := el.Tag
			switch(name) {

				case "sync-token":

				case "sync-level":

				case "prop":
					//response += handleReportCalendarProp(comm, child, cal, ics);
					// TODO
					fmt.Println("found: " + ics.Content)

				default:
					if name != "text" {
						fmt.Println("CAL-RSC: not handled: " + name)
					}
			}
		}
	}

	ms := dRet.FindElement("/multistatus")

	st := ms.CreateElement("sync-token")
	st.Space = "d"
	st.SetText("https://swordlord.org/ns/sync/" + fmt.Sprint(cal.Synctoken))

	handler.SendETreeDocument(w, http.StatusMultiStatus, dRet)
}

/*
function handleReportCalendarProp(comm, node, cal, ics)
{
    var response = "";

    var reqUrl = comm.getURL();
    reqUrl += reqUrl.match("\/$") ? "" : "/";

    response += "<d:response>";
    response += "<d:href>" + reqUrl + ics.pkey + ".ics</d:href>";
    response += "<d:propstat><d:prop>";

    var childs = node.childNodes();

    var date = Date.parse(ics.updatedAt);

    var len = childs.length;
    for (var i=0; i < len; ++i)
    {
        var child = childs[i];
        var name = child.name();
        switch(name)
        {
            case 'getetag':
                response += "<d:getetag>\"" + Number(date) + "\"</d:getetag>";
                break;

            case 'getcontenttype':
                response += "<d:getcontenttype>text/calendar; charset=utf-8; component=" + cal.supported_cal_component + "</d:getcontenttype>";
                break;

            default:
                if(name != 'text') log.warn("P-R: not handled: " + name);
                break;
        }
    }

    response += "</d:prop>";
    response += "<d:status>HTTP/1.1 200 OK</d:status>";
    response += "</d:propstat>";
    response += "</d:response>";

    return response;
}*/

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

		case "calendar-order":
			response += "" // "<xical:calendar-order xmlns:xical=\"http://apple.com/ns/ical/\">" + cal.Order + "</xical:calendar-order>";

		case "calendar-timezone":
			var timezone = cal.Timezone;
			//timezone = timezone.replace(/\r\n|\r|\n/g,"&#13;\r\n");
			response += "<cal:calendar-timezone>" + timezone + "</cal:calendar-timezone>";

		case "current-user-privilege-set":
			response += "" //getCurrentUserPrivilegeSet();

			//case "default-alarm-vevent-date":
			//case "default-alarm-vevent-datetime":

		case "displayname":
			response += "<d:displayname>" + cal.Displayname + "</d:displayname>";

			//case "language-code":
			//case "location-code":

		case "owner":
			response += "<d:owner><d:href>/p/" + user +"/</d:href></d:owner>";

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


func MakeCalendar(w http.ResponseWriter, req *http.Request){

	vars := mux.Vars(req)
	sUser := vars["user"]
	sCal := vars["calendar"]

	//var dm Xmlmakecalendar

	decoder := xml.NewDecoder(req.Body)
	sentCal := Xmlmakecalendar{}
	err := decoder.Decode(&sentCal)
	if err != nil {

		// TODO: or internal server error?
		handler.RespondWithMessage(w, http.StatusPreconditionFailed, err.Error())
		return
	}
	// fmt.Println(err)
	// fmt.Println(sentCal)

	// fmt.Println(sentCal.Set.Prop.Displayname)

	prop := sentCal.Set.Prop

	// TODO: set freebusyset only to yes if tag exists ->"YES"
	// TODO: set supported cal component to "VEVENT" if tag exists

	cal, err := tablemodule.AddCal(sUser, sCal, prop.Displayname, prop.CalendarColour, "YES", prop.CalendarOrder, "VEVENT", 0, prop.CalendarTimezone)
	if err != nil {

		handler.RespondWithMessage(w, http.StatusPreconditionFailed, err.Error())
		return
	}

	handler.RespondWithMessage(w, http.StatusCreated, "Make calendar: " + cal.Pkey + " for user: " + sUser)
}

func Put(w http.ResponseWriter, req *http.Request){

	vars := mux.Vars(req)
	sUser := vars["user"]
	sCal := vars["calendar"]
	sEvent := vars["event"]

	// var parser = require("../libs/parser");
	//    var pbody = parser.parseICS(body);
	//
	//    var dtStart = moment(pbody.VCALENDAR.VEVENT.DTSTART);
	//    var dtEnd = moment(pbody.VCALENDAR.VEVENT.DTEND);
	// dates toISOString

	bodyBuffer, _ := ioutil.ReadAll(req.Body)

	ics, err := tablemodule.AddIcs(sEvent, sUser, sCal, time.Now(), time.Now(), string(bodyBuffer))
	if err != nil {

		handler.RespondWithMessage(w, http.StatusPreconditionFailed, err.Error())
		return
	}

	// todo: cal increment sync token
	// todo: return e-tag
	handler.RespondWithMessage(w, http.StatusCreated, "ICS added: " + ics.Pkey)
}

func Get(w http.ResponseWriter, req *http.Request) {

	vars := mux.Vars(req)
	//sUser := vars["user"]
	//sCal := vars["calendar"]
	sEvent := vars["event"]

	ics, err := tablemodule.GetICS(sEvent)

	if err != nil {

		fmt.Println("err: could not find ics " + sEvent)
		// TODO send error
		handler.RespondWithMessage(w, http.StatusInternalServerError, err.Error())
		return
	}

	handler.RespondWithICS(w, http.StatusOK, ics.Content)
}

func Delete(w http.ResponseWriter, req *http.Request){

	vars := mux.Vars(req)
	//sUser := vars["user"]
	//sCal := vars["calendar"]
	sEvent := vars["event"]

	err := tablemodule.DeleteIcs(sEvent)

	if err != nil {
		log.Printf("Error with deleting Ics %q: %s\n", sEvent)

		handler.RespondWithMessage(w, http.StatusInternalServerError, err.Error())

		return
	}

	handler.RespondWithMessage(w, http.StatusOK, "Deleted")
}

func Move(w http.ResponseWriter, req *http.Request){

	handler.RespondWithMessage(w, http.StatusOK, "Not implemented yet")

}
