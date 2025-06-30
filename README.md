# Pakay

This library allows you to declare your secrets upfront and provides functions to list and retrieve them achieving a comprehensive security.

This is specially useful when you need to provide different sources for secrets depending on different environments. For example, locally you could use a password manager to retrieve secrets vs on CI or production where you use environment variables that are stored safely.

The main goal is to raise visibility to secrets:

1. The developer can declare the secrets and how to retrieve them from different sources
2. The user will know what secrets they need and how to pass them.
3. Machines will be able to check whether all secrets required are passed or not.

## Getting started

Declare your secrets' manifest:

```yaml
---
# secrets.yaml
- name: my_api_account
  description: The account to connect to the API
  sources:
  - type: env
    env: 
      key: MY_API_ACCOUNT
  - type: 1password
    1password: 
      ref: op://MY_APP_VAULT/my_api/username
- name: my_api_token
  sources:
  - type: env
    env: 
      key: MY_API_TOKEN
  - type: 1password
    1password:
      ref: op://MY_APP_VAULT/my_api/password
```

After that, retrieving the secret is as easy as:

```go
//go:embed secrets.yaml
var secretsConfig string

if err := pakay.LoadSecretsConfig([]byte(secretsConfig)); err != nil {
    return fmt.Errorf("loading secrets config: %w", err)
}

//...

token, found := pakay.GetSecret(ctx, "my_api_token")
```

Notice that order matters so the secrets are retrieved from sources in the same
order they are declared in the config.

You can see [more examples here](./examples).
