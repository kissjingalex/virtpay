package urls

import (
	cjson "github.com/kissjingalex/virtpay/internal/util/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

type urlsTestCase struct {
	RawUrl string
	Want   string
}

var urlsTestCases = []urlsTestCase{
	{RawUrl: "", Want: ""},
	{RawUrl: "https://example.com", Want: "example.com"},
	{RawUrl: "apis.example.com", Want: "apis.example.com"},
	{RawUrl: "apis.example.com/_/health", Want: "apis.example.com"},
	{RawUrl: "wss://example.com", Want: "example.com"},
	{RawUrl: "wss://example.com:8443", Want: "example.com:8443"},
	{RawUrl: "ftp://www.adobe.com", Want: "www.adobe.com"},
	{RawUrl: "ftp://root@www.adobe.com", Want: "www.adobe.com"},
	{RawUrl: "https://www.apple.com.cn/mac/", Want: "www.apple.com.cn"},
	{RawUrl: "https://stackoverflow.com/questions/62083272/parsing-url-with-port-and-without-scheme?t=1", Want: "stackoverflow.com"},
}

func TestTryParseRawUrl(t *testing.T) {
	for _, test := range urlsTestCases {
		scheme, host, path, query, err := TryParseRawUrl(test.RawUrl)
		if err != nil {
			t.Fatal(err)
		}
		assert.True(t, host == test.Want, "unexpected result of host, want %s but get %s", test.Want, host)
		t.Logf("raw: %#v, scheme: %s, host: %s, path: %s, query: %s", test, scheme, host, path, query)
	}
}

func TestParseRawUrl(t *testing.T) {
	for _, test := range urlsTestCases {
		u, err := ParseRawUrl(test.RawUrl)
		if err != nil {
			t.Fatal(err)
		}
		assert.True(t, u.Host == test.Want, "unexpected result of host, want %s but get %s", test.Want, u.Host)
		bs, _ := cjson.JSON.Marshal(u)
		t.Log(string(bs))
		if u.User != nil {
			t.Logf("%#v", u.User)
		}
	}
}
