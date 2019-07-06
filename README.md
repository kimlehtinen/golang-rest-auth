# golang-rest-auth
Golang rest jwt authentication rest api

## Routes
Register `POST /api/user/create/`:
``` 
{
	"email": "foo@bar.com",
	"password": "password"
}
``` 

Login `POST /api/user/login`:
``` 
{
	"email": "foo@bar.com",
	"password": "password"
}
```

## Authentication
Header:
``` 
Authorization : Bearer <token>
``` 

## Todo
- password reset