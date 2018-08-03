package auth
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
									"errors"
)


func ValidateCourier(uid string, pwd string) (error, string) {

	return errors.New("Method unknown"), ""

}

/*
function checkCourier(username, password, callback)
{
    log.debug("Authenticating user with courier method.");

    var socketPath = config.auth_method_courier_socket;
    log.debug("Using socket: " + socketPath);

    var client = net.createConnection({path: socketPath});

    client.on("connect", function() {
        //console.log('connect');
        var payload = 'service\nlogin\n' + username + '\n' + password;
        client.write('AUTH ' + payload.length + '\n' + payload);
    });

    var response = "";

    client.on("data", function(data) {
        //console.log('data: ' + data);
        response += data.toString();
    });

    client.on('end', function() {
        var result = response.indexOf('FAIL', 0);
        callback(result < 0);
    });
}
*/