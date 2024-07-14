package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	soqapi "github.com/mole-squad/soq-api/api"
	"github.com/mole-squad/soq-tui/pkg/config"
)

type Client struct {
	apiHost    string
	httpClient *http.Client

	token string
}

func NewClient() *Client {
	apiHost := config.APIHost

	c := &http.Client{
		Timeout: 10 * time.Second,
	}

	return &Client{
		apiHost:    apiHost,
		httpClient: c,
	}
}

func (c *Client) IsAuthenticated() bool {
	return c.token != ""
}

func (c *Client) LoadToken() error {
	tokenFilePath, err := c.getTokenFilePath()

	data, err := os.ReadFile(tokenFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			slog.Debug("token file does not exist")
			return nil
		}

		return fmt.Errorf("error reading token file: %w", err)
	}

	rawToken := string(data)
	cleanToken := strings.TrimRight(rawToken, "\r\n")
	cleanToken = strings.TrimRight(cleanToken, "\n")

	c.token = cleanToken

	return nil
}

func (c *Client) SetToken(token string) error {
	c.token = token

	tokenFilePath, err := c.getTokenFilePath()

	err = os.Mkdir(filepath.Dir(tokenFilePath), 0777)
	if err != nil {
		if !os.IsExist(err) {
			return fmt.Errorf("error creating token directory: %w", err)
		}
	}

	err = os.WriteFile(tokenFilePath, []byte(token), 0777)
	if err != nil {
		return fmt.Errorf("error writing token file: %w", err)
	}

	return nil
}

func (c *Client) ClearToken() error {
	c.token = ""

	tokenFilePath, err := c.getTokenFilePath()

	err = os.Remove(tokenFilePath)
	if err != nil {
		return fmt.Errorf("error removing token file: %w", err)
	}

	return nil
}

func (c *Client) Login(ctx context.Context, username, password string) (string, error) {
	var tokenResp soqapi.TokenResponseDTO

	reqUrl := url.URL{
		Scheme: "http",
		Host:   c.apiHost,
		Path:   "/auth/token",
	}

	req := soqapi.LoginRequestDTO{
		Username: username,
		Password: password,
	}

	body, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("error marshalling login request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, reqUrl.String(), bytes.NewBuffer(body))
	if err != nil {
		return "", fmt.Errorf("error building login request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	res, err := c.httpClient.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("error executing login request: %w", err)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	respBody, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("error reading login response: %w", err)
	}

	if err = json.Unmarshal(respBody, &tokenResp); err != nil {
		return "", fmt.Errorf("error unmarshalling login response: %w", err)
	}

	return tokenResp.Token, nil
}

func (c *Client) ListTasks(ctx context.Context) ([]soqapi.TaskDTO, error) {
	reqUrl := url.URL{
		Scheme: "http",
		Host:   c.apiHost,
		Path:   "/tasks",
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqUrl.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("error building list tasks request: %w", err)
	}

	req.Header = c.buildHeaders()

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error executing list tasks request: %w", err)
	}

	defer res.Body.Close()

	if res.StatusCode == http.StatusUnauthorized {
		c.ClearToken()
		return nil, fmt.Errorf("unauthorized")
	}

	respBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading list tasks response: %w", err)
	}

	var tasksResp []soqapi.TaskDTO
	if err = json.Unmarshal(respBody, &tasksResp); err != nil {
		return nil, fmt.Errorf("error unmarshalling list tasks response: %w", err)
	}

	return tasksResp, nil
}

