package response

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type envelope map[string]any

func WriteJson(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")

	js, err := json.Marshal(data)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte(`{"error": "internal server error"}`))
		if err != nil {
			return
		}
		return
	}

	js = append(js, '\n')

	w.WriteHeader(status)
	_, err = w.Write(js)
	if err != nil {
		return
	}
}

func ErrorResponse(w http.ResponseWriter, _ *http.Request, status int, message any) {
	env := envelope{"error": message}
	//logger.Log.Error("error in incoming request", zap.Int("status", status), zap.String("url", r.URL.String()))

	WriteJson(w, status, env)
}
func BadRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	ErrorResponse(w, r, http.StatusBadRequest, err.Error())
}

func ServerErrorResponse(w http.ResponseWriter, r *http.Request, _ error) {
	//logger.Log.Error("error in incoming request", zap.Error(err))

	message := "the server encountered a problem and could not process your request"
	ErrorResponse(w, r, http.StatusInternalServerError, message)
}

func NotFoundResponse(w http.ResponseWriter, _ *http.Request) {
	message := "the required resource could not be found"
	env := envelope{"error": message}

	WriteJson(w, http.StatusNotFound, env)
}

func MethodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("the %s method is not supported for this resource", r.Method)
	env := envelope{"error": message}

	WriteJson(w, http.StatusMethodNotAllowed, env)
}
