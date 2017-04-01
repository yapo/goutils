package goutils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Yapo/goutils/commonErrors"
	"github.com/Yapo/logger"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

// Response is a struct to generate a response from POST/PUT requests
type Response struct {
	Code int
	Body interface{}
}

// ErrorType allow have common error on diferent projects
type ErrorType struct {
	Code    int
	Message string
}

// WriteJSONResponse write to te response stream
func WriteJSONResponse(w http.ResponseWriter, response *Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.Code)
	fmt.Fprintf(w, "%s", response.Body)
}

// CreateJSON convert Body to json format
func CreateJSON(response *Response) {
	jsonResponse, err := json.Marshal(response.Body)
	if err != nil {
		logger.Info("CAN'T ENCODE \"%+v\" TO JSON", response.Body)
		response.Body = ""
		response.Code = http.StatusInternalServerError
		return
	}
	response.Body = jsonResponse
}

// SendHTTPRequest does that
func SendHTTPRequest(url, endpoint, method, query, body string, headers map[string]string, errors ErrorType) (int, string) {
	req, _ := http.NewRequest(method, url+endpoint+query, strings.NewReader(body))
	for key, value := range headers {
		req.Header.Add(key, value)
	}
	logger.Debug("sending HTTP request: %+v", req)
	resp, err := http.DefaultClient.Do(req)
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

	if req.Method == "GET" {
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
	} else {
		rawBody, _ := ioutil.ReadAll(req.Body)
		req.Body = ioutil.NopCloser(bytes.NewBuffer(rawBody))
		logger.Debug("%s params were:%+v", req.Method, req.Body)
		err := json.Unmarshal(rawBody, &jsonValues)
		if err != nil {
			resp.Code = commonErrors.CodeMalformedRequest
			resp.Body = commonErrors.GenericError{ErrorMessage: commonErrors.MsgMalformedRequest}
			return jsonValues, true
		}
	}
	for _, value := range dataValues {
		if jsonValues[value] == nil {
			resp.Code = commonErrors.CodeMissingParam
			resp.Body = commonErrors.GenericError{ErrorMessage: fmt.Sprintf(commonErrors.MsgMissingParam, value)}
			return jsonValues, true
		}
		if vals, ok := validators[value]; ok {
			for _, val := range vals {
				if vok, verrmsg := val(jsonValues[value]); !vok {
					resp.Code = commonErrors.CodeBadParam
					resp.Body = commonErrors.GenericError{ErrorMessage: fmt.Sprintf(commonErrors.MsgBadParam, value+" "+verrmsg)}
					return jsonValues, true
				}
			}
		}
	}
	return jsonValues, false
}
