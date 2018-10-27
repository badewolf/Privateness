package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/daemon"
	"github.com/skycoin/skycoin/src/readable"
	"github.com/skycoin/skycoin/src/util/useragent"
)

func TestConnection(t *testing.T) {
	tt := []struct {
		name                       string
		method                     string
		status                     int
		err                        string
		addr                       string
		gatewayGetConnectionResult *daemon.Connection
		gatewayGetConnectionError  error
		result                     *readable.Connection
	}{
		{
			name:   "405",
			method: http.MethodPost,
			status: http.StatusMethodNotAllowed,
			err:    "405 Method Not Allowed",
		},
		{
			name:                       "400 - empty addr",
			method:                     http.MethodGet,
			status:                     http.StatusBadRequest,
			err:                        "400 Bad Request - addr is required",
			addr:                       "",
			gatewayGetConnectionResult: nil,
			result:                     nil,
		},
		{
			name:   "200",
			method: http.MethodGet,
			status: http.StatusOK,
			err:    "",
			addr:   "addr",
			gatewayGetConnectionResult: &daemon.Connection{
				Addr: "127.0.0.1:6061",
				Gnet: daemon.GnetConnectionDetails{
					ID:           1,
					LastSent:     time.Unix(99999, 0),
					LastReceived: time.Unix(1111111, 0),
				},
				ConnectionDetails: daemon.ConnectionDetails{
					Outgoing:    true,
					ConnectedAt: time.Unix(222222, 0),
					State:       daemon.ConnectionStateIntroduced,
					Mirror:      6789,
					ListenPort:  9877,
					Height:      1234,
					UserAgent:   useragent.MustParse("skycoin:0.25.1(foo)"),
				},
			},
			result: &readable.Connection{
				Addr:         "127.0.0.1:6061",
				GnetID:       1,
				LastSent:     99999,
				LastReceived: 1111111,
				ConnectedAt:  222222,
				Outgoing:     true,
				State:        daemon.ConnectionStateIntroduced,
				Mirror:       6789,
				ListenPort:   9877,
				Height:       1234,
				UserAgent:    useragent.MustParse("skycoin:0.25.1(foo)"),
			},
		},

		{
			name:                      "500 - GetConnection failed",
			method:                    http.MethodGet,
			status:                    http.StatusInternalServerError,
			err:                       "500 Internal Server Error - GetConnection failed",
			addr:                      "addr",
			gatewayGetConnectionError: errors.New("GetConnection failed"),
		},

		{
			name:   "404",
			method: http.MethodGet,
			status: http.StatusNotFound,
			addr:   "addr",
			err:    "404 Not Found",
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			endpoint := "/api/v1/network/connection"
			gateway := &MockGatewayer{}
			gateway.On("GetConnection", tc.addr).Return(tc.gatewayGetConnectionResult, tc.gatewayGetConnectionError)

			v := url.Values{}
			if tc.addr != "" {
				v.Add("addr", tc.addr)
			}
			if len(v) > 0 {
				endpoint += "?" + v.Encode()
			}
			req, err := http.NewRequest(tc.method, endpoint, nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler := newServerMux(defaultMuxConfig(), gateway, &CSRFStore{}, nil)
			handler.ServeHTTP(rr, req)

			status := rr.Code
			require.Equal(t, tc.status, status, "got `%v` want `%v`", status, tc.status)

			if status != http.StatusOK {
				require.Equal(t, tc.err, strings.TrimSpace(rr.Body.String()), "got `%v`| %d, want `%v`",
					strings.TrimSpace(rr.Body.String()), status, tc.err)
			} else {
				var msg *readable.Connection
				err = json.Unmarshal(rr.Body.Bytes(), &msg)
				require.NoError(t, err)
				require.Equal(t, tc.result, msg)
			}
		})
	}
}

