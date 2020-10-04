package websocket

import (
    "fmt"
    "io"
    "log"
    "net/http"

    "github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    CheckOrigin: func(r *http.Request) bool { return true },
}

func Upgrade(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
    ws, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println(err)
        return ws, err
    }
    return ws, nil
}

func Reader(conn *websocket.Conn, a *App) {
    for {
        messageType, p, err := conn.ReadMessage()
        if err != nil {
            log.Println(err)
            return
        }

        fmt.Println(string(p))
        msg := string(p)
	row1 := &Menu{Name: msg, Description: msg}
        a.DB.Create(row1)

        if err := conn.WriteMessage(messageType, p); err != nil {
            log.Println(err)
            return
        }
    }
}

func Writer(conn *websocket.Conn) {
    for {
        fmt.Println("Sending")
        messageType, r, err := conn.NextReader()
        if err != nil {
            fmt.Println(err)
            return
        }
        w, err := conn.NextWriter(messageType)
        if err != nil {
            fmt.Println(err)
            return
        }
        if _, err := io.Copy(w, r); err != nil {
            fmt.Println(err)
            return
        }
        if err := w.Close(); err != nil {
            fmt.Println(err)
            return
        }
    }
}
func SendMsg(conn *websocket.Conn) {
      msg := []byte("Let's start to talk something.")
      err := conn.WriteMessage(websocket.TextMessage, msg)
      fmt.Printf("ffq")
      if err != nil {
        return
      }
      fmt.Printf("ff")
}
