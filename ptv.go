package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

const baseUrl = "https://timetableapi.ptv.vic.gov.au"
const apiVersion = "v3"

var devID = os.Getenv("PTV_DEV_ID")
var apiKey = os.Getenv("PTV_API_KEY")

func ptvRequest[T any](api string) (T, bool) {
	apiCall := fmt.Sprintf("/%s%s", apiVersion, api)
	if strings.ContainsRune(apiCall, '?') {
		apiCall += "&"
	} else {
		apiCall += "?"
	}
	apiCall += fmt.Sprintf("devid=%s", devID)

	mac := hmac.New(sha1.New, []byte(apiKey))
	mac.Write([]byte(apiCall))
	signature := hex.EncodeToString(mac.Sum(nil))

	finalUrl := fmt.Sprintf("%s%s&signature=%s", baseUrl, apiCall, signature)

	resp, err := http.Get(finalUrl)
	if err != nil {
		return *new(T), false
	}

	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return *new(T), false
	}

	var jsonData T
	err = json.Unmarshal(body, &jsonData)
	if err != nil {
		return *new(T), false
	}

	return jsonData, true
}
