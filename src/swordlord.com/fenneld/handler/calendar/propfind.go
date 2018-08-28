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
	"strconv"
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
	fillPropfindResponse(prop, sUser, model.CAL{}, propsQuery, true)

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
	fillPropfindResponse(prop, sUser, cal, propsQuery, false)

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

func fillPropfindResponse(node *etree.Element, user string, cal model.CAL, props []*etree.Element, isRoot bool) {

	// TODO
	token := ""

	for _, e := range props {

		//fmt.Println(e.Tag)
		name := e.Tag
		switch(name) {

		//case "add-member":

		case "allowed-sharing-modes":
			// <cs:allowed-sharing-modes><cs:can-be-shared/><cs:can-be-published/></cs:allowed-sharing-modes>";
			asm := node.CreateElement("allowed-sharing-modes")
			asm.Space = "cs"
			cbs := asm.CreateElement("can-be-shared")
			cbs.Space = "cs"
			cbp := asm.CreateElement("can-be-published")
			cbp.Space = "cs"

			//case "autoprovisioned":
			//case "bulk-requests":

		case "calendar-color":
			// <xical:calendar-color xmlns:xical=\"http://apple.com/ns/ical/\">" + cal.Colour + "</xical:calendar-color>";
			cc := node.CreateElement("calendar-color")
			cc.Space = "xical"
			cc.CreateAttr("xmlns:xical", "http://apple.com/ns/ical/")
			cc.SetText(cal.Colour)

			//case "calendar-description":
			//case "calendar-free-busy-set":
			//response += "<d:response><d:href>/</d:href></d:response>";

		case "calendar-order":
			// "<xical:calendar-order xmlns:xical=\"http://apple.com/ns/ical/\">" + cal.Order + "</xical:calendar-order>";
			co := node.CreateElement("calendar-order")
			co.Space = "xical"
			co.CreateAttr("xmlns:xical", "http://apple.com/ns/ical/")
			co.SetText(strconv.Itoa(cal.Order))

		case "calendar-timezone":
			var timezone = cal.Timezone;
			// TODO check why here we had a replace
			//timezone = timezone.replace(/\r\n|\r|\n/g,"&#13;\r\n");
			//"<cal:calendar-timezone>" + timezone + "</cal:calendar-timezone>";
			ct := node.CreateElement("calendar-timezone")
			ct.Space = "cal"
			ct.SetText(timezone)

		case "current-user-privilege-set":
			getCurrentUserPrivilegeSet(node)

		case "current-user-principal":
			// <d:current-user-principal><d:href>/p/" + username + "/</d:href></d:current-user-principal>
			cup := node.CreateElement("current-user-principal")
			cup.Space = "d"
			handler.AddURLElement(cup, "/p/" + user + "/")

			//case "default-alarm-vevent-date":
			//case "default-alarm-vevent-datetime":

		case "displayname":
			// "<d:displayname>" + cal.Displayname + "</d:displayname>"
			ct := node.CreateElement("displayname")
			ct.Space = "d"
			ct.SetText(cal.Displayname)

			//case "language-code":
			//case "location-code":

		case "owner":
			// "<d:owner><d:href>/p/" + user +"/</d:href></d:owner>"
			o := node.CreateElement("owner")
			o.Space = "d"
			handler.AddURLElement(o, "/p/" + user + "/")

		case "principal-collection-set":
			//"<d:principal-collection-set><d:href>/p/</d:href></d:principal-collection-set>"
			pcs := node.CreateElement("principal-collection-set")
			pcs.Space = "d"
			handler.AddURLElement(pcs, "/p/")

			// TODO Check if relative URL is acceptable. if so -> OK
		case "pre-publish-url":
			//"<cs:pre-publish-url><d:href>https://127.0.0.1/cal/" + user + "/" + cal.Pkey + "</d:href></cs:pre-publish-url>";
			pcs := node.CreateElement("pre-publish-url")
			pcs.Space = "cs"
			handler.AddURLElement(pcs, "/cal/" + user + "/" + cal.Pkey)

			//case "publish-url":
			//case "push-transports":
			//case "pushkey":
			//case "quota-available-bytes":
			//case "quota-used-bytes":
			//case "refreshrate":
			//case "resource-id":

		case "resourcetype":
			// "<d:resourcetype><d:collection/><cal:calendar/></d:resourcetype>";
			rt := node.CreateElement("resourcetype")
			rt.Space = "d"
			col := rt.CreateElement("collection")
			col.Space = "d"
			cal := rt.CreateElement("calendar")
			cal.Space = "cal"

		case "schedule-calendar-transp":
			// "<cal:schedule-calendar-transp><cal:opaque/></cal:schedule-calendar-transp>";
			sct := node.CreateElement("schedule-calendar-transp")
			sct.Space = "cal"
			o := sct.CreateElement("opaque")
			o.Space = "cal"

			//case "schedule-default-calendar-URL":
			//case "source":
			//case "subscribed-strip-alarms":
			//case "subscribed-strip-attachments":
			//case "subscribed-strip-todos":
			//case "supported-calendar-component-set":

		case "supported-calendar-component-sets":
			// "<cal:supported-calendar-component-set><cal:comp name=\"VEVENT\"/></cal:supported-calendar-component-set>";
			scc := node.CreateElement("supported-calendar-component-set")
			scc.Space = "cal"
			c := scc.CreateElement("comp")
			c.Space = "cal"
			c.CreateAttr("name", "VEVENT")

		case "supported-report-set":
			getSupportedReportSet(node, isRoot)

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
			getACL(node, user)

		case "getcontenttype":
			//response += "<d:getcontenttype>text/calendar;charset=utf-8</d:getcontenttype>";
			gct := node.CreateElement("getcontenttype")
			gct.Space = "d"
			gct.SetText("text/calendar;charset=utf-8")

		default:
			if name != "text" {
				fmt.Println("CAL-PF: not handled: " + name)
			}
		}
	}
}

