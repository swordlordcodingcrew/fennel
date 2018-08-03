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

func ListVCardsPerAddressbook(addressbook string) {

	var vcard [][]string

	err, rows := FindVcardByAddressbook(addressbook)
	if err != nil {

		log.Printf("Error with VCARD in Addressbook %q: %s\n", addressbook, err)
		return
	}

	for _, rec := range rows {

		vcard = append(vcard, []string{rec.Pkey, rec.CrtDat.Format("2006-01-02 15:04:05"), rec.UpdDat.Format("2006-01-02 15:04:05")})
	}

	fcdb.WriteTable([]string{"Id", "CrtDat", "UpdDat"}, vcard)
}

func AddVCard(vcardId string, owner string, addressbookId string, isGroup bool, content string) (model.VCARD, error) {

	db := fcdb.GetDB()

	vcard := model.VCARD{Pkey: vcardId}

	vcard.AddressbookId = addressbookId
	vcard.Owner = owner
	vcard.Content = content
	vcard.IsGroup = isGroup

	retDB := db.Create(&vcard)

	if retDB.Error != nil {
		log.Printf("Error with VCARD %q: %s\n", vcardId, retDB.Error)
		return model.VCARD{}, retDB.Error
	}

	fmt.Printf("VCARD %s for owner %s added.\n", vcardId, owner)

	return vcard, nil
}

func UpdateVCard(name string, password string) error {

	return nil
}

func GetVCard(vcardId string) (model.VCARD, error) {

	db := fcdb.GetDB()

	var vcard model.VCARD
	retDB := db.First(&vcard, "pkey = ?", vcardId)

	if retDB.Error != nil {
		log.Printf("Error with loading VCARD %q: %s\n", vcardId, retDB.Error)
		return model.VCARD{}, retDB.Error
	}

	return vcard, nil
}

func FindVcardByAddressbook(addressbookID string) (error, []*model.VCARD)  {

	var vcard model.VCARD

	db := fcdb.GetDB()
	db = db.Model(vcard).Where("addressbook_id = ?", addressbookID)

	var rows []*model.VCARD

	retDB := db.Find(&rows)

	if retDB.Error != nil {
		log.Printf("Error with loading VCARD %s\n", retDB.Error)
		return retDB.Error, rows
	}

	return nil, rows
}

func FindVCardsFromAddressbook(adbID string, vcardIDs []string) (error, []*model.VCARD) {

	var vcard model.VCARD

	db := fcdb.GetDB()
	db = db.Model(vcard)

	db = db.Where("pkey in (?)", vcardIDs).Where("addressbook_id = ?", adbID)

	var rows []*model.VCARD

	retDB := db.Find(&rows)

	if retDB.Error != nil {
		log.Printf("Error with VCARD from Addressbook %s: %s\n", adbID, retDB.Error)
		return retDB.Error, rows
	}

	return nil, rows
}


func DeleteVCard(vcardId string) error {

	db := fcdb.GetDB()

	vcard := &model.VCARD{}

	retDB := db.Where("pkey = ?", vcardId).First(&vcard)

	if retDB.Error != nil {
		log.Printf("Error with VCARD %q: %s\n", vcardId, retDB.Error)
		log.Fatal(retDB.Error)
		return retDB.Error
	}

	if retDB.RowsAffected <= 0 {
		log.Printf("VCARD not found: %s\n", vcardId)
		log.Fatal("VCARD not found: " + vcardId + "\n")
		return retDB.Error
	}

	log.Printf("Deleting VCARD: %s", &vcard.Pkey)

	ret := db.Delete(&vcard)

	return ret.Error
}
