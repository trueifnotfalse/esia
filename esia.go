package esia

type AuthCode struct {
	Nbf          string `json:"-"`
	Scope        string `json:"scope"`
	Iss          string `json:"iss"`
	UrnEsiaSid   string `json:"urn:esia:sid"`
	UrnEsiaSbjId int32  `json:"urn:esia:sbj_id"`
	Iat          int32  `json:"iat"`
	ClientId     string `json:"client_id"`
	Exp          int32  `json:"exp"`
}

type Token struct {
	AccessToken  string   `json:"access_token"`
	AuthCode     AuthCode `json:"-"`
	RefreshToken string   `json:"refresh_token"`
	State        string   `json:"state"`
	TokenType    string   `json:"token_type"`
	ExpiresIn    int32    `json:"expires_in"`
}

type Person struct {
	FirstName         string `json:"firstName"`
	LastName          string `json:"lastName"`
	MiddleName        string `json:"middleName"`
	Trusted           bool   `json:"trusted"`
	Citizenship       string `json:"citizenship"`
	Status            string `json:"status"`
	Verifying         bool   `json:"verifying"`
	RIdDoc            int32  `json:"rIdDoc"`
	ContainsUpCfmCode bool   `json:"containsUpCfmCode"`
	ETag              string `json:"eTag"`
}

type Docs struct {
	Id        int32  `json:"id"`
	Type      string `json:"type"`
	VrfStu    string `json:"vrfStu"`
	Series    string `json:"series"`
	Number    string `json:"number"`
	IssueDate string `json:"issueDate"`
	IssueId   string `json:"issueId"`
	IssuedBy  string `json:"issuedBy"`
	Etag      string `json:"eTag"`
}
