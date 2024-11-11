package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	c "github.com/Maxxxxxx-x/cicd-cloudflare-zero-trust-thing/config"
	"github.com/Maxxxxxx-x/cicd-cloudflare-zero-trust-thing/models"
	"io"
	"net/http"
)

const API_URL = "https://api.cloudflare.com/client/v4"

type IngressRoot struct {
	Config Config `json:"config"`
}

type Response struct {
	Errors   []string `json:"errors"`
	Messages []string `json:"messages"`
	Success  bool     `json:"success"`
	Result   Result   `json:"result"`
}

type Result struct {
	AccountID string `json:"account_id"`
	Config    Config `json:"config"`
	CreatedAt string `json:"created_at"`
	Source    string `json:"source"`
	TunnelID  string `json:"tunnel_id"`
	Version   int    `json:"version"`
}

type Config struct {
	Ingress       []IngressConfig `json:"ingress"`
	OriginRequest OriginRequest   `json:"origin_request"`
	WarpRouting   WarpRouting     `json:"warp-routing"`
}

type IngressConfig struct {
	Hostname      string        `json:"hostname"`
	Path          string        `json:"path"`
	Service       string        `json:"service"`
	OriginRequest OriginRequest `json:"origin_request"`
}

type OriginRequest struct {
	Access                 map[string]interface{} `json:"access,omitempty"`
	CaPool                 interface{}            `json:"ca_pool,omitempty"`
	ConnectTimeout         interface{}            `json:"connect_timeout,omitempty"`
	DisableChunkedEncoding bool                   `json:"disable_chunked_encoding,omitempty"`
	Http2Origin            bool                   `json:"http2_origin,omitempty"`
	HttpHostHeader         interface{}            `json:"http_host_header,omitempty"`
	KeepAliveConnections   interface{}            `json:"keep_alive_connections,omitempty"`
	KeepAliveTimeout       interface{}            `json:"keep_alive_timeout,omitempty"`
	NoHappyEyeballs        bool                   `json:"no_happy_eye_balls,omitempty"`
	NoTLSVerify            bool                   `json:"no_tls_verify,omitempty"`
	OriginServerName       interface{}            `json:"origin_server_name,omitempty"`
	ProxyType              interface{}            `json:"proxy_type,omitempty"`
	TcpKeepAlive           interface{}            `json:"tcp_keep_alive,omitempty"`
	TLSTimeout             interface{}            `json:"tls_timeout,omitempty"`
}

type AccessConfig struct {
	AudTag   []string `json:"aud_tag"`
	Required bool     `json:"required"`
	TeamName string   `json:"team_name"`
}

type WarpRouting struct {
	Enabled bool `json:"enabled"`
}

type ResultInfo struct {
	Count      int `json:"count"`
	Page       int `json:"page"`
	PerPage    int `json:"per_page"`
	TotalCount int `json:"total_count"`
}

type DomainRecord struct {
	Comment           string                 `json:"comment"`
	Name              string                 `json:"name"`
	Proxied           bool                   `json:"proxied"`
	Settings          map[string]interface{} `json:"settings"`
	Tags              []string               `json:"tags"`
	TTL               int                    `json:"ttl"`
	Content           string                 `json:"content"`
	Type              string                 `json:"type"`
	CommentModifiedOn string                 `json:"comment_modified_on"`
	CreatedOn         string                 `json:"created_on"`
	Id                string                 `json:"id"`
	Meta              map[string]interface{} `json:"meta"`
	ModifiedOn        string                 `json:"modified_on"`
	Proxiable         bool                   `json:"proxiable"`
	TagsModifiedOn    string                 `json:"tags_modified_on"`
}

type DNSResponse struct {
	Errors     []string       `json:"errors"`
	Messages   []string       `json:"messages"`
	Success    bool           `json:"success"`
	ResultInfo ResultInfo     `json:"result_info"`
	Result     []DomainRecord `json:"result"`
}

type CreateDNSResponse struct {
	Errors   []string        `json:"errors"`
	Messages []string        `json:"messages"`
	Success  bool            `json:"success"`
	Result   CreateDNSResult `json:"result"`
}

type CreateDNSResult struct {
	Comment           string                 `json:"comment"`
	Name              string                 `json:"name"`
	Proxied           bool                   `json:"proxied"`
	Settings          map[string]interface{} `json:"settings"`
	Tags              []string               `json:"tags"`
	TTL               int                    `json:"ttl"`
	Content           string                 `json:"content"`
	Type              string                 `json:"type"`
	CommentModifiedOn string                 `json:"comment_modified_on"`
	CreatedOn         string                 `json:"created_on"`
	ID                string                 `json:"id"`
	Meta              map[string]interface{} `json:"meta"`
	ModifiedOn        string                 `json:"modified_on"`
	Proxiable         bool                   `json:"proxiable"`
	TagsModifiedOn    string                 `json:"tags_modified_on"`
}

func getTunnelConfig(authKey string, email string, accountId string, tunnelId string) (Config, error) {
	var config Config
	client := &http.Client{}
	url := fmt.Sprintf("%s/accounts/%s/cfd_tunnel/%s/configurations", API_URL, accountId, tunnelId)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return config, err
	}
	req.Header.Add("X-Auth-Email", email)
	req.Header.Add("X-Auth-Key", authKey)
	req.Header.Add("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		return config, err
	}
	if res.StatusCode != 200 {
		return config, fmt.Errorf("%s", res.Status)
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	var Body Response
	err = json.Unmarshal(body, &Body)
	if err != nil {
		return config, err
	}
	config = Body.Result.Config
	return config, nil
}

