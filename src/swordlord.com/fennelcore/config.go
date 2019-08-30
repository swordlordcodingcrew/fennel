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
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
)

func InitConfig() {

	// Note: Viper does not require any initialization before using, unless we'll be dealing multiple different configurations.
	// check [working with multiple vipers](https://github.com/spf13/viper#working-with-multiple-vipers)

	// Set config file we want to read. 2 ways to do this.
	// 1. Set config file path including file name and extension
	//viper.SetConfigFile("./configs/config.json")

	// OR
	// 2. Register path to look for config files in. It can accept multiple paths.
	// It will search these paths in given order
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME")
	// And then register config file name (no extension)
	viper.SetConfigName("fennel.config")
	// Optionally we can set specific config type
	viper.SetConfigType("json")

	// viper allows watching of config files for changes (and potential reloads)
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
	//	fmt.Println("Config file changed:", e.Name)
	})

	// Find and read the config file
	if err := viper.ReadInConfig(); err != nil {

		// TODO: don't just overwrite, check for existence first, then write a standard config file and move on...
		WriteStandardConfig()

		if err := viper.ReadInConfig(); err != nil {
			// we tried it once, crash now
				log.Fatalf("Error reading config file, %s", err)
		}
	}

	// Confirm which config file is used
	// fmt.Printf("Using config: %s\n", viper.ConfigFileUsed())

	// Confirm which config file is used
	// fmt.Printf("Env set to: %s\n", env)

	//EnsureTemplateFilesExist()
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

//
func WriteStandardConfig() (error) {

	err := ioutil.WriteFile("fennel.config.json", defaultConfig, 0700)

	return err
}

var defaultConfig = []byte(`
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
    "dialect": "sqlite3",
    "args": "fennel.db",
    "logmode": "true"
  }
}
`)

