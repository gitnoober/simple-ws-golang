package handlers

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"

	"github.com/CloudyKit/jet/v6"
)

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

// WSJsonResponse Defines the response send back from the websocket
type WSJsonResponse struct {
	Action      string `json:"action"`
	Message     string `json:"message"`
	MessageType string `json:"message_type"`
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
	err = ws.WriteJSON(response) // Takes care of marshalling as well & responds with JSON
	if err != nil {
		log.Println(err)
	}
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
