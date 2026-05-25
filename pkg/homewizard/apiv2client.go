package homewizard

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type APIv2Client struct {
	client   *http.Client
	url      string
	username string
	token    string
}

type externalMeasurement struct {
	Type      string  `json:"type"`
	Timestamp string  `json:"timestamp"`
	Value     float64 `json:"value"`
	Unit      string  `json:"unit"`
}

type measurement struct {
	Tariff            int                   `json:"tariff"`
	Timestamp         string                `json:"timestamp"`
	EnergyImportkWH   float64               `json:"energy_import_kwh"`
	EnergyImportT1kWH float64               `json:"energy_import_t1_kwh"`
	EnergyImportT2kWH float64               `json:"energy_import_t2_kwh"`
	EnergyExportkWH   float64               `json:"energy_export_kwh"`
	EnergyExportT1kWH float64               `json:"energy_export_t1_kwh"`
	EnergyExportT2kWH float64               `json:"energy_export_t2_kwh"`
	PowerW            int                   `json:"power_w"`
	PowerL1W          int                   `json:"power_l1_w"`
	PowerL2W          int                   `json:"power_l2_w"`
	PowerL3W          int                   `json:"power_l3_w"`
	VoltageV          float64               `json:"voltage_v"`
	VoltageL1V        float64               `json:"voltage_l1_v"`
	VoltageL2V        float64               `json:"voltage_l2_v"`
	VoltageL3V        float64               `json:"voltage_l3_v"`
	CurrentA          float64               `json:"current_a"`
	CurrentL1A        float64               `json:"current_l1_a"`
	CurrentL2A        float64               `json:"current_l2_a"`
	CurrentL3A        float64               `json:"current_l3_a"`
	External          []externalMeasurement `json:"external"`
}

type Option func(*APIv2Client)

func WithToken(token string) Option {
	return func(c *APIv2Client) { c.token = token }
}

func NewAPIv2Client(url string, opts ...Option) *APIv2Client {
	log.Debug("Homewizard device URL=", url)
	c := &APIv2Client{
		url: url,
		client: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		},
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func (c *APIv2Client) CreateLocalUser(username string) error {
	body := fmt.Sprintf(`{ "name": "local/%s" }`, username)
	jsonBody := bytes.NewBufferString(body)
	req, err := http.NewRequest(http.MethodPost, c.urlWithPath("api/user"), jsonBody)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Api-Version", "2")
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusForbidden {
		log.Warn("Retry later")
	} else if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected server http response: %s", resp.Status)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	log.Info(string(b))
	return nil
}

func (c APIv2Client) GetMeasurement() (measurement, error) {
	var m measurement
	req, err := http.NewRequest(http.MethodGet, c.urlWithPath("api/measurement"), nil)
	if err != nil {
		return m, err
	}
	req.Header.Add("Authorization", "Bearer "+c.token)
	req.Header.Add("X-Api-Version", "2")
	resp, err := c.client.Do(req)
	if err != nil {
		return m, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return m, fmt.Errorf("unexpected server http response: %s", resp.Status)
	}

	err = json.NewDecoder(resp.Body).Decode(&m)
	if err != nil {
		return m, err
	}

	return m, nil
}

type User struct {
	Name    string
	Current bool
}

func (c APIv2Client) ListUsers() ([]User, error) {
	req, err := http.NewRequest(http.MethodGet, c.urlWithPath("api/user"), nil)
	if err != nil {
		return []User{}, err
	}
	req.Header.Add("Authorization", "Bearer "+c.token)
	req.Header.Add("X-Api-Version", "2")
	resp, err := c.client.Do(req)
	if err != nil {
		return []User{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return []User{}, fmt.Errorf("unexpected server http response: %s", resp.Status)
	}

	var users []User
	err = json.NewDecoder(resp.Body).Decode(&users)
	return users, err
}

func (c *APIv2Client) urlWithPath(path string) string {
	return fmt.Sprintf("%s/%s", c.url, path)
}
