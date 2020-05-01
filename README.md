# Rest API with golang, gin, gorm  
This project is example for rest api.

- List articles
`GET /api/articles?tag=AngularJS&author=jake&favorited=jake&limit=20&offset=0`  

- Feed articles  
`GET /api/articles/feed`  

- Get article  
`GET /api/articles/:slug`  

- Post article
`POST /api/articles`  

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

- Update article  
`PUT /api/articles/:slug`  
optional : title, description, body
```json
{
  "article": {
    "title": "Did you train your dragon?"
  }
}
```  

- Delete article
`DELETE /api/articles/:slug`  

- Post comment  
`POST /api/articles/:slug/comments`  
```json
{
  "comment": {
    "body": "His name was my name too."
  }
}
```

- Get comments  
`GET /api/articles/:slug/comments`  

- Delete comment  
`DELETE /api/articles/:slug/comments/:id`  

- Favorite article  
`POST /api/articles/:slug/favorite`  

- Unfavorite article  
`DELETE /api/articles/:slug/favorite`

- Get Tags  
`GET /api/tags`