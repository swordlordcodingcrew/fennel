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
	"bytes"
	"flag"
	"log"

	"github.com/spf13/viper"
)

var defaultConfig = `
{
  "log": {
    "level": "debug"
  },
  "www": {
      "host": "127.0.0.1",
      "port": "8888"
  },
  "auth": {
    "module": "htpasswd",
    "file": "demouser.htpasswd"
  },
  "folder": {
    "templates": "templates"
  },
  "db": {
    "dialect": "sqlite",
    "args": "fennel.db",
    "logmode": "true"
  }
}
`

func InitConfig() {
	configOverride := flag.String("config", "", "Configuration file path (Optional, will read from standard config locations otherwise)")
	flag.Parse()
	if *configOverride != "" {
		viper.SetConfigFile(*configOverride)
	} else {
		viper.SetConfigName("fennel")
		viper.AddConfigPath(".")
		viper.AddConfigPath("$HOME/.config/fennel")
		viper.AddConfigPath("$HOME/.config")
		viper.AddConfigPath("/etc/fennel")
		viper.AddConfigPath("/etc")
	}

	// Find and read the config file
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Error reading config %v\n Falling back to default %v", err, defaultConfig)
		viper.SetConfigType("json")
		err = viper.ReadConfig(bytes.NewBuffer([]byte(defaultConfig)))
		if err != nil {
			log.Fatal(err)
		}
	}

	// Confirm which config file is used
	log.Printf("config read from: %s\n", viper.ConfigFileUsed())
}

func GetBoolFromConfig(key string) bool {
	return viper.GetBool(key)
}

func GetStringFromConfig(key string) string {
	return viper.GetString(key)
}

func GetLogLevel() string {
	loglevel := viper.GetString("log.level")
	if loglevel == "" {
		return "warn"
	} else {

		return loglevel
	}
}

