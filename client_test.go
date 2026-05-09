package twitterapi

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"os"
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"
)

var apiKey string

const (
	testTweetID     = "1907859017325678727"
	testUserID      = "783214"
	testUserName    = "Twitter"
	testListID      = "1"
	testSpaceID     = "1PlJQbgjNpnJE"
	testCommunityID = "1"
)

func TestMain(m *testing.M) {
	apiKey = os.Getenv("TWITTERAPI_IO_KEY")
	os.Exit(m.Run())
}

func newTestClient(t *testing.T) *twitterApi {
	t.Helper()
	return New(apiKey)
}

func doGet(t *testing.T, client *twitterApi, url string) ([]byte, int, error) {
	t.Helper()
	if client == nil {
		return nil, 0, errors.New("client is nil")
	}
	ctx, cancel := context.WithTimeout(client.ctx, 10*time.Second)
	defer cancel()
	jsonData, resp, err := getDataWithHeader(ctx, client.httpClient, url, client.headers)
	if err != nil {
		return jsonData, 0, err
	}
	if resp == nil {
		return jsonData, 0, errors.New("response is nil")
	}
	return jsonData, resp.StatusCode, nil
}

func decodeJSONDisallowUnknowns(t *testing.T, raw []byte, target any) error {
	t.Helper()
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.DisallowUnknownFields()
	return decoder.Decode(target)
}

func logJSONFieldDiff(t *testing.T, raw []byte, target any) {
	t.Helper()
	actualKeys := jsonKeysFromRaw(t, raw)
	if len(actualKeys) == 0 {
		t.Logf("json key diff skipped: response is not a JSON object")
		return
	}
	expectedKeys := jsonKeysFromStruct(target)

	missing := make([]string, 0)
	for key := range expectedKeys {
		if _, ok := actualKeys[key]; !ok {
			missing = append(missing, key)
		}
	}
	extra := make([]string, 0)
	for key := range actualKeys {
		if _, ok := expectedKeys[key]; !ok {
			extra = append(extra, key)
		}
	}
	sort.Strings(missing)
	sort.Strings(extra)
	if len(missing) > 0 {
		t.Logf("missing JSON fields in response: %s", strings.Join(missing, ", "))
	}
	if len(extra) > 0 {
		t.Logf("extra JSON fields in response: %s", strings.Join(extra, ", "))
	}
}

func jsonKeysFromRaw(t *testing.T, raw []byte) map[string]struct{} {
	t.Helper()
	result := map[string]struct{}{}
	var object map[string]any
	if err := json.Unmarshal(raw, &object); err == nil {
		for key := range object {
			result[key] = struct{}{}
		}
		return result
	}
	var array []any
	if err := json.Unmarshal(raw, &array); err == nil && len(array) > 0 {
		if first, ok := array[0].(map[string]any); ok {
			for key := range first {
				result[key] = struct{}{}
			}
			return result
		}
	}
	return result
}

func jsonKeysFromStruct(value any) map[string]struct{} {
	result := map[string]struct{}{}
	if value == nil {
		return result
	}
	valueType := reflect.TypeOf(value)
	for valueType.Kind() == reflect.Ptr {
		valueType = valueType.Elem()
	}
	collectJSONTags(valueType, result)
	return result
}

func collectJSONTags(valueType reflect.Type, out map[string]struct{}) {
	if valueType.Kind() != reflect.Struct {
		return
	}
	for i := 0; i < valueType.NumField(); i++ {
		field := valueType.Field(i)
		if field.PkgPath != "" { // unexported
			continue
		}
		if field.Anonymous {
			collectJSONTags(field.Type, out)
			continue
		}
		tag := field.Tag.Get("json")
		if tag == "-" {
			continue
		}
		name := strings.Split(tag, ",")[0]
		if name == "" {
			name = field.Name
		}
		out[name] = struct{}{}
	}
}

func stringPtr(value string) *string {
	return &value
}

func intPtr(value int) *int {
	return &value
}

func int64Ptr(value int64) *int64 {
	return &value
}

func boolPtr(value bool) *bool {
	return &value
}

func last7DaysRangeUnix() (int64, int64) {
	until := time.Now().Unix()
	since := time.Now().Add(-7 * 24 * time.Hour).Unix()
	return since, until
}

func last7DaysRangeUnixInt() (int, int) {
	since, until := last7DaysRangeUnix()
	return int(since), int(until)
}
