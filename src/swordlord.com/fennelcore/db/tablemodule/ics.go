package tablemodule
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
	"fmt"
	"log"
	fcdb "swordlord.com/fennelcore"
	"swordlord.com/fennelcore/db/model"
	"time"
	"github.com/vjeantet/jodaTime"
)

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

	var dtmStart *time.Time
	var dtmEnd *time.Time

	json := fcdb.ParseICS(content)

	sStart, ok := json.Path("VCALENDAR.VEVENT.DTSTART").Data().(string)
	if ok {

		start, err := jodaTime.Parse("yMd'T'Hms'Z'", sStart)
		if err == nil {

			dtmStart = &start
		}
	}

	sEnd, ok := json.Path("VCALENDAR.VEVENT.DTEND").Data().(string)
	if ok {

		end, err := jodaTime.Parse("yMd'T'Hms'Z'", sEnd)
		if err == nil {

			dtmEnd = &end
		}
	}
	// value == 10.0, ok == true

	return AddIcsParsed(calId, user, calendar, dtmStart, dtmEnd, content)
}

func AddIcsParsed(calId string, user string, calendar string, start *time.Time, end *time.Time, content string) (model.ICS, error) {

	db := fcdb.GetDB()

	ics := model.ICS{Pkey: calId}

	ics.CalendarId = calendar

	if start != nil {
		ics.StartDate = *start
	}
	if end != nil {
		ics.EndDate = *end
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

func UpdateIcs(name string, password string) error {

	db := fcdb.GetDB()

	pwd, err := hashPassword(password)
	if err != nil {
		log.Printf("Error with hashing password %q: %s\n", password, err )
		return err
	}

	retDB := db.Model(&model.ADB{}).Where("Id=?", name).Update("Token", pwd)

	if retDB.Error != nil {
		log.Printf("Error with Device %q: %s\n", name, retDB.Error)
		return retDB.Error
	}

	fmt.Printf("Device %s updated.\n", name)

	return nil
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


func FindIcsByCalendar(calID string) ([]*model.ICS, error)  {

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

func FindIcsByTimeslot(calID string, start *time.Time, end *time.Time) ([]*model.ICS, error)  {

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

func FindIcsInList(arrICS []string) ([]*model.ICS, error)  {

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
