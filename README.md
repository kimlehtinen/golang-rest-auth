# golang-rest-auth

Golang rest jwt authentication rest api.

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

