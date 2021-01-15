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
	"io/ioutil"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
	"github.com/gorilla/sessions"
	"reflect"
)

//session_secret
const SESSION_SECRET_KEY = "xWYMMrMKTd6XD34e5MTOmrdQyXkCaIeQ"  // need to place somewhere in files


//database constants
const DATABASE = "sqlite3"
const DATABASE_NAME = "db/test.db"

//rest path(RP) constants
const RP_MENU = "/menu"
const RP_MENU_UPDATE = "/bakemenuupdate"
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
	ID		int64	`json:"dish_id"`
	Name		string	`gorm:"not null" json:"name"`
	Price		int64	`gorm:"not null" json:"price"`
	Description	string	`json:"description"`
	Type		string	`gorm:"default:snacks" json:"type"`
}

type Order struct {
	ID		int64	`json:"order_id"`
	ChefId		string	`gorm:"not null" json:"chef_id"`
	DishId		int64	`gorm:"not null" json:"dish_id"`
	OrderNumber	int64	`gorm:"not null" json:"order_number"`
	Type		string	`gorm:"default:dine_in" json:"type"`
	Status		string	`gorm:"default:in_progress" json:status"`
	PaymentType	string	`gorm:"default:cash" json:"payment_type"`
	Note		string	`json:"note"`
	Quantity	int64	`gorm:"default:1" json:"quantity"`
}

type OrderStatus struct {
	DishId		int64	`json: "dish_id"`
	DishPrice	int64	`json: "dish_price"`
	DishName	string	`json: "dish_name"`
	DishDescription	string	`json: "dish_description"`
	DishType	string	`json: "dish_type"`
	OrderId		int64	`json: "order_id"`
	OrderNumber	int64	`json: "order_number"`
	OrderQuantity	int64	`json: "order_quantity"`
	ChefId		string	`json: "chef_id"`
	OrderType	string	`json: "order_type"`
	OrderStatus	string	`json: "order_status"`
	OrderPaymentType string	`json: "order_payment_type"`
	OrderNote	string	`json: "order_note"`
}

type App struct {
	DB *gorm.DB
}

type Client struct {
	ID	string
	Conn	*websocket.Conn
	Pool	*Pool
}

type Message struct {
	Type	int	`json:"type"`
	Body	string	`json:"body"`
}

type Pool struct {
	Register	chan *Client
	Unregister	chan *Client
	Clients		map[*Client]bool
	Broadcast	chan Message
}

type OrderByNumber struct {
	OrderNumber	int64	`json: "order_number"`
	OrderList	[]*OrderStatus	`json: "order_list"`
}

func newPool() *Pool {
	return &Pool{
		Register:	make(chan *Client),
		Unregister:	make(chan *Client),
		Clients:	make(map[*Client]bool),
		Broadcast:	make(chan Message),
	}
}

