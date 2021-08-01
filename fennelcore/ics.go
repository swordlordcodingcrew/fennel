package fennelcore
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
	"github.com/Jeffail/gabs"
	"strings"
	"regexp"
)

var (
	rEOL *regexp.Regexp
)

func GetRegexEndOfLine() *regexp.Regexp {

	// lazy init and global member to not re-compile constantly...
	if rEOL == nil {

		r, err := regexp.Compile("[\r\n|\n\r|\n|\r]")
		if err != nil {

			panic(err)
		}

		rEOL = r
	}

	return rEOL
}

func ParseICS(file string) *gabs.Container {

	json := generateJSON(file)

	println(json)

	jsonParsed, err := gabs.ParseJSON([]byte(json))

	if err != nil {
		println(err)
		return nil
	}

	return jsonParsed
}

func generateJSON(file string) string {

	var result = ""

	r := GetRegexEndOfLine()
	linesUnfiltered := r.Split(file, -1)

	lines := linesUnfiltered[:0]

	// clean up lines
	for _, line := range linesUnfiltered {

		// remove empty lines
		if len(line) == 0 {
			continue
		}

		// Unfold the lines, if no : at any position, assume it is folded with previous
		if !strings.Contains(line, ":") {

			lines[len(lines) - 1] = lines[len(lines) - 1] + line
		}

		// line seems to be ok, add it to cleaned list
		lines = append(lines, line)
	}

	for _, line := range lines {

		if strings.HasPrefix(line, "BEGIN:") {

			if strings.HasSuffix(line, ".") {
				result += "\"" + line[6:len(line)-1] + "\": {"
			} else {
				result += "\"" + line[6:] + "\": {"
			}
		} else if strings.HasPrefix(line, "END:") {

			// TODO: terrible hack, fixme
			if strings.HasSuffix(result, ",") {
				result = result[0:len(result)-1]
			}

			result += "},"
		} else {

			arrKeyVal := strings.Split(line, ":")
			key := arrKeyVal[0]
			val := arrKeyVal[1]

			if strings.HasSuffix(val, ".") {
				// todo: the key.split is a terrible hack as well, we loose some information like that
				// example: DTEND;TZID=Europe/Zurich:20161210T010000Z. -> tzid will be lost
				// result += "\"" + key.split(";")[0] + "\":\"" + val.substr(0, val.length -1) + "\",";

				result += "\"" + strings.Split(key, ";")[0] + "\":\"" + val[0:len(val)-1] + "\","
			} else {
				// todo, see above
				result += "\"" + key + "\":\"" + val + "\","
			}
		}
	}

	// TODO: terrible hack, fixme
	if strings.HasSuffix(result, ",")	{
		result = result[0:len(result)-1]
	}

	result = "{" + result + "}"

	return result
}
