package main

import (
	"encoding/json"
	"os"
	"strings"
)

type Location string

func (l Location) ToString(lang string) string {
	return strings.Replace(string(l), ":lang", lang, -1)
}

type OsName string
type ProductName string

type LocationsMap map[OsName]Location
type ProductInfo struct {
	Locations LocationsMap
	SSLOnly   bool
}
type ProductLocationsMap map[ProductName]ProductInfo

type LocationsFile struct {
	Products []struct {
		Name      string
		SSLOnly   bool `json:"ssl_only"`
		Locations map[OsName]Location
	}
}

func ParseLocationsFile(file string) (ProductLocationsMap, error) {
	plm := make(ProductLocationsMap)
	f, err := os.Open(file)
	defer f.Close()
	if err != nil {
		return plm, err
	}

	lf := LocationsFile{}
	err = json.NewDecoder(f).Decode(&lf)
	if err != nil {
		return plm, err
	}
	for _, product := range lf.Products {
		pName := ProductName(strings.TrimSpace(strings.ToLower(product.Name)))
		plm[pName] = ProductInfo{
			SSLOnly:   product.SSLOnly,
			Locations: make(LocationsMap),
		}

		for os, path := range product.Locations {
			plm[pName].Locations[os] = path
		}
	}
	return plm, nil
}
