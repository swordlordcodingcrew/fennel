package auth
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
	"os"
	"encoding/csv"
	"swordlord.com/fennelcore"
	"strings"
	"crypto/sha1"
	"crypto/subtle"
	"encoding/base64"
	"bytes"
	"errors"
)

type usermap struct {
	Users map[string]string // The map of htpasswd User key value pairs
	IsInitialised bool
}

var um usermap // A reference to the singleton

func LoadHTPasswd(fromFile string) error {

	r, err := os.Open(fromFile)
	if err != nil {
		return err
	}

	defer r.Close()

	csv_reader := csv.NewReader(r)
	csv_reader.Comma = ':'
	csv_reader.Comment = '#'
	csv_reader.TrimLeadingSpace = true

	records, err := csv_reader.ReadAll()
	if err != nil {
		return err
	}

	// Create a straps object
	um = usermap{
		Users: make(map[string]string),
		IsInitialised: false,
	}

	for _, record := range records {

		um.Users[record[0]] = record[1]
		//println(record)
	}

	um.IsInitialised = true

	return nil
}


func ValidateUserHTPasswd(uid string, pwd string) (error, string) {

	// lazy initialisation
	if !um.IsInitialised {

		htpasswd := fennelcore.GetStringFromConfig("auth.file")
		err := LoadHTPasswd(htpasswd)
		if err != nil {
			return err, ""
		}
	}

	pwdHash := um.Users[uid]

	// check password
	if strings.HasPrefix(pwdHash, "{SHA}") {

		d := sha1.New()
		d.Write([]byte(pwd))
		if subtle.ConstantTimeCompare([]byte(pwdHash)[5:], []byte(base64.StdEncoding.EncodeToString(d.Sum(nil)))) != 1 {
			return errors.New("Password not correct"), ""
		}
	} else if strings.HasPrefix(pwdHash, "$apr1$"){

		err := compareMD5HashAndPassword([]byte(pwdHash), []byte(pwd))
		if err != nil {
			return err, ""
		} else {
			return nil, ""
		}
	} else {


	}

	return nil, ""
}

func compareMD5HashAndPassword(hashedPassword, password []byte) error {
	parts := bytes.SplitN(hashedPassword, []byte("$"), 4)
	if len(parts) != 4 {
		return errors.New("Password not correct")
	}
	magic := []byte("$" + string(parts[1]) + "$")
	salt := parts[2]

	if subtle.ConstantTimeCompare(hashedPassword, MD5Crypt(password, salt, magic)) != 1 {
		return errors.New("Password not correct")
	}
	return nil
}