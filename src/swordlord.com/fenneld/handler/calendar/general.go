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
			"fmt"
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

func Proppatch(w http.ResponseWriter, req *http.Request){

	handler.RespondWithMessage(w, http.StatusOK, "Not implemented yet")

}

func Options(w http.ResponseWriter, req *http.Request){

	handler.RespondWithStandardOptions(w, http.StatusOK, "")
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
	sCal := vars["calendar"]

	sUser, ok := req.Context().Value("auth_user").(string)
	if !ok {
		// TODO fail when there is no user, since this can't really happen!
		sUser = ""
	}
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
	sCal := vars["calendar"]
	sEvent := vars["event"]

	sUser, ok := req.Context().Value("auth_user").(string)
	if !ok {
		// TODO fail when there is no user, since this can't really happen!
		sUser = ""
	}

	bodyBuffer, _ := ioutil.ReadAll(req.Body)

	ics, err := tablemodule.AddIcs(sEvent, sUser, sCal, string(bodyBuffer))
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
