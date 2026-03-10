package curlbuilder

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

// BuildCurlCommand converts an http.Request into a curl command.
func BuildCurlCommand(req *http.Request) (string, error) {
	var cmd strings.Builder

	cmd.WriteString("curl ")
	cmd.WriteString("-X " + req.Method + " ")
	cmd.WriteString(req.URL.String() + " ")

	for k, v := range req.Header {
		for _, vv := range v {
			cmd.WriteString("-H '" + k + ": " + vv + "' ")
		}
	}

	if req.Body != nil && req.GetBody != nil {
		body, err := req.GetBody()
		if err != nil {
			return "", fmt.Errorf("failed to get request body: %w", err)
		}
		defer body.Close()

		bodyBytes, err := io.ReadAll(body)
		if err != nil {
			return "", fmt.Errorf("failed to read request body: %w", err)
		}

		bodyStr := string(bodyBytes)
		bodyStr = strings.ReplaceAll(bodyStr, "'", "'\"'\"'")
		cmd.WriteString("-d '" + bodyStr + "' ")
	}

	return cmd.String(), nil
}
