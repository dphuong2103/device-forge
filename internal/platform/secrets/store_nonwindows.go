//go:build !windows

package secrets

import "context"

type Store struct{}

func New(_ string) (*Store, error)                           { return nil, ErrUnsupported }
func (s *Store) Put(context.Context, string, []byte) error   { return ErrUnsupported }
func (s *Store) Get(context.Context, string) ([]byte, error) { return nil, ErrUnsupported }
func (s *Store) Delete(context.Context, string) error        { return ErrUnsupported }
