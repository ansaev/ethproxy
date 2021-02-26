package api

import (
	"encoding/json"
	"ethproxy/internal/domain"
	"fmt"

	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

const (
	BlockIDParam = "blockID"
	TxIDParam    = "txID"
	Latest       = "latest"
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

// TODO: add structs in errors and in responces
func (srv *service) txHandler(w http.ResponseWriter, r *http.Request) {
	// get and validate request's data
	blockID := chi.URLParam(r, BlockIDParam)
	txID := chi.URLParam(r, TxIDParam)
	if blockID == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("empty block id"))
		return
	}
	if txID == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("empty block id"))
		return
	}

	// find tx
	tx, errFindTx := srv.txService.FindTx(blockID, txID)
	// return response
	if errFindTx != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("failed to get tx: %v", errFindTx)))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tx)
}

func (srv *service) ListenAndServe() error {
	return srv.server.ListenAndServe()
}
