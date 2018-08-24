package auth

import (
	"github.com/urfave/negroni"
	"net/http"
	"strings"
	"encoding/base64"
				"github.com/pkg/errors"
	"swordlord.com/fenneld/handler"
	"swordlord.com/fennelcore"
	"context"
)

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

func NewFennelAuthentication() negroni.HandlerFunc {

	return func(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {

		// do not manage service detection
		if strings.HasPrefix(req.RequestURI, "/.well-known") {
			next(w, req)
			return
		}

		authHeader := req.Header.Get("Authorization")

		// not authenticated
		if len(authHeader) == 0 {

			handler.RespondWithUnauthenticated(w)
			return
		}

		err, uid, pwd := parseAuthHeader(authHeader)
		if err != nil {

			handler.RespondWithUnauthenticated(w)
			return
		}

		err, roles := ValidateUser(uid, pwd)
		if err != nil {

			handler.RespondWithUnauthenticated(w)
			return
		}

		println(roles)

		// env/context var when authenticated
		ctx := req.Context()
		ctx = context.WithValue(ctx, "auth_user", uid)

		next(w, req.WithContext(ctx))
	}
}

func ValidateUser(uid string, pwd string) (error, string) {

	authModule := fennelcore.GetStringFromConfig("auth.module")

	switch authModule {

		case "htpasswd":
			return ValidateUserHTPasswd(uid, pwd)
		case "ldap":
		case "db":
			return ValidateDB(uid, pwd)
		case "courier":
			return ValidateCourier(uid, pwd)
		default:
			return errors.New("Authentication Module unknown, can't authenticate"), ""
	}

	return errors.New("Authentication Module unknown, can't authenticate"), ""
}

func parseAuthHeader(header string) (error, string, string) {

	aHeader := strings.SplitN(header, " ", 2)

	if len(aHeader) != 2 || aHeader[0] != "Basic" {
		return errors.New("There is no Basic Header, or some other Header problem"), "", ""
	}

	sDecoded, err := base64.StdEncoding.DecodeString(aHeader[1])
	if err != nil {
		return errors.New("Can't decrypt BASE64"), "", ""
	}

	aUidPwd := strings.SplitN(string(sDecoded), ":", 2)
	if len(aUidPwd) != 2 {
		return errors.New("Can't split to username:password"), "", ""
	}

	return nil, aUidPwd[0], aUidPwd[1]
}
