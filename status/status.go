package status

import (
	"context"

	"github.com/jcchavezs/pakay"
	"github.com/jcchavezs/pakay/internal/secrets"
)

type SecretStatus struct {
	filterIn pakay.FilterIn
	secret   secrets.Secret
}

func (ss SecretStatus) Name() string {
	return ss.secret.Name
}

func (ss SecretStatus) Description() string {
	return ss.secret.Description
}

func (ss SecretStatus) Sources() []string {
	sources := make([]string, 0, len(ss.secret.Sources))
	for _, s := range ss.secret.Sources {
		if ss.filterIn != nil {
			if !ss.filterIn(pakay.Source{Type: s.Type, Labels: s.Labels}) {
				continue
			}
		}

		sources = append(sources, s.String())
	}

	return sources
}

func (ss SecretStatus) GetValue(ctx context.Context) (string, bool) {
	return pakay.GetSecretWithOptions(ctx, ss.secret.Name, pakay.SecretOptions{
		FilterIn: ss.filterIn,
	})
}

type CheckOptions struct {
	FilterIn pakay.FilterIn
}

func CheckSecrets(ctx context.Context) []SecretStatus {
	return CheckSecretsWithOptions(ctx, CheckOptions{})
}

func CheckSecretsWithOptions(ctx context.Context, opts CheckOptions) []SecretStatus {
	ss := make([]SecretStatus, 0, len(secrets.All))

	for _, s := range secrets.All {
		ss = append(ss, SecretStatus{
			filterIn: opts.FilterIn,
			secret:   s,
		})
	}

	return ss
}
