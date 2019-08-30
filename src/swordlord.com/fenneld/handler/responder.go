package handler
/*-----------------------------------------------------------------------------
 **
 ** - Wombag -
 **
 ** the alternative, native backend for your Wallabag apps
 **
 ** Copyright 2017-18 by SwordLord - the coding crew - http://www.swordlord.com
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
 ** LordLightningBolt@swordlord.com
 **
-----------------------------------------------------------------------------*/
import (
	"fmt"
	"net/http"
//	"swordlord.com/fenneld/render"
	"github.com/beevik/etree"
)

/*
func Render(w http.ResponseWriter, status int, r render.Render){

	w.WriteHeader(status)

	if !bodyAllowedForStatus(status) {

		r.WriteContentType(w)
		return
	}

	if err := r.Render(w); err != nil {
		panic(err)
	}
}
*/

// bodyAllowedForStatus is a copy of http.bodyAllowedForStatus non-exported function.
func bodyAllowedForStatus(status int) bool {
	switch {
	case status >= 100 && status <= 199:
		return false
	case status == 204:
		return false
	case status == 304:
		return false
	}
	return true
}


func NotImplementedYet(w http.ResponseWriter){

	w.WriteHeader(http.StatusNotImplemented)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "This function is not implemented yet\n")
}

func SetAllowHeader(w http.ResponseWriter){

	w.Header().Add("Allow", "OPTIONS, PROPFIND, HEAD, GET, REPORT, PROPPATCH, PUT, DELETE, POST, COPY, MOVE")
}

func SetDAVHeader(w http.ResponseWriter){

	w.Header().Set("DAV", "1, 3, extended-mkcol, calendar-access, calendar-schedule, calendar-proxy, calendarserver-sharing, calendarserver-subscribed, addressbook, access-control, calendarserver-principal-property-search")
}

func SetStandardXMLHeader(w http.ResponseWriter) {

	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	w.Header().Set("Server", "Fennel")
}

func SetStandardHTMLHeader(w http.ResponseWriter) {

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Server", "Fennel")
}

func RespondWithRedirect(w http.ResponseWriter, req *http.Request, uri string){

	http.Redirect(w, req, uri, http.StatusMovedPermanently)
}


func RespondWithMessage(w http.ResponseWriter, status int, message string){

	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, message + "\n")
}

func RespondWithICS(w http.ResponseWriter, status int, ics string){

	w.WriteHeader(status)
	w.Header().Set("Content-Type", "text/calendar")
	fmt.Fprintf(w, ics + "\n")
}

func RespondWithVCARD(w http.ResponseWriter, status int, vcard string){

	w.WriteHeader(status)
	w.Header().Set("Content-Type", "text/vcard; charset=utf-8")
	fmt.Fprintf(w, vcard + "\n")
}

func RespondWithUnauthenticated(w http.ResponseWriter){

	w.Header().Set("WWW-Authenticate", "Basic realm=\"Fennel\"")
	w.WriteHeader(http.StatusUnauthorized)
	fmt.Fprintf(w, "Not authorized\n")
}

func RespondWithStandardOptions(w http.ResponseWriter, status int, message string) {

	//log.debug("pushOptionsResponse called");

	SetStandardHTMLHeader(w)
	SetDAVHeader(w)
	SetAllowHeader(w)

	RespondWithMessage(w, http.StatusOK, "")
}

func SendMultiStatus(w http.ResponseWriter, httpStatus int, dRet *etree.Document, propstat *etree.Element) {

	SetStandardHTMLHeader(w)
	SetDAVHeader(w)
	SetAllowHeader(w)

	AddStatusToPropstat(httpStatus, propstat)

	w.WriteHeader(http.StatusMultiStatus)
	dRet.WriteTo(w)
}

func SendETreeDocument(w http.ResponseWriter, status int, dRet *etree.Document) {

	SetStandardHTMLHeader(w)
	SetDAVHeader(w)
	SetAllowHeader(w)

	w.WriteHeader(status)
	dRet.WriteTo(w)
}

func GetMultistatusDoc(sURL string) (*etree.Document, *etree.Element) {

	doc, ms := GetMultistatusDocWOResponseTag()

	response := ms.CreateElement("response")
	response.Space = "d"

	href := response.CreateElement("href")
	href.SetText(sURL)
	href.Space = "d"

	propstat := response.CreateElement("propstat")
	propstat.Space = "d"

	return doc, propstat
}

func GetMultistatusDocWOResponseTag() (*etree.Document, *etree.Element) {

	doc := etree.NewDocument()
	doc.Indent(2)
	doc.CreateProcInst("xml", `version="1.0" encoding="utf-8"`)

	ms := doc.CreateElement("multistatus")
	ms.Space = "d"

	/*
	<d:multistatus xmlns:d="DAV:"
				xmlns:s="http://swordlord.com/ns"
				xmlns:cal="urn:ietf:params:xml:ns:caldav"
				xmlns:cs="http://calendarserver.org/ns/"
				xmlns:card="urn:ietf:params:xml:ns:carddav">
	  <d:response>
		<d:href>/card/user/</d:href>
		<d:propstat>
		  <d:prop>
	*/
	ms.CreateAttr("xmlns:d", "DAV:")
	ms.CreateAttr("xmlns:d", "DAV:")
	ms.CreateAttr("xmlns:s", "http://swordlord.com/ns")
	ms.CreateAttr("xmlns:cal", "urn:ietf:params:xml:ns:caldav")
	ms.CreateAttr("xmlns:cs", "http://calendarserver.org/ns/")
	ms.CreateAttr("xmlns:card", "urn:ietf:params:xml:ns:carddav")

	return doc, ms
}


func AddResponseWStatusToMultistat(ms *etree.Element, uri string, httpStatus int) *etree.Element {

	propstat := AddResponseToMultistat(ms, uri)

	AddStatusToPropstat(httpStatus, propstat)

	return propstat
}

func AddResponseToMultistat(ms *etree.Element, uri string) *etree.Element {

	response := ms.CreateElement("response")
	response.Space = "d"

	href := response.CreateElement("href")
	href.SetText(uri)
	href.Space = "d"

	propstat := response.CreateElement("propstat")
	propstat.Space = "d"

	return propstat
}

func AddStatusToPropstat(httpStatus int, propstat *etree.Element) {

	status := propstat.CreateElement("status")
	status.Space = "d"

	switch(httpStatus) {

	case http.StatusNotFound:
		status.SetText("HTTP/1.1 404 Not Found")
	case http.StatusOK:
		status.SetText("HTTP/1.1 200 OK")
	default:
		status.SetText("HTTP/1.1 500 Internal Server Error")
	}
}
