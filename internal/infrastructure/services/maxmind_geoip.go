package services

import (
	"context"
	"fmt"
	"net"

	"github.com/oschwald/geoip2-golang"
	appContracts "github.com/victorotene80/authentication_api/internal/application/contracts"
)

type MaxMindGeoIPService struct {
	db *geoip2.Reader
}

func NewMaxMindGeoIPService(dbPath string) (*MaxMindGeoIPService, error) {
	db, err := geoip2.Open(dbPath)
	if err != nil {
		return nil, fmt.Errorf("open maxmind db: %w", err)
	}
	return &MaxMindGeoIPService{db: db}, nil
}

func (s *MaxMindGeoIPService) Close() error {
	return s.db.Close()
}

var _ appContracts.GeoIPService = (*MaxMindGeoIPService)(nil)

func (s *MaxMindGeoIPService) Lookup(
	_ context.Context,
	ipStr string,
) (string, string, error) {

	ip := net.ParseIP(ipStr)
	if ip == nil {
		return "", "", fmt.Errorf("invalid IP: %s", ipStr)
	}

	record, err := s.db.City(ip)
	if err != nil {
		return "", "", fmt.Errorf("geoip lookup: %w", err)
	}

	countryCode := record.Country.IsoCode 
	cityName := ""
	if name, ok := record.City.Names["en"]; ok {
		cityName = name
	}

	return countryCode, cityName, nil
}