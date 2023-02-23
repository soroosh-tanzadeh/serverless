package utils

import "rogchap.com/v8go"

func ThrowException(message string, isolate *v8go.Isolate) error {
	value, err := v8go.NewValue(isolate, message)
	if err != nil {
		return err
	}
	isolate.ThrowException(value)
	return nil
}
