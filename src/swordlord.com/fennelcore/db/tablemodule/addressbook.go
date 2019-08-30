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

func ListAddressbook() {

	db := fcdb.GetDB()

	var rows []*model.ADB

	db.Find(&rows)

	var adb [][]string

	for _, rec := range rows {

		adb = append(adb, []string{rec.Pkey, rec.CrtDat.Format("2006-01-02 15:04:05"), rec.UpdDat.Format("2006-01-02 15:04:05")})
	}

	//wombag.WriteTable([]string{"Id", "CrtDat", "UpdDat"}, adb)
}

func AddAddressbook(name string, password string, user string) (model.ADB, error) {

	db := fcdb.GetDB()

	_, err := hashPassword(password)
	if err != nil {
		log.Printf("Error with hashing password %q: %s\n", password, err )
		return model.ADB{}, err
	}

	adb := model.ADB{Pkey: name}
	retDB := db.Create(&adb)

	if retDB.Error != nil {
		log.Printf("Error with Device %q: %s\n", name, retDB.Error)
		log.Fatal(retDB.Error)
		return model.ADB{}, retDB.Error
	}

	fmt.Printf("Device %s for user %s added.\n", name, user)

	return adb, nil
}

func GetAddressbooksFromUser(user string) (error, []*model.ADB) {

	var adb model.ADB

	db := fcdb.GetDB()
	db = db.Model(adb)

	db = db.Where("owner = ?", user)

	var rows []*model.ADB

	retDB := db.Find(&rows)

	if retDB.Error != nil {
		log.Printf("Error with ADB from User %q: %s\n", user, retDB.Error)
		return retDB.Error, rows
	}

	return nil, rows
}

func GetAddressbookByName(name string) (error, *model.ADB) {

	var adb model.ADB

	db := fcdb.GetDB()
	db = db.Model(adb)

	db = db.Where("name = ?", name)

	var row *model.ADB

	retDB := db.First(&row)

	if retDB.Error != nil {
		log.Printf("Error with ADB by name %q: %s\n", name, retDB.Error)
		return retDB.Error, row
	}

	return nil, row
}

func GetOrCreateAddressbookByName(name string, owner string) (error, *model.ADB) {

	var adb model.ADB

	db := fcdb.GetDB()
	db = db.Model(adb)

	db = db.Where(model.ADB{Name: name}).Attrs(model.ADB{Owner: owner})

	var row *model.ADB

	retDB := db.FirstOrCreate(&row)

	if retDB.Error != nil {
		log.Printf("Error with ADB by name %q: %s\n", name, retDB.Error)
		return retDB.Error, row
	}

	return nil, row
}

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
}

func DeleteAddressbook(name string) {

	db := fcdb.GetDB()

	rec := &model.ADB{}

	retDB := db.Where("id = ?", name).First(&rec)

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

	log.Printf("Deleting Device: %s", &rec.Pkey)

	db.Delete(&rec)

	fmt.Printf("Device %s deleted.\n", name)
}