func (pool *Pool) start() {
	for {
		select {
		case client := <-pool.Register:
			pool.Clients[client] = true
			log.Println(client)
			log.Println(
				"Size of Connection Pool: ",
				len(pool.Clients))
			for client, _ := range pool.Clients {
				client.Conn.WriteJSON(
					Message{
						Type:	1,
						Body:	"New User Joined...",
					})
			}
			break
		case client := <-pool.Unregister:
			delete(pool.Clients, client)
			log.Println(
				"Size of Connection Pool: ",
				len(pool.Clients))
			for client, _ := range pool.Clients {
				client.Conn.WriteJSON(
					Message{
						Type: 1,
						Body: "User Disconnected...",
					})
			}
			break
		case message := <-pool.Broadcast:
			log.Println("Sending message to all clients in Pool")
			for client, _ := range pool.Clients {
				if err := client.Conn.WriteJSON(message);
				   err != nil {
					log.Println(err)
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
		message := Message{
				Type: messageType,
				Body: string(p)}
		//msg := string(p)
		c.Pool.Broadcast <- message
		log.Printf("Message Received: %+v\n", message)
	}
}

var upgrader = websocket.Upgrader {
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {return true},
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
	log.Printf(r.Host)
	for _, cookie := range r.Cookies() {
		log.Printf(cookie.Name)
		log.Printf(cookie.Value)
	}
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
	log.Printf("qwerty" + client.ID)
	log.Printf("Client Connected")
}

func (a *App) Initialize(dbDriver string, dbURI string) {
	db, err := gorm.Open(dbDriver, dbURI + "?_busy_timeout=5000")
	if err != nil {
		panic("failed to connect database")
	}
	a.DB = db
	//a.DB.AutoMigrate(&Menu{})
	//a.DB.AutoMigrate(&Order{})
}

func ProfileHandler(
		w http.ResponseWriter,
		r *http.Request) {

			p := ""
			log.Printf("cookie")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			if origin := r.Header.Get("Origin"); origin != "" {
				w.Header().Set("Access-Control-Allow-Origin", origin)
			}

			for _, cookie := range r.Cookies() {
				log.Printf(cookie.Name)
				log.Printf(cookie.Value)
			}

			store := sessions.NewCookieStore([]byte(SESSION_SECRET_KEY))
			session,err := store.Get(r, "project101_session")
			if err != nil {
				log.Println(err)
			}else if session.Values["google_id"] != nil {
				log.Printf(session.Values["google_id"].(string))
				p = session.Values["google_id"].(string)
			}
		w.WriteHeader(200)
		w.Write([]byte(p))
		}

func (a *App) ListHandler(
		tableId int,
		w http.ResponseWriter,
		r *http.Request) {
	var tableJSON []uint8
	switch tableId {
	case 1:	var menus []Menu
		log.Printf("cookie")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
//		w.Header().Set("Access-Control-Allow-Origin", "http://192.168.3.120,http://192.168.3.120:3000")
		if origin := r.Header.Get("Origin"); origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}
		for _, cookie := range r.Cookies() {
			log.Printf(cookie.Name)
			log.Printf(cookie.Value)
		}
		store := sessions.NewCookieStore([]byte(SESSION_SECRET_KEY))
		session,err := store.Get(r, "project101_session")
		if err != nil {
			log.Println(err)
		}else if session.Values["google_id"] != nil {
			log.Printf(session.Values["google_id"].(string))
		}
		a.DB.Find(&menus)
		tableJSON, _ = json.Marshal(menus)
	case 2:
		w.Header().Set("Access-Control-Allow-Origin", "*")
		var orders []Order
		a.DB.Find(&orders)
		tableJSON, _ = json.Marshal(orders)
	case 3:
		w.Header().Set("Access-Control-Allow-Origin", "*")
		var orders []*OrderStatus
		rows, err := a.DB.Table(
					"orders",
				).Select(
					"orders.id, orders.order_number, " +
					"orders.status, orders.chef_id, " +
					"orders.type, orders.payment_type, " +
					"orders.quantity, menus.name, " +
					"menus.description, menus.type, " +
					"menus.id, menus.price",
				).Joins(
					"left join menus on " +
					"menus.id = orders.dish_id",
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
		tableJSON, _ = json.Marshal(orders)
	case 4:
		w.Header().Set("Access-Control-Allow-Origin", "*")
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
				).Order(
					"order_number desc",
				).Rows()
		if err != nil {
			log.Printf("Error Code 1034a");
			log.Fatal(err)
		}
		defer order_numbers.Close()
		log.Printf("colected orderumbers")
		for order_numbers.Next() {
			var orders []*OrderStatus
			var order_number int64
			if err := order_numbers.Scan(&order_number,);
			   err != nil {
				log.Printf("Error Code 1034b");
				log.Fatal(err)
			}
			rows, err := a.DB.Table(
					"orders",
				).Select(
					"orders.id, orders.order_number, " +
					"orders.status, orders.chef_id, " +
					"orders.type, orders.payment_type, " +
					"orders.quantity, menus.name, " +
					"menus.description, menus.type, " +
					"menus.id, menus.price",
				).Joins(
					"left join menus on " +
					"menus.id = orders.dish_id",
				).Where(
					"order_number = ? AND status = ?",
					order_number,order_status,
				).Rows()
			if err != nil {
				log.Printf("Error Code 1034c");
				log.Fatal(err)
			}
			defer rows.Close()
			for rows.Next() {
				o := new(OrderStatus)
				if err := rows.Scan(
						&o.OrderId, &o.OrderNumber,
						&o.OrderStatus, &o.ChefId,
						&o.OrderType,
						&o.OrderPaymentType,
						&o.OrderQuantity, &o.DishName,
						&o.DishDescription, &o.DishType,
						&o.DishId, &o.DishPrice,
						); err != nil {
					log.Printf("Error Code 1034d");
					log.Fatal(err)
				}
				orders = append(orders,o)
			}
			if err := rows.Err(); err != nil {
				log.Printf("Error Code 1034e");
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

	log.Printf("collected orders in progress")
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
					"orders.id, orders.order_number, " +
					"orders.status, orders.chef_id, " +
					"orders.type, orders.payment_type, " +
					"orders.quantity, menus.name, " +
					"menus.description, menus.type, " +
					"menus.id, menus.price",
				).Joins(
					"left join menus on " +
					"menus.id = orders.dish_id",
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
		log.Printf("Error Code 1014a");
		log.Fatal(err)

	}
	defer order_numbers.Close()
	for order_numbers.Next() {
		var orders []*OrderStatus
		var order_number int64
		if err := order_numbers.Scan(
					&order_number,
				); err != nil {
				log.Printf("Error Code 1014b");
				log.Fatal(err)
			}
		rows, err := a.DB.Table(
					"orders",
				).Select(
					"orders.id, orders.order_number, " +
					"orders.status, orders.chef_id, " +
					"orders.type, orders.payment_type, " +
					"orders.quantity, menus.name, " +
					"menus.description, menus.type, " +
					"menus.id, menus.price",
				).Joins(
					"left join menus on " +
					"menus.id = orders.dish_id",
				).Where(
					"order_number = ?",
					order_number,
				).Rows()
		if err != nil {
			log.Printf("Error Code 1014c");
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
				log.Printf("Error Code 1014d");
				log.Fatal(err)
			}
			orders = append(orders,o)
		}

		if err := rows.Err(); err != nil {
			log.Printf("Error Code 1014e");
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
  w.Header().Set("Access-Control-Allow-Origin", "*")
  w.WriteHeader(200)
  w.Write([]byte(tableJSON))
}


func (a *App) CreateHandler(
		tableId int,
		w http.ResponseWriter,
		r *http.Request) {

	// Parse the POST body to populate r.PostForm.
	if err := r.ParseForm(); err != nil {
		panic("failed in ParseForm() call")
	}
	log.Printf("sdsd\n")
	body, _ := ioutil.ReadAll(r.Body)
	switch tableId {
		case 1:	me := Menu{}
			json.Unmarshal(body, &me)
			log.Printf(string(me.Name))
			a.DB.Create(&me)

			u, err := url.Parse(
				fmt.Sprintf(RP_MENU + "/%s", me.Name))
			if err != nil {
				panic("failed to form new Menu URL")
			}

			base, err := url.Parse(r.URL.String())
			if err != nil {
				panic("failed to parse request URL")
			}

			w.Header().Set(
				"Location",
				base.ResolveReference(u).String())
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.WriteHeader(201)
		case 2:	or := Order{}
			json.Unmarshal(body, &or)
			log.Println(or)
			or.ChefId = "1"
			a.DB.Create(&or)
			u, err := url.Parse(fmt.Sprintf(
						RP_MENU + "/%d",
						or.DishId))
			if err != nil {
				panic("failed to form new Menu URL")
			}

			base, err := url.Parse(r.URL.String())
			if err != nil {
				panic("failed to parse request URL")
			}

			w.Header().Set("Location",
					base.ResolveReference(u).String())
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Headers", "*")
			w.WriteHeader(201)
		case 3:
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Headers", "*")
			w.Header().Set("Access-Control-Allow-Methods", "*")
			w.WriteHeader(201)
	}
}

func (a *App) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// Parse the request body to populate r.PostForm.
	if err := r.ParseForm(); err != nil {
		panic("failed in ParseForm() call")
	}

        // Update the order number as completed 
	a.DB.Model(
		&Order{},
		).Where(
			"order_number = ?",
			vars["order_number"],
		).Update(
			"status", "completed",
		)

	// Write to HTTP response.
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(204)
	log.Printf("Updated Order: #%s as completed", vars["order_number"]);
	//sleep for some time to ensure database is not locked.
	//time.Sleep(2 * time.Second)
}

func (a *App) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// Delete the menu with the given name.
	a.DB.Where("name = ?", vars["name"]).Delete(Menu{})

	// Write to HTTP response.
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(204)
}

func dbUpdateCallback(scope *gorm.Scope) {
	if scope.HasError() {
		log.Printf("DB Update Error!!")
	}
	log.Printf("DB Updated")
}

func setupNewOuth() {
	maxAge := 86400 * 1
	isProd := false

	store := sessions.NewCookieStore([]byte(SESSION_SECRET_KEY))
	store.MaxAge(maxAge)
	store.Options.Path = "/"
	store.Options.HttpOnly = true
	store.Options.Secure = isProd

	gothic.Store = store

	goth.UseProviders(
		google.New(
			"294787796772-h5jpfgnq9rd3cf8roohi126f5em34q30.apps.googleusercontent.com",
			"xje35yKstxrwSDruOEmc9yDQ",
			"http://ckcserver.loca.lt/v1/auth/google/callback",
			"email",
			"profile"),
	)
}


func httpRoutingHandler(a *App, pool *Pool) {

	//For routing server for different paths -- rest approach
	//Using mux for routing
	r := mux.NewRouter()

	var r1 = r.PathPrefix("/v1").Subrouter()

	r1.HandleFunc(
		RP_MENU,
		func(w http.ResponseWriter, r *http.Request) {
			a.ListHandler(TABLEID["MENUS"], w, r)
		}).Methods("GET")
	r1.HandleFunc(
		"/profile",
		func(w http.ResponseWriter, r *http.Request) {
			ProfileHandler(w, r)
		}).Methods("GET")
	r1.HandleFunc(
		RP_BAKEMENU + "",
		func(w http.ResponseWriter, r *http.Request) {
			a.ListHandler(TABLEID["ORDERS"], w, r)
		}).Methods("GET")
	r1.HandleFunc(
		"/bakemenufull",
		func(w http.ResponseWriter, r *http.Request) {
			a.ListHandler(TABLEID["ORDERSTATUSFULL"], w, r)
		}).Methods("GET")
	r1.HandleFunc(
		"/bakemenubyorder",
		func(w http.ResponseWriter, r *http.Request) {
			a.ListHandler(TABLEID["ORDERSTATUSBYORDER"], w, r)
		}).Methods("GET")
	r1.HandleFunc(
		RP_MENU + "/{dish_id:.+}",
		func(w http.ResponseWriter, r *http.Request) {
			a.ViewHandler(TABLEID["MENUS"], w, r)
		}).Methods("GET")
	r1.HandleFunc(
		RP_BAKEMENU + "/{order_number:.+}",
		func(w http.ResponseWriter, r *http.Request) {
			a.ViewHandler(TABLEID["ORDERSTATUSBYORDER"], w, r)
		}).Methods("GET")
	r1.HandleFunc(
		RP_MENU,
		func(w http.ResponseWriter, r *http.Request) {
			a.CreateHandler(1, w, r)
		}).Methods("POST")
	r1.HandleFunc(
		"/order",
		func(w http.ResponseWriter, r *http.Request) {
			a.CreateHandler(2, w, r)
		}).Methods("POST")
	r1.HandleFunc(
		"/order",
		func(w http.ResponseWriter, r *http.Request) {
			a.CreateHandler(3, w, r)
		}).Methods("OPTIONS")
	r1.HandleFunc(
		RP_MENU_UPDATE + "/{order_number:.+}",
		a.UpdateHandler).Methods("PUT")
	r1.HandleFunc(
		RP_MENU_UPDATE + "/{order_number:.+}",
		func(w http.ResponseWriter, r *http.Request) {
			a.CreateHandler(3, w, r)
		}).Methods("OPTIONS")
	r1.HandleFunc(
		RP_MENU + "/{name:.+}",
		a.DeleteHandler,
		).Methods("DELETE")
	r1.HandleFunc(
		"/",
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "!!Backend Server running correctly!!")
		})
	r1.HandleFunc(
		"/ws",func(w http.ResponseWriter, r *http.Request) {
			a.serveWs(pool, w, r)
		})



	r1.HandleFunc("/auth/{provider}", func(res http.ResponseWriter, req *http.Request) {
		maxAge := 86400 * 1
		isProd := false

		store := sessions.NewCookieStore([]byte(SESSION_SECRET_KEY))
		store.MaxAge(maxAge)
		store.Options.Path = "/"
		store.Options.HttpOnly = true
		store.Options.Secure = isProd

		gothic.Store = store

		goth.UseProviders(
			google.New(
				"294787796772-h5jpfgnq9rd3cf8roohi126f5em34q30.apps.googleusercontent.com",
				"xje35yKstxrwSDruOEmc9yDQ",
				"http://ckcserver.loca.lt/v1/auth/google/callback",
				"email",
				"profile"),
		)
		log.Println("sdsds")
		log.Println(reflect.TypeOf(store))
		gothic.BeginAuthHandler(res, req)
	})

	r1.HandleFunc("/auth/{provider}/callback", func(res http.ResponseWriter, req *http.Request) {

		user, err := gothic.CompleteUserAuth(res, req)

		if err != nil {
			log.Fatal(err)
			return
		}

		log.Printf(user.Name)

		maxAge := 86400 * 1
		isProd := false

		store := sessions.NewCookieStore([]byte(SESSION_SECRET_KEY))
		store.MaxAge(maxAge)
		store.Options.Path = "/"
		store.Options.HttpOnly = true
		store.Options.Secure = isProd
		store.Options.Domain = "loca.lt"

		session,err := store.New(req, "project101_session")
		session.Values["google_id"] = user.Email
		session.Values["google_un"] = user.Name
		session.Values["google_at"] = user.AccessToken
		session.Save(req,res)
		for _, cookie := range req.Cookies() {
			log.Printf(cookie.Name)
			log.Printf(cookie.Value)
		}
		http.Redirect(
			res,
			req,
			"http://ckcclient.loca.lt",
			http.StatusFound)

	})

	r1.HandleFunc("/logout/{provider}", func(res http.ResponseWriter, req *http.Request) {

		gothic.Logout(res, req)

		store := sessions.NewCookieStore([]byte(SESSION_SECRET_KEY))
		store.Options.Domain = "loca.lt"

		session,err := store.Get(req, "project101_session")
		if err != nil {
			log.Println(err)
			log.Println("error|^")
		}
		session.Options.MaxAge = -1
		session.Values = make(map[interface{}]interface{})
		session.Save(req,res)

		http.Redirect(
			res,
			req,
			"http://ckcclient.loca.lt",
			http.StatusFound)
	})

	//go http library to handle http requests
	//http.Handle("/", r)

	//Starting backend server at 8080
	if err := http.ListenAndServe(":8080", r); err != nil {
		panic(err)
	}

}


func main() {

	//Initilise Database
	a := &App{}
	a.Initialize(DATABASE, DATABASE_NAME)
	//Callback functioin when and if database gets a new row -- gorm
	//a.DB.Callback().Create().Register(
	//			"gorm:after_create",
	//			dbUpdateCallback)
	defer a.DB.Close()

	//Running a thread for multiple websocket connections
	pool := newPool()
	go pool.start()

	//Setting required params for Google Authrisation
	setupNewOuth()

	//Setting up handlers for http api handling
	httpRoutingHandler(a, pool)

}

