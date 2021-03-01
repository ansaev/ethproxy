package api

import (
	"encoding/json"
	"ethproxy/internal/domain"
	"fmt"
	"log"

	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

const (
	BlockIDParam = "blockID"
	TxIDParam    = "txID"
)

func New(listeningAddress string, txService txFinder) *service {
	srv := &service{txService: txService}

	// set routing
	handler := srv.getServiceHandler()

	// init http server
	srv.server = &http.Server{
		Addr:    listeningAddress,
		Handler: handler,
	}

	return srv
}

type txFinder interface {
	FindTx(blockID string, txID string) (*domain.Transaction, error)
}

type service struct {
	server    *http.Server
	txService txFinder
}

func (srv *service) getServiceHandler() chi.Router {
	router := chi.NewRouter()
	router.Use(
		middleware.RequestID,
		middleware.Logger,
	)

	router.Route("/tx/", func(r chi.Router) {
		router.Use(
			middleware.Recoverer,
			middleware.NoCache,
			middleware.SetHeader("Content-Type", "application/json"),
		)

		router.Get(fmt.Sprintf("/block/{%s}/tx/{%s}", BlockIDParam, TxIDParam), srv.txHandler)
	})

	return router

}

func (srv *service) txHandler(w http.ResponseWriter, r *http.Request) {
	// get and validate request's data
	blockID := chi.URLParam(r, BlockIDParam)
	txID := chi.URLParam(r, TxIDParam)
	if blockID == "" {
		w.WriteHeader(http.StatusBadRequest)
		errWrite := json.NewEncoder(w).Encode(&TxResponse{Ok: false, Error: &ErrorForm{"empty block id"}})
		if errWrite != nil {
			log.Printf("failed to write response #1: %v\n", errWrite)
		}
		return
	}
	if txID == "" {
		w.WriteHeader(http.StatusBadRequest)
		errWrite := json.NewEncoder(w).Encode(&TxResponse{Ok: false, Error: &ErrorForm{"empty tx id"}})
		if errWrite != nil {
			log.Printf("failed to write response #2: %v\n", errWrite)
		}
		return
	}

	// find tx
	tx, errFindTx := srv.txService.FindTx(blockID, txID)
	// return response
	if errFindTx != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errWrite := json.NewEncoder(w).Encode(&TxResponse{
			Ok: false, Error: &ErrorForm{fmt.Sprintf("failed to get tx: %v", errFindTx)}})
		if errWrite != nil {
			log.Printf("failed to write response #3: %v\n", errWrite)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	errWrite := json.NewEncoder(w).Encode(&TxResponse{Ok: true, Data: tx})
	if errWrite != nil {
		log.Printf("failed to write response #4: %v\n", errWrite)
	}
}

func (srv *service) ListenAndServe() error {
	return srv.server.ListenAndServe()
}
