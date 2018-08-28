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
)

func ListCal() {

	db := fcdb.GetDB()

	var rows []*model.CAL

	db.Find(&rows)

	var cal [][]string

	for _, rec := range rows {

		cal = append(cal, []string{rec.Pkey, rec.CrtDat.Format("2006-01-02 15:04:05"), rec.UpdDat.Format("2006-01-02 15:04:05")})
	}

	fcdb.WriteTable([]string{"Id", "CrtDat", "UpdDat"}, cal)
}

func AddCal(user string, calId string, displayname string, colour string, freebusyset string, order int, supportedCalComponent string, synctoken int, timezone string) (model.CAL, error) {

	db := fcdb.GetDB()

	cal := model.CAL{Pkey: calId}

	cal.Owner = user
	cal.Displayname = displayname
	cal.Colour = colour
	cal.FreeBusySet = freebusyset
	cal.Order = order
	cal.SupportedCalComponent = supportedCalComponent
	cal.Synctoken = synctoken
	cal.Timezone = timezone

	retDB := db.Create(&cal)

	if retDB.Error != nil {
		log.Printf("Error with CAL %q: %s\n", calId, retDB.Error)
		return model.CAL{}, retDB.Error
	}

	fmt.Printf("CAL %s for user %s added.\n", calId, user)

	return cal, nil
}

func GetCal(calId string) (model.CAL, error) {

	db := fcdb.GetDB()

	var cal model.CAL
	retDB := db.First(&cal, "pkey = ?", calId)


	//	retDB := db.Model(&model.CAL{}).Where("pkey=?", calId)
	//retDB := db.Where("Pkey = ?", calId).First(&model.CAL{})
	if retDB.Error != nil {
		log.Printf("Error with loading CAL %q: %s\n", calId, retDB.Error)
		return model.CAL{}, retDB.Error
	}

	return cal, nil
}

/*
func UpdateAddressbook(name string, password string) error {

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
} */

func DeleteCal(name string) {

	db := fcdb.GetDB()

	cal := &model.CAL{}

	retDB := db.Where("id = ?", name).First(&cal)

	if retDB.Error != nil {
		log.Printf("Error with Device %q: %s\n", name, retDB.Error)
		log.Fatal(retDB.Error)
		return
	}

	if retDB.RowsAffected <= 0 {
		log.Printf("Device not found: %s\n", name)
		log.Fatal("Device not found: " + name + "\n")
		return
	}

	log.Printf("Deleting Calendar: %s", &cal.Pkey)

	db.Delete(&cal)

	fmt.Printf("Calendar %s deleted.\n", name)
}