func getDnsRecords(authKey string, email string, zoneId string) ([]DomainRecord, error) {
	var records []DomainRecord
	client := &http.Client{}
	url := fmt.Sprintf("%s/zones/%s/dns_records", API_URL, zoneId)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return records, err
	}
	req.Header.Add("X-Auth-Email", email)
	req.Header.Add("X-Auth-key", authKey)
	req.Header.Add("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		return records, err
	}
	if res.StatusCode != 200 {
		return records, fmt.Errorf("%s", res.Status)
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return records, err
	}
	var Body DNSResponse
	err = json.Unmarshal(body, &Body)
	if err != nil {
		return records, err
	}
	records = Body.Result
	return records, nil
}

func checkHostnameInUsed(tunnelConfig *Config, hostname string) *IngressConfig {
	for i, ingress := range tunnelConfig.Ingress {
		if hostname == ingress.Hostname {
			return &tunnelConfig.Ingress[i]
		}
	}
	var ingress IngressConfig
	return &ingress
}

func checkRecordExist(records []DomainRecord, hostname string) DomainRecord {
	for _, record := range records {
		if hostname == record.Name {
			return record
		}
	}
	var record DomainRecord
	return record
}

func checkServiceIsSame(ingress IngressConfig, ip string, port string) bool {
	return ingress.Service == fmt.Sprintf("http://%s:%s", ip, port)
}

func checkRecordIsSame(record DomainRecord, tunnelId string) bool {
	return record.Content == fmt.Sprintf("%s.cfargotunnel.com", tunnelId)
}

func addTunnelRecord(authKey string, email string, accountId string, tunnelId string, tunnelConfig Config) error {
	client := &http.Client{}
	url := fmt.Sprintf("%s/accounts/%s/cfd_tunnel/%s/configurations", API_URL, accountId, tunnelId)
	root := IngressRoot{tunnelConfig}
	config, err := json.MarshalIndent(root, "", "  ")
	if err != nil {
		return err
	}
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(config))
	if err != nil {
		return err
	}
	req.Header.Add("X-Auth-Email", email)
	req.Header.Add("X-Auth-Key", authKey)
	req.Header.Add("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		return fmt.Errorf("%s", res.Status)
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	var Body Response
	err = json.Unmarshal(body, &Body)
	return nil
}

func addDNSRecord(authKey string, email string, zoneId string, hostname string, tunnelId string) error {
	client := &http.Client{}
	url := fmt.Sprintf("%s/zones/%s/dns_records", API_URL, zoneId)
	target := fmt.Sprintf("%s.cfargotunnel.com", tunnelId)
	record := fmt.Sprintf(`{"name": "%s", "content": "%s", "proxiable":true, "proxied": true, "type": "CNAME"}`, hostname, target)
	req, err := http.NewRequest("POST", url, bytes.NewBufferString(record))
	if err != nil {
		return err
	}
	req.Header.Add("X-Auth-Email", email)
	req.Header.Add("X-Auth-Key", authKey)
	req.Header.Add("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		return fmt.Errorf("%s", res.Status)
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	var Body CreateDNSResponse
	err = json.Unmarshal(body, &Body)
	fmt.Println(Body.Success)
	return nil
}

func SetZeroTrust(config c.Config, input models.Input) error {
	tunnelConfig, err := getTunnelConfig(config.Cloudflare.AuthKey, config.Cloudflare.Email, config.Cloudflare.AccountId, config.Cloudflare.TunnelId)
	dnsRecords, err := getDnsRecords(config.Cloudflare.AuthKey, config.Cloudflare.Email, config.Cloudflare.ZoneId)
	if err != nil {
		return err
	}
	ingress := checkHostnameInUsed(&tunnelConfig, input.Hostname)
	domainRecord := checkRecordExist(dnsRecords, input.Hostname)
	if ingress.Hostname != input.Hostname {
		newIngress := IngressConfig{
			Hostname:      input.Hostname,
			Path:          "",
			Service:       fmt.Sprintf("http://%s:%s", input.DockerIP, input.DockerPort),
			OriginRequest: OriginRequest{},
		}
		tmp := tunnelConfig.Ingress[len(tunnelConfig.Ingress)-1]
		tunnelConfig.Ingress = append(tunnelConfig.Ingress[:len(tunnelConfig.Ingress)-1], newIngress, tmp)
		err = addTunnelRecord(config.Cloudflare.AuthKey, config.Cloudflare.Email, config.Cloudflare.AccountId, config.Cloudflare.TunnelId, tunnelConfig)
		if err != nil {
			return err
		}
	}
	if !checkServiceIsSame(*ingress, input.DockerIP, input.DockerPort) {
		ingress.Service = fmt.Sprintf("http://%s:%s", input.DockerIP, input.DockerPort)
		err = addTunnelRecord(config.Cloudflare.AuthKey, config.Cloudflare.Email, config.Cloudflare.AccountId, config.Cloudflare.TunnelId, tunnelConfig)
		if err != nil {
			return err
		}
	}
	if domainRecord.Name != input.Hostname {
		err = addDNSRecord(config.Cloudflare.AuthKey, config.Cloudflare.Email, config.Cloudflare.ZoneId, input.Hostname, config.Cloudflare.TunnelId)
		if err != nil {
			return err
		}
	}
	if checkRecordIsSame(domainRecord, config.Cloudflare.TunnelId) {
		return fmt.Errorf("This domain is used by other project")
	}
	return nil
}
