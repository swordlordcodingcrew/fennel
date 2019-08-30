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
	)

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

func Options(w http.ResponseWriter, req *http.Request){

	handler.RespondWithStandardOptions(w, http.StatusOK, "")
}
