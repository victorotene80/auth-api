package messaging

type Context struct {
	Aggregate string
	Action    string

	IPAddress string
	UserAgent string
	DeviceID  string
}

func (c Context) ToMetadata() map[string]string {
	meta := map[string]string{
		"aggregate": c.Aggregate,
		"action":    c.Action,
	}

	if c.IPAddress != "" {
		meta["ip_address"] = c.IPAddress
	}
	if c.UserAgent != "" {
		meta["user_agent"] = c.UserAgent
	}
	if c.DeviceID != "" {
		meta["device_id"] = c.DeviceID
	}

	return meta
}
