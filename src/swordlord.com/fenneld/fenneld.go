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
)

func main() {

	// Initialise env and params
	fennelcore.InitConfig()

	// Initialise database
	// todo: make sure database is working as expected, chicken out otherwise
	fennelcore.InitDatabase()
	defer fennelcore.CloseDB()

	env := fennelcore.GetEnv()

	// TODO add authentication
	// n := negroni.New(negroni.NewRecovery(), negroni.HandlerFunc(auth.OAuthMiddleware), negroni.NewLogger())
	n := negroni.New(negroni.NewRecovery(), negroni.NewLogger())

	gr := mux.NewRouter().StrictSlash(false)

	//gr.NotFoundHandler
	//gr.MethodNotAllowedHandler

	n.UseHandler(gr)

	// TODO: parameters in URLs need regular expressions to make sure no unwanted char is used

	// what to do when a user hits the root
	// 	crossroads.addRoute('/', onHitRoot);
	gr.HandleFunc("/", handler.OnRoot).Methods("GET")

	// 	crossroads.addRoute('/.well-known/:params*:', onHitWellKnown);
	gr.HandleFunc("/.well-known/{param}", handler.OnWellKnown).Methods("GET")

	// crossroads.addRoute('/p/:params*:', onHitPrincipal);
	sr_p := gr.PathPrefix("/p").Subrouter()
	//sr_p.HandleFunc("", handler.onPrincipal).Methods("GET") -> should not happen?
	sr_p.HandleFunc("/", principal.Options).Methods("OPTIONS")
	sr_p.HandleFunc("", principal.Report).Methods("REPORT")
	sr_p.HandleFunc("", principal.Propfind).Methods("PROPFIND")
	sr_p.HandleFunc("", principal.Proppatch).Methods("PROPPATCH")

	// crossroads.addRoute('/cal/:username:/:cal:/:params*:', onHitCalendar);
	sr_cal := gr.PathPrefix("/cal").Subrouter()
	sr_cal.HandleFunc("/{user}/", calendar.Options).Methods("OPTIONS")
	sr_cal.HandleFunc("/{user}/{calendar}/", calendar.MakeCalendar).Methods("MKCALENDAR")
	sr_cal.HandleFunc("/{user}/{calendar}/{event}.ics", calendar.Put).Methods("PUT")
	sr_cal.HandleFunc("/{user}/{calendar}/{event}.ics", calendar.Get).Methods("GET")

	sr_cal.HandleFunc("/{user}/", calendar.PropfindUser).Methods("PROPFIND")
	sr_cal.HandleFunc("/{user}/inbox/", calendar.PropfindInbox).Methods("PROPFIND")
	sr_cal.HandleFunc("/{user}/outbox/", calendar.PropfindOutbox).Methods("PROPFIND")
	sr_cal.HandleFunc("/{user}/notifications/", calendar.PropfindNotification).Methods("PROPFIND")
	sr_cal.HandleFunc("/{user}/{calendar}/", calendar.PropfindCalendar).Methods("PROPFIND")

	sr_cal.HandleFunc("/{user}/{calendar}", calendar.Report).Methods("REPORT")

	// crossroads.addRoute('/card/:username:/:card:/:params*:', onHitCard);
	sr_card := gr.PathPrefix("/card").Subrouter()
	sr_card.HandleFunc("/{user}/{addressbook}/", addressbook.Options).Methods("OPTIONS")
	sr_card.HandleFunc("/{user}/{addressbook}/", addressbook.Report).Methods("REPORT")
	sr_card.HandleFunc("/{user}/{addressbook}/{card}.vcf", addressbook.Put).Methods("PUT")

	// crossroads.bypassed.add(onBypass); -> 404
	/*
	api.HandleFunc("/entries{ext:(?:.json)?}", handler.OnRetrieveEntries).Methods("GET")
	//api.HandleFunc("/entries.json", handler.OnRetrieveEntries).Methods("GET")
	api.HandleFunc("/entries{ext:(?:.json)?}", handler.OnCreateEntry).Methods("POST")
	//api.HandleFunc("/entries.json", handler.OnCreateEntry).Methods("POST")
	api.HandleFunc("/entries/{entry:[0-9]+}{ext:(?:.json)?}", handler.OnDeleteEntry).Methods("DELETE")
	api.HandleFunc("/entries/{entry:[0-9]+}/export{ext:(?:.json)?}", handler.OnGetEntryFormatted).Methods("GET")
	api.HandleFunc("/entries/{entry:[0-9]+}/tags/{tag:[0-9]+}{ext:(?:.json)?}", handler.OnDeleteTagOnEntry).Methods("DELETE")
	api.HandleFunc("/tags{ext:(?:.json|.txt|.xml)?}", handler.OnRetrieveAllTags).Methods("GET")
	api.HandleFunc("/version{ext:(?:.json|.txt|.xml|.html)?}", handler.OnRetrieveVersionNumber).Methods("GET")
	*/

	host := fennelcore.GetStringFromConfig("www.host")
	port := fennelcore.GetStringFromConfig("www.port")

	if env == "dev" {

		// give the user the possibility to trace and profile the app
		/*
		TODO RE ADD
		r.GET("/debug/pprof/block", pprofHandler(pprof.Index))
		r.GET("/debug/pprof/heap", pprofHandler(pprof.Index))
		r.GET("/debug/pprof/profile", pprofHandler(pprof.Profile))
		r.POST("/debug/pprof/symbol", pprofHandler(pprof.Symbol))
		r.GET("/debug/pprof/symbol", pprofHandler(pprof.Symbol))
		r.GET("/debug/pprof/trace", pprofHandler(pprof.Trace))
		*/

		// give the user some hints on what URLs she could test
		fmt.Printf("fenneld running on %v:%v\n", host, port)

		// TODO, fix URLs
		fmt.Printf("** get token  : curl -X POST 'http://%s:%s/oauth/v2/token' --data 'client_id=1&client_secret=secret&grant_type=password&password=pwd&username=uid' -H 'Content-Type:application/x-www-form-urlencoded'\n", host, port)
		fmt.Printf("** add entry  : curl -X POST 'http://%s:%s/api/entries/' --data 'url=http://test' -H 'Content-Type:application/x-www-form-urlencoded' -H 'Authorization: Bearer (access token)'\n", host, port)
		fmt.Printf("** get entries: curl -X GET 'http://%s:%s/api/entries/?page=1&perPage=20' -H 'Authorization: Bearer (access token)\n", host, port)
		fmt.Printf("** get entry  : curl -X GET 'http://%s:%s/api/entries/1' -H 'Authorization: Bearer (access token)\n", host, port)
		fmt.Printf("** patch entry: curl -X PATCH 'http://%s:%s/api/entries/1' --data 'archive=1&starred=1' -H 'Content-Type:application/x-www-form-urlencoded' -H 'Authorization: Bearer (access token)\n", host, port)

	}

	// have fun with fennel
	// http.ListenAndServe(host + ":" + port, n)
	n.Run(host + ":" + port)
}

/*
func pprofHandler(h http.HandlerFunc) negroni.HandlerFunc {
handler := http.HandlerFunc(h)
return func(c *gin.Context) {
	handler.ServeHTTP(c.Writer, c.Request)
}

}
*/
