package goutils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/Yapo/logger"
)

// Response is a struct to generate a response from POST/PUT requests
type Response struct {
	Code    int
	Body    interface{}
	Headers http.Header
}

// ErrorType allow have common error on diferent projects
type ErrorType struct {
	Code    int
	Message string
}

// WriteJSONResponse write to te response stream
func WriteJSONResponse(w http.ResponseWriter, response *Response) {
	header := w.Header()
	if !reflect.ValueOf(response).IsNil() {
		for key, values := range response.Headers {
			for _, value := range values {
				header.Add(key, value)
			}
		}
	}
	header.Set("Content-Type", "application/json")
	w.WriteHeader(response.Code)
	fmt.Fprintf(w, "%s", response.Body)
}

// CreateJSON converts response.Body to json format
func CreateJSON(response *Response) {
	body := new(bytes.Buffer)
	encoder := json.NewEncoder(body)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(response.Body)

	if err != nil {
		logger.Info("CAN'T ENCODE \"%+v\" TO JSON", response.Body)
		response.Body = ""
		response.Code = http.StatusInternalServerError
		return
	}
	response.Body = body
}

// ParseJSONBody
func ParseJSONBody(r *http.Request, input interface{}) *Response {
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(input)
	r.Body.Close()
	if err != nil {
		return &Response{
			Code: http.StatusBadRequest,
			Body: GenericError{
				ErrorMessage: err.Error(),
			},
		}
	}
	return nil
}

// SendHTTPRequest does that
func SendHTTPRequest(url, endpoint, method, query, body string, headers map[string]string, timeout int, errors ErrorType) (int, string) {
	req, _ := http.NewRequest(method, url+endpoint+query, strings.NewReader(body))
	for key, value := range headers {
		req.Header.Add(key, value)
	}
	logger.Debug("sending HTTP request: %+v", req)
	// this makes a custom http client with a timeout in secs for each request
	var httpClient = &http.Client{
		Timeout: time.Second * time.Duration(timeout),
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		logger.Error("sending HTTP request: %+v", err)
		return errors.Code, fmt.Sprintf(errors.Message, err)
	}
	response, _ := ioutil.ReadAll(resp.Body)
	return resp.StatusCode, string(response)
}

// GetParamsFromRequest get requested params from http.Request on methods GET or OTHER and run Validations
func GetParamsFromRequest(
	req *http.Request,
	dataValues []string,
	resp *Response,
	defaultValues map[string]interface{},
	validators map[string][]ValidatorFunc,
) (map[string]interface{}, bool) {

	jsonValues := make(map[string]interface{})

	if req.Method != "GET" {
		rawBody, _ := ioutil.ReadAll(req.Body)
		req.Body = ioutil.NopCloser(bytes.NewBuffer(rawBody))
		logger.Debug("%s params were:%+v", req.Method, req.Body)
		err := json.Unmarshal(rawBody, &jsonValues)
		if err != nil {
			resp.Code = CodeMalformedRequest
			resp.Body = GenericError{ErrorMessage: MsgMalformedRequest}
			return jsonValues, true
		}
	}

	logger.Debug("GET params were:%+v", req.URL.Query())
	for _, value := range dataValues {
		param := req.URL.Query().Get(value)
		if param != "" {
			if v, err := strconv.Atoi(param); err == nil {
				jsonValues[value] = v
				logger.Debug("%s %T were:%d", value, v, v)
			} else {
				jsonValues[value] = param
				logger.Debug("%s %T were:%s", value, param, param)
			}
		} else {
			if defaultValues[value] != nil {
				jsonValues[value] = defaultValues[value]
			}
		}
	}
	for _, value := range dataValues {
		if jsonValues[value] == nil {
			resp.Code = CodeMissingParam
			resp.Body = GenericError{ErrorMessage: fmt.Sprintf(MsgMissingParam, value)}
			return jsonValues, true
		}
		if vals, ok := validators[value]; ok {
			for _, val := range vals {
				if vok, verrmsg := val(jsonValues[value]); !vok {
					resp.Code = CodeBadParam
					resp.Body = GenericError{ErrorMessage: fmt.Sprintf(MsgBadParam, value+" "+verrmsg)}
					return jsonValues, true
				}
			}
		}
	}
	return jsonValues, false
}