func (c *Client) CreateTask(ctx context.Context, t *soqapi.CreateTaskRequestDto) (soqapi.TaskDTO, error) {
	var task soqapi.TaskDTO

	reqUrl := url.URL{
		Scheme: "http",
		Host:   c.apiHost,
		Path:   "/tasks",
	}

	body, err := json.Marshal(t)
	if err != nil {
		return task, fmt.Errorf("error marshalling create task request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, reqUrl.String(), bytes.NewBuffer(body))
	if err != nil {
		return task, fmt.Errorf("error building create task request: %w", err)
	}

	req.Header = c.buildHeaders()

	req.Header.Set("Content-Type", "application/json")

	res, err := c.httpClient.Do(req)
	if err != nil {
		return task, fmt.Errorf("error executing create task request: %w", err)
	}

	defer res.Body.Close()

	if res.StatusCode == http.StatusUnauthorized {
		c.ClearToken()
		return task, fmt.Errorf("unauthorized")
	}

	respBody, err := io.ReadAll(res.Body)
	if err != nil {
		return task, fmt.Errorf("error reading create task response: %w", err)
	}

	if err = json.Unmarshal(respBody, &task); err != nil {
		return task, fmt.Errorf("error unmarshalling create task response: %w", err)
	}

	return task, nil
}

func (c *Client) UpdateTask(ctx context.Context, taskID uint, t *soqapi.UpdateTaskRequestDto) (soqapi.TaskDTO, error) {
	var task soqapi.TaskDTO

	reqUrl := url.URL{
		Scheme: "http",
		Host:   c.apiHost,
		Path:   fmt.Sprintf("/tasks/%d", taskID),
	}

	body, err := json.Marshal(t)
	if err != nil {
		return task, fmt.Errorf("error marshalling update task request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, reqUrl.String(), bytes.NewBuffer(body))
	if err != nil {
		return task, fmt.Errorf("error building update task request: %w", err)
	}

	req.Header = c.buildHeaders()

	res, err := c.httpClient.Do(req)
	if err != nil {
		return task, fmt.Errorf("error executing update task request: %w", err)
	}

	defer res.Body.Close()

	if res.StatusCode == http.StatusUnauthorized {
		c.ClearToken()
		return task, fmt.Errorf("unauthorized")
	}

	respBody, err := io.ReadAll(res.Body)
	if err != nil {
		return task, fmt.Errorf("error reading update task response: %w", err)
	}

	if err = json.Unmarshal(respBody, &task); err != nil {
		return task, fmt.Errorf("error unmarshalling update task response: %w", err)
	}

	return task, nil
}

func (c *Client) DeleteTask(ctx context.Context, taskID uint) error {
	reqUrl := url.URL{
		Scheme: "http",
		Host:   c.apiHost,
		Path:   fmt.Sprintf("/tasks/%d", taskID),
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, reqUrl.String(), nil)
	if err != nil {
		return fmt.Errorf("error building delete task request: %w", err)
	}

	req.Header = c.buildHeaders()

	res, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error executing delete task request: %w", err)
	}

	defer res.Body.Close()

	if res.StatusCode == http.StatusUnauthorized {
		c.ClearToken()
		return fmt.Errorf("unauthorized")
	}

	if res.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	return nil
}

func (c *Client) ListFocusAreas(ctx context.Context) ([]soqapi.FocusAreaDTO, error) {
	reqUrl := url.URL{
		Scheme: "http",
		Host:   c.apiHost,
		Path:   "/focusareas",
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqUrl.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("error building list focus areas request: %w", err)
	}

	req.Header = c.buildHeaders()

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error executing list focus areas request: %w", err)
	}

	defer res.Body.Close()

	if res.StatusCode == http.StatusUnauthorized {
		c.ClearToken()
		return nil, fmt.Errorf("unauthorized")
	}

	respBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading list focus areas response: %w", err)
	}

	var focusAreasResp []soqapi.FocusAreaDTO
	if err = json.Unmarshal(respBody, &focusAreasResp); err != nil {
		return nil, fmt.Errorf("error unmarshalling list focus areas response: %w", err)
	}

	return focusAreasResp, nil
}

func (c *Client) buildHeaders() http.Header {
	headers := http.Header{}

	if c.token != "" {
		headers.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))
	}

	headers.Set("Content-Type", "application/json")

	return headers
}

func (c *Client) getTokenFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("error getting user home directory: %w", err)
	}

	return filepath.Join(homeDir, ".soq", "token"), nil
}
