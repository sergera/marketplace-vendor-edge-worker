package service

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/sergera/marketplace-vendor-edge-worker/internal/conf"
	"github.com/sergera/marketplace-vendor-edge-worker/internal/domain"
)

type VendorAPIService struct {
	host        string
	port        string
	contentType string
	client      *http.Client
}

func NewVendorAPIService() *VendorAPIService {
	conf := conf.GetConf()
	return &VendorAPIService{
		conf.VendorAPIHost,
		conf.VendorAPIPort,
		"application/json; charset=UTF-8",
		&http.Client{},
	}
}

func (v VendorAPIService) Post(route string, jsonData []byte) error {
	request, err := http.NewRequest("POST", v.host+":"+v.port+"/"+route, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println("failed to create post request: " + err.Error())
		return err
	}

	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	response, err := v.client.Do(request)
	if err != nil {
		log.Println("failed to perform vendor api post request: " + err.Error())
		return err
	}

	defer response.Body.Close()
	return nil
}

func (v VendorAPIService) Put(route string, jsonData []byte) error {
	request, err := http.NewRequest("PUT", v.host+":"+v.port+"/"+route, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println("failed to create put request: " + err.Error())
		return err
	}

	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	response, err := v.client.Do(request)
	if err != nil {
		log.Println("failed to perform vendor api put request: " + err.Error())
		return err
	}

	defer response.Body.Close()
	return nil
}

func (v VendorAPIService) SendOrder(o domain.OrderModel) error {
	m, err := json.Marshal(o)
	if err != nil {
		log.Println("failed to marshal order model into json")
		return err
	}

	err = v.Post("send-order", m)
	if err != nil {
		return err
	}

	return nil
}
