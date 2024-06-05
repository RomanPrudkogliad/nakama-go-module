Nakama Go Module
=================

[Nakama](https://github.com/heroiclabs/nakama) is an open source game engine for backend written using golang and supports typescript in runtime.

## Features

Go module implements RPC function accept payload with type, version, hash (all parameters
are optional, defaults: type=core, version=1.0.0, hash=null).
Like `{type: "core", version: "1.0.0", hash: "d6e4677dc8987b7b140ad75384bb"}`
function read file from the disk path= `<type>/<version>.json`
calculate content file hash and if If hashes are not equal, then content will be null, response `{type: "core", version: "1.0.0", hash: "123123", content: "jsonFileContent"}`
If file doesn't exist, error returned.

### RPC Endpoint
Accepting POST with payload

    `curl -X POST "http://127.0.0.1:7351/v2/console/api/endpoints/rpc/getcontent" \
        -d '{"userId":"","body":"{\"version\": \"1.0.0\", \"hash\": \"123\", \"type\": \"core\"}"}' \
        -H 'Content-Type: application/json' \
        -H 'Accept: application/json' \
        -H 'Authorization: Bearer <AUTH_JWT_TOKEN>'`

Response json

`{"body":"{\"type\":\"core\",\"version\":\"1.0.0\",\"hash\":\"\",\"content\":\"\"}","error_message":""}`


### Requirements
    Docker
    GoLang 1.21

## Usage 
Build Go module image 
`docker compose build`

Run image 
`docker compose up`

Open in browse 
`http://127.0.0.1:7351/#/login` and login using admin:admin credentials 

Go to API Explorer tab 
`http://127.0.0.1:7351/#/apiexplorer` 
choose `getcontent` endpoint and send request
Done you are awesome!!

## Tests
As I failed to run nktest on my mac, due to connectivity [issue](https://github.com/ascii8/nktest/issues/1)
before running the tests you need to past active JWT token to the .env file
then run

`go test`

## Thoughts and leftovers

1. Refactor tests on [nktest](https://github.com/ascii8/nktest)
2. As it's not safe to load anything from the disc by parameters in request, I've added a validation to the Type and Version parameters but it's not agile enough if new type or version appears we need to update application.
3. This file loading logic written that files are relatively small, and we can load them in memory. For large files we need another approach. 
4. Storing requests: I decided no to create new table in db, because of nakama developers recommendations, all data saved nakama collections storage. It's also point for future development.    









 




