package main

import (
  "fmt"
  "log"
  "net/url"
  "net/http"
  "encoding/json"
  "github.com/jinzhu/gorm"
  _ "github.com/jinzhu/gorm/dialects/sqlite"
  "github.com/gorilla/mux"
  "github.com/gorilla/websocket"
//  "net/http/httputil"
  "io/ioutil"
)

var web_conn *websocket.Conn

type Menu struct {
  ID int64 `json:"dish_id"`
  Name string `gorm:"not null" json:"name"`
  Description string `json:"description"`
  Type string `gorm:"default:snacks" json:"type"`

}
type App struct {
  DB *gorm.DB
}

type Client struct {
    ID   string
    Conn *websocket.Conn
    Pool *Pool
}

type Message struct {
    Type int    `json:"type"`
    Body string `json:"body"`
}

type Pool struct {
    Register   chan *Client
    Unregister chan *Client
    Clients    map[*Client]bool
    Broadcast  chan Message
}

func newPool() *Pool {
    return &Pool{
        Register:   make(chan *Client),
        Unregister: make(chan *Client),
        Clients:    make(map[*Client]bool),
        Broadcast:  make(chan Message),
    }
}

func (pool *Pool) start() {
    for {
        select {
        case client := <-pool.Register:
            pool.Clients[client] = true
            fmt.Println("Size of Connection Pool: ", len(pool.Clients))
            for client, _ := range pool.Clients {
                fmt.Println(client)
                client.Conn.WriteJSON(Message{Type: 1, Body: "New User Joined..."})
            }
            break
        case client := <-pool.Unregister:
            delete(pool.Clients, client)
            fmt.Println("Size of Connection Pool: ", len(pool.Clients))
            for client, _ := range pool.Clients {
                client.Conn.WriteJSON(Message{Type: 1, Body: "User Disconnected..."})
            }
            break
        case message := <-pool.Broadcast:
            fmt.Println("Sending message to all clients in Pool")
            for client, _ := range pool.Clients {
                if err := client.Conn.WriteJSON(message); err != nil {
                    fmt.Println(err)
                    return
                }
            }
        }
    }
}

func (c *Client) read(a *App) {
    defer func() {
        c.Pool.Unregister <- c
        c.Conn.Close()
    }()

    for {
        messageType, p, err := c.Conn.ReadMessage()
        if err != nil {
            log.Println(err)
            return
        }
        message := Message{Type: messageType, Body: string(p)}
//	msg := string(p)
        c.Pool.Broadcast <- message
        fmt.Printf("Message Received: %+v\n", message)
    }
}

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    CheckOrigin: func(r *http.Request) bool { return true },
}

func upgrade(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
    ws, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println(err)
        return ws, err
    }
    return ws, nil
}


// define our WebSocket endpoint
func (a *App) serveWs(pool *Pool, w http.ResponseWriter, r *http.Request) {
    fmt.Printf(r.Host)
    conn, err := upgrade(w, r)
    if err != nil {
        log.Println(err)
    }
    web_conn = conn

    client := &Client{
        Conn: conn,
        Pool: pool,
    }

    pool.Register <- client
    client.read(a)


    fmt.Printf("Client Connected")
}

func (a *App) Initialize(dbDriver string, dbURI string) {
  db, err := gorm.Open(dbDriver, dbURI)
  if err != nil {
    panic("failed to connect database")
  }
  a.DB = db
  a.DB.AutoMigrate(&Menu{})
}

func (a *App) ListHandler(w http.ResponseWriter, r *http.Request) {
  var menus []Menu

  // Select all menus and convert to JSON.
  a.DB.Find(&menus)
  menusJSON, _ := json.Marshal(menus)

  // Write to HTTP response.
  w.WriteHeader(200)
  w.Write([]byte(menusJSON))
}

func (a *App) ViewHandler(w http.ResponseWriter, r *http.Request) {
  var menu []Menu
  vars := mux.Vars(r)

  // Select the menu with the given name, and convert to JSON.
  a.DB.First(&menu, "name = ?", vars["name"])

  menuJSON, _ := json.Marshal(menu)

  // Write to HTTP response.
  w.WriteHeader(200)
  w.Write([]byte(menuJSON))
}

func (a *App) CreateHandler(w http.ResponseWriter, r *http.Request) {
  // Parse the POST body to populate r.PostForm.
  if err := r.ParseForm(); err != nil {
    panic("failed in ParseForm() call")
  }
  fmt.Printf("sdsd\n")
  body, err := ioutil.ReadAll(r.Body)
  me := Menu{}
  json.Unmarshal(body, &me)
  fmt.Printf(string(me.Name))

  // Create a new menu from the request body.
  a.DB.Create(&me)

  // Form the URL of the newly created menu.
  u, err := url.Parse(fmt.Sprintf("/menu/%s", me.Name))
  if err != nil {
    panic("failed to form new Menu URL")
  }
  base, err := url.Parse(r.URL.String())
  if err != nil {
    panic("failed to parse request URL")
  }

  // Write to HTTP response.
  w.Header().Set("Location", base.ResolveReference(u).String())
  w.WriteHeader(201)
}

func (a *App) UpdateHandler(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)

  // Parse the request body to populate r.PostForm.
  if err := r.ParseForm(); err != nil {
    panic("failed in ParseForm() call")
  }

  // Set new menu values from the request body.
  menu := &Menu{
    Name: r.PostFormValue("name"),
    Description: r.PostFormValue("description"),
  }

  // Update the menu with the given name.
  a.DB.Model(&menu).Where("name = ?", vars["name"]).Updates(&menu)

  // Write to HTTP response.
  w.WriteHeader(204)
}

func (a *App) DeleteHandler(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)

  // Delete the menu with the given name.
  a.DB.Where("name = ?", vars["name"]).Delete(Menu{})

  // Write to HTTP response.
  w.WriteHeader(204)
}


func oneCallback(scope *gorm.Scope) {
    if !scope.HasError() {
        fmt.Printf("jjjfjf")
    }
    fmt.Printf("23")

}



func main() {
  a := &App{}
  a.Initialize("sqlite3", "test.db")
  fmt.Printf("fff")
  a.DB.Callback().Create().Register("gorm:after_create", oneCallback)

  pool := newPool()
  go pool.start()

  r := mux.NewRouter()
  r.HandleFunc("/menu", a.ListHandler).Methods("GET")
  r.HandleFunc("/menu/{name:.+}", a.ViewHandler).Methods("GET")
  r.HandleFunc("/menu", a.CreateHandler).Methods("POST")
  r.HandleFunc("/menu/{name:.+}", a.UpdateHandler).Methods("PUT")
  r.HandleFunc("/menu/{name:.+}", a.DeleteHandler).Methods("DELETE")
  r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Simple Server")
    })
  r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
                         a.serveWs(pool, w, r)
  })
  http.Handle("/", r)
  if err := http.ListenAndServe(":8080", nil); err != nil {
    panic(err)
  }

  defer a.DB.Close()
}
