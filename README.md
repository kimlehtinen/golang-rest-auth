# Table of Contents:
1. [Getting started](#getting-started)
2. [Features](#features)
3. [Stack](#stack)
4. [Routes](#routes)
5. [Email example](#email-example)

  
## Getting started
Create a `.env` file and change this so it works for you:
``` 
db_name = myapp
db_pass = password
db_user = postgres
db_type = postgres
db_host = localhost
db_port = 5434
jwt_secret = mysecret
database_url = postgres://postgres:password@localhost:5434
mail_server = smtp.gmail.com
mail_port = 465
mail_user = mailuser
mail_password = mailpass
```
 
### Run application
 If you don't want to build executable you can just `go run main.go`

### Build and run executable
#### Build
If you have go1.12 installed: <br>
`make generate-bin` <br><br>
If you don't have go1.12 you can do it with docker as well: <br>
`make docker-build`
#### Run executable
`bin/server`

## Features:

- user register

- user login

- protected routes (middleware)

- password reset email

  

## Stack

- golang server

- postgres db

  

## Routes

  

### User register `POST /api/user/create/`

#### Description

A route for registering a new user.

#### Params:

```
{
	"email": "foo@bar.com",
	"password": "password"
}
```

  

### User login `POST /api/user/login`

  

#### Description

User login route. As a response user gets jwt token that can be used to access protected routes that requires user login.

  

#### Params:

```
{
	"email": "foo@bar.com",
	"password": "password"
}
```

  

### User forgot password: `POST /api/user/forgot-password/:email`

  

#### Description

This route can be used to send a reset password email to user, it stores a reset token in database that expires in 1h. This token is sent to user's email.

Emails are sent using smtp, it can be configured in `.env` file.

Only `smtp.gmail.com` has been tested so far. In order for gmail to work, less secure apps has to be enabled in gmail settings, you can read more about this [here](https://support.google.com/accounts/answer/6010255?hl=en) .

  

#### Params

Url param :email **string**

  

### User reset password check `GET /api/user/reset-psw-check/:reset-token`

  

#### Description

This route can be used to check the status of a reset password token.

A link to this route can also be found in reset password email.

  

#### Params

Url param :reset-token **string**

  

### User reset password `POST /api/user/reset-password`

  

#### Description:

This changes the users password, using the token/link that was sent to users email.

Reset token status can be checked with endpoint `/api/user/reset-psw-check/:reset-token`

  

#### Params:

```
{
	"token_reset": "<reset-token>",
	"new_password": "newpassword"
}
```

  

## Authentication

  

Header:

  

```
Authorization : Bearer <token>
```

## Email example

<img src="email_example.jpg" alt="Example email" width="200"/>