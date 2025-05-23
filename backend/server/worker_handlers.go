package server

import (
	"encoding/json"
	"eqweqr/bdkurach/controllers"
	"log"
	"net/http"
)

// посмотреть все заказы которые ожидают работы (на фронте надо чтобы разделялось какие предложения адресованные и нет)
func (server *Server) GetAllSuggestionsWorkerHandler(w http.ResponseWriter, r *http.Request) {
	// user_id := r.Request.Context().Value("id")
	id := r.URL.Query().Get("id")
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
	id := r.URL.Query().Get("id")
	order_id := r.URL.Query().Get("order_id")

	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	summary := r.FormValue("summary")
	term := r.FormValue("term")

	err := controllers.CreateSuggestin(order_id, id, summary, term, server.DB)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

// посмотреть какие выплаты(прогнозируемые на основе заказов которые выполняются, и сумма которая есть на основе )
func (server *Server) GetAllOrdersWorkerHandler(w http.ResponseWriter, r *http.Request) {
	days := r.URL.Query().Get("days")
	id := r.URL.Query().Get("id")
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
	days := r.URL.Query().Get("days")
	id := r.URL.Query().Get("id")

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
	id := r.URL.Query().Get("id")
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
	id := r.URL.Query().Get("id")
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
	id := r.URL.Query().Get("id")
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
	id := r.URL.Query().Get("id")
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
