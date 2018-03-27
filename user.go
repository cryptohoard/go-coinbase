package coinbase

import (
	"net/http"
	"time"
)

const userURL = "/user"

type User struct {
	AvatarURL   string `json:"avatar_url"`
	BitcoinUnit string `json:"bitcoin_unit"`
	Country     struct {
		Code string `json:"code"`
		Name string `json:"name"`
	} `json:"country"`
	CreatedAt       time.Time `json:"created_at"`
	Email           string    `json:"email"`
	ID              string    `json:"id"`
	Name            string    `json:"name"`
	NativeCurrency  string    `json:"native_currency"`
	ProfileBio      string    `json:"profile_bio"`
	ProfileLocation string    `json:"profile_location"`
	ProfileURL      string    `json:"profile_url"`
	Resource        string    `json:"resource"`
	ResourcePath    string    `json:"resource_path"`
	State           string    `json:"state"`
	Tiers           struct {
		Body                 string `json:"body"`
		CompletedDescription string `json:"completed_description"`
		Header               string `json:"header"`
		UpgradeButtonText    string `json:"upgrade_button_text"`
	} `json:"tiers"`
	TimeZone string `json:"time_zone"`
	Username string `json:"username"`
}

type userResp struct {
	Data User `json:"data"`
}

func (c *client) User() (*User, error) {
	r := userResp{}
	_, err := c.Request(http.MethodGet, userURL, nil, &r)
	return &r.Data, err
}
