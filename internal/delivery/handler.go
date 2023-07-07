package delivery

import (
	"github.com/sirupsen/logrus"
	"log"
	"mailService/internal/service"
	"net/http"
)

type Handler struct {
	// access to business logic
	service *service.Service
	logger  *logrus.Logger
}

func NewHandler(service *service.Service) *Handler {
	logger := logrus.New()
	return &Handler{
		service: service,
		logger:  logger,
	}
}

func (h *Handler) Init() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/add-user-mail", h.addUserMail)
	mux.HandleFunc("/get-user-zip", h.returnZIP)

	return mux
}

func (h *Handler) addUserMail(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		msg := "method not provide!"
		_, err := w.Write([]byte(msg))
		if err != nil {
			return
		}
	}

	username := r.URL.Query().Get("mail")

	ps, err := h.service.EmailService.AddUser(username)
	if err != nil {
		h.logger.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println(ps)
	_, err = w.Write([]byte(ps.UniqueCode))
	if err != nil {
		h.logger.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) returnZIP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		msg := "method not provide!"
		_, err := w.Write([]byte(msg))
		if err != nil {
			return
		}
	}

	key := r.URL.Query().Get("key")

	b, err := h.service.EmailService.CheckUserByKeyword(key)
	if err != nil {
		h.logger.Error(err)
		http.Error(w, err.Error(), 200)
		return
	}

	w.Write(b)
	// implement the logic
}
