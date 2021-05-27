package utils

type PubData struct {
	PubID  int    `json:"pubId"`
	PubStr string `json:"pubStr"`
}

type TokenRequest struct {
	AppID  string `json:"appId"`
	AppKey string `json:"appKey"`
}

type TokenInfo struct {
	Token string `json:"token"`
	To    int    `json:"timeout"`
}

type ReqTokenResponse struct {
	Code    int       `json:"code"`
	Message string    `json:"message"`
	Data    TokenInfo `json:"data"`
}

type ResponseFromBd struct {
	Code       int    `json:"code"`
	Message    string `json:"message"`
	DataDetail struct {
		SignLen   int    `json:"signLen"`
		Data      string `json:"data"`
		DataIndex string `json:"dataIndex"`
		Sign      string `json:"sign"`
		Cert      string `json:"cert"`
		CertLen   int    `json:"certLen"`
	} `json:"data"`
}
