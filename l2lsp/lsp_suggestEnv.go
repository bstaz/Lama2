// lsp_suggestEnv.go
package l2lsp

import (
	"errors"
	"path/filepath"
	"strings"

	l2envpackege "github.com/HexmosTech/lama2/l2env"
	"github.com/rs/zerolog/log"
)

func getSearchQueryString(request JSONRPCRequest) string {
	if request.Params.SearchQuery != nil {
		return *request.Params.SearchQuery
	}
	return ""
}

func getRequestURI(request JSONRPCRequest) (string, int, error) {
	if request.Params.TextDocument.Uri == nil {
		return "", ErrInvalidURI, errors.New("URI cannot be empty. Ex: 'file:///path/to/workspace/myapi.l2'")
	}
	uri := *request.Params.TextDocument.Uri

	// Handle local files
	if strings.HasPrefix(uri, "file://") {
		return uri[len("file://"):], 0, nil

		// Handle remote files (example for SSH)
	} else if strings.HasPrefix(uri, "ssh://") {
		return uri[len("ssh://"):], 0, nil

		// Handle WSL files
	} else if strings.HasPrefix(uri, "wsl://") {
		return uri[len("wsl://"):], 0, nil

		// Handle Windows files
	} else if strings.Contains(uri, "\\") {
		return "", ErrUnsupportedFeature, errors.New("windows is not supported as of now. To contribute visit here: https://github.com/HexmosTech/Lama2")

	} else {
		// Log the unexpected URI scheme
		log.Warn().Str("URI", uri).Msg("Encountered unexpected URI scheme.")
		return "", ErrUnexpectedURIScheme, errors.New("encountered unexpected URI scheme. Ex: 'file:///path/to/workspace/myapi.l2'")
	}
}

func SuggestEnvironmentVariables(request JSONRPCRequest) JSONRPCResponse {
	/*
		{
			"jsonrpc": "2.0",
			"id": 2,
			"method": "suggest/environmentVariables",
			"params": {
				"textDocument": {
					"uri": "file:///home/lovestaco/repos/Lama2/elfparser/ElfTestSuite/root_variable_override/api/y_0020_root_override.l2"
				},
				"position": {
					"line": 1,
					"character": 2
				},
				"relevantSearchString": ""
			}
		}
	*/

	log.Info().Msg("L2 LSP environment variables suggestion requested")
	log.Info().Str("Method", request.Method).Interface("Params", request.Params)

	relevantSearchString := getSearchQueryString(request)
	uri, errorCode, err := getRequestURI(request)
	if err != nil {
		return ErrorResp(request, errorCode, err.Error())
	}
	parentFolder := filepath.Dir(uri)
	res := l2envpackege.ProcessEnvironmentVariables(relevantSearchString, parentFolder)
	return createEnvironmentVariablesResp(request, res)
}
