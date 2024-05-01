package handlers

import (
	"fmt"
	"log"
	"net/http"
	"sort"

	"github.com/gorilla/websocket"

	"github.com/CloudyKit/jet/v6"
)

// User Flow:
//User comes to the webpage
//JS within the webpage creates a websocket connection(by upgrading the connection from HTTP to WS protocol
//Now we are listening to WS(go routine) and if a payload comes it sends it to a channel
//Now the job is to broadcast it to all online users (send the response to every client you know about)

var wsChan = make(chan WSPayload)

var clients = make(map[WebSocketConnection]string)

var views = jet.NewSet(
	jet.NewOSFileSystemLoader("./html"),
	jet.InDevelopmentMode(),
)

var upgradeConnection = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// Home renders the home page
func Home(w http.ResponseWriter, r *http.Request) {
	err := RenderPage(w, "home.jet", nil)
	if err != nil {
		log.Println(err)
	}
}

type WebSocketConnection struct {
	*websocket.Conn
}

// ListenForWS listens for WS, once it gets payload it sends it off to the channel, it also has auto recovery
func ListenForWS(conn *WebSocketConnection) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("Error", fmt.Sprintf("%v", err))
		}
	}()

	var payload WSPayload

	for {
		err := conn.ReadJSON(&payload)
		if err != nil {
			// do nothing
		} else {
			payload.Conn = *conn
			wsChan <- payload
		}
	}
}

// ListenToWSChannel listens to the channel for payload and broadcasts it to all users.
func ListenToWSChannel() {
	var response WSJsonResponse
	for {
		e := <-wsChan

		switch e.Action {
		case "username":
			//get a list of all users and send it back
			clients[e.Conn] = e.UserName
			allUsers := GetUserList()
			response.Action = "list_users"
			response.ConnectedUsers = allUsers
			BroadcastToAll(response)

		case "left":
			//delete user from map
			response.Action = "list_users"
			delete(clients, e.Conn)
			users := GetUserList()
			response.ConnectedUsers = users
			BroadcastToAll(response)

		case "broadcast":
			// broadcast message to all users
			response.Action = "broadcast"
			response.Message = "<strong>" + e.UserName + "</strong>: " + e.Message
			BroadcastToAll(response)

		}
		//response.Action = "Got here"
		//response.Message = fmt.Sprintf("Some message, and action was %s", e.Action)
		//BroadcastToAll(response)
	}
}

func GetUserList() []string {
	var UserList []string
	for _, client := range clients {
		if client != "" {
			UserList = append(UserList, client)
		}
	}
	sort.Strings(UserList)
	return UserList
}

func BroadcastToAll(response WSJsonResponse) {
	for client := range clients {
		err := client.WriteJSON(response)
		if err != nil {
			log.Println("websocket err")
			_ = client.Close()
			delete(clients, client)
		}

	}
}

// WSJsonResponse Defines the response send back from the websocket
type WSJsonResponse struct {
	Action         string   `json:"action"`
	Message        string   `json:"message"`
	MessageType    string   `json:"message_type"`
	ConnectedUsers []string `json:"connected_users"`
}

type WSPayload struct {
	Action   string              `json:"action"`
	UserName string              `json:"username"`
	Message  string              `json:"message"`
	Conn     WebSocketConnection `json:"-"`
}

// WSEndpoint Upgrades connection from HTTP protocol to websocket
func WSEndpoint(w http.ResponseWriter, r *http.Request) {
	ws, err := upgradeConnection.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	log.Println("Client connected to endpoint")
	var response WSJsonResponse
	response.Message = `<em><small>Connected to server</small></em>`

	conn := WebSocketConnection{ws}
	clients[conn] = ""

	err = ws.WriteJSON(response) // Takes care of marshalling as well & responds with JSON
	if err != nil {
		log.Println(err)
	}

	go ListenForWS(&conn) // go routine
}

func RenderPage(w http.ResponseWriter, tmpl string, data jet.VarMap) error {
	view, err := views.GetTemplate(tmpl)
	if err != nil {
		log.Println(err)
		return err
	}

	err = view.Execute(w, data, nil)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
