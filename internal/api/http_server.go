package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ethereum_parser/internal/types"
)

type HTTPServer struct {
	parser types.Parser
}

func NewHTTPServer(p types.Parser) *HTTPServer {
	return &HTTPServer{parser: p}
}

func (s *HTTPServer) Start(addr string) error {
	http.HandleFunc("/subscribe", s.handleSubscribe)
	http.HandleFunc("/transactions", s.handleGetTransactions)
	http.HandleFunc("/current-block", s.handleGetCurrentBlock)

	log.Printf("Starting HTTP server on %s", addr)
	return http.ListenAndServe(addr, nil)
}

// subscribe the address
func (s *HTTPServer) handleSubscribe(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Address string `json:"address"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	success := s.parser.Subscribe(req.Address)
	json.NewEncoder(w).Encode(map[string]bool{"success": success})
}

// get transction for given address
func (s *HTTPServer) handleGetTransactions(w http.ResponseWriter, r *http.Request) {
	address := r.URL.Query().Get("address")
	if address == "" {
		http.Error(w, "address is required", http.StatusBadRequest)
		return
	}

	txs, err := s.parser.GetTransactions(address)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(txs)
}

// get the current block
func (s *HTTPServer) handleGetCurrentBlock(w http.ResponseWriter, r *http.Request) {
	block, err := s.parser.GetCurrentBlock()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]int64{"block": block})
}
