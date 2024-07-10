package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ppreis1976/multithreading/model"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

func main() {
	cep := "03373100" // Substitua por qualquer CEP que você queira consultar

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(2)

	resultChan := make(chan string, 2)
	errorChan := make(chan error, 2)

	go func() {
		defer wg.Done()
		brasilApi, err := BrasilApi(ctx, cep)
		if err != nil {
			errorChan <- err
			return
		}
		resultChan <- fmt.Sprintf("-=- BRASIL API -=- \nCEP: %s\nEstado: %s\nCidade: %s\nBairro: %s\nRua: %s\nServiço: %s\n",
			brasilApi.Cep, brasilApi.State, brasilApi.City, brasilApi.Neighborhood, brasilApi.Street, brasilApi.Service)
	}()

	go func() {
		defer wg.Done()
		viaCep, err := ViaCep(ctx, cep)
		if err != nil {
			errorChan <- err
			return
		}
		resultChan <- fmt.Sprintf("-=- VIA CEP -=- \nCEP: %s\nLogradouro: %s\nComplemento: %s\nUnidade: %s\nBairro: %s\nLocalidade: %s\nUF: %s\nIBGE: %s\nGIA: %s\nDDD: %s\nSIAFI: %s\n",
			viaCep.Cep, viaCep.Logradouro, viaCep.Complemento, viaCep.Unidade, viaCep.Bairro, viaCep.Localidade, viaCep.Uf, viaCep.Ibge, viaCep.Gia, viaCep.Ddd, viaCep.Siafi)
	}()

	go func() {
		wg.Wait()
		close(resultChan)
		close(errorChan)
	}()

	select {
	case result := <-resultChan:
		fmt.Println(result)
	case err := <-errorChan:
		fmt.Println("Erro ao consultar o CEP:", err)
	case <-ctx.Done():
		fmt.Println("Timeout ao consultar o CEP")
	}
}

func BrasilApi(ctx context.Context, cep string) (model.BrasilApi, error) {
	var brasilApi model.BrasilApi

	url := fmt.Sprintf("https://brasilapi.com.br/api/cep/v1/%s", cep)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return brasilApi, err
	}

	client := http.Client{
		Timeout: time.Second,
	}

	resp, err := client.Do(req.WithContext(ctx))
	if err != nil {
		return brasilApi, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return brasilApi, err
	}

	err = json.Unmarshal(body, &brasilApi)
	if err != nil {
		return brasilApi, err
	}

	return brasilApi, nil
}

func ViaCep(ctx context.Context, cep string) (model.ViaCep, error) {
	var viaCep model.ViaCep

	url := fmt.Sprintf("http://viacep.com.br/ws/%s/json/", cep)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return viaCep, err
	}

	client := http.Client{
		Timeout: time.Second,
	}

	resp, err := client.Do(req.WithContext(ctx))
	if err != nil {
		return viaCep, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return viaCep, err
	}

	err = json.Unmarshal(body, &viaCep)
	if err != nil {
		return viaCep, err
	}

	return viaCep, nil
}
