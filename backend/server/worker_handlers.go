package server

import (
	"encoding/json"
	"eqweqr/bdkurach/controllers"
	"log"
	"net/http"
)

// посмотреть все заказы которые ожидают работы (на фронте надо чтобы разделялось какие предложения адресованные и нет)
func (server *Server) GetAllSuggestionsWorkerHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	// user_id := r.Request.Context().Value("id")
	// id := r.URL.Query().Get("id")
	// id := r.URL.Query().Get("id")
	id := r.Context().Value("id").(string)
	// user_id := r.Request.URL.Query().Get("user_id")
	sugs, err := controllers.GetAllSuggestionsWorker(id, server.DB)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(sugs)
}

// посмотреть все заказы которые взял пользователь(можно 'processing', 'done')
// func (server *Server) GetAllOrdersWorkerHandler(w http.ResponseWriter, r *http.Request) {
// 	id := r.URL.Query().Get("id")
// 	status := r.URL.Query().Get("status")
// 	orders, err := controllers.GetAllStatusOrdersWorker(id, status, server.DB)
// 	if err != nil {
// 		log.Println(err)
// 		w.WriteHeader(http.StatusBadRequest)
// 		return
// 	}
// 	json.NewEncoder(w).Encode(orders)

// }

// сделать пользователю предложение
func (server *Server) MakeSuggestionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	// id := r.URL.Query().Get("id")
	// id := r.URL.Query().Get("id")
	id := r.Context().Value("id").(string)
	order_id := r.URL.Query().Get("order_id")

	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	summary := r.FormValue("summary")
	term := r.FormValue("term")

	log.Println(id, order_id, summary, term)

	err := controllers.CreateSuggestin(order_id, id, summary, term, server.DB)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

// посмотреть какие выплаты(прогнозируемые на основе заказов которые выполняются, и сумма которая есть на основе )
func (server *Server) GetAllOrdersWorkerHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	days := r.URL.Query().Get("days")
	// id := r.URL.Query().Get("id")
	// id := r.URL.Query().Get("id")
	id := r.Context().Value("id").(string)
	// status := r.URL.Query().Get("status")

	if days == "all" {
		orders, err := controllers.GetAllStatusOrdersWorker(id, server.DB)

		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if err := json.NewEncoder(w).Encode(orders); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		return
	}
	orders, err := controllers.GetAllOrdersByTime(id, days, server.DB)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(orders)
}

// получить информацию о выплатах
func (server *Server) TotalSalaryHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	days := r.URL.Query().Get("days")
	// id := r.URL.Query().Get("id")
	// id := r.URL.Query().Get("id")
	id := r.Context().Value("id").(string)

	if days == "all" {
		sum, err := controllers.GetTotalSummary(id, server.DB)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if err := json.NewEncoder(w).Encode(sum); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		return
	}

	sum, err := controllers.GetTotalSummaryByTime(id, days, server.DB)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := json.NewEncoder(w).Encode(sum); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	return
}

func (server *Server) GetAllOrdersByStatusWorkerHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	// id := r.URL.Query().Get("id")
	// id := r.URL.Query().Get("id")
	id := r.Context().Value("id").(string)
	status := r.URL.Query().Get("status")
	orders, err := controllers.GetAllOrderStatusWorker(id, status, server.DB)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err = json.NewEncoder(w).Encode(orders); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func (server *Server) GetAllsugavaitHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	// id := r.URL.Query().Get("id")
	id := r.Context().Value("id").(string)
	suggestions, err := controllers.GetAllSugessiongWorker(id, server.DB)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err = json.NewEncoder(w).Encode(suggestions); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func (server *Server) AllSalaryHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	// id := r.URL.Query().Get("id")
	id := r.Context().Value("id").(string)
	sum, err := controllers.GetTotalSummary(id, server.DB)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err = json.NewEncoder(w).Encode(sum); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func (server *Server) SalaryByTimeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	id := r.Context().Value("id").(string)
	log.Println(id)
	// id := r.URL.Query().Get("id")
	days := r.URL.Query().Get("days")
	sum, err := controllers.GetTotalSummaryByTime(id, days, server.DB)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err = json.NewEncoder(w).Encode(sum); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func (server *Server) HandleMySuggestionsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	id := r.Context().Value("id").(string)

	// id := r.URL.Query().Get("id")
	// id := r.URL.Query().Get("id")
	sum, err := controllers.GetAllWorkerOwnSuggestions(id, server.DB)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err = json.NewEncoder(w).Encode(sum); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func (server *Server) HandleAllWorkerSuggestionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	// id := r.URL.Query().Get("id")
	// id := r.URL.Query().Get("id")
	id := r.Context().Value("id").(string)
	log.Println(id)

	sum, err := controllers.GetAllWorkerSuggestions(id, server.DB)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err = json.NewEncoder(w).Encode(sum); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func (server *Server) CreateNewSuggestHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	// /api/worker/all/suggestion/send?id=2&orderid=${orderid}&cost=${cost}&term=${term}
	// id := r.URL.Query().Get("id")
	// id := r.URL.Query().Get("id")
	id := r.Context().Value("id").(string)

	orderid := r.URL.Query().Get("orderid")
	cost := r.URL.Query().Get("cost")
	term := r.URL.Query().Get("term")
	log.Println(id, orderid, cost, term)
	err := controllers.CreateSuggest(id, orderid, cost, term, server.DB)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}
