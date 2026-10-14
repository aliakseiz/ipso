package ipso_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/aliakseiz/ipso"
)

const listPath = "/Object"

const testObjectXML = `<LWM2M><Object ObjectType="MODefinition">
<Name>Temperature</Name><Description1>desc</Description1><ObjectID>3303</ObjectID>
<ObjectURN>urn:oma:lwm2m:ext:3303</ObjectURN><MultipleInstances>Multiple</MultipleInstances>
<Mandatory>Optional</Mandatory>
<Resources><Item ID="5700"><Name>Sensor Value</Name><Operations>R</Operations>
<MultipleInstances>Single</MultipleInstances><Mandatory>Mandatory</Mandatory><Type>Float</Type>
<Units>Cel</Units></Item></Resources></Object></LWM2M>`

func TestReg_ImportFromAPI(t *testing.T) {
	t.Parallel()

	reg := newAPIServer(t, map[string]http.HandlerFunc{
		listPath:  serveList(`[{"ObjectID":3303,"ObjectLink":"%[1]s/o/3303"}]`),
		"/o/3303": serveXML(testObjectXML),
	})

	objects, err := reg.ImportFromAPI()
	require.NoError(t, err)
	require.Len(t, objects, 1)

	obj := objects[0]
	require.Equal(t, int32(3303), obj.ObjectID)
	require.Equal(t, "Temperature", obj.Name)
	require.Equal(t, ipso.DefaultObjectVersion, obj.ObjectVersion, "missing object version gets default")
	require.Equal(t, ipso.DefaultLwM2MVersion, obj.LWM2MVersion, "missing LwM2M version gets default")
	require.Len(t, obj.Resources.Item, 1)
	require.Equal(t, int32(5700), obj.Resources.Item[0].ID)
	require.Equal(t, ipso.ResourceTypeFloat, obj.Resources.Item[0].Type)
}

func TestReg_ImportFromAPIListErrors(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		list    http.HandlerFunc
		wantErr error
	}{
		"HTTP 404": {list: status(http.StatusNotFound), wantErr: ipso.ErrHTTPStatus},
		"not JSON": {list: serveXML("<html>maintenance</html>")},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			reg := newAPIServer(t, map[string]http.HandlerFunc{listPath: tt.list})
			_, err := reg.ImportFromAPI()
			require.Error(t, err)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
			}
		})
	}
}

func TestReg_ImportFromAPINoObjects(t *testing.T) {
	t.Parallel()

	reg := newAPIServer(t, map[string]http.HandlerFunc{
		listPath: serveList(`[{"ObjectID":1,"ObjectLink":"%[1]s/bad"},{"ObjectID":2,"ObjectLink":""}]`),
		"/bad":   status(http.StatusInternalServerError),
	})
	reg.Config.SkipInitErrors = true

	_, err := reg.ImportFromAPI()
	require.ErrorIs(t, err, ipso.ErrNoObjects)
}

func TestNew_Sanitizer(t *testing.T) {
	t.Parallel()

	const dirty = "<LWM2M><Object><Name>N</Name><Description1>a#b\tc</Description1><ObjectID>1</ObjectID></Object></LWM2M>"

	tests := map[string]struct {
		cfg  ipso.Configuration
		want string
	}{
		"custom sanitizer":  {cfg: ipso.Configuration{Sanitize: true, Sanitizer: []string{"#"}}, want: "ab\tc"},
		"default sanitizer": {cfg: ipso.Configuration{Sanitize: true}, want: "a#b c"},
		"no sanitize":       {cfg: ipso.Configuration{Sanitizer: []string{"#"}}, want: "a#b\tc"},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			srv := newAPIServer(t, map[string]http.HandlerFunc{
				listPath: serveList(`[{"ObjectID":1,"ObjectLink":"%[1]s/o"}]`),
				"/o":     serveXML(dirty),
			})

			tt.cfg.InitOnNew = true
			tt.cfg.BaseURL = srv.Config.BaseURL
			tt.cfg.HTTPClient = srv.Config.HTTPClient

			reg, err := ipso.New(tt.cfg)
			require.NoError(t, err)
			require.Len(t, reg.Objects, 1)
			require.Equal(t, tt.want, reg.Objects[0].Description1)
		})
	}
}

func TestReg_ImportFromAPIUnreachable(t *testing.T) {
	t.Parallel()

	reg := &ipso.Reg{Config: ipso.Configuration{BaseURL: "http://127.0.0.1:1"}}

	_, err := reg.ImportFromAPI()
	require.Error(t, err)
}

