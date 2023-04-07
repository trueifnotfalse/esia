# esia

```Golang
Esia openId

signer := esia.NewCliSigner(&esia.CliSignerConfig{
    CertPath:       "/path/to/esia_auth.pem",
    PrivateKeyPath: "/path/to/esia_auth.key",
})
openId := esia.NewOpenId(&esia.OpenIdConfig{
    MnemonicsSystem: "000000",
    RedirectUrl:     "https://your-site/esia/callback",
    PortalUrl:       "https://esia-portal1.test.gosuslugi.ru/",
    Scope:           "fullname id_doc",
}, signer)

//Get auth url
url, _ := esiaOpenId.GetUrl()

//Get persone info
var person esia.EsiaPerson
esiaOpenId.GetInfoByPath("", &person)

//Get docs info
var docs esia.EsiaDocs
esiaOpenId.GetInfoByPath("/docs/" + fmt.Sprint(person.RIdDoc), &docs)
```
