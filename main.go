package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	plist "github.com/DHowett/go-plist"
	"github.com/bitrise-io/go-utils/command"
	"github.com/gorilla/mux"
	"github.com/trapacska/certificate-info/pkcs"
)

func getCertsJSON(p12 []byte) (string, error) {
	certs, err := pkcs.DecodeAllCerts(p12, "")
	if err != nil {
		return "", err
	}

	b, err := json.Marshal(certs)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func getProfileJSON(profile []byte) (string, error) {
	cmd := command.New("openssl", "smime", "-inform", "der", "-verify")
	cmd.SetStdin(strings.NewReader(string(profile)))

	var b bytes.Buffer
	cmd.SetStdout(&b)

	err := cmd.Run()
	if err != nil {
		return "", err
	}

	var intf interface{}
	dec := plist.NewDecoder(bytes.NewReader(b.Bytes()))

	err = dec.Decode(&intf)
	if err != nil {
		return "", err
	}

	str, err := json.Marshal(intf)
	if err != nil {
		return "", err
	}

	return string(str), nil
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", index).Methods("GET")
	router.HandleFunc("/certificate/url", certFromURL).Methods("POST")
	router.HandleFunc("/certificate", certFromContent).Methods("POST")
	router.HandleFunc("/profile", profFromContent).Methods("POST")
	router.HandleFunc("/profile", profFromURL).Methods("POST")

	if err := http.ListenAndServe(":"+os.Getenv("PORT"), router); err != nil {
		fmt.Printf("Failed to listen, error: %s\n", err)
	}
}

func certFromContent(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write([]byte(fmt.Sprintf(`{"error":"Failed to read body: %s"}`, err)))
		if err != nil {
			fmt.Printf("Failed to write response, error: %s\n", err)
			return
		}
		return
	}

	certsJSON, err := getCertsJSON(body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write([]byte(`{"error":"Failed to get certificate info"}`))
		if err != nil {
			fmt.Printf("Failed to write response, error: %s\n", err)
			return
		}
		return
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(certsJSON))
	if err != nil {
		fmt.Printf("Failed to write response, error: %s\n", err)
		return
	}
}

func profFromContent(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write([]byte(fmt.Sprintf(`{"error":"Failed to read body: %s"}`, err)))
		if err != nil {
			fmt.Printf("Failed to write response, error: %s\n", err)
		}
		return
	}

	profJSON, err := getProfileJSON(body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write([]byte(fmt.Sprintf(`{"error":"Failed to get profile info: %s"}`, err)))
		if err != nil {
			fmt.Printf("Failed to write response, error: %s\n", err)
		}
		fmt.Printf("Failed to get profile info, error: %s\n", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(profJSON))
	if err != nil {
		fmt.Printf("Failed to write response, error: %s\n", err)
		return
	}
}

func profFromURL(w http.ResponseWriter, r *http.Request) {
	url, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write([]byte(fmt.Sprintf(`{"error":"Failed to read body: %s"}`, err)))
		if err != nil {
			fmt.Printf("Failed to write response, error: %s\n", err)
			return
		}
		return
	}

	response, err := http.Get(string(url))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write([]byte(`{"error":"Failed to create request for the given URL"}`))
		if err != nil {
			fmt.Printf("Failed to write response, error: %s\n", err)
			return
		}
		return
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write([]byte(fmt.Sprintf(`{"error":"Failed to read body: %s"}`, err)))
		if err != nil {
			fmt.Printf("Failed to write response, error: %s\n", err)
		}
		return
	}

	profJSON, err := getProfileJSON(body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write([]byte(fmt.Sprintf(`{"error":"Failed to get profile info: %s"}`, err)))
		if err != nil {
			fmt.Printf("Failed to write response, error: %s\n", err)
		}
		fmt.Printf("Failed to get profile info, error: %s\n", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(profJSON))
	if err != nil {
		fmt.Printf("Failed to write response, error: %s\n", err)
		return
	}
}

func certFromURL(w http.ResponseWriter, r *http.Request) {
	url, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write([]byte(fmt.Sprintf(`{"error":"Failed to read body: %s"}`, err)))
		if err != nil {
			fmt.Printf("Failed to write response, error: %s\n", err)
			return
		}
		return
	}

	response, err := http.Get(string(url))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write([]byte(`{"error":"Failed to create request for the given URL"}`))
		if err != nil {
			fmt.Printf("Failed to write response, error: %s\n", err)
			return
		}
		return
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write([]byte(fmt.Sprintf(`{"error":"Failed to get data from the given url: %s"}`, err)))
		if err != nil {
			fmt.Printf("Failed to write response, error: %s\n", err)
			return
		}
		return
	}

	certsJSON, err := getCertsJSON(body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write([]byte(`{"error":"Failed to get certificate info"}`))
		if err != nil {
			fmt.Printf("Failed to write response, error: %s\n", err)
			return
		}
		return
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(certsJSON))
	if err != nil {
		fmt.Printf("Failed to write response, error: %s\n", err)
		return
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome!")
}