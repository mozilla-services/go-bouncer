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
type AliasName string

func NewProductName(s string) ProductName {
	return ProductName(strings.TrimSpace(strings.ToLower(s)))
}

func NewAliasName(s string) AliasName {
	return AliasName(strings.TrimSpace(strings.ToLower(s)))
}

type LocationsMap map[OsName]Location
type ProductInfo struct {
	Locations LocationsMap
	SSLOnly   bool
}

type BouncerMap struct {
	ProductLocationMap map[ProductName]ProductInfo
	Aliases            map[AliasName]ProductName
}
type ProductLocationsMap map[ProductName]ProductInfo

type LocationsFile struct {
	Products []struct {
		Name      string
		SSLOnly   bool `json:"ssl_only"`
		Locations map[OsName]Location
	}
	Aliases map[string]string
}

func ParseLocationsFile(file string) (BouncerMap, error) {
	bouncerMap := BouncerMap{
		ProductLocationMap: make(ProductLocationsMap),
		Aliases:            make(map[AliasName]ProductName),
	}

	f, err := os.Open(file)
	defer f.Close()
	if err != nil {
		return bouncerMap, err
	}

	lf := LocationsFile{}
	err = json.NewDecoder(f).Decode(&lf)
	if err != nil {
		return bouncerMap, err
	}

	plm := bouncerMap.ProductLocationMap
	for _, product := range lf.Products {
		pName := NewProductName(product.Name)
		plm[pName] = ProductInfo{
			SSLOnly:   product.SSLOnly,
			Locations: make(LocationsMap),
		}

		for os, path := range product.Locations {
			plm[pName].Locations[os] = path
		}
	}

	aliases := bouncerMap.Aliases
	for alias, product := range lf.Aliases {
		aliasName := NewAliasName(alias)
		productName := NewProductName(product)
		aliases[aliasName] = productName
	}
	return bouncerMap, nil
}
