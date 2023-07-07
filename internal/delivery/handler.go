package delivery

import (
	"github.com/sirupsen/logrus"
	"log"
	"mailService/internal/models"
	"net/http"
)

type EmailService interface {
	AddUser(mail string) (models.Email, error)
	CheckUserByKeyword(keyword string) ([]byte, error) // check if user exists
}

type Handler struct {
	// access to business logic
	userEmailService EmailService
	logger           *logrus.Logger
}

func NewHandler(userEmailService EmailService) *Handler {
	logger := logrus.New()
	return &Handler{
		userEmailService: userEmailService,
		logger:           logger,
	}
}

func (h *Handler) Init() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/add-user-mail", h.addUserMail)
	mux.HandleFunc("/get-user-zip", h.returnUserEmailZip)

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

	ps, err := h.userEmailService.AddUser(username)
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

func (h *Handler) returnUserEmailZip(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		msg := "method not provide!"
		_, err := w.Write([]byte(msg))
		if err != nil {
			return
		}
	}

	key := r.URL.Query().Get("key")

	b, err := h.userEmailService.CheckUserByKeyword(key)
	if err != nil {
		h.logger.Error(err)
		http.Error(w, err.Error(), 200)
		return
	}

	w.Header().Set("Content-Type", "application/zip")
	w.Write(b)
}
