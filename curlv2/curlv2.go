package curlv2

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/http/httputil"

	"github.com/codegangsta/cli"
)

func CliDomainAction(c *cli.Context) {
	domain := c.String("domain")
	if len(domain) == 0 {
		fmt.Println("Please provide cloudfoundry sysdomain")
	}
	t := Target{
		Domain: domain,
	}
	https := c.Bool("https")
	var urlBuilder Next
	if https {
		urlBuilder = HttpsUrlBuilder
	} else {
		urlBuilder = HttpUrlBuilder
	}

	err := urlBuilder.Next(BuildRequest).Request(t)
	if err != nil {
		fmt.Println(err)
	}
}

func CliRouterAction(c *cli.Context) {
	domain := c.String("domain")
	if len(domain) == 0 {
		fmt.Println("Please provide cloudfoundry sysdomain")
	}
	host := c.String("ip")
	if len(host) == 0 {
		fmt.Println("Please provide router IP")
	}
	t := Target{
		Domain: domain,
		Host:   host,
	}

	err := RouterUrlBuilder.Next(BuildRequest).Next(AddDomainHeader).Request(t)
	if err != nil {
		fmt.Println(err)
	}
}

type Target struct {
	Domain string
	Host   string
	Url    string
	Req    *http.Request
}

type Next func(t Target) (Target, error)

var HttpUrlBuilder Next = func(t Target) (Target, error) {
	return Target{
		Url: fmt.Sprintf("http://api.%s/v2/info", t.Domain),
	}, nil
}

var HttpsUrlBuilder Next = func(t Target) (Target, error) {
	return Target{
		Url: fmt.Sprintf("https://api.%s/v2/info", t.Domain),
	}, nil
}

var RouterUrlBuilder Next = func(t Target) (Target, error) {
	return Target{
		Url:    fmt.Sprintf("http://%s/v2/info", t.Host),
		Domain: t.Domain,
	}, nil
}

var BuildRequest Next = func(t Target) (target Target, err error) {
	req, err := http.NewRequest("GET", t.Url, nil)
	if err != nil {
		return
	}
	return Target{
		Domain: t.Domain,
		Host:   t.Host,
		Url:    t.Url,
		Req:    req,
	}, nil

}

var AddDomainHeader Next = func(t Target) (target Target, err error) {
	t.Req.Host = fmt.Sprintf("api.%s", t.Domain)
	return t, nil
}

func (current Next) Next(next Next) Next {
	return func(t Target) (target Target, err error) {
		target, err = current(t)
		if err != nil {
			return
		}
		return next(target)
	}
}

func (n Next) Request(t Target) error {
	target, err := n(t)
	if err != nil {
		return err
	}
	return request(target.Req)
}

func request(req *http.Request) (err error) {
	fmt.Printf("Perform Request to: %s\n", req.URL.String())
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	dump, _ := httputil.DumpRequest(req, true)
	fmt.Println(string(dump))

	response, err := transport.RoundTrip(req)
	if err != nil {
		return
	}

	dump, _ = httputil.DumpResponse(response, true)
	fmt.Println(string(dump))
	return
}
