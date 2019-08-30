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

func Report(w http.ResponseWriter, req *http.Request) {

	doc := etree.NewDocument()
	size, err := doc.ReadFrom(req.Body)
	if err != nil || size == 0 {

		fmt.Printf("Error reading XML Body. Error: %s Size: %v", err, size)

		handler.RespondWithMessage(w, http.StatusInternalServerError, "")
		return
	}

	root := doc.Root()
	name := root.Tag

	/*
	<?xml version="1.0" encoding="UTF-8"?>
	<A:principal-search-property-set xmlns:A="DAV:"/>
	*/
	
	switch name {

	case "principal-search-property-set":
		handleSearchPropertySet(w, req.RequestURI, root)

	default:
		if name != "text" {
			fmt.Println("Principal-Report: not handled: " + name)
		}
	}
}

func handleSearchPropertySet(w http.ResponseWriter, url string, root *etree.Element){

	/*
	<?xml version="1.0"?>
	<d:principal-search-property-set xmlns:d="DAV:" ...>
	  <d:principal-search-property>
		<d:prop>
		  <d:displayname/>
		</d:prop>
		<d:description xml:lang="en">Display name</d:description>
	  </d:principal-search-property>
	  <d:principal-search-property>
		<d:prop>
		  <s:email-address/>
		</d:prop>
		<d:description xml:lang="en">Email address</d:description>
	  </d:principal-search-property>
	  <d:principal-search-property>
		<d:prop>
		  <cal:calendar-user-address-set/>
		</d:prop>
		<d:description xml:lang="en">Calendar address</d:description>
	  </d:principal-search-property>
	</d:principal-search-property-set>
	*/

	dRet := etree.NewDocument()
	dRet.Indent(2)
	dRet.CreateProcInst("xml", `version="1.0" encoding="utf-8"`)

	psps := dRet.CreateElement("principal-search-property-set")
	psps.Space = "d"

	psps.CreateAttr("xmlns:d", "DAV:")
	psps.CreateAttr("xmlns:d", "DAV:")
	psps.CreateAttr("xmlns:s", "http://swordlord.com/ns")
	psps.CreateAttr("xmlns:cal", "urn:ietf:params:xml:ns:caldav")
	psps.CreateAttr("xmlns:cs", "http://calendarserver.org/ns/")
	psps.CreateAttr("xmlns:card", "urn:ietf:params:xml:ns:carddav")

	addPrincipalSearchProperty(psps, "displayname", "d", "Display Name")
	addPrincipalSearchProperty(psps, "email-address", "s", "Email address")
	addPrincipalSearchProperty(psps, "calendar-user-address-set", "cal", "Calendar address")

	handler.SendETreeDocument(w, http.StatusOK, dRet)
}

func addPrincipalSearchProperty(psps *etree.Element, prop string, ns string, desc string) {

	/*
	<d:principal-search-property>
		<d:prop>
		  <cal:calendar-user-address-set/>
		</d:prop>
		<d:description xml:lang="en">Calendar address</d:description>
	  </d:principal-search-property>
	*/

	psp := psps.CreateElement("principal-search-property")
	psp.Space = "d"

	p := psp.CreateElement("prop")
	p.Space = "d"

	el := p.CreateElement(prop)
	el.Space = ns

	d := psp.CreateElement("description")
	d.Space = "d"
	d.CreateAttr("xml:lang", "en")
	d.SetText(desc)
}