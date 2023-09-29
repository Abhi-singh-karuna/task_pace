package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"
	"os/exec"
)

type Group struct {
	Id    string   `json:"id"`
	Name  string   `json:"name"`
	Users []string `json:"users"`
}

func main() {

	// initialize the gin
	//  imkshdiuwdfgiuwrgf
	router := gin.New()

	server := socketio.NewServer(nil)

	// load the environment file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var users = make(map[string]interface{})
	var groupname Group

	// for connecting the users
	server.OnConnect("/", func(s socketio.Conn) error {
		fmt.Println("connected:", s.ID())
		return nil
	})

	// assing the username to all the connected users
	server.OnEvent("/", "username", func(s socketio.Conn, name string) {

		fmt.Println("new-users:", name)
		users[s.ID()] = name
		s.Emit("newuser", users[s.ID()])
		fmt.Println("allconectedusers : ", users)
		for key, _ := range users {
			server.BroadcastToRoom("/", key, "allconectedusers", users)
		}

	})

	// for creating the room (during creating the room admin shoud be added another user to the room)
	server.OnEvent("/", "creategroup", func(s socketio.Conn, grp Group) {

		if len(grp.Users) > 4 {
			s.Emit("warning", "maximum 5 users shuld you only add in group ... ")
			return
		}
		newUUID, err := exec.Command("uuidgen").Output()
		if err != nil {
			log.Fatal(err)
		}

		grp.Id = string(newUUID[:36])
		grp.Users = append(grp.Users, s.ID())
		s.SetContext(grp)
		fmt.Println("new group created :  ", grp.Name+" from  ID :"+grp.Id+"  users are now  :")

		fmt.Println("groupname : ", groupname)
		fmt.Println("group : ", grp)

		for _, key := range grp.Users {
			if key == s.ID() {
				server.BroadcastToRoom("/", key, "notification", "you successfully created a new group : "+grp.Name)
			} else {
				server.BroadcastToRoom("/", key, "notification", "you are added in new group : "+grp.Name+" from :"+(users[s.ID()]).(string))
			}
		}
		groupname = grp
	})

	// for joining the room
	server.OnEvent("/", "join", func(s socketio.Conn, grpid string) {
		for _, key := range groupname.Users {
			if key == s.ID() {
				server.BroadcastToRoom("/", key, "notification", "you are already in group : "+groupname.Name)
				return
			}
		}
		if grpid == groupname.Id {

			if len(groupname.Users) > 4 {
				s.Emit("warning", "maximum 5 users should be in the group ... ")
				return
			} else {
				groupname.Users = append(groupname.Users, s.ID())
				s.SetContext(groupname)
				s.Join(grpid)
				fmt.Println("new user joined : ", groupname.Users)
				for _, key := range groupname.Users {
					if key == s.ID() {
						server.BroadcastToRoom("/", key, "notification", "you successfully joined in group : "+groupname.Name)
					} else {
						server.BroadcastToRoom("/", key, "notification", "you are added in new group : "+groupname.Name+" from :"+(users[s.ID()]).(string))
					}
				}
			}
		} else {
			s.Emit("warning", "you are not allowed to join this group")
		}
	})

	// for disconnecting
	server.OnDisconnect("/", func(s socketio.Conn, msg string) {
		fmt.Println("disconnected:", users[s.ID()])
		s.Emit("disconnecteduser", users[s.ID()])
		delete(users, s.ID())
		for key, _ := range users {
			fmt.Println("key : ", key)
			server.BroadcastToRoom("/", key, "allconectedusers", users)
		}
		fmt.Println("closed", msg)
	})

	go server.Serve()
	defer server.Close()

	router.GET("/socket.io/*any", gin.WrapH(server))
	router.POST("/socket.io/*any", gin.WrapH(server))

	router.Run(os.Getenv("PORT"))
}
