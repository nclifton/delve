### Account Type

```
type Account struct {
	ID        OID       `bson:"_id" json:"_id"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`

	Deleted            bool   `json:"deleted" bson:"deleted"`
	Country            string `json:"country" bson:"country"`

	AccountDetails `bson:",inline"`

	APIKeys  []AccountAPIKey `json:"api_keys" bson:"api_keys"`

	FeatureFlags []FeatureFlag `json:"feature_flags" bson:"feature_flags"`
}

type FeatureFlag struct {
	Domain   string   `json:"domain" bson:"domain"`
	Features []string `json:"features" bson:"features"`
}

type AccountDetails struct {
	Name           string `bson:"name" json:"name"`
	ShortURLDomain string `json:"short_url_domain" bson:"short_url_domain"`
}

type AccountAPIKey struct {
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
	KeyName string `json:"key_name" bson:"key_name"`
	Key               string `json:"key" bson:"key"`
}
```

### RPC account.FindByAPIKey

Request

```
type FindByAPIKeyParams struct {
	Key string `valid:"required"`
}

```

Reply

```
type FindByAPIKeyReply struct {
	Account *Account
}
```