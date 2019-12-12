package messaging

import (
	"context"
	"errors"
	phttp "github.com/beatlabs/patron/trace/http"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mkcl struct{}

func (m *mkcl) Do(ctx context.Context, req *http.Request) (*http.Response, error) {
	return nil, errors.New("wrong")
}

func TestHTTPClient_Do(t *testing.T) {
	bcl, err := phttp.New()
	assert.Nil(t, err)

	rs := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(200)
		rw.Write([]byte("hello"))
	}))
	cl := NewHTTPClient(bcl, "testprovider")
	rq, _ := http.NewRequest("GET", rs.URL, nil)
	rsp, err := cl.Do(context.Background(), rq)
	assert.Nil(t, err)
	assert.Equal(t, 200, rsp.StatusCode)
	bd, err := ioutil.ReadAll(rsp.Body)
	assert.Nil(t, err)
	assert.Equal(t, "hello", string(bd))
}

func TestClientError(t *testing.T) {
	cl := NewHTTPClient(&mkcl{}, "testprovider")
	rq, _ := http.NewRequest("GET", "", nil)
	rsp, err := cl.Do(context.Background(), rq)
	assert.Nil(t, rsp)
	assert.EqualError(t, err, "wrong")
}
