package entity

type AnyAddress struct {
	Api        string
	Cep        string
	Bairro     string
	Localidade string
	UF         string
	Logradouro string
}

type AddressViaCep struct {
	Api        string
	Cep        string `json:"cep"`
	Bairro     string `json:"bairro"`
	Localidade string `json:"localidade"`
	UF         string `json:"uf"`
	Logradouro string `json:"logradouro"`
}

type AddressApiCep struct {
	Api        string
	Cep        string `json:"code"`
	Bairro     string `json:"district"`
	Localidade string `json:"city"`
	UF         string `json:"state"`
	Logradouro string `json:"address"`
}

type Address struct {
	Cep        string `json:"cep"`
	Bairro     string `json:"bairro"`
	Localidade string `json:"localidade"`
	UF         string `json:"uf"`
	Logradouro string `json:"logradouro"`
	Erro       bool   `json:"erro"`
}
