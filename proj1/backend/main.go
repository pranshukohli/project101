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

)

// We'll need to define an Upgrader
// this will require a Read and Write buffer size
var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
  WriteBufferSize: 1024,

  // We'll need to check the origin of our connection
  // this will allow us to make requests from our React
  // development server to here.
  // For now, we'll do no checking and just allow any connection
  CheckOrigin: func(r *http.Request) bool { return true },
}


type Menu struct {
  ID int64 `json:"dish_id"`
  Name string `gorm:"not null" json:"name"`
  Description string `json:"description"`
  Type string `gorm:"default:snacks" json:"type"`

}
type App struct {
  DB *gorm.DB
}

// define a reader which will listen for
// new messages being sent to our WebSocket
// endpoint
func reader(conn *websocket.Conn, a *App) {
    for {
    // read in a message
        messageType, p, err := conn.ReadMessage()
        if err != nil {
            log.Println(err)
            return
        }
    // print out that message for clarity
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

// define our WebSocket endpoint
func (a *App) serveWs(w http.ResponseWriter, r *http.Request) {
    fmt.Println(r.Host)

  // upgrade this connection to a WebSocket
  // connection
    ws, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println(err)
  }
  // listen indefinitely for new messages coming
  // through on our WebSocket connection
    reader(ws, a)
}
func (a *App) Initialize(dbDriver string, dbURI string) {
  db, err := gorm.Open(dbDriver, dbURI)
  if err != nil {
    panic("failed to connect database")
  }
  a.DB = db

  // Migrate the schema.
  a.DB.AutoMigrate(&Menu{})
}
/*
func handler(w http.ResponseWriter, r *http.Request) {
  w.WriteHeader(200)
  w.Write([]byte("hello world"))
}
*/

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

  // Create a new menu from the request body.
  menu := &Menu{
    Name: r.PostFormValue("name"),
    Description: r.PostFormValue("description"),
  }
  a.DB.Create(menu)

  // Form the URL of the newly created menu.
  u, err := url.Parse(fmt.Sprintf("/menu/%s", menu.Name))
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


  r := mux.NewRouter()
  r.HandleFunc("/menu", a.ListHandler).Methods("GET")
  r.HandleFunc("/menu/{name:.+}", a.ViewHandler).Methods("GET")
  r.HandleFunc("/menu", a.CreateHandler).Methods("POST")
  r.HandleFunc("/menu/{name:.+}", a.UpdateHandler).Methods("PUT")
  r.HandleFunc("/menu/{name:.+}", a.DeleteHandler).Methods("DELETE")
  r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Simple Server")
    })
  r.HandleFunc("/ws", a.serveWs)
  http.Handle("/", r)
  if err := http.ListenAndServe(":8080", nil); err != nil {
    panic(err)
  }

  defer a.DB.Close()
}
