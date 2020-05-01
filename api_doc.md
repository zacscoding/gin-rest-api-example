# Api documents  
This api spec is extend api of [real world api spec](https://github.com/gothinkster/realworld/tree/master/api)  

Most of api is same to real world but updated user's profile for followers.

- [User API](#User-API)

---  

# User API

## Overview  
User API provide registration and authentication, resource of users.

---  

## Operation

- <a href="user_registration">Registration</a>
- <a href="user_authentication">Authentication</a>

---    

<div id="user_registration"></div>  

### Registration  
`POST /api/users`  

Save a new user and not required authentication.  

#### Request body

| **Parameter** | **Type** | **Description**       |
|---------------|----------|-----------------------|
| user          | Object   | Information of a user |
| user.username | String   | Name of a user        |
| user.email    | String   | Email address         |
| user.password | String   | Password of a user    |

> #### Request example

```json
{
  "user":{
    "username": "Jacob",
    "email": "jake@jake.jake",
    "password": "jakejake"
  }
}
```  

#### Response  

```json
{
  "user": {
    "email": "jake@jake.jake",
    "token": "jwt.token.here",
    "username": "jake",
    "bio": "I work at statefarm",
    "image": null
  }
}
```




---  

API Template

# API Title

## Overview

## Operation

##