func TestConnections(t *testing.T) {
	tt := []struct {
		name                                 string
		method                               string
		status                               int
		err                                  string
		gatewayGetSolicitedConnectionsResult []daemon.Connection
		gatewayGetSolicitedConnectionsError  error
		result                               Connections
	}{
		{
			name:   "405",
			method: http.MethodPost,
			status: http.StatusMethodNotAllowed,
			err:    "405 Method Not Allowed",
		},
		{
			name:   "200",
			method: http.MethodGet,
			status: http.StatusOK,
			err:    "",
			gatewayGetSolicitedConnectionsResult: []daemon.Connection{
				{
					Addr: "127.0.0.1:6061",
					Gnet: daemon.GnetConnectionDetails{
						ID:           1,
						LastSent:     time.Unix(99999, 0),
						LastReceived: time.Unix(1111111, 0),
					},
					ConnectionDetails: daemon.ConnectionDetails{
						Outgoing:    true,
						State:       daemon.ConnectionStateIntroduced,
						ConnectedAt: time.Unix(222222, 0),
						Mirror:      9876,
						ListenPort:  9877,
						Height:      1234,
						UserAgent:   useragent.MustParse("skycoin:0.25.1(foo)"),
					},
				},
			},
			result: Connections{
				Connections: []readable.Connection{
					{
						Addr:         "127.0.0.1:6061",
						GnetID:       1,
						LastSent:     99999,
						LastReceived: 1111111,
						ConnectedAt:  222222,
						Outgoing:     true,
						State:        daemon.ConnectionStateIntroduced,
						Mirror:       9876,
						ListenPort:   9877,
						Height:       1234,
						UserAgent:    useragent.MustParse("skycoin:0.25.1(foo)"),
					},
				},
			},
		},

		{
			name:                                "500 - GetOutgoingConnections failed",
			method:                              http.MethodGet,
			status:                              http.StatusInternalServerError,
			err:                                 "500 Internal Server Error - GetOutgoingConnections failed",
			gatewayGetSolicitedConnectionsError: errors.New("GetOutgoingConnections failed"),
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			endpoint := "/api/v1/network/connections"
			gateway := &MockGatewayer{}
			gateway.On("GetOutgoingConnections").Return(tc.gatewayGetSolicitedConnectionsResult, tc.gatewayGetSolicitedConnectionsError)

			req, err := http.NewRequest(tc.method, endpoint, nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler := newServerMux(defaultMuxConfig(), gateway, &CSRFStore{}, nil)
			handler.ServeHTTP(rr, req)

			status := rr.Code
			require.Equal(t, tc.status, status, "got `%v` want `%v`", status, tc.status)

			if status != http.StatusOK {
				require.Equal(t, tc.err, strings.TrimSpace(rr.Body.String()), "got `%v`| %d, want `%v`",
					strings.TrimSpace(rr.Body.String()), status, tc.err)
			} else {
				var msg Connections
				err = json.Unmarshal(rr.Body.Bytes(), &msg)
				require.NoError(t, err)
				require.Equal(t, tc.result, msg)
			}
		})
	}
}

func TestDefaultConnections(t *testing.T) {
	tt := []struct {
		name                               string
		method                             string
		status                             int
		err                                string
		gatewayGetDefaultConnectionsResult []string
		result                             []string
	}{
		{
			name:   "405",
			method: http.MethodPost,
			status: http.StatusMethodNotAllowed,
			err:    "405 Method Not Allowed",
		},
		{
			name:                               "200",
			method:                             http.MethodGet,
			status:                             http.StatusOK,
			err:                                "",
			gatewayGetDefaultConnectionsResult: []string{"44.33.22.11", "11.44.66.88"},
			result:                             []string{"11.44.66.88", "44.33.22.11"},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			endpoint := "/api/v1/network/defaultConnections"
			gateway := &MockGatewayer{}
			gateway.On("GetDefaultConnections").Return(tc.gatewayGetDefaultConnectionsResult)

			req, err := http.NewRequest(tc.method, endpoint, nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler := newServerMux(defaultMuxConfig(), gateway, &CSRFStore{}, nil)
			handler.ServeHTTP(rr, req)

			status := rr.Code
			require.Equal(t, tc.status, status, "got `%v` want `%v`", status, tc.status)

			if status != http.StatusOK {
				require.Equal(t, tc.err, strings.TrimSpace(rr.Body.String()), "got `%v`| %d, want `%v`",
					strings.TrimSpace(rr.Body.String()), status, tc.err)
			} else {
				var msg []string
				err = json.Unmarshal(rr.Body.Bytes(), &msg)
				require.NoError(t, err)
				require.Equal(t, tc.result, msg)
			}
		})
	}
}

