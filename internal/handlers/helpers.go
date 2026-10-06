package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	gonertia "github.com/romsar/gonertia/v3"
)

// setValidationErrors stores gonertia validation errors in the request context.
// gonertia.SetValidationErrors returns a new context; this helper mutates the
// request in-place so that i.Back() picks the errors up correctly.
func setValidationErrors(r *http.Request, errors gonertia.ValidationErrors) {
	ctx := gonertia.SetValidationErrors(r.Context(), errors)
	*r = *r.WithContext(ctx)
}

func parseJSONOrForm(r *http.Request, target any) error {
	contentType := r.Header.Get("Content-Type")
	if strings.Contains(contentType, "application/json") {
		return json.NewDecoder(r.Body).Decode(target)
	}
	if err := r.ParseForm(); err != nil {
		return err
	}
	data := make(map[string]any)
	for key, values := range r.PostForm {
		if len(values) == 1 {
			val := values[0]
			if val == "true" {
				data[key] = true
			} else if val == "false" {
				data[key] = false
			} else if num, err := strconv.ParseFloat(val, 64); err == nil && !strings.HasPrefix(val, "0") {
				data[key] = num
			} else {
				data[key] = val
			}
		} else if len(values) > 1 {
			data[key] = values
		}
	}
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, target)
}

func getIntQuery(r *http.Request, key string, defaultVal int) int {
	valStr := r.URL.Query().Get(key)
	if valStr == "" {
		return defaultVal
	}
	val, err := strconv.Atoi(valStr)
	if err != nil || val <= 0 {
		return defaultVal
	}
	return val
}