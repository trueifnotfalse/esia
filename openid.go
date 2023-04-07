package esia

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	defaultCodeUrlPath  = "aas/oauth2/ac"
	defaultTokenUrlPath = "aas/oauth2/te"
)

type OpenId struct {
	Config *OpenIdConfig
	Token  string
	Oid    int32
	signer Signer
}

func NewOpenId(config *OpenIdConfig, signer Signer) *OpenId {
	if 0 == len(config.CodeUrl) {
		config.CodeUrl = defaultCodeUrlPath
	}
	if 0 == len(config.TokenUrl) {
		config.TokenUrl = defaultTokenUrlPath
	}

	return &OpenId{
		Config: config,
		signer: signer,
	}
}

func (c *OpenId) getState() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return uuid, nil
}

func (c *OpenId) getTimeStamp() string {
	return time.Now().Format("2006.01.02 15:04:05 Z0700")
}

func (c *OpenId) GetUrl() (string, error) {
	state, err := c.getState()
	if err != nil {
		return "", err
	}

	timestamp := c.getTimeStamp()
	clientSecret := c.Config.Scope + timestamp + c.Config.MnemonicsSystem + state
	clientSecret, err = c.sign(clientSecret)
	if err != nil {
		return "", err
	}

	var Url *url.URL
	Url, err = url.Parse(c.Config.PortalUrl)
	if err != nil {
		return "", err
	}

	Url.Path += c.Config.CodeUrl

	params := &url.Values{
		"client_id":     []string{c.Config.MnemonicsSystem},
		"client_secret": []string{clientSecret},
		"redirect_uri":  []string{c.Config.RedirectUrl},
		"scope":         []string{c.Config.Scope},
		"response_type": []string{"code"},
		"state":         []string{state},
		"access_type":   []string{"offline"},
		"timestamp":     []string{timestamp},
	}

	Url.RawQuery = params.Encode()

	return Url.String(), nil
}

func (c *OpenId) GetInfoByPath(path string, item interface{}) error {
	if c.Oid <= 0 {
		return errors.New("oid empty")
	}
	if len(c.Token) <= 0 {
		return errors.New("token empty")
	}
	client := &http.Client{}

	req, err := http.NewRequest("GET", c.Config.PortalUrl+"rs/prns/"+fmt.Sprint(c.Oid)+path, nil)
	req.Header.Add("Authorization", "Bearer "+c.Token)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, &item)
	if err != nil {
		return err
	}

	return nil
}

func (c *OpenId) GetTokenState(code string) (Token, error) {
	var esiaToken Token
	state, err := c.getState()
	if err != nil {
		return esiaToken, err
	}

	timestamp := c.getTimeStamp()
	clientSecret := c.Config.Scope + timestamp + c.Config.MnemonicsSystem + state
	clientSecret, err = c.sign(clientSecret)
	if err != nil {
		return esiaToken, err
	}

	params := url.Values{
		"client_id":     []string{c.Config.MnemonicsSystem},
		"code":          []string{code},
		"grant_type":    []string{"authorization_code"},
		"client_secret": []string{clientSecret},
		"state":         []string{state},
		"redirect_uri":  []string{c.Config.RedirectUrl},
		"scope":         []string{c.Config.Scope},
		"timestamp":     []string{timestamp},
		"token_type":    []string{"Bearer"},
		"refresh_token": []string{state},
	}
	resp, err := http.PostForm(c.Config.PortalUrl+c.Config.TokenUrl, params)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return esiaToken, err
	}

	err = json.Unmarshal(body, &esiaToken)
	if err != nil {
		return esiaToken, err
	}

	chunks := strings.Split(esiaToken.AccessToken, ".")
	data, err := base64.URLEncoding.DecodeString(chunks[1])
	if err != nil {
		data, err = base64.URLEncoding.DecodeString(chunks[1] + "==")
		if err != nil {
			return esiaToken, err
		}
	}

	err = json.Unmarshal([]byte(string(data)), &esiaToken.AuthCode)
	if err != nil {
		return esiaToken, err
	}

	c.Token = esiaToken.AccessToken
	c.Oid = esiaToken.AuthCode.UrnEsiaSbjId

	return esiaToken, nil
}

func (c *OpenId) sign(message string) (string, error) {
	return c.signer.Sign(message)
}
