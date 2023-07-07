package delivery

import (
	"log"
	"mailService/internal/service"
	"net/http"
)

type Handler struct {
	// access to business logic
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{
		service: service,
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
		// TODO use logging
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return
	}

	log.Println(ps)
	_, err = w.Write([]byte(ps.UniqueCode))
	if err != nil {
		// TODO USE logger
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
		log.Println(err)
		http.Error(w, err.Error(), 200)
		return
	}

	w.Write(b)
	// implement the logic
}
