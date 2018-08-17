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
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"github.com/urfave/negroni"
	"swordlord.com/fennelcore"
	"swordlord.com/fenneld/handler"
	"swordlord.com/fenneld/handler/principal"
	"swordlord.com/fenneld/handler/calendar"
	"swordlord.com/fenneld/handler/addressbook"
	"net/http/pprof"
	"swordlord.com/fenneld/auth"
)

func main() {

	// Initialise env and params
	fennelcore.InitConfig()

	// Initialise database
	// if there is an error, this function will quit the app
	fennelcore.InitDatabase()
	defer fennelcore.CloseDB()

	env := fennelcore.GetEnv()

	n := negroni.New(negroni.NewRecovery(), negroni.NewLogger(), auth.NewFennelAuthentication())

	gr := mux.NewRouter().StrictSlash(false)

	n.UseHandler(gr)

	// TODO add these handlers
	//gr.NotFoundHandler
	//gr.MethodNotAllowedHandler

	// what to do when a user hits the root
	gr.HandleFunc("/", handler.OnRoot).Methods("GET")

	// ******************* SERVICE DISCOVERY
	gr.HandleFunc("/.well-known", handler.OnWellKnownNoParam).Methods("GET")
	gr.HandleFunc("/.well-known/{param:[0-9a-zA-Z-]+}", handler.OnWellKnown).Methods("GET")

	// ******************* PRINCIPAL
	sr_p := gr.PathPrefix("/p").Subrouter()
	//sr_p.HandleFunc("", handler.onPrincipal).Methods("GET") -> should not happen?
	sr_p.HandleFunc("/", principal.Options).Methods("OPTIONS")
	sr_p.HandleFunc("/", principal.Report).Methods("REPORT")
	sr_p.HandleFunc("/{user:[0-9a-zA-Z-]+}/", principal.Propfind).Methods("PROPFIND")
	sr_p.HandleFunc("", principal.Proppatch).Methods("PROPPATCH")

	// ******************* CALENDAR
	sr_cal := gr.PathPrefix("/cal").Subrouter()
	sr_cal.HandleFunc("/{user:[0-9a-zA-Z-]+}/", calendar.Options).Methods("OPTIONS")
	sr_cal.HandleFunc("/{user:[0-9a-zA-Z-]+}/{calendar:[0-9a-zA-Z-]+}/", calendar.MakeCalendar).Methods("MKCALENDAR")
	sr_cal.HandleFunc("/{user:[0-9a-zA-Z-]+}/{calendar:[0-9a-zA-Z-]+}/{event:[0-9a-zA-Z-]+}.ics", calendar.Put).Methods("PUT")
	sr_cal.HandleFunc("/{user:[0-9a-zA-Z-]+}/{calendar:[0-9a-zA-Z-]+}/{event:[0-9a-zA-Z-]+}.ics", calendar.Get).Methods("GET")

	sr_cal.HandleFunc("/{user:[0-9a-zA-Z-]+}/", calendar.PropfindUser).Methods("PROPFIND")
	sr_cal.HandleFunc("/{user:[0-9a-zA-Z-]+}/inbox/", calendar.PropfindInbox).Methods("PROPFIND")
	sr_cal.HandleFunc("/{user:[0-9a-zA-Z-]+}/outbox/", calendar.PropfindOutbox).Methods("PROPFIND")
	sr_cal.HandleFunc("/{user:[0-9a-zA-Z-]+}/notifications/", calendar.PropfindNotification).Methods("PROPFIND")
	sr_cal.HandleFunc("/{user:[0-9a-zA-Z-]+}/{calendar:[0-9a-zA-Z-]+}/", calendar.PropfindCalendar).Methods("PROPFIND")

	sr_cal.HandleFunc("/{user:[0-9a-zA-Z-]+}/{calendar:[0-9a-zA-Z-]+}", calendar.Report).Methods("REPORT")

	// ******************* ADDRESSBOOK
	sr_card := gr.PathPrefix("/card").Subrouter()
	sr_card.HandleFunc("/{user:[0-9a-zA-Z-]+}/", addressbook.Propfind).Methods("PROPFIND")
	sr_card.HandleFunc("/{user:[0-9a-zA-Z-]+}/{addressbook:[0-9a-zA-Z-]+}/", addressbook.Options).Methods("OPTIONS")
	sr_card.HandleFunc("/{user:[0-9a-zA-Z-]+}/{addressbook:[0-9a-zA-Z-]+}/", addressbook.Report).Methods("REPORT")
	sr_card.HandleFunc("/{user:[0-9a-zA-Z-]+}/{addressbook:[0-9a-zA-Z-]+}/{card:[0-9a-zA-Z-]+}.vcf", addressbook.Put).Methods("PUT")

	// get settings
	host := fennelcore.GetStringFromConfig("www.host")
	port := fennelcore.GetStringFromConfig("www.port")

	// check if user wants to mount debug urls
	if env == "dev" {

		// give the user the possibility to trace and profile the app
		sr_debug := gr.PathPrefix("/debug/pprof").Subrouter()
		sr_debug.HandleFunc("/block", pprof.Index).Methods("GET")
		sr_debug.HandleFunc("/heap", pprof.Index).Methods("GET")
		sr_debug.HandleFunc("/profile", pprof.Profile).Methods("GET")
		sr_debug.HandleFunc("/symbol", pprof.Symbol).Methods("POST")
		sr_debug.HandleFunc("/symbol", pprof.Symbol).Methods("GET")
		sr_debug.HandleFunc("/trace", pprof.Trace).Methods("GET")

		// give the user some hints on what URLs she could test
		fmt.Printf("fenneld running on %v:%v\n", host, port)

		fmt.Printf("** get  options : curl -X OPTIONS 'http://%s:%s/cal/demo/'\n\n", host, port)

		fmt.Printf("** get  block  	: go tool pprof 'http://%s:%s/debug/pprof/block'\n", host, port)
		fmt.Printf("** get  heap  	: go tool pprof 'http://%s:%s/debug/pprof/heap'\n", host, port)
		fmt.Printf("** get  profile : go tool pprof 'http://%s:%s/debug/pprof/profile'\n", host, port)
		fmt.Printf("** post symbol  : go tool pprof 'http://%s:%s/debug/pprof/symbol'\n", host, port)
		fmt.Printf("** get  symbol  : go tool pprof 'http://%s:%s/debug/pprof/symbol'\n", host, port)
		fmt.Printf("** get  trace  	: go tool pprof 'http://%s:%s/debug/pprof/trace'\n", host, port)
	}

	// have fun with fennel
	// http.ListenAndServe(host + ":" + port, n)
	n.Run(host + ":" + port)
}