func getCurrentUserPrivilegeSet(node *etree.Element) {

	/*
	response += "<d:current-user-privilege-set>";
    response += "<d:privilege xmlns:d=\"DAV:\"><cal:read-free-busy/></d:privilege>";
    response += "<d:privilege xmlns:d=\"DAV:\"><d:write/></d:privilege>";
    response += "<d:privilege xmlns:d=\"DAV:\"><d:write-acl/></d:privilege>";
    response += "<d:privilege xmlns:d=\"DAV:\"><d:write-content/></d:privilege>";
    response += "<d:privilege xmlns:d=\"DAV:\"><d:write-properties/></d:privilege>";
    response += "<d:privilege xmlns:d=\"DAV:\"><d:bind/></d:privilege>";
    response += "<d:privilege xmlns:d=\"DAV:\"><d:unbind/></d:privilege>";
    response += "<d:privilege xmlns:d=\"DAV:\"><d:unlock/></d:privilege>";
    response += "<d:privilege xmlns:d=\"DAV:\"><d:read/></d:privilege>";
    response += "<d:privilege xmlns:d=\"DAV:\"><d:read-acl/></d:privilege>";
    response += "<d:privilege xmlns:d=\"DAV:\"><d:read-current-user-privilege-set/></d:privilege>";
    response += "</d:current-user-privilege-set>";
	*/

	cups := node.CreateElement("current-user-privilege-set")
	cups.Space = "d"

	addPrivilegeToPrivilegeSet(cups, "cal", "read-free-busy")

	addPrivilegeToPrivilegeSet(cups, "d", "write")
	addPrivilegeToPrivilegeSet(cups, "d", "write-acl")
	addPrivilegeToPrivilegeSet(cups, "d", "write-content")
	addPrivilegeToPrivilegeSet(cups, "d", "write-properties")
	addPrivilegeToPrivilegeSet(cups, "d", "bind")
	addPrivilegeToPrivilegeSet(cups, "d", "unbind")
	addPrivilegeToPrivilegeSet(cups, "d", "unlock")
	addPrivilegeToPrivilegeSet(cups, "d", "read")
	addPrivilegeToPrivilegeSet(cups, "d", "read-acl")
	addPrivilegeToPrivilegeSet(cups, "d", "read-current-user-privilege-set")
}

func addPrivilegeToPrivilegeSet(cups *etree.Element, namespace string, privilege string) {

	p := cups.CreateElement("privilege")
	p.Space = "d"
	p.CreateAttr("xmlns:d", "DAV")

	e := p.CreateElement(privilege)
	e.Space = namespace
}

func getSupportedReportSet(node *etree.Element, isRoot bool) {

	/*
	response += "<d:supported-report-set>";

	if(!isRoot)
	{
		response += "<d:supported-report><d:report><cal:calendar-multiget/></d:report></d:supported-report>";
		response += "<d:supported-report><d:report><cal:calendar-query/></d:report></d:supported-report>";
		response += "<d:supported-report><d:report><cal:free-busy-query/></d:report></d:supported-report>";
	}

	response += "<d:supported-report><d:report><d:sync-collection/></d:report></d:supported-report>";
	response += "<d:supported-report><d:report><d:expand-property/></d:report></d:supported-report>";
	response += "<d:supported-report><d:report><d:principal-property-search/></d:report></d:supported-report>";
	response += "<d:supported-report><d:report><d:principal-search-property-set/></d:report></d:supported-report>";
	response += "</d:supported-report-set>";
	*/
	srs := node.CreateElement("supported-report-set")
	srs.Space = "d"

	if isRoot {

		addSupportedReport(srs, "calendar-multiget")
		addSupportedReport(srs, "calendar-query")
		addSupportedReport(srs, "free-busy-query")
	}

	addSupportedReport(srs, "sync-collection")
	addSupportedReport(srs, "expand-property")
	addSupportedReport(srs, "principal-property-search")
	addSupportedReport(srs, "principal-search-property-set")
}

