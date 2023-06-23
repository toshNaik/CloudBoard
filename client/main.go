package main

import (
	"fmt"
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"context"
	"syscall"

	"golang.design/x/clipboard"
	"github.com/gorilla/websocket"
	"golang.org/x/term"
)

var (
	URL string = "https://cloudboard-app-gtcylmuena-uc.a.run.app/"
	cookies []*http.Cookie
	clipboardChangeByServer bool = false
)

func init() {
	err := clipboard.Init()
	if err != nil {
		panic(err)
	}

	var email string
	fmt.Println("Enter email:")
	fmt.Scanln(&email)
	fmt.Println("Enter password:")
	password, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		panic(err)
	}
	signUp(email, string(password))
	cookies = login(email, string(password))
}

func main() {
	connectToWS()
}

func signUp(email string, password string) {
	// make json payload
	jsonPayload, err := json.Marshal(map[string]interface{}{
		"email": email,
		"password": password,
	})
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
	}

	// send request
	req, err := http.Post(URL + "signup", "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		fmt.Println("Error sending request:", err)
	}
	defer req.Body.Close()

	if req.StatusCode == http.StatusConflict {
		// user already exists return and call login
		return
	} else if req.StatusCode != http.StatusOK {
		fmt.Println("Error:", req.Body)
		panic("Error signing up")
	}
	fmt.Println("Signed up successfully!")
}

func login(email string, password string) []*http.Cookie {
	// make json payload
	jsonPayload, err := json.Marshal(map[string]interface{}{
		"email": email,
		"password": password,
	})
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
	}

	// send request
	req, err := http.Post(URL + "login", "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		fmt.Println("Error sending request:", err)
	}
	defer req.Body.Close()
	if req.StatusCode != http.StatusOK {
		fmt.Println("Error:", req.Body)
		panic("Error logging in")
	}

	// get the cookies
	return req.Cookies()
}

func connectToWS() {
	wsURL := url.URL{Scheme: "wss", Host: "cloudboard-app-gtcylmuena-uc.a.run.app", Path: "/cloudboard/ws"}
	
	// create a new http header and add the cookies to it
	header := http.Header{}
	for _, cookie := range cookies {
		if cookie.Path == "/" {
			header.Add("Cookie", cookie.Name + "=" + cookie.Value)
		}
	}

	// connect to the websocket
	conn, _, err := websocket.DefaultDialer.Dial(wsURL.String(), header)
	if err != nil {
		fmt.Println("Error connecting to websocket:", err)
	}
	defer conn.Close()

	clipboardChanged := clipboard.Watch(context.Background(), clipboard.FmtText)

	go handleClipboardChanges(conn, clipboardChanged)
	handleServerMessages(conn)
}

func handleClipboardChanges(conn *websocket.Conn, changed <-chan []byte) {
	for {
		select {
		case text := <-changed:
			if clipboardChangeByServer {
				clipboardChangeByServer = false
				continue
			}
			fmt.Println("Clipboard changed!")
			err := conn.WriteMessage(websocket.TextMessage, text)
			if err != nil {
				fmt.Println("Error writing to websocket:", err)
			}
		}
	}
}

func handleServerMessages(conn *websocket.Conn) {
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Error reading from websocket:", err)
			return
		}
		fmt.Println("Message received:", string(message))
		// current := clipboard.Read(clipboard.FmtText)
		// if !bytes.Equal(current, message) {
		clipboard.Write(clipboard.FmtText, message)
		clipboardChangeByServer = true
		// }
	}
}
