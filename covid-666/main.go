package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type Patient struct {
	Id           string `json:"id"`
	City         string `json:"city"`
	State        string `json:"state"`
	Country      string `json:"country"`
	Status       string `json:"status"`
	InfectedDate int    `json:"infected_date"`
	UpdatedAt    int    `json:"updated_at"`
}

func main() {
	allPatients := map[string]Patient{}
	patientCounter := 0
	http.HandleFunc("/ping", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(200)
		fmt.Fprintf(writer, "pong")
	})

	http.HandleFunc("/add", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		if request.Method != http.MethodPost {
			writer.WriteHeader(http.StatusBadRequest)
			return
		}

		body, err := ioutil.ReadAll(request.Body)
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			return
		}

		p := Patient{}
		json.Unmarshal(body, &p)

		patientCounter += 1
		p.Id = fmt.Sprintf("%d", patientCounter)

		allPatients[p.Id] = p

		s, _ := json.Marshal(p)
		fmt.Fprintf(writer, string(s))
	})

	http.HandleFunc("/get/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		splits := strings.Split(request.URL.Path, "/")
		id := splits[len(splits)-1]

		p, ok := allPatients[id]
		if !ok {
			writer.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(writer, "{}")
			return
		}

		s, _ := json.Marshal(p)
		fmt.Fprintf(writer, string(s))
	})

	http.HandleFunc("/update/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		splits := strings.Split(request.URL.Path, "/")
		id := splits[len(splits)-1]

		p, ok := allPatients[id]
		if !ok {
			writer.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(writer, "{}")
			return
		}

		requestBody, err := ioutil.ReadAll(request.Body)
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			return
		}

		status := struct {
			Status string `json:"status"`
		}{}

		json.Unmarshal(requestBody, &status)
		p.Status = status.Status

		allPatients[p.Id] = p
		s, _ := json.Marshal(p)
		fmt.Fprintf(writer, string(s))
	})

	http.HandleFunc("/count", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		country := request.URL.Query().Get("country")
		state := request.URL.Query().Get("state")
		city := request.URL.Query().Get("city")
		status := request.URL.Query().Get("status")
		count := 0
		for _, p := range allPatients {
			flag := true

			if len(country) > 0 {
				flag = p.Country == country
			}
			if len(state) > 0 {
				flag = p.State == state
			}
			if len(city) > 0 {
				flag = p.City == city
			}
			if len(status) > 0 {
				flag = p.Status == status
			}

			if flag {
				count++
			}
		}

		res := struct {
			Count int `json:"count"`
		}{
			Count: count,
		}
		temp, _ := json.Marshal(res)
		fmt.Fprintf(writer, string(temp))

	})

	http.HandleFunc("/reset", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		allPatients = map[string]Patient{}
		fmt.Fprintf(writer, "ok")
	})

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
