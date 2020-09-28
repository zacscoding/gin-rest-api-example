# Api references  

- [API Overview](#API-Overview)
- [User API](#User-API)  
    - [Authentication](#Authentication)
    - [User registration](#User-Registration)  
    - [Get current user](#Get-current-user)  
    - [Update user](#Update-user)
- [Article API](#Article-API)  
    - [Create a article](#Create-a-article)
    - [Get a article](#Get-a-article)  
    - [List articles](#List-Articles)  
    - [Delete a article](#Delete-a-article)
- [Comment API](#Comment-API)  
    - [Create a comment](#Create-a-comment)  
    - [List Comments from an Article](#List-Comments-from-an-Article)

## API Overview

---  
    
## User API  

### Authentication  

`POST /v1/api/users/login`  

#### Request body  

| **Parameter** | **Type** | **Description** | **Required** |
|---------------|----------|-----------------|--------------|
| user          | Object   | a user          | yes          |
| user.email    | String   | email address   | yes          |
| user.password | String   | password        | yes          |

```json
{
  "user":{
    "email": "zaccoding@github.com",
    "password": "zaccoding"
  }
}
```  

#### Response  

`Status: 200 OK`  

```json
{
    "code": 200,
    "expire": "2020-09-23T00:09:36.1524+09:00",
    "meta": "meta",
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDA3ODczNzYsImlkIjoiYWRtaW5AZW1haWwuY29tIiwib3JpZ19pYXQiOjE2MDA3ODcyODl9.AIoE8jlM68l05l4be0irTo6-0fw1nbBIINDJhLDJ0IY"
}
```  

<br />

### User Registration  
  
`POST /v1/api/users`  

#### Request body  

| **Parameter** | **Type** | **Description** | **Required** |
|---------------|----------|-----------------|--------------|
| user          | Object   | a user          | yes          |
| user.username | String   | name            | yes          |
| user.email    | String   | email address   | yes          |
| user.password | String   | password        | yes          |

```json
{
  "user":{
    "username": "zaccoding",
    "email": "zaccoding@github.com",
    "password": "zaccoding"
  }
}
```  

#### Response  

`Status: 201 Created`  

```json
{
  "user": {
    "username": "zaccoding",
    "email": "zaccoding@github.com",
    "bio": "",
    "image": ""
  }
}
```  

<br />

### Get current user  

`GET /v1/api/user/me`  

Authentication required.  

#### Response  

`Status: 200 OK`  

```json
{
  "user": {
    "username": "zaccoding",
    "email": "zaccoding@github.com",
    "bio": "",
    "image": ""
  }
}
```  

### Update user  

`PUT /v1/api/user`  

Authentication required.  

#### Request Body  

| **Parameter** | **Type** | **Description** | **Required** |
|---------------|----------|-----------------|--------------|
| user          | Object   | User's object   | no           |
| user.username | String   | user name       | no           |
| user.password | String   | password        | no           |
| user.bio      | String   | biography       | no           |
| user.image    | String   | image url       | no           |  

```json
{
  "user": {
    "bio": "I like coding"
  }
}
```  

#### Response  

`Status: 200 OK`  

```json
{
  "user": {
    "username": "zaccoding",
    "email": "zaccoding@github.com",
    "bio": "I like coding",
    "image": ""
  }
}
```

---  

## Article API  

### Create a article  

`POST /v1/api/articles`  

Authentication required.  

#### Request Body    

| **Parameter** | **Type** | **Description**  | **Required** |
|---------------|----------|------------------|--------------|
| article       | Object   | article's object | yes          |
| article.title | String   | title            | yes          |
| article.body  | String   | body             | yes          |
| article.tags  | Array    | article's tags   | no           |  

```json
{
  "article": {
    "title": "How to train your dragon",
    "description": "Ever wonder how?",
    "body": "You have to believe",
    "tagList": ["reactjs", "angularjs", "dragons"]
  }
}
```  

#### Response  

`Status: 201 Created`  

```json
{
  "article": {
    "slug": "how-to-train-your-dragon",
    "title": "How to train your dragon",
    "body": "It takes a Jacobian",
    "tagList": ["dragons", "training"],
    "createdAt": "2016-02-18T03:22:56.637Z",
    "updatedAt": "2016-02-18T03:48:35.824Z",
    "author": {
      "username": "jake",
      "bio": "I work at statefarm",
      "image": "https://i.stack.imgur.com/xHWG8.jpg"
    }
  }
}
```  

<br />

## Get a article  

`GET /v1/api/articles/:slug`  

Authentication required.

#### Path parameter

| **Parameter** | **Description** |
|---------------|-----------------|
| slug          | article's slug  |

#### Response  

`Status: 200 OK`  

```
{
  "article": {
    "slug": "how-to-train-your-dragon",
    "title": "How to train your dragon",
    "body": "It takes a Jacobian",
    "tagList": ["dragons", "training"],
    "createdAt": "2016-02-18T03:22:56.637Z",
    "updatedAt": "2016-02-18T03:48:35.824Z",
    "author": {
      "username": "jake",
      "bio": "I work at statefarm",
      "image": "https://i.stack.imgur.com/xHWG8.jpg"
    }
  }
}
``` 

<br />

## List Articles  

`GET /v1/api/articles?tag=AngularJS&author=zaccoding&limit=20&offset=0`  

#### Request parameter  

| **Parameter** | **Type** | **Description**          | **Default** |
|---------------|----------|--------------------------|-------------|
| tag           | Array    | filter by tag            | none        |
| author        | String   | filter by author         | none        |
| limit         | Numeric  | limit number of articles | 5           |
| offset        | Numeric  | skip number of articles  | 0           |

#### Response  

`Status: 200 OK`  

```json
{
  "articles":[{
    "slug": "how-to-train-your-dragon",
    "title": "How to train your dragon",
    "body": "It takes a Jacobian",
    "tagList": ["dragons", "training"],
    "createdAt": "2016-02-18T03:22:56.637Z",
    "updatedAt": "2016-02-18T03:48:35.824Z",
    "author": {
      "username": "jake",
      "bio": "I work at statefarm",
      "image": "https://i.stack.imgur.com/xHWG8.jpg"
    }
  }, {
    "slug": "how-to-train-your-dragon-2",
    "title": "How to train your dragon 2",
    "body": "It a dragon",
    "tagList": ["dragons", "training"],
    "createdAt": "2016-02-18T03:22:56.637Z",
    "updatedAt": "2016-02-18T03:48:35.824Z",
    "author": {
      "username": "jake",
      "bio": "I work at statefarm",
      "image": "https://i.stack.imgur.com/xHWG8.jpg"
    }
  }],
  "articlesCount": 2
}
```  

<br />

## Delete a article  

`DELETE /v1/api/articles/:slug`  

Authentication required

#### Path parameter

| **Parameter** | **Description** |
|---------------|-----------------|
| slug          | article's slug  |  

#### Response  

`Status: 200 OK`  

---  

## Comment API

### Create a comment  

`POST /v1/api/articles/:slug/comments`  

Authentication required.

#### Path parameter

| **Parameter** | **Description** |
|---------------|-----------------|
| slug          | article's slug  |

#### Request Body  

| **Parameter** | **Type** | **Description** | **Required** |
|---------------|----------|-----------------|--------------|
| comment       | Object   | comment object  | yes          |
| comment.body  | String   | content         | yes          |

```json
{
  "comment": {
    "body": "His name was my name too."
  }
}
```  

#### Response

`Status: 201 Created`  

```json
{
  "comment": {
    "id": 1,
    "createdAt": "2016-02-18T03:22:56.637Z",
    "updatedAt": "2016-02-18T03:22:56.637Z",
    "body": "It takes a Jacobian",
    "author": {
      "username": "jake",
      "bio": "I work at statefarm",
      "image": "https://i.stack.imgur.com/xHWG8.jpg"
    }
  }
}
```  

<br />

### List Comments from an Article

`GET /v1/api/articles/:slug/comments`  

#### Path parameter

| **Parameter** | **Description** |
|---------------|-----------------|
| slug          | article's slug  |

#### Response  

`Status: 200 OK`  

```json
{
  "comments": [{
    "id": 1,
    "createdAt": "2016-02-18T03:22:56.637Z",
    "updatedAt": "2016-02-18T03:22:56.637Z",
    "body": "It takes a Jacobian",
    "author": {
      "username": "jake",
      "bio": "I work at statefarm",
      "image": "https://i.stack.imgur.com/xHWG8.jpg"
    }
  }]
}
```  

<br />

### Delete a comment  

`DELETE /v1/api/articles/:slug/comments/:id`  

Authentication required.

#### Path parameter

| **Parameter** | **Description** |
|---------------|-----------------|
| slug          | article's slug  |
| id            | comment's id    |  

#### Response  

`Status: 200 OK`  

---  

