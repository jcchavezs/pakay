package view

import (
	"context"

	"github.com/jcchavezs/pakay"
	"github.com/jcchavezs/pakay/internal/secrets"
)

type Secret interface {
	Name() string
	Description() string
	Sources() []string
	GetValue(ctx context.Context) (string, bool)
}

type secret struct {
	filterIn pakay.FilterIn[pakay.Source]
	secrets.Secret
}

var _ Secret = secret{}

func (ss secret) Name() string {
	return ss.Secret.Name
}

func (ss secret) Description() string {
	return ss.Secret.Description
}

func (ss secret) Sources() []string {
	sources := make([]string, 0, len(ss.Secret.Sources))
	for _, s := range ss.Secret.Sources {
		if ss.filterIn != nil {
			if !ss.filterIn(pakay.Source{Type: s.Type, Labels: s.Labels}) {
				continue
			}
		}

		sources = append(sources, s.String())
	}

	return sources
}

func (ss secret) GetValue(ctx context.Context) (string, bool) {
	return pakay.GetSecretWithOptions(ctx, ss.Secret.Name, pakay.SecretOptions{
		FilterIn: ss.filterIn,
	})
}

type ListOptions struct {
	FilterIn pakay.FilterIn[pakay.Source]
}

// ListSecrets returns the status of all secrets
func ListSecrets(ctx context.Context) []Secret {
	return ListSecretsWithOptions(ctx, ListOptions{})
}

// ListSecretsWithOptions returns the status of secrets, applying the provided filter if any.
func ListSecretsWithOptions(ctx context.Context, opts ListOptions) []Secret {
	ss := make([]Secret, 0, len(secrets.All))

	for _, s := range secrets.All {
		ss = append(ss, secret{
			filterIn: opts.FilterIn,
			Secret:   s,
		})
	}

	return ss
}
