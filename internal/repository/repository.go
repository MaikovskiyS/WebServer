package repository

import (
	"context"
	"database/sql"
	"fmt"
	"stadium/internal/models"
	"stadium/internal/service"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
)

type mySQL struct {
	log     *logrus.Logger // здесь не надо но пока не дойдем до логера можно
	db      *sql.DB        // должно быть не экспоритируемо db
	timeOut time.Duration
}

// NewRepository package constructor
func NewRepository(connString string, log *logrus.Logger, timeOut time.Duration) (service.Repository, error) {
	db, err := sql.Open("mysql", connString)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeOut)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return &mySQL{
		db:      db,
		log:     log,
		timeOut: timeOut,
	}, nil
}

// Close ...
func (m *mySQL) Close() error {
	return m.db.Close()
}

// GetCities ...
func (m *mySQL) GetCities(ctx context.Context) (arr []*models.City, err error) {
	m.log.Info("GetSities in repositiry")
	execCTX, cancel := context.WithTimeout(ctx, m.timeOut)
	defer cancel()
	r, err := m.db.QueryContext(execCTX, "SELECT `name`,`population`FROM world.city WHERE `population`>1000000")
	if err != nil {
		m.log.Info("Cant send request to db")
		return nil, err
	}

	for r.Next() {
		var city models.City
		err := r.Scan(&city.Name, &city.Population)
		if err != nil {
			m.log.Info("cant write data from db to model city")
			return nil, err
		}
		arr = append(arr, &city)
	}
	return arr, nil
}

func (m *mySQL) GetCitiesByCountry(ctx context.Context, country *models.Country) (arr []*models.City, err error) {
	m.log.Info("Get Cities By Country From Repository")
	execCTX, cancel := context.WithTimeout(ctx, m.timeOut)
	defer cancel()
	rows, err := m.db.QueryContext(execCTX, "SELECT city.Name,city.countryCode,city.Population FROM city WHERE city.CountryCode=? ORDER BY city.Name", country.Code)
	fmt.Println(country.Code)
	if err != nil {
		m.log.Info("cant get rows from db")
		return nil, err
	}
	for rows.Next() {
		var city models.City
		err = rows.Scan(&city.Name, &city.CountryCode, &city.Population)
		if err != nil {
			m.log.Info("cant scan data from db to city")
			return nil, err
		}
		arr = append(arr, &city)
	}
	return arr, nil
}
func (m *mySQL) GetCapital(ctx context.Context, country *models.Country) (capital *models.City, err error) {
	execCTX, cancel := context.WithTimeout(ctx, m.timeOut)
	defer cancel()
	rows, err := m.db.QueryContext(execCTX, "SELECT city.`name`,city.`Countrycode`,city.`population` from city join country on city.id=country.Capital where country.Code=?", country.Code)
	if err != nil {
		m.log.Info("cant get data from db")
	}
	for rows.Next() {
		var city models.City
		err := rows.Scan(&city.Name, &city.CountryCode, &city.Population)
		if err != nil {
			m.log.Info("cant scan data from db")
		}
		capital = &city
	}
	return capital, nil
}

// GetAllCountries ...
func (m *mySQL) GetAllCountries(ctx context.Context) (arr []*models.Country, err error) {
	m.log.Info("Get Countries From Repository")
	execCTX, cancel := context.WithTimeout(ctx, m.timeOut)
	defer cancel()

	rows, err := m.db.QueryContext(execCTX, "SELECT `name`,`code`,`capital` FROM world.country WHERE Population>100000000")
	if err != nil {
		m.log.Info("cant get rows from db")
		return nil, err
	}
	for rows.Next() {
		var country models.Country

		err = rows.Scan(&country.Name, &country.Code, &country.Capital)
		if err != nil {
			return nil, err
		}
		arr = append(arr, &country)
	}
	return arr, nil
}
