github.com/rhtps/gochat/
# gochat
A nifty chat client written in Go.

## Usage
### Building
```
go get gitub.com/rhtps/gochat
cd $GOPATH/src/github.com/rhtps/gochat
go build
```
### Help
```
./gochat --help
Usage of ./gochat:
  -avatarPath="avatars/": The path to the folder for the avatar images  This is relative to the location from which "gochat" is executed.  Can be absolute.
  -callBackHost="localhost:8080": The host address of the application.
  -facebookProviderKey="12345": The FaceBook OAuth provider key.
  -facebookProviderSecretKey="12345": The FaceBook OAuth provider secret key.
  -githubProviderKey="12345": The GitHub OAuth provider key.
  -githubProviderSecretKey="12345": The GitHub OAuth provider secret key.
  -googleProviderKey="12345": The Google OAuth provider key.
  -googleProviderSecretKey="12345": The Google OAuth provider secret key.
  -host=":8080": The host address of the application.
  -securityKey="12345": The OAuth security key.
  -templatePath="templates/": The path to the HTML templates.  This is relative to the location from which "gochat" is executed.  Can be absolute.
```

### Basic Execution
To execute, binding to any interface on the host, execute:
```
./gochat -host=0.0.0.0:8080
```

## Resources
"gochat" is based on an example from [Go Programming Blueprints](https://github.com/matryer/goblueprints) by Mat Ryer

["go-http-auth"](https://github.com/abbot/go-http-auth) provides the basic authentication
## License
Copyright (c) 2015 Kenneth D. Evensen

    This program is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.

    You should have received a copy of the GNU General Public License
    along with this program.  If not, see <http://www.gnu.org/licenses/>.
