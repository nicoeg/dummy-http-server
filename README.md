# Dummy HTTP Server

A Simple HTTP server which will respond with what you tell it to through a JSON file.  
It will stop on the first match in the config JSON file.

Project is forked from [jasonrm/dummy-server](https://github.com/jasonrm/dummy-server).

## Usage

```sh
$ docker run -it -p 8080:8080 nicoeg/dummy-http-server
``` 

> Note: the custom config file needs to be mounted into the docker container for the HTTP server to use it.  
```sh
$ docker run -v $PWD/config.json:/app/config.json -it -p 8080:8080 nicoeg/dummy-http-server
``` 

### Options

```
Usage of ./dummy-http-server:
  -port int
        port number (default 8080)
  -config string
        path to config file (default config.json)
```

## Request configuration

Use the config to match different requests. Configuration is loaded on each request right now to allow for easy testing.  
When match keys are not present they will be ignored.

Example configuration
```JSON
[
  {
    "match": {
      "url": "/data",
      "method": "GET"
    },
    "response": {
      "status": 200,
      "body": "success"
    }
  },
  {
    "match": {
      "url": "/match-any-method",
    },
    "response": {
      "status": 403,
      "body": "Forbidden"
    }
  }
]
```

## TODO:

- Point output of response to file
- Match on body, query, and headers
- Wildcards when matching
- Tests ofc...
- Listen for config file changes instead of loading on each request
