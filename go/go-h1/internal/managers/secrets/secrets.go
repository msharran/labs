package secrets

import (
	"context"
	"log/slog"
	"sync"
)

type SecretsManager struct {
	secrets map[string]string
	ctx     context.Context
	log     *slog.Logger
	actionC chan func()
	wg      sync.WaitGroup
}

func NewManager(ctx context.Context, log *slog.Logger) *SecretsManager {
	s := &SecretsManager{
		secrets: make(map[string]string),
		ctx:     ctx,
		log:     log.WithGroup("secrets"),
		actionC: make(chan func()),
	}

	return s
}

func (s *SecretsManager) Run() {
	s.log.Info("starting secrets manager")
	defer s.log.Info("secrets manager stopped")

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		s.reconcileSecrets()
	}()

	s.wg.Wait()
}

func (s *SecretsManager) Set(key, value string) error {
	err := make(chan error)

	s.actionC <- func() {
		s.secrets[key] = value
		err <- nil
	}

	return <-err
}

func (s *SecretsManager) Keys() []string {
	keys := make(chan []string)

	s.actionC <- func() {
		var k []string
		for key := range s.secrets {
			k = append(k, key)
		}
		keys <- k
	}

	return <-keys
}

func (s *SecretsManager) Get(key string) (string, error) {
	result := make(chan struct {
		value string
		err   error
	})

	s.actionC <- func() {
		result <- struct {
			value string
			err   error
		}{s.secrets[key], nil}
	}

	res := <-result
	return res.value, res.err
}

func (s *SecretsManager) reconcileSecrets() {
	s.log.Info("reconciling secrets")
	defer s.log.Info("stopped reconciling secrets")

	for {
		select {
		case <-s.ctx.Done():
			s.log.Info("context done, exiting", "err", s.ctx.Err())
			return
		case action := <-s.actionC:
			action()
		}
	}
}