func addSupportedReport(srs *etree.Element, report string) {

	sr := srs.CreateElement("supported-report")
	sr.Space = "d"

	r := sr.CreateElement("report")
	r.Space = "d"

	e := r.CreateElement(report)
	e.Space = "d"
}

func getACL(node *etree.Element, user string) {

	/*
	    response += "<d:acl>";
    response += "    <d:ace>";
    response += "        <d:principal><d:href>/p/" + username + "</d:href></d:principal>";
    response += "        <d:grant><d:privilege><d:read/></d:privilege></d:grant>";
    response += "        <d:protected/>";
    response += "    </d:ace>";

    response += "    <d:ace>";
    response += "        <d:principal><d:href>/p/" + username + "</d:href></d:principal>";
    response += "        <d:grant><d:privilege><d:write/></d:privilege></d:grant>";
    response += "        <d:protected/>";
    response += "    </d:ace>";

    response += "    <d:ace>";
    response += "        <d:principal><d:href>/p/" + username + "/calendar-proxy-write/</d:href></d:principal>";
    response += "        <d:grant><d:privilege><d:read/></d:privilege></d:grant>";
    response += "        <d:protected/>";
    response += "    </d:ace>";

    response += "    <d:ace>";
    response += "        <d:principal><d:href>/p/" + username + "/calendar-proxy-write/</d:href></d:principal>";
    response += "        <d:grant><d:privilege><d:write/></d:privilege></d:grant>";
    response += "        <d:protected/>";
    response += "    </d:ace>";

    response += "    <d:ace>";
    response += "        <d:principal><d:href>/p/" + username + "/calendar-proxy-read/</d:href></d:principal>";
    response += "        <d:grant><d:privilege><d:read/></d:privilege></d:grant>";
    response += "        <d:protected/>";
    response += "    </d:ace>";

    response += "    <d:ace>";
    response += "        <d:principal><d:authenticated/></d:principal>";
    response += "        <d:grant><d:privilege><cal:read-free-busy/></d:privilege></d:grant>";
    response += "        <d:protected/>";
    response += "    </d:ace>";

    response += "    <d:ace>";
    response += "        <d:principal><d:href>/p/system/admins/</d:href></d:principal>";
    response += "        <d:grant><d:privilege><d:all/></d:privilege></d:grant>";
    response += "        <d:protected/>";
    response += "    </d:ace>";
*/
	acl := node.CreateElement("acl")
	acl.Space = "d"

	addACEwURL(acl, "/p/" + user, "read")
	addACEwURL(acl, "/p/" + user, "write")

	addACEwURL(acl, "/p/" + user + "/calendar-proxy-write/", "read")
	addACEwURL(acl, "/p/" + user + "/calendar-proxy-write/", "write")

	addACEwURL(acl, "/p/" + user + "/calendar-proxy-read/", "read")

	addACEFreeBusy(acl)

	addACEwURL(acl, "/p/system/admins/", "all")
}

func addACEwURL(acl *etree.Element, url string, privilege string)  {

	//    <d:ace>";
	//        <d:principal><d:href>/p/" + username + "</d:href></d:principal>";
	//        <d:grant><d:privilege><d:read/></d:privilege></d:grant>";
	//        <d:protected/>";
	//    </d:ace>";

	ace := acl.CreateElement("ace")
	ace.Space = "d"

	princ := ace.CreateElement("principal")
	princ.Space = "d"

	href := princ.CreateElement("href")
	href.Space = "d"
	href.SetText(url)

	g := ace.CreateElement("grant")
	g.Space = "d"

	priv := g.CreateElement("privilege")
	priv.Space = "d"

	rw := priv.CreateElement(privilege)
	rw.Space = "d"

	prot := ace.CreateElement("protected")
	prot.Space = "d"
}

func addACEFreeBusy(acl *etree.Element)  {

	//    <d:ace>";
	//        <d:principal><d:authenticated/></d:principal>";
	//        <d:grant><d:privilege><cal:read-free-busy/></d:privilege></d:grant>";
	//        <d:protected/>";
	//    </d:ace>";

	ace := acl.CreateElement("ace")
	ace.Space = "d"

	princ := ace.CreateElement("principal")
	princ.Space = "d"

	a := princ.CreateElement("authenticated")
	a.Space = "d"

	g := ace.CreateElement("grant")
	g.Space = "d"

	priv := g.CreateElement("privilege")
	priv.Space = "d"

	rfb := priv.CreateElement("read-free-busy")
	rfb.Space = "cal"

	prot := ace.CreateElement("protected")
	prot.Space = "d"
}
