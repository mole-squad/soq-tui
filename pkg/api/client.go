package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	soqapi "github.com/mole-squad/soq-api/api"
	"github.com/mole-squad/soq-tui/pkg/config"
	"github.com/mole-squad/soq-tui/pkg/logger"
)

type Client struct {
	apiHost    string
	httpClient *http.Client
	logger     *logger.Logger

	configDir string

	token string
}

func NewClient(logger *logger.Logger, configDir string) *Client {
	apiHost := config.APIHost

	c := &http.Client{
		Timeout: 10 * time.Second,
	}

	return &Client{
		apiHost:    apiHost,
		logger:     logger,
		configDir:  configDir,
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
			c.logger.Debug("token file does not exist")
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

	dto := soqapi.LoginRequestDTO{
		Username: username,
		Password: password,
	}

	err := c.doRequest(ctx, http.MethodPost, "/auth/token", dto, &tokenResp)
	if err != nil {
		return "", fmt.Errorf("error logging in: %w", err)
	}

	return tokenResp.Token, nil
}

func (c *Client) ListTasks(ctx context.Context) ([]soqapi.TaskDTO, error) {
	var tasks []soqapi.TaskDTO

	err := c.doRequest(ctx, http.MethodGet, "/tasks", nil, &tasks)
	if err != nil {
		return nil, fmt.Errorf("error listing tasks: %w", err)
	}

	return tasks, nil
}

func (c *Client) CreateTask(ctx context.Context, t *soqapi.CreateTaskRequestDTO) (soqapi.TaskDTO, error) {
	var task soqapi.TaskDTO

	err := c.doRequest(ctx, http.MethodPost, "/tasks", t, &task)
	if err != nil {
		return task, fmt.Errorf("error creating task: %w", err)
	}

	return task, nil
}

func (c *Client) UpdateTask(ctx context.Context, taskID uint, t *soqapi.UpdateTaskRequestDTO) (soqapi.TaskDTO, error) {
	var task soqapi.TaskDTO

	err := c.doRequest(ctx, http.MethodPatch, fmt.Sprintf("/tasks/%d", taskID), t, &task)
	if err != nil {
		return task, fmt.Errorf("error updating task: %w", err)
	}

	return task, nil
}

func (c *Client) ResolveTask(ctx context.Context, taskID uint) (soqapi.TaskDTO, error) {
	var task soqapi.TaskDTO

	err := c.doRequest(ctx, http.MethodPatch, fmt.Sprintf("/tasks/%d/resolve", taskID), nil, &task)
	if err != nil {
		return task, fmt.Errorf("error resolving task: %w", err)
	}

	return task, nil
}

func (c *Client) DeleteTask(ctx context.Context, taskID uint) error {
	err := c.doRequest(ctx, http.MethodDelete, fmt.Sprintf("/tasks/%d", taskID), nil, nil)
	if err != nil {
		return fmt.Errorf("error deleting task: %w", err)
	}

	return nil
}

func (c *Client) ListFocusAreas(ctx context.Context) ([]soqapi.FocusAreaDTO, error) {
	var focusAreas []soqapi.FocusAreaDTO

	err := c.doRequest(ctx, http.MethodGet, "/focusareas", nil, &focusAreas)
	if err != nil {
		return nil, fmt.Errorf("error listing focus areas: %w", err)
	}

	return focusAreas, nil
}

func (c *Client) CreateFocusArea(ctx context.Context, f *soqapi.CreateFocusAreaRequestDTO) (soqapi.FocusAreaDTO, error) {
	var focusArea soqapi.FocusAreaDTO

	err := c.doRequest(ctx, http.MethodPost, "/focusareas", f, &focusArea)
	if err != nil {
		return focusArea, fmt.Errorf("error creating focus area: %w", err)
	}

	return focusArea, nil
}

func (c *Client) UpdateFocusArea(ctx context.Context, focusAreaID uint, f *soqapi.UpdateFocusAreaRequestDTO) (soqapi.FocusAreaDTO, error) {
	var focusArea soqapi.FocusAreaDTO

	err := c.doRequest(ctx, http.MethodPatch, fmt.Sprintf("/focusareas/%d", focusAreaID), f, &focusArea)
	if err != nil {
		return focusArea, fmt.Errorf("error updating focus area: %w", err)
	}

	return focusArea, nil
}

func (c *Client) DeleteFocusArea(ctx context.Context, focusAreaID uint) error {
	err := c.doRequest(ctx, http.MethodDelete, fmt.Sprintf("/focusareas/%d", focusAreaID), nil, nil)
	if err != nil {
		return fmt.Errorf("error deleting focus area: %w", err)
	}

	return nil
}

func (c *Client) doRequest(ctx context.Context, method, path string, dto interface{}, respBody interface{}) error {
	var req *http.Request
	var err error

	c.logger.Debug("Request", "method", method, "url", path)

	reqUrl := url.URL{
		Scheme: "http",
		Host:   c.apiHost,
		Path:   path,
	}

	if dto == nil {
		req, err = http.NewRequestWithContext(ctx, method, reqUrl.String(), nil)
	} else {
		serializedDto, err := json.Marshal(dto)
		if err != nil {
			return fmt.Errorf("error marshalling request: %w", err)
		}

		req, err = http.NewRequestWithContext(ctx, method, reqUrl.String(), bytes.NewBuffer(serializedDto))
	}

	if err != nil {
		return fmt.Errorf("error building request: %w", err)
	}

	req.Header = c.buildHeaders()
	req.Header.Set("Content-Type", "application/json")

	res, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error executing request: %w", err)
	}

	defer res.Body.Close()

	if res.StatusCode == http.StatusUnauthorized {
		c.ClearToken()
		return fmt.Errorf("unauthorized")
	}

	respBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("error reading response: %w", err)
	}

	var badStatusCode bool

	switch method {
	case http.MethodGet:
		badStatusCode = res.StatusCode != http.StatusOK

	case http.MethodPost:
		badStatusCode = res.StatusCode != http.StatusCreated

	case http.MethodPatch:
		badStatusCode = res.StatusCode != http.StatusOK

	case http.MethodDelete:
		badStatusCode = res.StatusCode != http.StatusNoContent
	}

	if badStatusCode {
		return fmt.Errorf("unexpected status code %d. Error: %s", res.StatusCode, respBytes)
	}

	if len(respBytes) == 0 {
		return nil
	}

	if err = json.Unmarshal(respBytes, respBody); err != nil {
		return fmt.Errorf("error unmarshalling response: %w", err)
	}

	return nil
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
	return filepath.Join(c.configDir, "token"), nil
}
