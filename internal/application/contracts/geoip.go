package contracts

import "context"

type GeoIPService interface {
	Lookup(ctx context.Context, ipAddress string) (countryCode string, city string, err error)
}