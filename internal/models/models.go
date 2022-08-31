package models

type City struct {
	Name        string `json:"name"`
	CountryCode string `json:"countrycode"`
	Population  int    `json:"population"`
}

type Country struct {
	Code    string `json:"code"`
	Name    string `json:"name"`
	Capital int
}