func TestReg_ImportFromAPIObjectErrors(t *testing.T) {
	t.Parallel()

	tests := map[string]http.HandlerFunc{
		"HTTP 500": status(http.StatusInternalServerError),
		"bad XML":  serveXML("<LWM2M><Object>"),
	}

	for name, broken := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			reg := newAPIServer(t, map[string]http.HandlerFunc{
				listPath: serveList(`[{"ObjectID":1,"ObjectLink":"%[1]s/bad"},` +
					`{"ObjectID":3303,"ObjectLink":"%[1]s/ok"},{"ObjectID":2,"ObjectLink":""}]`),
				"/bad": broken,
				"/ok":  serveXML(testObjectXML),
			})

			// Strict mode fails on the first broken object.
			_, err := reg.ImportFromAPI()
			require.Error(t, err)

			// Lenient mode skips broken objects and objects without a link.
			reg.Config.SkipInitErrors = true
			objects, err := reg.ImportFromAPI()
			require.NoError(t, err)
			require.Len(t, objects, 1)
			require.Equal(t, int32(3303), objects[0].ObjectID)
		})
	}
}

func TestReg_ImportFromAPIContextCanceled(t *testing.T) {
	t.Parallel()

	for _, skip := range []bool{false, true} {
		t.Run(fmt.Sprintf("SkipInitErrors=%t", skip), func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithCancel(t.Context())

			reg := newAPIServer(t, map[string]http.HandlerFunc{
				listPath: serveList(`[{"ObjectID":1,"ObjectLink":"%[1]s/o"},{"ObjectID":2,"ObjectLink":"%[1]s/o"}]`),
				"/o": func(w http.ResponseWriter, r *http.Request) {
					cancel() // the caller gives up while the first object is being fetched
					serveXML(testObjectXML)(w, r)
				},
			})
			reg.Config.SkipInitErrors = skip

			objects, err := reg.ImportFromAPIContext(ctx)
			require.ErrorIs(t, err, context.Canceled)
			require.Nil(t, objects)
		})
	}
}

func TestNewContext(t *testing.T) {
	t.Parallel()

	srv := newAPIServer(t, map[string]http.HandlerFunc{
		listPath: serveList(`[{"ObjectID":3303,"ObjectLink":"%[1]s/o"}]`),
		"/o":     serveXML(testObjectXML),
	})
	cfg := srv.Config
	cfg.InitOnNew = true

	reg, err := ipso.NewContext(t.Context(), cfg)
	require.NoError(t, err)
	require.Len(t, reg.Objects, 1)

	ctx, cancel := context.WithCancel(t.Context())
	cancel()

	_, err = ipso.NewContext(ctx, cfg)
	require.ErrorIs(t, err, context.Canceled)
}

func TestReg_ImportFromAPIDefaults(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		baseURL string
		want    string
	}{
		"default base URL":        {baseURL: "", want: ipso.DefaultBaseURLForTest + ipso.ObjectsPathForTest},
		"custom base URL":         {baseURL: "https://example.test/api", want: "https://example.test/api/Object"},
		"base URL trailing slash": {baseURL: "https://example.test/api/", want: "https://example.test/api/Object"},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			var urls []string

			client := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
				urls = append(urls, r.URL.String())

				return nil, errTransport
			})}

			reg := &ipso.Reg{Config: ipso.Configuration{HTTPClient: client, BaseURL: tt.baseURL}}

			_, err := reg.ImportFromAPI()
			require.ErrorIs(t, err, errTransport, "request goes through the configured client")
			require.Equal(t, []string{tt.want}, urls)
		})
	}
}

// newAPIServer serves the handlers and returns a registry whose base URL and client point at the server.
func newAPIServer(t *testing.T, handlers map[string]http.HandlerFunc) *ipso.Reg {
	t.Helper()

	mux := http.NewServeMux()
	for path, h := range handlers {
		mux.HandleFunc(path, h)
	}

	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)

	return &ipso.Reg{Config: ipso.Configuration{BaseURL: srv.URL, HTTPClient: srv.Client()}}
}

// serveList serves a JSON objects list, replacing %[1]s in the format with the server base URL.
func serveList(format string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, format, "http://"+r.Host) //nolint:gosec // test server, no HTML rendered
	}
}

func serveXML(body string) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		_, _ = fmt.Fprint(w, body)
	}
}

var errTransport = errors.New("transport stub")

// roundTripFunc adapts a function to http.RoundTripper.
type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func status(code int) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "failure", code)
	}
}
