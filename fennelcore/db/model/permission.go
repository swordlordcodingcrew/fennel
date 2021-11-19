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
	"gorm.io/gorm"
	"time"
)

type Permission struct {
	Pkey        string `gorm:"primary_key"`
	GroupId     string `sql:"NOT NULL"`
	Description string
	Permission  string
	CrtDat      time.Time `sql:"DEFAULT:current_timestamp"`
	UpdDat      time.Time `sql:"DEFAULT:current_timestamp"`
}

func (m *Permission) BeforeCreate(tx *gorm.DB) (err error) {
    m.CrtDat = time.Now()
    return nil
}

func (m *Permission) BeforeSave(scope *gorm.DB) (err error) {
    m.UpdDat = time.Now()
    return nil
}
