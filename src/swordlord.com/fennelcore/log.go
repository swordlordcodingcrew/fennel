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
	log "github.com/sirupsen/logrus"
)

func InitLog() {

	level, err := log.ParseLevel(GetLogLevel())
	if err != nil {

		log.SetLevel(log.WarnLevel)

	} else {

		log.SetLevel(level)
	}
}

func LogTrace(msg string, fields log.Fields) {

	if fields == nil {

		log.Trace(msg)

	} else {
		log.WithFields(fields).Trace(msg)
	}
}

func LogDebug(msg string, fields log.Fields) {

	if fields == nil {

		log.Debug(msg)

	} else {
		log.WithFields(fields).Debug(msg)
	}
}

func LogDebugFmt(err string, a ...interface{}) {

	log.Debugf(err, a...)
}

func LogInfo(msg string, fields log.Fields) {

	if fields == nil {

		log.Info(msg)

	} else {
		log.WithFields(fields).Info(msg)
	}
}

func LogInfoFmt(err string, a ...interface{}) {

	log.Infof(err, a...)
}

func LogWarn(msg string, fields log.Fields) {

	if fields == nil {

		log.Warn(msg)

	} else {
		log.WithFields(fields).Warn(msg)
	}
}

func LogError(msg string, fields log.Fields) {

	if fields == nil {

		log.Error(msg)

	} else {
		log.WithFields(fields).Error(msg)
	}
}

func LogErrorFmt(err string, a ...interface{}) {

	log.Errorf(err, a...)
}


func LogFatal(msg string, fields log.Fields) {

	if fields == nil {

		log.Fatal(msg)

	} else {
		log.WithFields(fields).Fatal(msg)
	}
}

func LogPanic(msg string, fields log.Fields) {

	if fields == nil {

		log.Panic(msg)

	} else {
		log.WithFields(fields).Panic(msg)
	}
}
