package main

import (
	"context"
	"database/sql"
	"net"
	"proto"

	_ "github.com/go-sql-driver/mysql"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func dbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "divyang"
	dbPass := "divyangk1998"
	dbName := "golang"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	return db
}

type server struct{}

var queue []string

func main() {
	listener, err := net.Listen("tcp", ":4040")
	if err != nil {
		panic(err)
	}

	srv := grpc.NewServer()
	proto.RegisterAddServiceServer(srv, &server{})
	reflection.Register(srv)

	if e := srv.Serve(listener); e != nil {
		panic(e)
	}

}

func (s *server) Login(ctx context.Context, request *proto.Request) (*proto.Response, error) {
	var a string = request.GetUsername()
	var b string = request.GetPassword()
	//var result string = a + b
	var jobString string
	db := dbConn()
	selDB, err := db.Query("SELECT uname,psw FROM user WHERE uname=?", a)
	var uname, psw string
	for selDB.Next() {

		err = selDB.Scan(&uname, &psw)

		if err != nil {
			panic(err.Error())
		}
		//fmt.Print("hello")
		if uname == a && psw == b {
			jobString = "Hello " + a
			//return &proto.Response{Status: jobString}, nil
		} else {
			jobString = "User not found"
			//return &proto.Response{Status: jobString}, nil
		}

	}
	//fmt.Print(a + b)
	return &proto.Response{Result: jobString}, nil
}

func (s *server) Signin(ctx context.Context, request1 *proto.Request1) (*proto.Response1, error) {
	var a string = request1.GetUsername()
	var b string = request1.GetPassword()
	var c string = request1.GetFirstname()
	var d string = request1.GetLastname()
	var result string
	// Make a queue of ints.

	//queue = dequeue(queue)
	db := dbConn()
	selDB, err := db.Query("SELECT uname FROM user")
	var uname string
	for selDB.Next() {

		err = selDB.Scan(&uname)

		if err != nil {
			panic(err.Error())
		}
		//fmt.Print("hello")
		if uname == a {
			result = "Already have username try different username"
			//return &proto.Response{Status: jobString}, nil
		} else {
			insForm, err := db.Prepare("INSERT INTO user(uname, psw, fname, lname) VALUES(?,?,?,?)")
			if err != nil {
				panic(err.Error())
			}
			insForm.Exec(a, b, c, d)
			result = "Successfully signup"
		}

	}

	//fmt.Print(a + b + c + d)
	/*insForm, err := db.Prepare("INSERT INTO user(uname, psw, fname, lname) VALUES(?,?,?,?)")
	if err != nil {
		panic(err.Error())
	}
	insForm.Exec(a, b, c, d)
	*/
	return &proto.Response1{Result: result}, nil
}
