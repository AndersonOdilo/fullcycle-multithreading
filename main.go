package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

func main() {

	cepConsulta := "85603260"

	chViaCep := make(chan ViaCepResponse)

	chBrasilApi := make(chan BrasilApiResponse)

	go consultaCepViaCep(cepConsulta, chViaCep)

	go consultaCepBrasilApi(cepConsulta, chBrasilApi)

	select {
	case retorno := <-chViaCep:
		fmt.Println("Dados retornado da API ViaCep")
		fmt.Printf("Cep: %s, Logradouro: %s, Complemento: %s, Bairro: %s, Cidade: %s, Estado: %s ", retorno.Cep, retorno.Logradouro, retorno.Complemento, retorno.Bairro, retorno.Cidade, retorno.Estado)
	case retorno := <-chBrasilApi:
		fmt.Println("Dados retornado da API BrasilApi")
		fmt.Printf("Cep: %s, Logradouro: %s, Bairro: %s, Cidade: %s, Estado: %s ", retorno.Cep, retorno.Logradouro, retorno.Bairro, retorno.Cidade, retorno.Estado)
	case <-time.After(time.Second):
		fmt.Print("timeout")
	}
}

type ViaCepResponse struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Cidade      string `json:"localidade"`
	Estado      string `json:"estado"`
}

type BrasilApiResponse struct {
	Cep        string `json:"cep"`
	Logradouro string `json:"street"`
	Bairro     string `json:"neighborhood"`
	Cidade     string `json:"city"`
	Estado     string `json:"state"`
}

func consultaCepViaCep(cep string, ch chan<- ViaCepResponse) {

	res, err := http.Get("http://viacep.com.br/ws/" + cep + "/json/")

	if err != nil {
		fmt.Println(err)
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	if err != nil {
		fmt.Println(err)
	}

	viaCepResponse := ViaCepResponse{}
	err = json.Unmarshal(body, &viaCepResponse)

	if err != nil {
		fmt.Println(err)
	}

	ch <- viaCepResponse

}

func consultaCepBrasilApi(cep string, ch chan<- BrasilApiResponse) {

	res, err := http.Get("https://brasilapi.com.br/api/cep/v1/" + cep)

	if err != nil {
		fmt.Println(err)
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	if err != nil {
		fmt.Println(err)
	}

	brasilApiResponse := BrasilApiResponse{}
	err = json.Unmarshal(body, &brasilApiResponse)

	if err != nil {
		fmt.Println(err)
	}

	ch <- brasilApiResponse

}
