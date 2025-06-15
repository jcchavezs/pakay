# Pakay

This library allows you to declare your secrets upfront and how they can be retrieved to achieve comprehensive security.

This is specially useful when you need to provide different sources depending on different environments. For example, locally you could use a password manager to retrieve secrets vs on CI or production where you use environment variables that are stored safely.

## Getting started

Declare a secret manifest:

```yaml
---
- name: my_api_account
  description: The account to connect to the API
  sources:
  - type: env
    env: 
      name: MY_API_ACCOUNT
  - type: 1password
    1password: 
      ref: op://MY_APP_VAULT/my_api/username
- name: my_api_token
  sources:
  - type: env
    env: 
      name: MY_API_TOKEN
  - type: 1password
    1password:
      ref: op://MY_APP_VAULT/my_api/password
```

This, retrieving the secret is as easy as:

```go
if err := pakay.LoadSecretsFromBytes([]byte(config), pakay.LoadOptions{}); err != nil {
    return fmt.Errorf("loading secrets: %w", err)
}

//...

token, found := pakay.GetSecret(ctx, "my_api_token")
```

Notice that order matters so the secrets are retrieved from sources in the same
order they are declared in the config.

You can see [more examples here](./examples).
