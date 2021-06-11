package go_utils

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/asaskevich/govalidator"
	"moul.io/http2curl"
)

// OAUTHResponse holds OAuth2 tokens and scope, to be referred to when composing Authentication headers
// and when checking permissions
type OAUTHResponse struct {
	Scope        string `json:"scope"`
	ExpiresIn    int    `json:"expires_in"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
}

// Client describes the interface that an API  client should offer
// It is used extensively in tests (to mock responses)
type Client interface {
	IsInitialized() bool
	Refresh() error
	Authenticate() error
	MakeRequest(method string, url string, body io.Reader) (*http.Response, error)
	APIScheme() string
	APIHost() string
	HTTPClient() *http.Client
	AccessToken() string
	TokenType() string
	RefreshToken() string
	AccessScope() string
	ExpiresIn() int
	RefreshAt() time.Time
	MeURL() (string, error)
	ClientID() string
	ClientSecret() string
	APITokenURL() string
	GrantType() string
	Username() string
	Password() string
	UpdateAuth(authResp *OAUTHResponse)
	SetInitialized(b bool)
}

// DefaultServerClient initializes a server client using default environment variables
func DefaultServerClient() (*ServerClient, error) {
	clientID, err := GetEnvVar(ClientIDEnvVarName)
	if err != nil {
		return nil, err
	}

	clientSecret, err := GetEnvVar(ClientSecretEnvVarName)
	if err != nil {
		return nil, err
	}

	username, err := GetEnvVar(UsernameEnvVarName)
	if err != nil {
		return nil, err
	}

	password, err := GetEnvVar(PasswordEnvVarName)
	if err != nil {
		return nil, err
	}

	grantType, err := GetEnvVar(GrantTypeEnvVarName)
	if err != nil {
		return nil, err
	}

	apiScheme, err := GetEnvVar(APISchemeEnvVarName)
	if err != nil {
		return nil, err
	}

	apiTokenURL, err := GetEnvVar(TokenURLEnvVarName)
	if err != nil {
		return nil, err
	}

	apiHost, err := GetEnvVar(APIHostEnvVarName)
	if err != nil {
		return nil, err
	}

	workstationID, err := GetEnvVar(WorkstationEnvVarName)
	if err != nil {
		// this is optional
		if IsDebug() {
			log.Printf("%s env var not found", WorkstationEnvVarName)
		}
	}

	customHeaders := map[string]string{
		WorkstationHeaderName: workstationID,
	}
	return NewServerClient(
		clientID,
		clientSecret,
		apiTokenURL,
		apiHost,
		apiScheme,
		grantType,
		username,
		password,
		customHeaders,
	)
}

// NewServerClient initializes a generic OAuth2 + HTTP server client
func NewServerClient(
	clientID string,
	clientSecret string,
	apiTokenURL string,
	apiHost string,
	apiScheme string,
	grantType string,
	username string,
	password string,
	extraHeaders map[string]string,
) (*ServerClient, error) {
	c := ServerClient{
		clientID:     clientID,
		clientSecret: clientSecret,
		apiTokenURL:  apiTokenURL,
		apiHost:      apiHost,
		apiScheme:    apiScheme,
		grantType:    grantType,
		username:     username,
		password:     password,
	}
	err := c.Initialize()
	if err != nil {
		return nil, fmt.Errorf("unable to initialize server client: %w", err)
	}
	if extraHeaders != nil {
		c.extraHeaders = extraHeaders
	}
	return &c, nil
}

// ServerClient is a general purpose client for interacting with servers that:
//
//  1. Offer a HTTP API (it need not be RESTful)
//  2. Support OAuth2 authentication with the password grant type
//
// ServerClient MUST be configured by calling the `Initialize` method.
type ServerClient struct {
	// key connec
	clientID     string
	clientSecret string
	apiTokenURL  string
	apiHost      string
	apiScheme    string
	grantType    string
	username     string
	password     string
	extraHeaders map[string]string // optional extra headers

	// these fields are set by the constructor upon successful initialization
	httpClient   *http.Client
	accessToken  string
	tokenType    string
	refreshToken string
	accessScope  string
	expiresIn    int
	refreshAt    time.Time

	// sentinel value to simplify later checks
	isInitialized bool
}

// MeURL calculates and returns the user profile URL
func (c *ServerClient) MeURL() (string, error) {
	parsedTokenURL, parseErr := url.Parse(c.apiTokenURL)
	if parseErr != nil {
		return "", parseErr
	}
	meURL := fmt.Sprintf("%s://%s/%s", parsedTokenURL.Scheme, parsedTokenURL.Host, meURLFragment)
	return meURL, nil
}

// Refresh uses the refresh token to obtain a fresh access token
func (c *ServerClient) Refresh() error {
	if !c.IsInitialized() {
		return fmt.Errorf("cannot Refresh API tokens on an uninitialized client")
	}
	refreshData := url.Values{}
	refreshData.Set("client_id", c.clientID)
	refreshData.Set("client_secret", c.clientSecret)
	refreshData.Set("grant_type", "refresh_token")
	refreshData.Set("refresh_token", c.refreshToken)
	encodedRefreshData := strings.NewReader(refreshData.Encode())
	resp, err := c.httpClient.Post(c.apiTokenURL, "application/x-www-form-urlencoded", encodedRefreshData)
	if err != nil {
		return err
	}
	if resp != nil && (resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices) {
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			msg := fmt.Sprintf("server error status: %d", resp.StatusCode)
			return fmt.Errorf(msg)
		}
		msg := fmt.Sprintf(
			"server error status: %d\nraw response: %s",
			resp.StatusCode,
			string(data),
		)
		return fmt.Errorf(msg)
	}
	authResp, decodeErr := DecodeOAUTHResponseFromJSON(resp)
	if decodeErr != nil {
		return decodeErr
	}
	c.UpdateAuth(authResp)
	return nil
}

// UpdateAuth updates the tokens stored on the  API client after successful authentication or refreshes
func (c *ServerClient) UpdateAuth(authResp *OAUTHResponse) {
	c.accessToken = authResp.AccessToken
	c.tokenType = authResp.TokenType
	c.accessScope = authResp.Scope
	c.refreshToken = authResp.RefreshToken
	c.expiresIn = authResp.ExpiresIn

	// wait out most of the token's duration to expiry before attempting to Refresh
	secondsToRefresh := int(float64(c.expiresIn) * TokenExpiryRatio)
	c.refreshAt = time.Now().Add(time.Second * time.Duration(secondsToRefresh))
	c.isInitialized = true
}

// Authenticate uses client credentials stored on the client to log in to OAuth2 authentication server
// and update stored credentials
func (c *ServerClient) Authenticate() error {
	if err := CheckAPIClientPreconditions(c); err != nil {
		return fmt.Errorf("failing API client preconditions: %w", err)
	}
	credsData := url.Values{}
	credsData.Set("client_id", c.clientID)
	credsData.Set("client_secret", c.clientSecret)
	credsData.Set("grant_type", c.grantType)
	credsData.Set("username", c.username)
	credsData.Set("password", c.password)
	encodedCredsData := strings.NewReader(credsData.Encode())

	authResp, authErr := c.httpClient.Post(c.apiTokenURL, "application/x-www-form-urlencoded", encodedCredsData)
	if authErr != nil {
		return authErr
	}
	if authResp != nil && (authResp.StatusCode < http.StatusOK || authResp.StatusCode >= http.StatusMultipleChoices) {
		data, err := ioutil.ReadAll(authResp.Body)
		if err != nil {
			msg := fmt.Sprintf("server error status: %d", authResp.StatusCode)
			return fmt.Errorf(msg)
		}
		msg := fmt.Sprintf(
			"server error status: %d\nraw response: %s",
			authResp.StatusCode,
			string(data),
		)
		return fmt.Errorf(msg)
	}
	decodedAuthResp, decodeErr := DecodeOAUTHResponseFromJSON(authResp)
	if decodeErr != nil {
		return decodeErr
	}
	c.UpdateAuth(decodedAuthResp)
	return nil // no error
}

// MakeRequest composes an authenticated  request that has the correct content type
func (c *ServerClient) MakeRequest(method string, url string, body io.Reader) (*http.Response, error) {
	if time.Now().UnixNano() > c.refreshAt.UnixNano() {
		refreshErr := c.Refresh()
		if refreshErr != nil {
			return nil, refreshErr
		}
	}
	req, reqErr := http.NewRequest(method, url, body)
	if reqErr != nil {
		return nil, reqErr
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.accessToken)

	// set extra headers e.g X-Workstation header
	if c.extraHeaders != nil {
		for k, v := range c.extraHeaders {
			req.Header.Set(k, v)
		}
	}

	if IsDebug() {
		command, _ := http2curl.GetCurlCommand(req)
		log.Printf("\nCurl command:\n%s\n", command)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	contentType := resp.Header.Get("Content-Type")
	if contentType != "application/json" {
		bs, err := ioutil.ReadAll(resp.Body)
		if IsDebug() {
			log.Printf("Mismatched content type error: %s\n", err)
			log.Printf("Mismatched content type body: %s\b", string(bs))
		}
		return nil, fmt.Errorf("expected application/json Content-Type, got " + contentType)
	}
	return resp, nil
}

// Initialize MUST be used to set up a working  Client
func (c *ServerClient) Initialize() error {
	err := CheckAPIClientPreconditions(c)
	if err != nil {
		return fmt.Errorf("server client precondition check error: %w", err)
	}

	// the timeout is half an hour, to match the timeout of a Cloud Run function
	// and to support somewhat long lived data "crawls"
	c.httpClient = &http.Client{Timeout: time.Second * 60 * 30}

	err = c.Authenticate()
	if err != nil {
		return fmt.Errorf("server client authentication error: %w", err)
	}

	err = CheckAPIClientPostConditions(c)
	if err != nil {
		return fmt.Errorf("server client postcondition check error: %w", err)
	}

	c.SetInitialized(true)
	return nil
}

// HTTPClient returns a properly configured HTTP client
func (c *ServerClient) HTTPClient() *http.Client {
	return c.httpClient
}

// AccessToken returns the latest access token
func (c *ServerClient) AccessToken() string {
	return c.accessToken
}

// TokenType returns the latest OAuth access token's token type
func (c *ServerClient) TokenType() string {
	return c.tokenType
}

// RefreshToken returns the latest refresh token
func (c *ServerClient) RefreshToken() string {
	return c.refreshToken
}

// AccessScope returns the latest access scope
func (c *ServerClient) AccessScope() string {
	return c.accessScope
}

// ExpiresIn returns the expiry seconds value returned after the last authentication
func (c *ServerClient) ExpiresIn() int {
	return c.expiresIn
}

// RefreshAt returns the target refresh time
func (c *ServerClient) RefreshAt() time.Time {
	return c.refreshAt
}

// APIScheme exports the configured  API scheme
func (c *ServerClient) APIScheme() string {
	return c.apiScheme
}

// APIHost returns the configured  API host
func (c *ServerClient) APIHost() string {
	return c.apiHost
}

// ClientID returns the configured client ID
func (c *ServerClient) ClientID() string {
	return c.clientID
}

// ClientSecret returns the configured client secret
func (c *ServerClient) ClientSecret() string {
	return c.clientSecret
}

// APITokenURL returns the configured API token URL on the client
func (c *ServerClient) APITokenURL() string {
	return c.apiTokenURL
}

// GrantType returns the configured grant type on the client
func (c *ServerClient) GrantType() string {
	return c.grantType
}

// Username returns the configured Username on the client
func (c *ServerClient) Username() string {
	return c.username
}

// Password returns the configured Password on the client
func (c *ServerClient) Password() string {
	return c.password
}

// IsInitialized returns true if the  httpClient is correctly initialized
func (c *ServerClient) IsInitialized() bool {
	return c.isInitialized
}

// SetInitialized sets the value of the isInitialized bool
func (c *ServerClient) SetInitialized(isInitialized bool) {
	c.isInitialized = isInitialized
}

// CheckAPIClientPreconditions ensures that all the parameters passed into `Initialize` make sense
func CheckAPIClientPreconditions(client Client) error {
	clientID := client.ClientID()
	if !govalidator.IsAlphanumeric(clientID) || len(clientID) < tokenMinLength {
		errMsg := fmt.Sprintf("%s is not a valid clientId, expected a non-blank alphanumeric string of at least %d characters", clientID, tokenMinLength)
		return fmt.Errorf(errMsg)
	}

	clientSecret := client.ClientSecret()
	if !govalidator.IsAlphanumeric(clientSecret) || len(clientSecret) < tokenMinLength {
		errMsg := fmt.Sprintf("%s is not a valid clientSecret, expected a non-blank alphanumeric string of at least %d characters", clientSecret, tokenMinLength)
		return fmt.Errorf(errMsg)
	}

	apiTokenURL := client.APITokenURL()
	if !govalidator.IsRequestURL(apiTokenURL) {
		errMsg := fmt.Sprintf("%s is not a valid apiTokenURL, expected an http(s) URL", apiTokenURL)
		return fmt.Errorf(errMsg)
	}

	apiHost := client.APIHost()
	if !govalidator.IsHost(apiHost) {
		errMsg := fmt.Sprintf("%s is not a valid apiHost, expected a valid IP or domain name", apiHost)
		return fmt.Errorf(errMsg)
	}

	apiScheme := client.APIScheme()
	if apiScheme != "http" && apiScheme != "https" {
		errMsg := fmt.Sprintf("%s is not a valid apiScheme, expected http or https", apiScheme)
		return fmt.Errorf(errMsg)
	}

	grantType := client.GrantType()
	if grantType != "password" {
		return fmt.Errorf("the only supported OAuth grant type for now is 'password'")
	}

	username := client.Username()
	if !govalidator.IsEmail(username) {
		return fmt.Errorf("the username `%s` is not a valid email", username)
	}

	password := client.Password()
	if len(password) < apiPasswordMinLength {
		msg := fmt.Sprintf("the Password should be a string of at least %d characters", apiPasswordMinLength)
		return fmt.Errorf(msg)
	}

	return nil
}

// CheckAPIClientPostConditions performs sanity checks on a freshly initialized  client
func CheckAPIClientPostConditions(client Client) error {
	accessToken := client.AccessToken()
	if !govalidator.IsAlphanumeric(accessToken) || len(accessToken) < tokenMinLength {
		return fmt.Errorf("invalid access token after APIClient initialization")
	}

	tokenType := client.TokenType()
	if tokenType != "Bearer" {
		return fmt.Errorf("invalid token type after APIClient initialization, expected 'Bearer'")
	}

	refreshToken := client.RefreshToken()
	if !govalidator.IsAlphanumeric(refreshToken) || len(refreshToken) < tokenMinLength {
		return fmt.Errorf("invalid Refresh token after APIClient initialization")
	}

	accessScope := client.AccessScope()
	if !govalidator.IsASCII(accessScope) || len(accessScope) < tokenMinLength {
		return fmt.Errorf("invalid access scope text after APIClient initialization")
	}

	expiresIn := client.ExpiresIn()
	if expiresIn < 1 {
		return fmt.Errorf("invalid expiresIn after APIClient initialization")
	}

	refreshAt := client.RefreshAt()
	if refreshAt.UnixNano() < time.Now().UnixNano() {
		return fmt.Errorf("invalid past refreshAt after APIClient initialization")
	}

	return nil // no errors found
}

// CheckAPIInitialization returns and error if the  httpClient was not correctly initialized by calling `.Initialize()`
func CheckAPIInitialization(client Client) error {
	if client == nil || !client.IsInitialized() {
		return fmt.Errorf("the API httpClient is not correctly initialized. Please use the `.Initialize` constructor")
	}
	return nil
}

// ComposeAPIURL assembles an  URL string for the supplied path and query string
func ComposeAPIURL(client Client, path string, query string) string {
	apiURL := url.URL{
		Scheme:   client.APIScheme(),
		Host:     client.APIHost(),
		Path:     path,
		RawQuery: query,
	}
	return apiURL.String()
}

// DecodeOAUTHResponseFromJSON extracts auth server OAUth response from the supplied HTTP response
func DecodeOAUTHResponseFromJSON(resp *http.Response) (*OAUTHResponse, error) {
	defer CloseRespBody(resp)
	var decodedAuthResp OAUTHResponse
	respBytes, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		return nil, readErr
	}
	decodeErr := json.Unmarshal(respBytes, &decodedAuthResp)
	if decodeErr != nil {
		return nil, decodeErr
	}
	return &decodedAuthResp, nil
}

//ClientServerOptions - required to compose a server client config payload
type ClientServerOptions struct {
	ClientID     string
	ClientSecret string
	APITokenURL  string
	APIHost      string
	APIScheme    string
	GrantType    string
	Username     string
	Password     string
	ExtraHeaders map[string]string
}

// GetAccessToken - given a generic OAuth2 + HTTP server client retrieve its access token
func GetAccessToken(clientConfig *ClientServerOptions) (string, error) {
	newServerClient, err := NewServerClient(
		clientConfig.ClientID,
		clientConfig.ClientSecret,
		clientConfig.APITokenURL,
		clientConfig.APIHost,
		clientConfig.APIScheme,
		clientConfig.GrantType,
		clientConfig.Username,
		clientConfig.Password,
		clientConfig.ExtraHeaders,
	)
	if err != nil {
		return "", fmt.Errorf("can't get server client for access token: %w", err)
	}
	return newServerClient.AccessToken(), nil
}

// NewPostRequest - http post request
func NewPostRequest(url string, values url.Values, headers map[string]string, timeoutDuration int) (*http.Response, error) {
	reader := strings.NewReader(values.Encode())

	req, err := http.NewRequest(http.MethodPost, url, reader)
	if err != nil {
		return nil, err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}
	req.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: time.Duration(timeoutDuration) * time.Second}
	return client.Do(req)
}

// RespondWithError writes an error response
func RespondWithError(w http.ResponseWriter, code int, err error) {
	errMap := ErrorMap(err)
	errBytes, err := json.Marshal(errMap)
	if err != nil {
		errBytes = []byte(fmt.Sprintf("error: %s", err))
	}
	RespondWithJSON(w, code, errBytes)
}

// RespondWithJSON writes a JSON response
func RespondWithJSON(w http.ResponseWriter, code int, payload []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err := w.Write(payload)
	if err != nil {
		log.Printf(
			"unable to write payload `%s` to the http.ResponseWriter: %s",
			string(payload),
			err,
		)
	}
}