func TestGetTrustConnections(t *testing.T) {
	tt := []struct {
		name                             string
		method                           string
		status                           int
		err                              string
		gatewayGetTrustConnectionsResult []string
		result                           []string
	}{
		{
			name:   "405",
			method: http.MethodPost,
			status: http.StatusMethodNotAllowed,
			err:    "405 Method Not Allowed",
		},
		{
			name:                             "200",
			method:                           http.MethodGet,
			status:                           http.StatusOK,
			err:                              "",
			gatewayGetTrustConnectionsResult: []string{"44.33.22.11", "11.44.66.88"},
			result:                           []string{"11.44.66.88", "44.33.22.11"},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			endpoint := "/api/v1/network/connections/trust"
			gateway := &MockGatewayer{}
			gateway.On("GetTrustConnections").Return(tc.gatewayGetTrustConnectionsResult)

			req, err := http.NewRequest(tc.method, endpoint, nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler := newServerMux(defaultMuxConfig(), gateway, &CSRFStore{}, nil)
			handler.ServeHTTP(rr, req)

			status := rr.Code
			require.Equal(t, tc.status, status, "got `%v` want `%v`", status, tc.status)

			if status != http.StatusOK {
				require.Equal(t, tc.err, strings.TrimSpace(rr.Body.String()), "got `%v`| %d, want `%v`",
					strings.TrimSpace(rr.Body.String()), status, tc.err)
			} else {
				var msg []string
				err = json.Unmarshal(rr.Body.Bytes(), &msg)
				require.NoError(t, err)
				require.Equal(t, tc.result, msg)
			}
		})
	}
}

func TestGetExchgConnection(t *testing.T) {
	tt := []struct {
		name                            string
		method                          string
		status                          int
		err                             string
		gatewayGetExchgConnectionResult []string
		result                          []string
	}{
		{
			name:   "405",
			method: http.MethodPost,
			status: http.StatusMethodNotAllowed,
			err:    "405 Method Not Allowed",
		},
		{
			name:                            "200",
			method:                          http.MethodGet,
			status:                          http.StatusOK,
			err:                             "",
			gatewayGetExchgConnectionResult: []string{"44.33.22.11", "11.44.66.88"},
			result:                          []string{"11.44.66.88", "44.33.22.11"},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			endpoint := "/api/v1/network/connections/exchange"
			gateway := &MockGatewayer{}
			gateway.On("GetExchgConnection").Return(tc.gatewayGetExchgConnectionResult)

			req, err := http.NewRequest(tc.method, endpoint, nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler := newServerMux(defaultMuxConfig(), gateway, &CSRFStore{}, nil)
			handler.ServeHTTP(rr, req)

			status := rr.Code
			require.Equal(t, tc.status, status, "got `%v` want `%v`", status, tc.status)

			if status != http.StatusOK {
				require.Equal(t, tc.err, strings.TrimSpace(rr.Body.String()), "got `%v`| %d, want `%v`",
					strings.TrimSpace(rr.Body.String()), status, tc.err)
			} else {
				var msg []string
				err = json.Unmarshal(rr.Body.Bytes(), &msg)
				require.NoError(t, err)
				require.Equal(t, tc.result, msg)
			}
		})
	}
}
