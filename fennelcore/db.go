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
	"log"

	"github.com/swordlordcodingcrew/fennel/fennelcore/db/model"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var db gorm.DB

//
func InitDatabase() {
	dialect := GetStringFromConfig("db.dialect")
	args := GetStringFromConfig("db.args")
	var dialector gorm.Dialector
	switch dialect {
	case "sqlite":
		dialector = sqlite.Open(args)
	case "postgres":
		dialector = postgres.Open(args)
	case "mysql":
		dialector = mysql.Open(args)
	default:
		log.Fatalf("Unsupported database dialect %v", dialect)
	}
	
	database, err := gorm.Open(dialector, &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: "fen_", // Avoid naming collisions with reserved tables like 'user'
			SingularTable: true,
		},
	})
	if err != nil {
		log.Fatalf("failed to connect database, %s", err)
		panic("failed to connect database")
	}

	db = *database
	db.AutoMigrate(&model.User{})
	db.AutoMigrate(&model.Group{})
	db.AutoMigrate(&model.UserGroup{})
	db.AutoMigrate(&model.Permission{})
	db.AutoMigrate(&model.CAL{})
	db.AutoMigrate(&model.ADB{})
	db.AutoMigrate(&model.ICS{})
	db.AutoMigrate(&model.VCARD{})
}

//
func CloseDB() {
	// db.Close()
}

//
func GetDB() *gorm.DB {
	return &db
}
