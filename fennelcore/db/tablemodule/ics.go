package tablemodule

/*-----------------------------------------------------------------------------
 **
 ** - Fennel -
 **
 ** your lightweight CalDAV and CardDAV server
 **
 ** Copyright 2019 by SwordLord - the coding crew - http://www.swordlord.com
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
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Jeffail/gabs"
	fcdb "github.com/swordlordcodingcrew/fennel/fennelcore"
	"github.com/swordlordcodingcrew/fennel/fennelcore/db/model"
	"github.com/vjeantet/jodaTime"
)

const DATEPARSER string = "yMd'T'Hms"
const UTCMARKER string = "'Z'"

func parseDatesAndFill(json *gabs.Container, ics *model.ICS) error {

	//	if with TZID -> parse timezone, directly, otherwise...
	//	DTSTART;TZID=Europe/Vatican:20190823T180200
	//	DTEND;TZID=Europe/Vatican:20190823T190200

	// TODO respect the TZID
	childMap, err := json.Search("VCALENDAR", "VEVENT").ChildrenMap()
	if err != nil {
		return err
	}
	for key, child := range childMap {
		if strings.HasPrefix(key, "DTSTART") {

			format := DATEPARSER

			sStart := child.Data().(string)
			if sStart[len(sStart)-1:] == "Z" {
				format += UTCMARKER
			}

			start, err := jodaTime.Parse(format, sStart)
			if err != nil {
				fcdb.LogErrorFmt("Error when parsing start date '%s' from VEVENT: %s", sStart, err.Error())
				return err
			}

			ics.StartDate = start

		} else if strings.HasPrefix(key, "DTEND") {

			format := DATEPARSER

			sEnd := child.Data().(string)
			if sEnd[len(sEnd)-1:] == "Z" {
				format += UTCMARKER
			}

			end, err := jodaTime.Parse(format, sEnd)
			if err != nil {
				fcdb.LogErrorFmt("Error when parsing end date '%s' from VEVENT: %s", sEnd, err.Error())
				return err
			}

			ics.EndDate = end
		}
	}

	return nil
}

func ListIcsPerCal(calendar string) {

	db := fcdb.GetDB()

	var rows []*model.ICS

	db.Find(&rows)

	var ics [][]string

	for _, rec := range rows {

		ics = append(ics, []string{rec.Pkey, rec.CrtDat.Format("2006-01-02 15:04:05"), rec.UpdDat.Format("2006-01-02 15:04:05")})
	}

	fcdb.WriteTable([]string{"Id", "CrtDat", "UpdDat"}, ics)
}

func AddIcs(calId string, user string, calendar string, content string) (model.ICS, error) {

	json := fcdb.ParseICS(content)

	db := fcdb.GetDB()

	ics := model.ICS{Pkey: calId}

	ics.CalendarId = calendar

	err := parseDatesAndFill(json, &ics)
	if err != nil {
		return model.ICS{}, err
	}

	ics.Content = content

	retDB := db.Create(&ics)

	if retDB.Error != nil {
		log.Printf("Error with ICS %q: %s\n", calId, retDB.Error)
		return model.ICS{}, retDB.Error
	}

	fmt.Printf("ICS %s for user %s added.\n", calId, user)

	return ics, nil
}

func UpdateIcs(calId string, content string) (model.ICS, error) {

	db := fcdb.GetDB()

	ics, err := GetICS(calId)
	if err != nil {
		return model.ICS{}, err
	}

	// todo update changed fields only
	ics.Content = content
	ics.UpdDat = time.Now()

	json := fcdb.ParseICS(content)

	err = parseDatesAndFill(json, &ics)
	if err != nil {
		return model.ICS{}, err
	}

	retDB := db.Save(&ics)

	if retDB.Error != nil {
		log.Printf("Error with ICS %q: %s\n", calId, retDB.Error)
		return model.ICS{}, retDB.Error
	}

	fmt.Printf("ICS %s updated.\n", calId)

	return ics, nil
}

func GetICS(icsId string) (model.ICS, error) {

	db := fcdb.GetDB()

	var ics model.ICS
	retDB := db.First(&ics, "pkey = ?", icsId)

	if retDB.Error != nil {
		log.Printf("Error with loading ICS %q: %s\n", icsId, retDB.Error)
		return model.ICS{}, retDB.Error
	}

	return ics, nil
}

func FindIcsByCalendar(calID string) ([]*model.ICS, error) {

	var ics model.ICS

	db := fcdb.GetDB()
	db = db.Model(ics).Where("calendar_id = ?", calID)

	var rows []*model.ICS

	retDB := db.Find(&rows)

	if retDB.Error != nil {
		log.Printf("Error with loading ICS %s\n", retDB.Error)
		return rows, retDB.Error
	}

	return rows, nil
}

func FindIcsByTimeslot(calID string, start *time.Time, end *time.Time) ([]*model.ICS, error) {

	var ics model.ICS

	db := fcdb.GetDB()
	db = db.Model(ics)

	if len(calID) > 0 {

		db = db.Where("calendar_id = ?", calID)
	}

	if start != nil && !start.IsZero() {

		db = db.Where("start_date >= ?", start)
	}

	if end != nil && !end.IsZero() {

		db = db.Where("end_date <= ?", end)
	}

	var rows []*model.ICS

	retDB := db.Find(&rows)

	if retDB.Error != nil {
		log.Printf("Error with loading ICS %s\n", retDB.Error)
		return rows, retDB.Error
	}

	return rows, nil
}

func FindIcsInList(arrICS []string) ([]*model.ICS, error) {

	var ics model.ICS

	db := fcdb.GetDB()
	db = db.Model(ics).Where("pkey in (?)", arrICS)

	var rows []*model.ICS

	retDB := db.Find(&rows)

	if retDB.Error != nil {
		log.Printf("Error with loading ICS %s\n", retDB.Error)
		return rows, retDB.Error
	}

	return rows, nil
}

func DeleteIcs(icsId string) error {

	db := fcdb.GetDB()

	ics := &model.ICS{}

	retDB := db.Where("pkey = ?", icsId).First(&ics)

	if retDB.Error != nil {
		log.Printf("Error with Ics %q: %s\n", icsId, retDB.Error)
		log.Fatal(retDB.Error)
		return retDB.Error
	}

	if retDB.RowsAffected <= 0 {
		log.Printf("ICS not found: %s\n", icsId)
		log.Fatal("ICS not found: " + icsId + "\n")
		return retDB.Error
	}

	log.Printf("Deleting ICS: %s", &ics.Pkey)

	ret := db.Delete(&ics)

	return ret.Error
}
