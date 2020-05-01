package main

import (
	"fmt"
	"gin-rest-api-example/database/models"
	"gin-rest-api-example/repository"
	"github.com/jinzhu/gorm"
	"log"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func main() {
	db, err := gorm.Open("mysql", "root:password@tcp(127.0.0.1:13306)/my_database?charset=utf8&parseTime=True")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	//db.DropTable(&models.Follow{}, &models.User{})
	//db.AutoMigrate(&models.Follow{}, &models.User{})
	//testUsers(db)

	//db.DropTable(&models.Comment{}, &models.ArticleTag{}, &models.Tag{}, &models.ArticleFavorite{},
	//	&models.Article{}, &models.Follow{}, &models.User{})
	//db.AutoMigrate(&models.Follow{}, &models.User{}, &models.Article{}, &models.ArticleFavorite{}, &models.Tag{},
	//	&models.ArticleTag{}, &models.Comment{})
	testArticles(db)
}

func testArticles(db *gorm.DB) {
	userRepo := repository.NewUserRepository(db)
	articleRepo := repository.NewArticleRepository(db)

	fmt.Println("Try to save user1")
	u1 := &models.User{
		Email:    "user1@email.com",
		Username: "user1",
		Password: "user1",
		Bio:      "user1 bio",
		Image:    "user1 image",
	}
	_ = userRepo.Save(u1)

	fmt.Println("Try to save article1")
	a := &models.Article{
		Title:       "title",
		Description: "description",
		Body:        "body",
		Author:      *u1,
		AuthorID:    u1.ID,
		Tags:        []models.Tag{
			{
				Name : "Tag1",
			},
			{
				Name : "Tag2",
			},
		},
		Comment:     nil,
	}
	a.UpdateSlug()
	_ = articleRepo.SaveArticle(a)

	fmt.Println("Try to save article1")
	a2 := &models.Article{
		Title:       "title2",
		Description: "description2",
		Body:        "body2",
		AuthorID:    u1.ID,
		Tags:        []models.Tag{
			{
				Name : "Tag1",
			},
			{
				Name : "Tag3",
			},
		},
		Comment:     nil,
	}
	a2.UpdateSlug()
	_ = articleRepo.SaveArticle(a2)

	fmt.Println("Try to save comment1")
	c := &models.Comment{
		Body:      "comment",
		ArticleID: a.ID,
		AuthorID:  u1.ID,
	}
	_ = articleRepo.SaveOne(c)
}

func testUsers(db *gorm.DB) {
	userRepo := repository.NewUserRepository(db)
	u1 := &models.User{
		Email:    "user1@email.com",
		Username: "user1",
		Password: "user1",
		Bio:      "user1 bio",
		Image:    "user1 image",
	}
	_ = userRepo.Save(u1)
	u2 := &models.User{
		Email:    "user2@email.com",
		Username: "user2",
		Password: "user2",
		Bio:      "user2 bio",
		Image:    "user2 image",
	}
	_ = userRepo.Save(u2)
	u3 := &models.User{
		Email:    "user3@email.com",
		Username: "user3",
		Password: "user3",
		Bio:      "user3 bio",
		Image:    "user3 image",
	}
	_ = userRepo.Save(u3)

	find1, err := userRepo.FindByEmail(u1.Email)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Find by email :", find1)

	find1, err = userRepo.FindByUsername(u1.Username)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Find by username :", find1)

	u1.Username = "Updated user1"
	_ = userRepo.Update(u1)
	find1, err = userRepo.FindByUsername(u1.Username)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Find by updated username :", find1)

	find1, err = userRepo.FindByUsername(u1.Username + "$Unknown$")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Find by updated username :", find1)

	b, _ := userRepo.IsFollowing(u2, u1)
	fmt.Println("u2 follow u1 ?", b)
	_ = userRepo.UpdateFollow(u2, u1)
	b, _ = userRepo.IsFollowing(u2, u1)
	fmt.Println("u2 follow u1 (after follow)?", b)
	_ = userRepo.UpdateFollow(u3, u1)
	_ = userRepo.UpdateFollow(u1, u3)

	followingCnt, followerCnt, err := userRepo.CountFollows(u1)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Find follows count following :", followingCnt, ", follower :", followerCnt)
	followers, _ := userRepo.FindFollowers(u1)
	fmt.Println("Find u1's followers")
	for _, f := range followers {
		fmt.Println("> ", f)
	}
	following, _ := userRepo.FindFollowing(u1)
	fmt.Println("Find u1's following")
	for _, f := range following {
		fmt.Println("> ", f)
	}

	_ = userRepo.UpdateUnFollow(u2, u1)
	followingCnt, followerCnt, err = userRepo.CountFollows(u1)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Find follows count after un follow. following :", followingCnt, ", follower :", followerCnt)
}
