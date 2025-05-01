package redis

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/redis/go-redis/v9"
	tracerRedis "gopkg.in/DataDog/dd-trace-go.v1/contrib/redis/go-redis.v9"
)

// errors
var (
	ErrInputNotSlice = errors.New("input is not a slice")
	ErrInputEmpty    = errors.New("input is Nil")
	ErrParsingType   = errors.New("parsing type error")
)

// Redis Command
const (
	FTSearchCmd string = "FT.SEARCH"
	SortAscArg  string = "ASC"
	SortByArg   string = "SORTBY"
)

type WrapperClient struct {
	Client *redis.Client
}

type RedisConfig struct {
	Host      string
	Port      string
	Username  string
	Password  string
	Index     string
	KeyPrefix string
}

func NewRedisConnection(ctx context.Context, cfg *RedisConfig) (*WrapperClient, error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.Host + ":" + cfg.Port,
		Username: cfg.Username,
		Password: cfg.Password,
		Protocol: 2,
	})

	err := redisClient.Do(ctx, "FT.INFO", cfg.Index).Err()
	if err != nil {
		query := fmt.Sprintf("FT.CREATE %s ON JSON PREFIX 1 %s: SCHEMA $.id AS id TEXT $.surname AS surname TEXT $.lastname AS lastname TEXT", cfg.Index, cfg.KeyPrefix)
		argS := strings.Split(query, " ")
		argI, _ := sliceToInterface(argS)
		err := redisClient.Do(ctx, argI...).Err()
		if err != nil {
			return nil, err
		}
	}

	pong, err := redisClient.Ping(ctx).Result()
	if err != nil || pong == "" {
		return nil, err
	}

	tracerRedis.WrapClient(redisClient)
	return &WrapperClient{
		redisClient,
	}, nil
}

// SearchJSONData searches for documents within a Redis full-text index using the FT.SEARCH command, allowing filtering and result limitation. More details can be found at: https://redis.io/commands/ft.search/.
// This function returns a slice of keys and values of the documents that match the search query, as well as an error.
// Note that go-redis still does not support the built-in FT.SEARCH command for now (v9.5.1). Therefore, we use Do() to execute the command.
func (r *WrapperClient) SearchJSONDataTag(c context.Context, index, field, value string, args ...interface{}) ([]string, []string, error) {
	query := fmt.Sprintf("'@%s:{%s}'", field, escapeString(value))
	allArgs := append([]interface{}{FTSearchCmd, index, query}, args...)

	reply, err := r.Client.Do(c, allArgs...).Result()
	if err != nil {
		return nil, nil, err
	}

	return wrapSearchOutput(reply)
}

// SearchJSONData searches for documents within a Redis full-text index using the FT.SEARCH command, allowing filtering and result limitation. More details can be found at: https://redis.io/commands/ft.search/.
// This function returns a slice of keys and values of the documents that match the search query, as well as an error.
// Note that go-redis still does not support the built-in FT.SEARCH command for now (v9.5.1). Therefore, we use Do() to execute the command.
func (r *WrapperClient) SearchJSONDataText(c context.Context, index, field, value string, args ...interface{}) ([]string, []string, error) {
	query := fmt.Sprintf("'@%s:(%s)'", field, escapeString(value))
	allArgs := append([]interface{}{FTSearchCmd, index, query}, args...)

	reply, err := r.Client.Do(c, allArgs...).Result()
	if err != nil {
		return nil, nil, err
	}

	return wrapSearchOutput(reply)
}

// SearchJSONData searches for documents within a Redis full-text index using the FT.SEARCH command, allowing filtering and result limitation. More details can be found at: https://redis.io/commands/ft.search/.
// This function returns a slice of keys and values of the documents that match the search query, as well as an error.
// Note that go-redis still does not support the built-in FT.SEARCH command for now (v9.5.1). Therefore, we use Do() to execute the command.
func (r *WrapperClient) SearchJSONDataByQuery(c context.Context, index, query string, args ...interface{}) ([]string, []string, error) {
	allArgs := append([]interface{}{FTSearchCmd, index, query}, args...)

	reply, err := r.Client.Do(c, allArgs...).Result()
	if err != nil {
		return nil, nil, err
	}

	return wrapSearchOutput(reply)
}

