package contorller

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"stadium/internal/models"
	"stadium/internal/service"

	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
)

type Contorller struct {
	log     *logrus.Logger
	service service.Service
	chi     *chi.Mux
	ctx     context.Context
	cancel  context.CancelFunc
}

// NewController ...
func NewController(logger *logrus.Logger, svc service.Service) *Contorller {
	c := chi.NewRouter()
	ctx, cancel := context.WithCancel(context.Background())
	return &Contorller{
		log:     logger,
		service: svc,
		chi:     c,
		ctx:     ctx,
		cancel:  cancel,
	}
}

func (c *Contorller) CreatingEndPoints() *Contorller {
	c.chi.Get("/cities", c.GetCities)
	c.chi.Post("/getcapital", c.GetCapital)
	c.chi.Get("/", c.GetAllCountries)
	c.chi.Post("/getcitiesbycountry", c.GetCitiesByCountry)
	return c
}

func (c *Contorller) Start() error {
	return http.ListenAndServe(":8080", c.chi)
}

func (c *Contorller) Stop() {
	fmt.Println("cancel context")
	c.cancel()
}

// GET
func (c *Contorller) GetCities(w http.ResponseWriter, r *http.Request) {
	c2, err := c.service.GetSities(c.ctx)
	if err != nil {
		c.log.Info("cant get cities from repository")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Connect-Type", "application/json")
	writeJSON(w, c2)

}

// Get capital
func (c *Contorller) GetCapital(w http.ResponseWriter, r *http.Request) {
	country := &models.Country{}
	err := json.NewDecoder(r.Body).Decode(country)
	if err != nil {
		c.log.Info("cant decode request body ")
	}
	capital, err := c.service.GetCapital(c.ctx, country)
	if err != nil {
		c.log.Info("cant get capital from service")
		return
	}
	writeJSON(w, capital)
}

// Get countries
func (c *Contorller) GetAllCountries(w http.ResponseWriter, r *http.Request) {

	data, err := c.service.GetAllCountries(c.ctx)
	fmt.Println(data)
	if err != nil {
		logrus.Info(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Ошибка на сервере"))
		return
	}
	writeJSON(w, data)
}

// Get cities by country POST
func (c *Contorller) GetCitiesByCountry(w http.ResponseWriter, r *http.Request) {
	logrus.Info("Parse request")
	country := &models.Country{}
	json.NewDecoder(r.Body).Decode(country)
	c.log.Info("decode request")
	defer r.Body.Close()

	citiesArr, err := c.service.GetCitiesByCountry(c.ctx, country)
	if err != nil {
		c.log.Info("cant get cities array")
	}
	writeJSON(w, citiesArr)

}

// UTILS
// write JSON
func writeJSON(w http.ResponseWriter, v interface{}) {
	data, err := json.Marshal(v)
	if err != nil {
		logrus.Debug()
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Connect-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func ParseJSONFromRequest(r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}
