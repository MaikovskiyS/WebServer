package service

import (
	"context"
	"stadium/internal/models"

	"github.com/sirupsen/logrus"
)

type Repository interface {
	GetCapital(ctx context.Context, country *models.Country) (capital *models.City, err error)
	GetCities(ctx context.Context) ([]*models.City, error)
	GetAllCountries(ctx context.Context) ([]*models.Country, error)
	GetCitiesByCountry(ctx context.Context, country *models.Country) (arr []*models.City, err error)
	Close() error
}

type Service interface {
	GetCapital(ctx context.Context, country *models.Country) (capital *models.City, err error)
	GetSities(ctx context.Context) ([]*models.City, error)
	GetAllCountries(ctx context.Context) ([]*models.Country, error)
	GetCitiesByCountry(ctx context.Context, country *models.Country) (arr []*models.City, err error)
	// DeleteCountry()
	// GetCapitalByCountry()
}

type service struct {
	log        logrus.Logger
	repository Repository
}

// NewService service constractor
func NewService(log *logrus.Logger, r Repository) *service {
	return &service{
		log:        *log,
		repository: r,
	}
}

// GetSities get cities
func (s *service) GetSities(ctx context.Context) ([]*models.City, error) {
	s.log.Info("Get sities inservice")
	return s.repository.GetCities(ctx)
}

// GetAllCountries get all countries
func (s *service) GetAllCountries(ctx context.Context) ([]*models.Country, error) {
	s.log.Info("Get countries in service")

	return s.repository.GetAllCountries(ctx)
}
func (s *service) GetCitiesByCountry(ctx context.Context, country *models.Country) (arr []*models.City, err error) {
	s.log.Info("GetCitiesByCountry in servise")
	return s.repository.GetCitiesByCountry(ctx, country)
}
func (s *service) GetCapital(ctx context.Context, country *models.Country) (capital *models.City, err error) {
	s.log.Info("Get capital in service")
	return s.repository.GetCapital(ctx, country)
}
