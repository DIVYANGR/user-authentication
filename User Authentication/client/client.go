package main

import (
	"fmt"
	"log"
	"net/http"
	"proto"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

type userlogin struct {
	Username string `json:username`
	Password string `json:password`
}

type usersignin struct {
	Username  string `json:username`
	Password  string `json:password`
	Firstname string `json:firstname`
	Lastname  string `json:lastname`
}

func main() {
	conn, err := grpc.Dial("localhost:4040", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	client := proto.NewAddServiceClient(conn)

	g := gin.Default()
	g.GET("/login", func(ctx *gin.Context) {
		var u userlogin
		if ctx.BindJSON(&u) == nil {

			ctx.JSON(200, gin.H{
				"username": u.Username,
				"password": u.Password,
			})

			req := &proto.Request{Username: string(u.Username), Password: string(u.Password)}
			if response, err := client.Login(ctx, req); err == nil {
				ctx.JSON(http.StatusOK, gin.H{
					"result": fmt.Sprint(response.Result),
				})
			} else {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}

			fmt.Printf("%s %s %s %s", u.Username, u.Password)
		}

	})

	g.POST("/signup", func(ctx *gin.Context) {

		var u usersignin
		if ctx.BindJSON(&u) == nil {

			ctx.JSON(200, gin.H{
				"username":  u.Username,
				"password":  u.Password,
				"firstname": u.Firstname,
				"lastname":  u.Lastname,
			})

			req := &proto.Request1{Username: string(u.Username), Password: string(u.Password), Firstname: string(u.Firstname), Lastname: string(u.Lastname)}
			if response, err := client.Signin(ctx, req); err == nil {
				ctx.JSON(http.StatusOK, gin.H{
					"result": fmt.Sprint(response.Result),
				})
			} else {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}

			fmt.Printf("%s %s %s %s", u.Username, u.Password, u.Firstname, u.Lastname)
		}

	})

	if err := g.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}

}
