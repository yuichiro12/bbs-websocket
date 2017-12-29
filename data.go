package main

import (
	"encoding/json"
	"log"
)

type Data struct {
	Ids     []int  `json:"ids"`
	Message string `json:"message"`
	Icon    string `json:"icon"`
	Url     string `json:"url"`
}

type Notification struct {
	Message string `json:"message"`
	Icon    string `json:"icon"`
	Url     string `json:"url"`
}

func jsonDecode(j []byte) *Data {
	data := new(Data)
	if err := json.Unmarshal(j, data); err != nil {
		log.Fatal("JSON Unmarshal error:", err)
	}
	return data
}

func jsonEncode(d *Data) []byte {
	n := Notification{
		Message: d.Message,
		Icon:    d.Icon,
		Url:     d.Url,
	}
	j, err := json.Marshal(n)
	if err != nil {
		log.Fatal("JSON Marshal error:", err)
	}
	return j
}
