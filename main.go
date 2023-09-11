package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type expectedOutputStruct struct {
	// check for status code match
	StatusCode int `compscore:"status_code"`

	// check for substring match in body
	SubstringMatch string `compscore:"substring_match"`

	// check for regex match in body
	RegexMatch string `compscore:"regex_match"`

	// check for exact match in body
	Match string `compscore:"match"`
}

func (e *expectedOutputStruct) Unmarshall(in string) error {
	structLookup := make(map[string]string)
	testSplit := strings.Split(in, ";")
	for _, item := range testSplit {
		itemSplit := strings.Split(item, "=")
		if len(itemSplit) != 2 {
			return fmt.Errorf("invalid parameter string: %s", item)
		}
		structLookup[strings.TrimSpace(itemSplit[0])] = strings.TrimSpace(itemSplit[1])
	}

	statusCodeStr, ok := structLookup["status_code"]
	if ok {
		statusCodeInt, err := strconv.Atoi(statusCodeStr)
		if err != nil {
			return fmt.Errorf("invalid status code: %s", statusCodeStr)
		}

		e.StatusCode = statusCodeInt
	}

	substringMatch, ok := structLookup["substring_match"]
	if ok {
		e.SubstringMatch = substringMatch
	}

	regexMatch, ok := structLookup["regex_match"]
	if ok {
		e.RegexMatch = regexMatch
	}

	match, ok := structLookup["match"]
	if ok {
		e.Match = match
	}

	return nil
}

func (e *expectedOutputStruct) Compare(response *http.Response) error {
	if e.StatusCode != 0 && e.StatusCode != response.StatusCode {
		return fmt.Errorf("status code mismatch: expected \"%d\", got \"%d\"", e.StatusCode, response.StatusCode)
	}

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("encountered error while reading response body: %v", err.Error())
	}

	body := string(bodyBytes)
	fmt.Println(body)

	if e.SubstringMatch != "" && !strings.Contains(body, e.SubstringMatch) {
		return fmt.Errorf("substring match mismatch: expected \"%s\"", e.SubstringMatch)
	}

	if e.Match != "" && e.Match != body {
		return fmt.Errorf("match mismatch: expected \"%s\", got \"%s\"", e.Match, body)
	}

	if e.RegexMatch != "" {
		pattern, err := regexp.Compile(e.RegexMatch)
		if err != nil {
			return fmt.Errorf("invalid regex pattern: \"%s\"", e.RegexMatch)
		}

		if !pattern.MatchString(body) {
			return fmt.Errorf("regex match mismatch: expected \"%s\"", e.RegexMatch)
		}
	}

	return nil
}

func Run(ctx context.Context, target string, command string, expectedOutput string, username string, password string) (bool, string) {
	var requestType string

	switch strings.ToUpper(command) {
	case "GET":
		requestType = http.MethodGet
	case "POST":
		requestType = http.MethodPost
	case "PUT":
		requestType = http.MethodPut
	case "DELETE":
		requestType = http.MethodDelete
	case "PATCH":
		requestType = http.MethodPatch
	case "HEAD":
		requestType = http.MethodHead
	case "OPTIONS":
		requestType = http.MethodOptions
	case "CONNECT":
		requestType = http.MethodConnect
	case "TRACE":
		requestType = http.MethodTrace
	default:
		return false, "provided invalid command/http verb: " + command
	}

	req, err := http.NewRequestWithContext(ctx, requestType, target, nil)
	if err != nil {
		return false, fmt.Sprintf("encounted error while creating request: %v", err.Error())
	}

	if username != "" && password != "" {
		req.SetBasicAuth(username, password)
	} else if password != "" && username == "" {
		req.Header.Add("Authorization", password)
	}

	errChan := make(chan error)

	go func() {
		defer close(errChan)
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			errChan <- fmt.Errorf("encounted error while making request: %v", err.Error())
			return
		}
		defer resp.Body.Close()

		var output expectedOutputStruct
		err = output.Unmarshall(expectedOutput)
		if err != nil {
			errChan <- fmt.Errorf("encounted error while parsing expected output: %v", err.Error())
			return
		}

		err = output.Compare(resp)
		if err != nil {
			errChan <- fmt.Errorf("encounted error while comparing expected output: %v", err.Error())
			return
		}
	}()

	select {
	case <-ctx.Done():
		return false, "Timeout exceeded; err %v" + ctx.Err().Error()
	case err := <-errChan:
		if err != nil {
			return false, fmt.Sprintf("Encountered error: %s", err)
		}

		return true, ""
	}
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	fmt.Println(Run(ctx, "https://pastebin.com/raw/LpnAL5jF", "GET", "status_code=200;regex_match=\\d", "", ""))
}
