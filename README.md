[![Build Status](https://travis-ci.org/maprost/restclient.svg)](https://travis-ci.org/maprost/restclient)
[![Coverage Status](https://coveralls.io/repos/github/maprost/restclient/badge.svg)](https://coveralls.io/github/maprost/restclient)
[![GoDoc](https://godoc.org/github.com/maprost/restclient?status.svg)](https://godoc.org/github.com/maprost/restclient)
[![Go Report Card](https://goreportcard.com/badge/github.com/maprost/restclient)](https://goreportcard.com/report/github.com/maprost/restclient)

# RestClient

## Install
```
go get github.com/maprost/restclient
```

## Supported Methods
- Get
- Post
- Put
- Delete

## Supported Format
- Json
- XML

## Features
- custom logger
- query builder

## Usage
```go
var users []User
result := restclient.Get(serverUrl + "/user").
            AddQueryParam("limit", 1).
            AddQueryParam("email", "example@gmail.com").
            SendAndGetJsonResponse(&users)
            
// check internal rest client error 
if result.Err != nil {
   return result.Err
}
// check response error
if result.StatusCode != 200 {
   return errors.New(result.ResponseError)
}

// or check both at once
if err := result.Error(); err != nil {
    return err
}
```

```go
var users []User
result := restclient.Get(serverUrl + "/user" + rcquery.New().Add("limit", 1).Get()).
            SendAndGetJsonResponse(&users)
if err := result.Error(); err != nil {
    return err
}
```

```go
var user User{/* init */}
result := restclient.Post(serverUrl + "/user").
            AddJsonBody(user).
            Send()
if err := result.Error(); err != nil {
    return err
}
```












