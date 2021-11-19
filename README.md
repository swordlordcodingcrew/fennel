Fennel
======

![Fennel](https://raw.github.com/swordlordcodingcrew/fennel/master/fennel_logo.png)

**Fennel** (c) 2014-19 by [SwordLord - the coding crew](http://www.swordlord.com/)

## Introduction ##

**Fennel** is a lightweight CardDAV / CalDAV server. It is written in Go and based on the proof of concept [Fennel.js](https://github.com/LordEidi/fennel.js) (which is written in JavaScript and running on NodeJS).

If you are looking for a lightweight CalDAV / CardDAV server, **Fennel** might be for you:

- hassle free installation. Drop a binary, start it, that's it.
- authentication is meant to be pluggable. While we concentrate on CourierAuth and .htaccess, you can add whatever can check a username and a password.
- authorisation is meant to be pluggable as well.
- the data storage backend is meant to be pluggable as well. While we start with SQLite3, we do use an ORM. Whatever database can be used with **Gorm** can be used as storage backend for **Fennel**.
- and after all, **Fennel** is OSS and is written in Go. Whatever you do not like, you are free to replace / rewrite. Just respect the licence and give back.

## Status ##

![Build Status](https://travis-ci.org/swordlordcodingcrew/fennel.svg?branch=master)

**Fennel** is beta software and should be handled as such:

- The CalDAV part is still work in progress.
- The CardDAV part is still work in progress.

**Fennel** is tested on Calendar on iOS > v10.0 and on OSX Calendar as well as with Mozilla Lightning. If you run
**Fennel** with another client your mileage may vary.

What's missing:

- different clients (we will somewhen test with other clients, but we did not do thoroughly yet)
- Test cases for everything. We would love to have test cases for as many scenarios and features as possible. It is a pain in the neck to test **Fennel** otherwise. If you wonder how we test, have a look at the testing code in Fennel.js as well as at the [Project Spoon](https://github.com/swordlordcodingcrew/spoon)
- While **Fennel**'s goal is to have an RBAC based authorisation system, **Fennel** does currently only know global permissions without groups. But we are working on it.

## Installation ##

### From source ###

Dependencies: [golang](https://golang.org/dl/), [GNU Make](https://www.gnu.org/software/make/)

```
git clone https://github.com/swordlordcodingcrew/fennel
cd fennel
make
```
Executables fenneld and fennelcli are built into the bin/ folder.

### How to set up transport security ###

Since **Fennel** does not bring it's own crypto, you may need to install a TLS server in front of **Fennel**. You can do so
with nginx, which is a lightweight http server and proxy.

First prepare your /etc/apt/sources.list file (or just install the standard Debian package, your choice):

    deb http://nginx.org/packages/debian/ stretch nginx
    deb-src http://nginx.org/packages/debian/ stretch nginx

Update apt-cache and install nginx to your system.

    sudo update
    sudo apt-get install nginx

Now configure a proxy configuration so that your instance of nginx will serve / prox the content of / for the
**Fennel** server. To do so, you will need a configuration along this example:

    server {
        listen   443;
        server_name  fennel.yourdomain.tld;

        access_log  /var/www/logs/fennel_access.log combined;
        error_log  /var/www/logs/fennel_error.log;

        root /var/www/pages/;
        index  index.html index.htm;

        error_page   500 502 503 504  /50x.html;
        location = /50x.html {
            root   /var/www/nginx-default;
        }

        location / {
            proxy_pass         http://127.0.0.1:8888;
            proxy_redirect     off;
            proxy_set_header   Host             $host;
            proxy_set_header   X-Real-IP        $remote_addr;
            proxy_set_header   X-Forwarded-For  $proxy_add_x_forwarded_for;
            proxy_buffering    off;
        }

        ssl  on;
        ssl_certificate  /etc/nginx/certs/yourdomain.tld.pem;
        ssl_certificate_key  /etc/nginx/certs/yourdomain.tld.pem;
        ssl_session_timeout  5m;

        # modern configuration. tweak to your needs.
        ssl_protocols TLSv1.1 TLSv1.2;
        ssl_ciphers 'ECDHE-RSA-AES128-GCM-SHA256:ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-RSA-AES256-GCM-SHA384:ECDHE-ECDSA-AES256-GCM-SHA384:DHE-RSA-AES128-GCM-SHA256:DHE-DSS-AES128-GCM-SHA256:kEDH+AESGCM:ECDHE-RSA-AES128-SHA256:ECDHE-ECDSA-AES128-SHA256:ECDHE-RSA-AES128-SHA:ECDHE-ECDSA-AES128-SHA:ECDHE-RSA-AES256-SHA384:ECDHE-ECDSA-AES256-SHA384:ECDHE-RSA-AES256-SHA:ECDHE-ECDSA-AES256-SHA:DHE-RSA-AES128-SHA256:DHE-RSA-AES128-SHA:DHE-DSS-AES128-SHA256:DHE-RSA-AES256-SHA256:DHE-DSS-AES256-SHA:DHE-RSA-AES256-SHA:!aNULL:!eNULL:!EXPORT:!DES:!RC4:!3DES:!MD5:!PSK';
        ssl_prefer_server_ciphers on;
    
        # HSTS (ngx_http_headers_module is required) (15768000 seconds = 6 months)
        add_header Strict-Transport-Security max-age=15768000;
    }

Please check this site for updates on what TLS settings currently make sense:

[https://mozilla.github.io/server-side-tls/ssl-config-generator](https://mozilla.github.io/server-side-tls/ssl-config-generator)

Now run or reset your nginx and start your instance of **Fennel**.

Thats it, your instance of **Fennel** should run now. All logs are sent to stdout for now. Have a look at */libs/log.js* if
you want to change the options.

## Configuration ##

All parameters which can be configured right now are in the file *fennel.config.js*. There are not much parameters yet, indeed.
But **Fennel** is not ready production anyway. And you are welcome to help out in adding parameters and configuration
options.

## Contribution ##

If you happen to know how to write Go, documentation or can help out with something else, drop us a note at *contact at swordlord dot com*. As more helping hands we have, as quicker this server gets up and feature complete.

If some feature is missing, just remember that this is an Open Source Project. If you need something, think about contributing it yourself...

## Dependencies ##

When compiling, have a look at the vendor folder. Binaries have no direct Dependency whatsoever. You might want to have your database backend ready though.

## License ##

**Fennel** is published under the GNU Affero General Public Licence version 3. See the LICENCE file for details.
