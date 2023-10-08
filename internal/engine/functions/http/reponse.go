package http

import (
	"encoding/json"
	"net/http"
	"serverless/internal/engine/utils"

	"rogchap.com/v8go"
)

type Header map[string]interface{}

type Response struct {
	Status  int
	Content string
	Headers Header
}

func ResponseFunction(isolate *v8go.Isolate, channelChannel chan Response) *v8go.FunctionTemplate {
	return v8go.NewFunctionTemplate(isolate, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
		args := info.Args()

		var headers Header
		var status = 200
		var content string

		if len(args) < 1 {
			err := utils.ThrowException("Response content is required", isolate)
			if err != nil {
				return nil
			}
			return nil
		}
		if len(args) >= 1 {
			content = args[0].String()
		}

		if len(args) >= 2 {
			status = int(args[1].Integer())
			if http.StatusText(status) == "" {
				err := utils.ThrowException("Invalid Http Status code", isolate)
				if err != nil {
					return nil
				}
				return nil
			}
		}

		if len(args) >= 3 {
			headersJson, err := args[2].MarshalJSON()
			if err != nil {
				return nil
			}

			err = json.Unmarshal(headersJson, &headers)
			if err != nil {
				return nil
			}
		}
		channelChannel <- Response{
			Headers: headers,
			Status:  status,
			Content: content,
		}

		return nil
	})
}
