package lib

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
)

// CaptureResponseBody keeps a copy of the response body in body and gives the SDK an identical reader to decode. The Go
// SDK decodes some fields into strings with omitempty, which cannot tell an empty string from null.
func CaptureResponseBody(body *[]byte) files_sdk.RequestResponseOption {
	return files_sdk.ResponseOption(func(response *http.Response) error {
		if response.Body == nil {
			return nil
		}
		data, err := io.ReadAll(response.Body)
		if closeErr := response.Body.Close(); err == nil {
			err = closeErr
		}
		if err != nil {
			// The original body is closed, so later cleanup must not close it again.
			response.Body = http.NoBody
			return err
		}
		*body = data
		response.Body = io.NopCloser(bytes.NewReader(data))
		return nil
	})
}

// ResponseProperty returns the named top-level property of the JSON object in body exactly as the API sent it. A
// missing property is nil, like a null one. A body that is not a JSON object is an error.
func ResponseProperty(attributePath path.Path, body []byte, name string) (any, diag.Diagnostics) {
	var object map[string]json.RawMessage
	if err := json.Unmarshal(body, &object); err != nil || object == nil {
		detail := "The response body is not a JSON object."
		if err != nil {
			detail = "Could not decode the response body: " + err.Error()
		}
		return nil, diag.Diagnostics{diag.NewAttributeErrorDiagnostic(attributePath, "Invalid API response", detail)}
	}
	var value any
	if raw, ok := object[name]; ok {
		if err := json.Unmarshal(raw, &value); err != nil {
			return nil, diag.Diagnostics{diag.NewAttributeErrorDiagnostic(attributePath, "Invalid API response", "Could not decode "+name+": "+err.Error())}
		}
	}
	return value, nil
}
