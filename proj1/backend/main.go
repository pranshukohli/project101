package main

import (
  "fmt"
  "log"
  "time"
  "strings"
  "strconv"
  "net/url"
  "net/http"
  "encoding/json"
  "github.com/jinzhu/gorm"
  _ "github.com/jinzhu/gorm/dialects/sqlite"
  "github.com/gorilla/mux"
  "github.com/gorilla/websocket"
  "io/ioutil"
)

//database constants
const DATABASE = "sqlite3"
const DATABASE_NAME = "test.db"

//rest path(RP) constants
const RP_MENU = "/menu"
const RP_BAKEMENU = "/bakemenu"

//Tables with the ID
var TABLEID = map[string]int {
	"MENUS" : 1,
	"ORDERS" : 2,
	"ORDERSTATUSFULL" : 3,
	"CHEFS" : 4,
	"ORDERSTATUSBYORDER": 4,
}


type Menu struct {
  ID int64 `json:"dish_id"`
  Name string `gorm:"not null" json:"name"`
  Price int64 `gorm:"not null" json:"price"`
  Description string `json:"description"`
  Type string `gorm:"default:snacks" json:"type"`
}

type Order struct {
  ID int64 `json:"order_id"`
  ChefId string `gorm:"not null" json:"chef_id"`
  DishId int64 `gorm:"not null" json:"dish_id"`
  OrderNumber int64 `gorm:"not null" json:"order_number"`
  Type string `gorm:"default:dine_in" json:"type"`
  Status string `gorm:"default:in_progress" json:status"`
  PaymentType string `gorm:"default:cash" json:"payment_type"`
  Note string `json:"note"`
  Quantity int64 `gorm:"default:1" json:"quantity"`
}

