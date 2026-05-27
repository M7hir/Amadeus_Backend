package reccobeats

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

type ResponseError struct {
	StatusCode int
	Message    string
}

func (e *ResponseError) Error() string {
	return e.Message
}

type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

func validationErrorf(format string, args ...interface{}) error {
	return &ValidationError{Message: fmt.Sprintf(format, args...)}
}

type queryParameterValueKind string

const (
	queryParameterValueKindInt   queryParameterValueKind = "int"
	queryParameterValueKindFloat queryParameterValueKind = "float"
	queryParameterValueKindList  queryParameterValueKind = "list"
)

type queryParameterQuerySpec struct {
	kind     queryParameterValueKind
	min      float64
	max      float64
	required bool
	minItems int
	maxItems int
}

func NewClient(baseURL string) *Client {
	return &Client{
		BaseURL:    baseURL,
		HTTPClient: &http.Client{},
	}
}

func (c *Client) TrackRecommendation(q url.Values) ([]byte, error) {
	if err := validateRecommendationQuery(q); err != nil {
		return nil, err
	}

	return c.doGet("/v1/track/recommendation", q.Encode())
}

func (c *Client) TrackMultiple(q url.Values) ([]byte, error) {
	if err := validateMultipleTrack(q); err != nil {
		return nil, err
	}

	return c.doGet("/v1/track", q.Encode())
}

func (c *Client) TrackDetail(id string) ([]byte, error) {
	return c.doGet(fmt.Sprintf("/v1/track/%s", id), "")
}

func (c *Client) TrackAlbum(id string) ([]byte, error) {
	return c.doGet(fmt.Sprintf("/v1/track/%s/album", id), "")
}

func (c *Client) TrackAudioFeatures(id string) ([]byte, error) {
	return c.doGet(fmt.Sprintf("/v1/track/%s/audio-features", id), "")
}

func (c *Client) doGet(path string, rawQuery string) ([]byte, error) {
	fullURL, err := url.JoinPath(c.BaseURL, path)
	if err != nil {
		return nil, err
	}

	if rawQuery != "" {
		fullURL = fullURL + "?" + rawQuery
	}

	httpClient := c.HTTPClient
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	req, err := http.NewRequest(http.MethodGet, fullURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/json")

	res, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, &ResponseError{StatusCode: res.StatusCode, Message: http.StatusText(res.StatusCode)}
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func validateMultipleTrack(q url.Values) error {
	specs := map[string]queryParameterQuerySpec{
		"ids": {kind: queryParameterValueKindList, minItems: 1, maxItems: 40, required: true},
	}

	for key, spec := range specs {
		if spec.required && q.Get(key) == "" {
			return validationErrorf("missing required query parameter: %s", key)
		}
	}

	for key := range q {
		spec, ok := specs[key]
		if !ok {
			return validationErrorf("unsupported query parameter: %s", key)
		}

		if err := validateListItems(q, spec.minItems, spec.maxItems, key); err != nil {
			return err
		}
	}

	return nil
}

func validateRecommendationQuery(q url.Values) error {
	specs := map[string]queryParameterQuerySpec{
		"size":             {kind: queryParameterValueKindInt, min: 1, max: 100, required: true},
		"seeds":            {kind: queryParameterValueKindList, minItems: 1, maxItems: 5, required: true},
		"negativeSeeds":    {kind: queryParameterValueKindList, minItems: 0, maxItems: 5},
		"acousticness":     {kind: queryParameterValueKindFloat, min: 0, max: 1},
		"danceability":     {kind: queryParameterValueKindFloat, min: 0, max: 1},
		"energy":           {kind: queryParameterValueKindFloat, min: 0, max: 1},
		"instrumentalness": {kind: queryParameterValueKindFloat, min: 0, max: 1},
		"key":              {kind: queryParameterValueKindInt, min: -1, max: 11},
		"liveness":         {kind: queryParameterValueKindFloat, min: 0, max: 1},
		"loudness":         {kind: queryParameterValueKindFloat, min: -60, max: 2},
		"mode":             {kind: queryParameterValueKindInt, min: 0, max: 1},
		"speechiness":      {kind: queryParameterValueKindFloat, min: 0, max: 1},
		"tempo":            {kind: queryParameterValueKindFloat, min: 0, max: 250},
		"valence":          {kind: queryParameterValueKindFloat, min: 0, max: 1},
		"popularity":       {kind: queryParameterValueKindInt, min: 0, max: 100},
		"featureWeight":    {kind: queryParameterValueKindFloat, min: 1, max: 5},
	}

	for key := range q {
		spec, ok := specs[key]
		if !ok {
			return validationErrorf("unsupported query parameter: %s", key)
		}

		if err := validateRecommendationValue(key, q, spec); err != nil {
			return err
		}
	}

	return nil
}

func validateRecommendationValue(key string, q url.Values, spec queryParameterQuerySpec) error {
	raw := q.Get(key)
	if raw == "" {
		if spec.required {
			return validationErrorf("missing required query parameter: %s", key)
		}
		return nil
	}

	switch spec.kind {
	case queryParameterValueKindInt:
		return validateIntRange(raw, int(spec.min), int(spec.max), key)
	case queryParameterValueKindFloat:
		return validateFloatRange(raw, spec.min, spec.max, key)
	case queryParameterValueKindList:
		return validateListItems(q, spec.minItems, spec.maxItems, key)
	default:
		return validationErrorf("unsupported query parameter type for %s", key)
	}
}

func parseCSVQuery(q url.Values, key string) []string {
	values := q[key]
	if len(values) == 0 {
		return nil
	}

	result := make([]string, 0, len(values))
	for _, value := range values {
		parts := strings.Split(value, ",")
		for _, part := range parts {
			trimmed := strings.TrimSpace(part)
			if trimmed == "" {
				continue
			}
			result = append(result, trimmed)
		}
	}

	return result
}

func validateIntRange(raw string, min int, max int, field string) error {
	if raw == "" {
		return nil
	}

	value, err := strconv.Atoi(raw)
	if err != nil {
		return validationErrorf("invalid %s: must be an integer", field)
	}

	if value < min || value > max {
		return validationErrorf("invalid %s: must be between %d and %d", field, min, max)
	}

	return nil
}

func validateFloatRange(raw string, min float64, max float64, field string) error {
	if raw == "" {
		return nil
	}

	value, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		return validationErrorf("invalid %s: must be a number", field)
	}

	if value < min || value > max {
		return validationErrorf("invalid %s: must be between %g and %g", field, min, max)
	}

	return nil
}

func validateListItems(q url.Values, minItems int, maxItems int, field string) error {
	items := parseCSVQuery(q, field)
	if len(items) < minItems || len(items) > maxItems {
		if minItems == 0 {
			return validationErrorf("invalid %s: provide at most %d values", field, maxItems)
		}
		return validationErrorf("invalid %s: provide between %d and %d values", field, minItems, maxItems)
	}
	for _, item := range items {
		if item == "" {
			return validationErrorf("invalid %s: empty track ID is not allowed", field)
		}
	}
	return nil
}
