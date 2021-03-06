package restapi

import (
  "log"
  "testing"
  "net/http"
  "time"
)

var api_client_server *http.Server

func TestAPIClient(t *testing.T) {
  debug := false

  if debug { log.Println("client_test.go: Starting HTTP server") }
  setup_api_client_server()

  /* Notice the intentional trailing / */
  client := NewAPIClient ("http://127.0.0.1:8080/", false, "", "", make(map[string]string, 0), 2, "id", make([]string, 0), false, false, debug)

  var res string
  var err error

  log.Printf("api_client_test.go: Testing standard OK request\n")
  res, err = client.send_request("GET", "/ok", "")
  if err != nil { t.Fatalf("client_test.go: %s", err) }
  if res != "It works!" {
    t.Fatalf("client_test.go: Got back '%s' but expected 'It works!'\n", res)
  }

  /* Verify timeout works */
  log.Printf("api_client_test.go: Testing timeout aborts requests\n")
  res, err = client.send_request("GET", "/slow", "")
  if err == nil { t.Fatalf("client_test.go: Timeout did not trigger on slow request") }

  if debug { log.Println("client_test.go: Stopping HTTP server") }
  shutdown_api_client_server()
  if debug { log.Println("client_test.go: Done") }
}

func setup_api_client_server () {
  serverMux := http.NewServeMux()
  serverMux.HandleFunc("/ok", func(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("It works!"))
  })
  serverMux.HandleFunc("/slow", func(w http.ResponseWriter, r *http.Request) {
    time.Sleep(9999 * time.Second)
    w.Write([]byte("This will never return!!!!!"))
  })


  api_client_server = &http.Server{
    Addr: "127.0.0.1:8080",
    Handler: serverMux,
  }
  go api_client_server.ListenAndServe()
  /* let the server start */
  time.Sleep(1 * time.Second)
}

func shutdown_api_client_server () {
  api_client_server.Close()
}
