package http

import (
	"encoding/binary"

	"gitlab.com/mjwhitta/win/winhttp"
)

// NewClient will return a pointer to a new Client instnace that
// simply wraps the net/http.Client type.
func NewClient() (*Client, error) {
	var c = &Client{}
	var e error

	// Create session
	c.hndl, e = winhttp.Open(
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.85 Safari/537.36",
		winhttp.WinhttpAccessTypeAutomaticProxy,
		"",
		"",
		0,
	)
	if e != nil {
		return nil, e
	}

	return c, nil
}

// Get will make a GET request using Winhttp.dll.
func (c *Client) Get(
	dst string,
	headers map[string]string,
) (*Response, error) {
	return c.request(MethodGet, dst, headers, nil)
}

// Head will make a HEAD request using Winhttp.dll.
func (c *Client) Head(
	dst string,
	headers map[string]string,
) (*Response, error) {
	return c.request(MethodHead, dst, headers, nil)
}

// Post will make a POST request using Winhttp.dll.
func (c *Client) Post(
	dst string,
	headers map[string]string,
	data []byte,
) (*Response, error) {
	return c.request(MethodPost, dst, headers, data)
}

// Put will make a PUT request using Winhttp.dll.
func (c *Client) Put(
	dst string,
	headers map[string]string,
	data []byte,
) (*Response, error) {
	return c.request(MethodPut, dst, headers, data)
}

func (c *Client) request(
	method string,
	dst string,
	headers map[string]string,
	data []byte,
) (*Response, error) {
	var e error
	var reqHndl uintptr
	var res *Response
	var tlsIgnore uintptr
	var tmp []byte

	if reqHndl, e = buildRequest(c.hndl, method, dst); e != nil {
		return nil, e
	}

	if c.Timeout > 0 {
		tmp = make([]byte, 4)
		binary.LittleEndian.PutUint32(
			tmp,
			uint32(c.Timeout.Milliseconds()),
		)

		e = winhttp.SetOption(
			reqHndl,
			winhttp.WinhttpOptionConnectTimeout,
			tmp,
			len(tmp),
		)
		if e != nil {
			return nil, e
		}

		e = winhttp.SetOption(
			reqHndl,
			winhttp.WinhttpOptionReceiveResponseTimeout,
			tmp,
			len(tmp),
		)
		if e != nil {
			return nil, e
		}

		e = winhttp.SetOption(
			reqHndl,
			winhttp.WinhttpOptionReceiveTimeout,
			tmp,
			len(tmp),
		)
		if e != nil {
			return nil, e
		}

		e = winhttp.SetOption(
			reqHndl,
			winhttp.WinhttpOptionResolveTimeout,
			tmp,
			len(tmp),
		)
		if e != nil {
			return nil, e
		}

		e = winhttp.SetOption(
			reqHndl,
			winhttp.WinhttpOptionSendTimeout,
			tmp,
			len(tmp),
		)
		if e != nil {
			return nil, e
		}
	}

	if c.TLSClientConfig.InsecureSkipVerify {
		tlsIgnore |= winhttp.SecurityFlagIgnoreUnknownCa
		tlsIgnore |= winhttp.SecurityFlagIgnoreCertDateInvalid
		tlsIgnore |= winhttp.SecurityFlagIgnoreCertCnInvalid
		tlsIgnore |= winhttp.SecurityFlagIgnoreCertWrongUsage

		tmp = make([]byte, 4)
		binary.LittleEndian.PutUint32(tmp, uint32(tlsIgnore))

		e = winhttp.SetOption(
			reqHndl,
			winhttp.WinhttpOptionSecurityFlags,
			tmp,
			len(tmp),
		)
		if e != nil {
			return nil, e
		}
	}

	if e = sendRequest(reqHndl, headers, data); e != nil {
		return nil, e
	}

	if res, e = buildResponse(reqHndl); e != nil {
		return nil, e
	}

	return res, nil
}
