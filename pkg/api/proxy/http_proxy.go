
// RAINBOND, Application Management Platform
// Copyright (C) 2014-2017 Goodrain Co., Ltd.
 
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version. For any non-GPL usage of Rainbond,
// one or multiple Commercial Licenses authorized by Goodrain Co., Ltd.
// must be obtained first.
 
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
 
// You should have received a copy of the GNU General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package proxy

import (
	"fmt"
	"net/http"
	"net/http/httputil"

	"github.com/Sirupsen/logrus"
)

//HTTPProxy HTTPProxy
type HTTPProxy struct {
	name      string
	endpoints EndpointList
	lb        LoadBalance
}

//Proxy 代理
func (h *HTTPProxy) Proxy(w http.ResponseWriter, r *http.Request) {
	endpoint := h.lb.Select(r, h.endpoints)
	logrus.Info(fmt.Sprintf("select a endpoint %s", endpoint))
	director := func(req *http.Request) {
		req = r
		req.URL.Scheme = "http"
		req.URL.Host = endpoint.String()
	}
	proxy := &httputil.ReverseProxy{Director: director}
	proxy.ServeHTTP(w, r)
}

//UpdateEndpoints 更新端点
func (h *HTTPProxy) UpdateEndpoints(endpoints ...string) {
	h.endpoints = CreateEndpoints(endpoints)
}

func createHTTPProxy(name string, endpoints []string) *HTTPProxy {
	return &HTTPProxy{name, CreateEndpoints(endpoints), NewRoundRobin()}
}
