package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/consul/api"
	"github.com/labstack/echo"
)

// Controller is the struct all API actions are registered on
type Controller struct {
	consulKV *api.KV
}

// New creates a new Controller and keeps the DB around for requests
func New(kv *api.KV) *Controller {
	return &Controller{consulKV: kv}
}

func (c *Controller) Consul(ctx echo.Context) error {
	prefix := ctx.QueryParam("prefix")
	if prefix != "" {
		return c.KVByPrefix(ctx, prefix)
	}
	key := ctx.QueryParam("key")
	if key != "" {
		return c.KVByKey(ctx, key)
	}
	return ctx.String(http.StatusNotFound, "")
}

func (c *Controller) KVByPrefix(ctx echo.Context, prefix string) error {

	result, _, err := c.consulKV.List(prefix, nil)
	if err != nil {
		panic(err)
	}
	upper := ctx.QueryParam("upper")
	var buf strings.Builder
	out := map[string]interface{}{}
	for _, pair := range result {
		pair.Key = strings.TrimPrefix(pair.Key, prefix)
		if !strings.HasPrefix(pair.Key, "/") {
			// We do this because consul's List() will return all matching prefixes.
			// We want an exact match. So if the key starts with anything but '/'
			// after trimming the prefix off, then it's not an exact match.
			continue
		}
		pair.Key = strings.TrimLeft(pair.Key, "/")
		if upper == "1" {
			pair.Key = strings.ToUpper(pair.Key)
		}

		s := strings.Split(pair.Key, "/")
		val, err := NestedMapLookup(out, string(pair.Value), s...)
		fmt.Println(pair.Key, " => ", pair.Value)
		fmt.Println("\t val => ", val)
		fmt.Println("\t err => ", err)
	}
	pretty, _ := json.MarshalIndent(out, "->", "\t")
	fmt.Println(string(pretty))

	//buf.WriteString(fmt.Sprintf("%v: %s\n", pair.Key, pair.Value))
	return ctx.String(http.StatusOK, buf.String())
}

func NestedMapLookup(m map[string]interface{}, val string, ks ...string) (rval interface{}, err error) {
	var ok bool
	if len(ks) == 0 { // degenerate input
		return nil, fmt.Errorf("NestedMapLookup needs at least one key")
	}
	if rval, ok = m[ks[0]]; !ok {
		// create key
		m[ks[0]] = map[string]interface{}{}
		return NestedMapLookup(m, val, ks[1:]...)
		//return nil, fmt.Errorf("key not found; remaining keys: %v", ks)
	} else if len(ks) == 1 { // we've reached the final key
		if rval == "" {
			rval = val
		}
		return rval, nil
	} else if m, ok = rval.(map[string]interface{}); !ok {
		return nil, fmt.Errorf("malformed structure at %#v", rval)
	} else { // 1+ more keys
		return NestedMapLookup(m, val, ks[1:]...)
	}
}

func (c *Controller) KVByKey(ctx echo.Context, key string) error {

	pair, _, err := c.consulKV.Get(key, nil)
	if err != nil {
		panic(err)
	}
	override := ctx.QueryParam("override")
	if override != "" {
		return ctx.String(http.StatusOK, fmt.Sprintf("%v: %s\n", override, pair.Value))
	}

	k := strings.Split(pair.Key, "/")
	keyFormatted := k[len(k)-1]

	upper := ctx.QueryParam("upper")
	if upper == "1" {
		keyFormatted = strings.ToUpper(keyFormatted)
	}

	return ctx.String(http.StatusOK, fmt.Sprintf("%v: %s\n", keyFormatted, pair.Value))
}
