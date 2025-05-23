package server

import (
	"encoding/json"
	"eqweqr/bdkurach/controllers"
	"fmt"
	"log"
	"net/http"
)

// создать новый заказ(для прода переделать, чтобы id брался из контекста.)
func (server *Server) CreateOrderHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")

	log.Println("new order")
	// if r.Method != http.MethodPost {
	// 	log.Println("not post")
	// 	w.WriteHeader(http.StatusServiceUnavailable)
	// 	return
	// }
	if err := r.ParseForm(); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	type t struct {
		Name    string `json:"modelName"`
		Warning bool   `json:"isWarranty"`
		Comment string `json:"comment"`
		Worker  string `json:"selectedRepairman"`
		Device  string `json:"deviceType"`
	}

	id := r.URL.Query().Get("id")

	var tmp t
	err := d.Decode(&tmp)
	log.Println(&tmp)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// id_c := r.Context().Value("id").(string)
	iid, err := controllers.CreateNewOrder(tmp.Name, tmp.Comment, id, "pending", tmp.Warning, tmp.Device, tmp.Worker, server.DB)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	err = json.NewEncoder(w).Encode(iid)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// посмотреть все предложения для пользователя
// дописать проверку по id через контекст
func (server *Server) ShowSuggestionByOrderHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("we")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	// проверка что заказ принадлежт этому человеку
	// id := r.Context().Value("id").(string)
	user_id := r.URL.Query().Get("id")
	// if err := controllers.CheckOrderBelong(order_id, id, server.DB); err != nil {
	// w.WriteHeader(http.StatusOK)
	// json.NewEncoder(w).Encode(nil)
	// return
	// }

	suggestions, err := controllers.GetAllSuggestions(user_id, server.DB)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(suggestions)
}

// посмотреть все свои заказы pending
// 'pending', 'processing', 'done'
func (server *Server) ShowByStatusOrdersHandler(w http.ResponseWriter, r *http.Request) {
	//id_c := r.Context().Value("id").(string)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	id := r.URL.Query().Get("id")
	status := r.URL.Query().Get("status")
	log.Println(id, status)
	orders, err := controllers.GetAllStatusOrders(id, status, server.DB)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(orders)
}

// cancel order only it has pending
func (server *Server) CancelOrderHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	// id := r.Context().Value("id").(string)
	user_id := r.URL.Query().Get("user")

	order_id := r.URL.Query().Get("id")
	err := controllers.CheckOrderBelong(order_id, user_id, server.DB)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = controllers.CancelOrder(order_id, server.DB)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

// assign sugest to your order
func (server *Server) AssignWorkerToOrder(w http.ResponseWriter, r *http.Request) {
	log.Println("re")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	sug_id := r.URL.Query().Get("id")
	user_id := r.URL.Query().Get("user")
	err := controllers.SuggestionOrder(sug_id, user_id, server.DB)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = controllers.AssignWorkerToOrder(sug_id, server.DB)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

// посмотреть всех доступных работников
func (server *Server) GetAllWorkersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	workers, err := controllers.GetAllWorkersName(server.DB)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(workers)
}

// выбрать профиль работы
func (server *Server) GetAllTypesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	types, err := controllers.GetAllTypes(server.DB)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(types)
}

// подтвердить что работник выполнил заказ
func (server *Server) ApproveWorkHandler(w http.ResponseWriter, r *http.Request) {
	// id := r.Context().Value("id")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	id := r.URL.Query().Get("user_id")
	order_id := r.URL.Query().Get("id")

	err := controllers.CheckOrderBelong(order_id, id, server.DB)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = controllers.ApproveWork(order_id, server.DB)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}