// SearchJSONDataNumeric searches for documents within a specific numeric range using an index. More details can be found at: https://redis.io/commands/ft.search/.
// This function returns a slice of keys and values of the documents that match the search query, as well as an error.
// Note: go-redis still does not support the built-in FT.SEARCH command for now (v9.5.1). Therefore, we use Do() to execute the command.
func (r *WrapperClient) SearchJSONDataNumeric(c context.Context, index, field string, v1, v2 int, args ...interface{}) ([]string, []string, error) {
	value1 := strconv.Itoa(v1)
	value2 := strconv.Itoa(v2)

	query := fmt.Sprintf("'@%s:[%s %s]", field, escapeString(value1), escapeString(value2))
	allArgs := append([]interface{}{FTSearchCmd, index, query}, args...)

	reply, err := r.Client.Do(c, allArgs...).Result()
	if err != nil {
		return nil, nil, err
	}

	return wrapSearchOutput(reply)
}

// SearchJSONDataAll get all document of specific index using wildcard *: https://redis.io/commands/ft.search/.
// This function returns a slice of keys and values of the documents that match the search query, as well as an error.
// Note: go-redis still does not support the built-in FT.SEARCH command (v9.5.x). So, we use Do() to execute the command.
func (r *WrapperClient) SearchJSONDataAll(c context.Context, index string, args ...interface{}) ([]string, []string, error) {
	query := "*"
	allArgs := append([]interface{}{FTSearchCmd, index, query}, args...)

	reply, err := r.Client.Do(c, allArgs...).Result()
	if err != nil {
		return nil, nil, err
	}
	return wrapSearchOutput(reply)
}

// sliceToInterface is a helper function that converts an interface to an slice of interface. The purpose of this function is to convert the string to a slice of interface to be used as command in the Do() function of the go-redis library.
func sliceToInterface(inp interface{}) ([]interface{}, error) {
	s := reflect.ValueOf(inp)
	if s.Kind() != reflect.Slice {
		return nil, ErrInputNotSlice
	}
	if s.IsNil() {
		return nil, ErrInputEmpty
	}

	res := make([]interface{}, s.Len())
	for i := 0; i < s.Len(); i++ {
		res[i] = s.Index(i).Interface()
	}

	return res, nil
}

// escapeString is a helper function that escapes special characters in a string. This is required by Redis. Read more details at: https://redis.io/docs/latest/develop/interact/search-and-query/advanced-concepts/tags/
func escapeString(inp string) string {
	var (
		out  []string
		regx = regexp.MustCompile("[^A-Za-z0-9]")
	)

	for i := 0; i < len(inp); i++ {
		if regx.MatchString(string(inp[i])) {
			v := fmt.Sprintf("\\%s", string(inp[i]))
			out = append(out, v)
		} else {
			out = append(out, string(inp[i]))
		}
	}
	return strings.Join(out, "")
}

// wrapSearchOutput is a helper function that wraps the output of the FT.SEARCH command into a map[string]string where the key of the map is the Redis key and the value is the Redis value.
func wrapSearchOutput(input interface{}) ([]string, []string, error) {
	res, ok := input.([]interface{})
	if !ok {
		return nil, nil, ErrParsingType
	}

	totalResults, ok := res[0].(int64)
	if !ok {
		return nil, nil, ErrParsingType
	}

	keys := make([]string, 0, totalResults)
	values := make([]string, 0, totalResults)

	if totalResults != 0 {
		for i := 1; i < len(res); i += 2 {
			// Get redis key.
			k := res[i].(string)
			if !ok {
				return nil, nil, ErrParsingType
			}

			// Get redis value.
			p := res[i+1].([]interface{})
			v, ok := p[len(p)-1].(string)
			if !ok {
				return nil, nil, ErrParsingType
			}

			// Append to keys and values.
			values = append(values, v)
			keys = append(keys, k)
		}
	}

	return keys, values, nil
}
