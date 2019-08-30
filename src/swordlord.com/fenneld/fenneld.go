package main
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
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"github.com/urfave/negroni"
	"net/http/pprof"
	"swordlord.com/fennelcore"
	"swordlord.com/fenneld/auth"
	"swordlord.com/fenneld/handler"
	"swordlord.com/fenneld/handler/addressbook"
	"swordlord.com/fenneld/handler/calendar"
	"swordlord.com/fenneld/handler/principal"
)

func main() {

	// Initialise env and params
	fennelcore.InitConfig()
	fennelcore.InitLog()

	// Initialise database
	// if there is an error, this function will quit the app
	fennelcore.InitDatabase()
	defer fennelcore.CloseDB()

	logLevel := fennelcore.GetLogLevel()

	// TODO write our own logger using logrus
	n := negroni.New(negroni.NewRecovery(), negroni.NewLogger(), auth.NewFennelAuthentication())

	gr := mux.NewRouter().StrictSlash(false)

	n.UseHandler(gr)

	// TODO add these handlers
	//gr.NotFoundHandler
	//gr.MethodNotAllowedHandler

	// what to do when a user hits the root
	gr.HandleFunc("/", handler.OnRoot).Methods("GET")
	gr.HandleFunc("/", handler.OnRoot).Methods("PROPFIND")

	// ******************* SERVICE DISCOVERY
	gr.HandleFunc("/.well-known", handler.OnWellKnownNoParam).Methods("GET", "PROPFIND")
	gr.HandleFunc("/.well-known/", handler.OnWellKnownNoParam).Methods("GET", "PROPFIND")
	gr.HandleFunc("/.well-known/{param:[0-9a-zA-Z-]+}", handler.OnWellKnown).Methods("GET", "PROPFIND")
	gr.HandleFunc("/.well-known/{param:[0-9a-zA-Z-]+}/", handler.OnWellKnown).Methods("GET", "PROPFIND")

	// ******************* PRINCIPAL
	srP := gr.PathPrefix("/p").Subrouter()
	//srP.HandleFunc("", handler.onPrincipal).Methods("GET") -> should not happen?
	srP.HandleFunc("/", principal.Options).Methods("OPTIONS")
	srP.HandleFunc("/", principal.Report).Methods("REPORT")
	srP.HandleFunc("/", principal.Propfind).Methods("PROPFIND")
	srP.HandleFunc("/{user:[0-9a-zA-Z-]+}/", principal.Propfind).Methods("PROPFIND")
	srP.HandleFunc("/{user:[0-9a-zA-Z-]+}/", principal.Options).Methods("OPTIONS")
	srP.HandleFunc("", principal.Proppatch).Methods("PROPPATCH")

	// ******************* CALENDAR
	srCal := gr.PathPrefix("/cal").Subrouter()
	srCal.HandleFunc("/{user:[0-9a-zA-Z-]+}/", calendar.PropfindRoot).Methods("PROPFIND")
	srCal.HandleFunc("/{user:[0-9a-zA-Z-]+}/", calendar.Options).Methods("OPTIONS")
	srCal.HandleFunc("/{user:[0-9a-zA-Z-]+}/{calendar:[0-9a-zA-Z-]+}/", calendar.MakeCalendar).Methods("MKCALENDAR")
	srCal.HandleFunc("/{user:[0-9a-zA-Z-]+}/{calendar:[0-9a-zA-Z-]+}/{event:[0-9a-zA-Z-]+}.ics", calendar.Put).Methods("PUT")
	srCal.HandleFunc("/{user:[0-9a-zA-Z-]+}/{calendar:[0-9a-zA-Z-]+}/{event:[0-9a-zA-Z-]+}.ics", calendar.Get).Methods("GET")

	srCal.HandleFunc("/{user:[0-9a-zA-Z-]+}/", calendar.PropfindUser).Methods("PROPFIND")
	srCal.HandleFunc("/{user:[0-9a-zA-Z-]+}/inbox/", calendar.PropfindInbox).Methods("PROPFIND")
	srCal.HandleFunc("/{user:[0-9a-zA-Z-]+}/outbox/", calendar.PropfindOutbox).Methods("PROPFIND")
	srCal.HandleFunc("/{user:[0-9a-zA-Z-]+}/notifications/", calendar.PropfindNotification).Methods("PROPFIND")
	srCal.HandleFunc("/{user:[0-9a-zA-Z-]+}/{calendar:[0-9a-zA-Z-]+}/", calendar.PropfindCalendar).Methods("PROPFIND")

	srCal.HandleFunc("/{user:[0-9a-zA-Z-]+}/{calendar:[0-9a-zA-Z-]+}/", calendar.Report).Methods("REPORT")

	// ******************* ADDRESSBOOK
	srCard := gr.PathPrefix("/card").Subrouter()
	srCard.HandleFunc("/", addressbook.PropfindRoot).Methods("PROPFIND") //todo find out when this happens...
	srCard.HandleFunc("/{user:[0-9a-zA-Z-]+}/", addressbook.PropfindUser).Methods("PROPFIND")
	//srCard.HandleFunc("/{user:[0-9a-zA-Z-]+}/{addressbook:[0-9a-zA-Z-]+}/", addressbook.PropfindAddressbook).Methods("PROPFIND")
	srCard.HandleFunc("/{user:[0-9a-zA-Z-]+}/{addressbook:[0-9a-zA-Z-]+}/", addressbook.Options).Methods("OPTIONS")
	srCard.HandleFunc("/{user:[0-9a-zA-Z-]+}/{addressbook:[0-9a-zA-Z-]+}/", addressbook.Report).Methods("REPORT")
	srCard.HandleFunc("/{user:[0-9a-zA-Z-]+}/{addressbook:[0-9a-zA-Z-]+}/{card:[0-9a-zA-Z-]+}.vcf", addressbook.Put).Methods("PUT")

	// get settings
	host := fennelcore.GetStringFromConfig("www.host")
	port := fennelcore.GetStringFromConfig("www.port")

	// check if user wants to mount debug urls
	if logLevel == "debug" {

		// give the user the possibility to trace and profile the app
		srDebug := gr.PathPrefix("/debug/pprof").Subrouter()
		srDebug.HandleFunc("/block", pprof.Index).Methods("GET")
		srDebug.HandleFunc("/heap", pprof.Index).Methods("GET")
		srDebug.HandleFunc("/profile", pprof.Profile).Methods("GET")
		srDebug.HandleFunc("/symbol", pprof.Symbol).Methods("POST")
		srDebug.HandleFunc("/symbol", pprof.Symbol).Methods("GET")
		srDebug.HandleFunc("/trace", pprof.Trace).Methods("GET")

		// give the user some hints on what URLs she could test
		fennelcore.LogDebugFmt("get options : curl -X OPTIONS 'http://%s:%s/cal/demo/'", host, port)
		fennelcore.LogDebugFmt("get block: go tool pprof 'http://%s:%s/debug/pprof/block'", host, port)
		fennelcore.LogDebugFmt("get heap: go tool pprof 'http://%s:%s/debug/pprof/heap'", host, port)
		fennelcore.LogDebugFmt("get profile: go tool pprof 'http://%s:%s/debug/pprof/profile'", host, port)
		fennelcore.LogDebugFmt("post symbol: go tool pprof 'http://%s:%s/debug/pprof/symbol'", host, port)
		fennelcore.LogDebugFmt("get symbol: go tool pprof 'http://%s:%s/debug/pprof/symbol'", host, port)
		fennelcore.LogDebugFmt("get trace: go tool pprof 'http://%s:%s/debug/pprof/trace'", host, port)
	}

	fennelcore.LogInfoFmt("fenneld running on %v:%v.", host, port)

	// have fun with fennel
	n.Run(host + ":" + port)
}
