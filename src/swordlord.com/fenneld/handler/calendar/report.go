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
	"time"
	"github.com/beevik/etree"
	"fmt"
	"strings"
	"swordlord.com/fennelcore/db/model"
)

func Report(w http.ResponseWriter, req *http.Request) {

	vars := mux.Vars(req)
	sCalId := vars["calendar"]

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
			handleReportSyncCollection(w, req.RequestURI, root, sCalId)

		case "calendar-multiget":
			// TODO add handleReportCalendarMultiget
			handleReportCalendarMultiget(w, req.RequestURI, root, sCalId);

		case "calendar-query":
			handleReportCalendarQuery(w, req.RequestURI, root, sCalId)

		default:
			if name != "text" {
				fmt.Println("CAL-Report: not handled: " + name)
			}
	}
}

func handleReportCalendarQuery(w http.ResponseWriter, uri string, nodeQuery *etree.Element, sCalId string) {

	dRet, ms := handler.GetMultistatusDocWOResponseTag()

	cal, err := tablemodule.GetCal(sCalId)
	if err != nil {

		fmt.Println(err)

		propstat := handler.AddResponseToMultistat(ms, uri)

		handler.SendMultiStatus(w, http.StatusNotFound, dRet, propstat)
		return
	}

	props := nodeQuery.FindElements("./prop/*")

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
	//timerange := nodeQuery.FindElement("./filter/comp-filter[name='VCALENDAR']/comp-filter[name='VEVENT']/time-range")
	var dtmStart *time.Time
	var dtmEnd *time.Time

	filter := nodeQuery.FindElement("./filter/comp-filter[@name='VCALENDAR']/comp-filter[@name='VEVENT']/time-range")

	// get everything if filter is not supplied
	if filter != nil {
		start := filter.SelectAttr("start")
		if start != nil {

			// TODO warning, can have either whole day or time
			s, err := time.Parse("20060102T150405Z", start.Value)
			if err != nil {

				fmt.Println("Could not parse time format: " + start.Value)
			} else {
				dtmStart = &s
			}
		}

		end := filter.SelectAttr("end")
		if end != nil {

			// TODO warning, can have either whole day or time
			e, err := time.Parse("20060102T150405Z", end.Value)
			if err != nil {

				fmt.Println("Could not parse time format: " + end.Value)
			} else {
				dtmEnd = &e
			}
		}
	}

	rows, err := tablemodule.FindIcsByTimeslot(sCalId, dtmStart, dtmEnd)
	for _, row := range rows {

		appendIcsAsResponseToMultistat(ms, row, props, uri, cal.SupportedCalComponent)
	}

	handler.SendETreeDocument(w, http.StatusMultiStatus, dRet)
}

func appendIcsAsResponseToMultistat(ms *etree.Element, ics *model.ICS, props []*etree.Element, uri string, supportedCalComponent string){

	propstat := handler.AddResponseToMultistat(ms, uri + ics.Pkey + ".ics")

	// values to return from: /B:calendar-query/A:prop
	for _, prop := range props {

		propName := prop.Tag
		switch(propName) {

		// TODO missing:
		// <C:created-by xmlns:C="http://calendarserver.org/ns/"/>
		//    <B:schedule-tag/>
		//    <C:schedule-changes xmlns:C="http://calendarserver.org/ns/"/>
		//    <C:updated-by xmlns:C="http://calendarserver.org/ns/"/>

		case "getetag":
			getETag := propstat.CreateElement("getetag")
			getETag.Space = "d"
			getETag.SetText("etag")
			//response += "<d:getetag>\"" + Number(date) + "\"</d:getetag>";

		case "getcontenttype":
			getCT := propstat.CreateElement("getcontenttype")
			getCT.Space = "d"
			getCT.SetText("text/calendar; charset=utf-8; component=" + supportedCalComponent)
			//response += "<d:getcontenttype>text/calendar; charset=utf-8; component=" + cal.supported_cal_component + "</d:getcontenttype>";

		case "calendar-data":
			getCD := propstat.CreateElement("calendar-data")
			getCD.Space = "cal"
			getCD.SetText(ics.Content)
			//response += "<cal:calendar-data>" + ics.content + "</cal:calendar-data>"; // has to be cal: since a few lines below the namespace is cal: not c:

		default:
			if propName != "text" {
				fmt.Println("CAL-Query: not handled: " + propName)
			}
		}

	}

	handler.AddStatusToPropstat(http.StatusOK, propstat)
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
					// TODO

				case "sync-level":
					// TODO

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

	// TODO there is a solution which is much more elegant, find it :)
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


func handleReportCalendarMultiget(w http.ResponseWriter, uri string, mg *etree.Element, sCalId string) {

	dRet, ms := handler.GetMultistatusDocWOResponseTag()

	cal, err := tablemodule.GetCal(sCalId)
	if err != nil {

		fmt.Println(err)

		propstat := handler.AddResponseToMultistat(ms, uri)

		handler.SendMultiStatus(w, http.StatusNotFound, dRet, propstat)
		return
	}

	// get all href into an array for lookup/query
	hrefs := mg.FindElements("./href")

	arrHref := make([]string, len(hrefs))

	for i, href := range hrefs {

		arrHref[i] = parseHrefToIcsId(href.Text())
	}

	// get all root/prop entities for query and return
	props := mg.FindElements("./prop/*")

	// send these together with array to format response
	rows, err := tablemodule.FindIcsInList(arrHref)
	for _, row := range rows {

		appendIcsAsResponseToMultistat(ms, row, props, uri, cal.SupportedCalComponent)
	}

	handler.SendETreeDocument(w, http.StatusMultiStatus, dRet)
}

func parseHrefToIcsId(href string) string {

	arrPath := strings.Split(href, "/")
	pathCount := len(arrPath)

	if pathCount < 2 {
		// TODO continue is suboptimal, we should add a bool to the array and send a 404 for said file.
		return ""
	}

	filename := arrPath[pathCount - 1]

	arrFile := strings.Split(filename, ".")
	if len(arrFile) < 2 {
		// TODO continue is suboptimal, we should add a bool to the array and send a 404 for said file.
		return ""
	}

	return arrFile[0]
}
