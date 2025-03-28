// формат запроса https://yandex.ru/dev/dialogs/alice/doc/ru/request
package main

import (
	"net/http"
)

// go build -o skill
// ./skill -a :8081
// RUN_ADDR=:8081 ./skill
// RUN_ADDR=:8082 ./skill -a :8081
func main() {
	parseFlags()

	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	return http.ListenAndServe(`:8080`, http.HandlerFunc(webhook))
}

// to run this service:
// curl -X POST http://localhost:8080 -H "Content-Type: application/json" -d '{}'
func webhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(`
      {
        "response": {
          "text": "Извините, я пока ничего не умею"
        },
        "version": "1.0"
      }
    `))
}
