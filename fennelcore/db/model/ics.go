package model
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
	"github.com/jinzhu/gorm"
	"time"
)

type ICS struct {
	Pkey    	string `gorm:"primary_key"`
	CalendarId  string `sql:"NOT NULL"`
	StartDate	time.Time `sql:"NOT NULL; DEFAULT:current_timestamp"`
	EndDate		time.Time `sql:"NOT NULL; DEFAULT:current_timestamp"`
	Content		string `sql:"type:blob"`
	CrtDat		time.Time `sql:"DEFAULT:current_timestamp"`
	UpdDat		time.Time `sql:"DEFAULT:current_timestamp"`
}

func (m *ICS) BeforeUpdate(scope *gorm.Scope) (err error) {

	scope.SetColumn("UpdDat", time.Now())
	return  nil
}

/*
func (u *User) BeforeSave(scope *gorm.Scope) (err error) {

	scope.SetColumn("upddat", time.Now())
	return nil
}
*/
