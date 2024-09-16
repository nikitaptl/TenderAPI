package main

import (
	"avitoTestTask/internal/database"
	"avitoTestTask/internal/handlers/bid_handler"
	"avitoTestTask/internal/handlers/feedback_handler"
	"avitoTestTask/internal/handlers/tender_handler"
	"avitoTestTask/internal/migrations"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

func main() {
	addr := os.Getenv("SERVER_ADDRESS")
	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("Не удалось подключиться к базе данных: %v", err)
		return
	}
	defer func() {
		if err := database.CloseDB(); err != nil {
			log.Fatalf("Ошибка при закрытии соединения с базой данных: %v", err)
		}
	}()
	log.Println("База данных успешно запущена!")

	if err = migrations.RunMigrations(db); err != nil {
		log.Fatalf("Ошибка при миграции базы данных: %v", err)
		return
	}
	log.Println("Миграция базы данных успешно завершена!")

	apiRouter := setupRouter()
	http.ListenAndServe(addr, apiRouter)
}

func setupRouter() *mux.Router {
	apiRouter := mux.NewRouter().PathPrefix("/api").Subrouter()

	apiRouter.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "ok")
	})

	// Тендеры
	apiRouter.HandleFunc("/tenders", tender_handler.Tenders).Methods(http.MethodGet)
	apiRouter.HandleFunc("/tenders/my", tender_handler.TendersMy).Methods(http.MethodGet)
	apiRouter.HandleFunc("/tenders/new", tender_handler.NewTender).Methods(http.MethodPost)
	apiRouter.HandleFunc("/tenders/{tenderId}/status", tender_handler.GetTenderStatus).Methods(http.MethodGet)
	apiRouter.HandleFunc("/tenders/{tenderId}/status", tender_handler.UpdateTenderStatus).Methods(http.MethodPut)
	apiRouter.HandleFunc("/tenders/{tenderId}/edit", tender_handler.UpdateTenderHandler).Methods(http.MethodPatch)
	apiRouter.HandleFunc("/tenders/{tenderId}/rollback/{version}", tender_handler.RollbackTender).Methods(http.MethodPut)

	// Предложения
	apiRouter.HandleFunc("/bids/new", bid_handler.NewBid).Methods(http.MethodPost)
	apiRouter.HandleFunc("/bids/my", bid_handler.BidsMy).Methods(http.MethodGet)
	apiRouter.HandleFunc("/bids/{tenderId}/list", bid_handler.BidsList).Methods(http.MethodGet)
	apiRouter.HandleFunc("/bids/{bidId}/status", bid_handler.GetBidStatus).Methods(http.MethodGet)
	apiRouter.HandleFunc("/bids/{bidId}/status", bid_handler.UpdateBidStatus).Methods(http.MethodPut)

	apiRouter.HandleFunc("/bids/{bidId}/submit_decision", bid_handler.SubmitDecisionBid).Methods(http.MethodPut)
	apiRouter.HandleFunc("/bids/{bidId}/edit", bid_handler.UpdateBid).Methods(http.MethodPatch)
	apiRouter.HandleFunc("/bids/{bidId}/rollback/{version}", bid_handler.RollbackBid).Methods(http.MethodPut)

	// Отзывы
	apiRouter.HandleFunc("/bids/{bidId}/feedback", feedback_handler.CreateNewFeedback).Methods(http.MethodPut)
	apiRouter.HandleFunc("/bids/{tenderId}/reviews", feedback_handler.GetFeedbackList).Methods(http.MethodGet)

	return apiRouter
}