type OrderStatus struct {
	DishId int64 `json: "dish_id"`
  DishPrice int64 `json: "dish_price"`
  DishName string `json: "dish_name"`
  DishDescription string `json: "dish_description"`
  DishType string `json: "dish_type"`
  OrderId int64 `json: "order_id"`
  OrderNumber int64 `json: "order_number"`
  OrderQuantity int64 `json: "order_quantity"`
  ChefId string `json: "chef_id"`
  OrderType string `json: "order_type"`
  OrderStatus string `json: "order_status"`
  OrderPaymentType string `json: "order_payment_type"`
  OrderNote string `json: "order_note"`
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

type OrderByNumber struct {
	OrderNumber int64 `json: "order_number"`
	OrderList []*OrderStatus `json: "order_list"`
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
  a.DB.AutoMigrate(&Order{})
}

func (a *App) ListHandler(tableId int, w http.ResponseWriter, r *http.Request) {

  var tableJSON []uint8
  switch tableId {
  case 1:
    var menus []Menu
    a.DB.Find(&menus)
    tableJSON, _ = json.Marshal(menus)
  case 2:
    var orders []Order
    a.DB.Find(&orders)
    tableJSON, _ = json.Marshal(orders)
  case 3:
    var orders []*OrderStatus
    rows, err := a.DB.Table("orders").Select("orders.id, orders.order_number, orders.status, orders.chef_id, orders.type, orders.payment_type, orders.quantity, menus.name, menus.description, menus.type, menus.id, menus.price").Joins("left join menus on menus.id = orders.dish_id").Rows()
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()
    for rows.Next() {
	    o := new(OrderStatus)
            if err := rows.Scan(&o.OrderId, &o.OrderNumber, &o.OrderStatus, &o.ChefId, &o.OrderType, &o.OrderPaymentType, &o.OrderQuantity, &o.DishName, &o.DishDescription, &o.DishType, &o.DishId, &o.DishPrice); err != nil {
                    log.Fatal(err)
            }
	    orders = append(orders,o)
    }
    if err := rows.Err(); err != nil {
            log.Fatal(err)
    }
    tableJSON, _ = json.Marshal(orders)
  case 4:
	var orders_list [][]*OrderStatus
	var onl []*OrderByNumber
	order_status := "in_progress"
	order_numbers, err := a.DB.Table(
					"orders",
				).Select(
					"distinct order_number",
				).Where(
					"status = ?",
					order_status,
				).Rows()
	if err != nil {
		log.Fatal(err)
	}
	defer order_numbers.Close()
	for order_numbers.Next() {
		var orders []*OrderStatus
		var order_number int64
		if err := order_numbers.Scan(
					&order_number,
				); err != nil {
				log.Fatal(err)
			}
		rows, err := a.DB.Table(
					"orders",
				).Select(
					"orders.id, orders.order_number, orders.status, orders.chef_id, orders.type, orders.payment_type, orders.quantity, menus.name, menus.description, menus.type, menus.id, menus.price",
				).Joins(
					"left join menus on menus.id = orders.dish_id",
				).Where(
					"order_number = ? AND status = ?",
					order_number,order_status,
				).Rows()
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()
		for rows.Next() {
			o := new(OrderStatus)
			if err := rows.Scan(
					&o.OrderId, &o.OrderNumber,
					&o.OrderStatus, &o.ChefId,
					&o.OrderType, &o.OrderPaymentType,
					&o.OrderQuantity, &o.DishName,
					&o.DishDescription, &o.DishType,
					&o.DishId, &o.DishPrice,
				); err != nil {
				log.Fatal(err)
			}
			orders = append(orders,o)
		}

		if err := rows.Err(); err != nil {
			log.Fatal(err)
		}
		on :=new(OrderByNumber)
		on.OrderNumber = order_number
		on.OrderList = orders
		onl = append(onl,on)
		orders_list = append(orders_list,orders)
	}
	//tableJSON, _ = json.MarshalIndent(onl, "", "")
	tableJSON, _ = json.Marshal(onl)

  }

  w.WriteHeader(200)
  w.Write([]byte(tableJSON))
}


func (a *App) ViewHandler(tableId int, w http.ResponseWriter, r *http.Request) {
	var tableJSON []uint8
	vars := mux.Vars(r)
	switch tableId {
	case 1:
		var menus []Menu
		a.DB.First(&menus, "id = ?", vars["dish_id"])
		tableJSON, _ = json.Marshal(menus)
	case 2:
		var orders []*OrderStatus
		rows, err := a.DB.Table(
					"orders",
				).Select(
					"orders.id, orders.order_number, orders.status, orders.chef_id, orders.type, orders.payment_type, orders.quantity, menus.name, menus.description, menus.type, menus.id, menus.price",
				).Joins(
					"left join menus on menus.id = orders.dish_id",
				).Where(
					"order_number = ?",
					vars["order_number"],
				).Rows()
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()
		for rows.Next() {
			o := new(OrderStatus)
			if err := rows.Scan(&o.OrderId, &o.OrderNumber,
					    &o.OrderStatus, &o.ChefId,
					    &o.OrderType, &o.OrderPaymentType,
					    &o.OrderQuantity,
					    &o.DishName, &o.DishDescription,
					    &o.DishType, &o.DishId,
					    &o.DishPrice);
			err != nil {
				log.Fatal(err)
			}
			orders = append(orders,o)
		}
		if err := rows.Err(); err != nil {
			log.Fatal(err)
		}
		tableJSON, _ = json.Marshal(orders)
  case 4:
	var orders_list [][]*OrderStatus
	var onl []*OrderByNumber
	order_numbers, err := a.DB.Table(
					"orders",
				).Select(
					"distinct order_number",
				).Where(
					"order_number = ?",
					vars["order_number"],
				).Rows()
	if err != nil {
		log.Fatal(err)
	}
	defer order_numbers.Close()
	for order_numbers.Next() {
		var orders []*OrderStatus
		var order_number int64
		if err := order_numbers.Scan(
					&order_number,
				); err != nil {
				log.Fatal(err)
			}
		rows, err := a.DB.Table(
					"orders",
				).Select(
					"orders.id, orders.order_number, orders.status, orders.chef_id, orders.type, orders.payment_type, orders.quantity, menus.name, menus.description, menus.type, menus.id, menus.price",
				).Joins(
					"left join menus on menus.id = orders.dish_id",
				).Where(
					"order_number = ?",
					order_number,
				).Rows()
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()
		for rows.Next() {
			o := new(OrderStatus)
			if err := rows.Scan(
					&o.OrderId, &o.OrderNumber,
					&o.OrderStatus, &o.ChefId,
					&o.OrderType, &o.OrderPaymentType,
					&o.OrderQuantity, &o.DishName,
					&o.DishDescription, &o.DishType,
					&o.DishId, &o.DishPrice,
				); err != nil {
				log.Fatal(err)
			}
			orders = append(orders,o)
		}

		if err := rows.Err(); err != nil {
			log.Fatal(err)
		}
		on :=new(OrderByNumber)
		on.OrderNumber = order_number
		on.OrderList = orders
		onl = append(onl,on)
		orders_list = append(orders_list,orders)
	}
	//tableJSON, _ = json.MarshalIndent(onl, "", "")
	tableJSON, _ = json.Marshal(onl)
	  }

  w.WriteHeader(200)
  w.Write([]byte(tableJSON))
}


func (a *App) CreateHandler(tableId int, w http.ResponseWriter, r *http.Request) {
  // Parse the POST body to populate r.PostForm.
  if err := r.ParseForm(); err != nil {
    panic("failed in ParseForm() call")
  }
  fmt.Printf("sdsd\n")
  body, _ := ioutil.ReadAll(r.Body)

  switch tableId {
  case 1:
	me := Menu{}
	json.Unmarshal(body, &me)
	fmt.Printf(string(me.Name))
	a.DB.Create(&me)
	u, err := url.Parse(fmt.Sprintf(RP_MENU + "/%s", me.Name))
	if err != nil {
		panic("failed to form new Menu URL")
	}
	base, err := url.Parse(r.URL.String())
	if err != nil {
		panic("failed to parse request URL")
	}
	w.Header().Set("Location", base.ResolveReference(u).String())
	w.WriteHeader(201)
  case 2:
	or := Order{}
	json.Unmarshal(body, &or)
	fmt.Println(or)
	current_time := time.Now()
	ct := current_time.Format("2006-01-02 15:04:05")
	ct4 := strings.Replace(strings.Replace(strings.Replace(ct,"-","",3)," ","",1),":","",3)
	n, err := strconv.ParseInt(ct4, 10, 64)
	if err == nil {
		fmt.Printf("%d of type %T", n, n)
	}
	or.ChefId = "1"
//	or.OrderNumber = n
	fmt.Printf(string(or.DishId))
	a.DB.Create(&or)
	u, err := url.Parse(fmt.Sprintf(RP_MENU + "/%d", or.DishId))
	if err != nil {
		panic("failed to form new Menu URL")
	}
	base, err := url.Parse(r.URL.String())
	if err != nil {
		panic("failed to parse request URL")
	}

	w.Header().Set("Location", base.ResolveReference(u).String())
	w.WriteHeader(201)
  }


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


func dbUpdateCallback(scope *gorm.Scope) {
    if !scope.HasError() {
        fmt.Printf("DB Update Error!!")
    }
    fmt.Printf("DB Updated")

}



func main() {

	//Initilise Database
	a := &App{}
	a.Initialize(DATABASE, DATABASE_NAME)
	//Callback functioin when and if database gets a new row -- gorm
	a.DB.Callback().Create().Register("gorm:after_create", dbUpdateCallback)
	defer a.DB.Close()

	//Running a thread for multiple websocket connections
	pool := newPool()
	go pool.start()


	//For routing server for different paths -- rest approach
	//Using mux for routing
	r := mux.NewRouter()
	r.HandleFunc(
		RP_MENU,
		func(w http.ResponseWriter, r *http.Request) {
			a.ListHandler(TABLEID["MENUS"], w, r)
		}).Methods("GET")
	r.HandleFunc(
		RP_BAKEMENU + "",
		func(w http.ResponseWriter, r *http.Request) {
			a.ListHandler(TABLEID["ORDERS"], w, r)
		}).Methods("GET")
	r.HandleFunc(
		"/bakemenufull",
		func(w http.ResponseWriter, r *http.Request) {
			a.ListHandler(TABLEID["ORDERSTATUSFULL"], w, r)
		}).Methods("GET")
	r.HandleFunc(
		"/bakemenubyorder",
		func(w http.ResponseWriter, r *http.Request) {
			a.ListHandler(TABLEID["ORDERSTATUSBYORDER"], w, r)
		}).Methods("GET")
	r.HandleFunc(
		RP_MENU + "/{dish_id:.+}",
		func(w http.ResponseWriter, r *http.Request) {
			a.ViewHandler(TABLEID["MENUS"], w, r)
		}).Methods("GET")
	r.HandleFunc(
		RP_BAKEMENU + "/{order_number:.+}",
		func(w http.ResponseWriter, r *http.Request) {
			a.ViewHandler(TABLEID["ORDERSTATUSBYORDER"], w, r)
		}).Methods("GET")
	r.HandleFunc(
		RP_MENU,
		func(w http.ResponseWriter, r *http.Request) {
			a.CreateHandler(1, w, r)
		}).Methods("POST")
	r.HandleFunc(
		"/order",
		func(w http.ResponseWriter, r *http.Request) {
			a.CreateHandler(2, w, r)
		}).Methods("POST")
	r.HandleFunc(
		RP_MENU + "/{name:.+}",
		a.UpdateHandler).Methods("PUT")
	r.HandleFunc(
		RP_MENU + "/{name:.+}",
		a.DeleteHandler).Methods("DELETE")
	r.HandleFunc(
		"/",
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Simple Server")
		})
	r.HandleFunc(
		"/ws",func(w http.ResponseWriter, r *http.Request) {
			a.serveWs(pool, w, r)
		})


	//go http library to handle http requests
	http.Handle("/", r)

	//Starting backend server at 8080
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
