package app

import (
	"net/http"

	"github.com/Sonicspeedfly/flymarket/cmd/app/middleware"
	"github.com/Sonicspeedfly/flymarket/pkg/accounts"
	"github.com/Sonicspeedfly/flymarket/pkg/product"
	"github.com/gorilla/mux"
)

type Server struct {
	mux        *mux.Router
	accountSvc *accounts.Service
	productSvc *product.Service
}

func NewServer(mux *mux.Router, accountSvc *accounts.Service, productSvc *product.Service) *Server {
	return &Server{mux: mux, accountSvc: accountSvc, productSvc: productSvc}
}

func (s *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	s.mux.ServeHTTP(writer, request)
}

func (s *Server) Init() {
	accountsSubrouter := s.mux.PathPrefix("/api/accounts").Subrouter()
	accountsSubrouter.Use(middleware.Authenticate(s.accountSvc.IDByToken))
	accountsSubrouter.Use(middleware.CheckRole(s.accountSvc.HasAnyRole, "CUSTOMER", "ADMIN"))
	accountsSubrouter.HandleFunc("", s.handleSaveAccount).Methods(http.MethodGet)
	accountsSubrouter.HandleFunc("/accounts.ByID", s.handleByIdAccount).Methods(http.MethodGet)
	accountsSubrouter.HandleFunc("/accounts.Remove", s.handleRemoveAccount).Methods(http.MethodGet)
	accountsSubrouter.HandleFunc("/accounts.All", s.handleAllAccount).Methods(http.MethodGet)
	accountsSubrouter.HandleFunc("/login", s.handleLoginAccount).Methods(http.MethodGet)
	accountsSubrouter.HandleFunc("/autification", s.handleAutificationAccount).Methods(http.MethodGet)
	
	productSubrouter := s.mux.PathPrefix("/api/products").Subrouter()
	productSubrouter.HandleFunc("/save", s.handleSaveProduct).Methods(http.MethodPost)
	productSubrouter.HandleFunc("/product.edit", s.handleEditProduct).Methods(http.MethodPost)
	productSubrouter.HandleFunc("/product.ByID", s.handleByIdProduct).Methods(http.MethodGet)
	productSubrouter.HandleFunc("/All", s.handleAllProduct).Methods(http.MethodGet)
	productSubrouter.HandleFunc("/product.Remove", s.handleRemoveProduct).Methods(http.MethodGet)
	productSubrouter.HandleFunc("/Buy", s.handleBuyProduct).Methods(http.MethodGet)
	productSubrouter.HandleFunc("/", notFound).Methods(http.MethodGet)
}

//notFound запрос notFound
func notFound(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusNotFound)
	writer.Write([]byte("sorry, page not found"))
}

