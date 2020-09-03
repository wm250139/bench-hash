package client_inproc

import "bench-hash/hasher"

type InProcessClient struct{}

func NewInProcessClient() *InProcessClient {
	return &InProcessClient{}
}

func (c InProcessClient) Hash(input string) (string, error) {
	return hasher.String(input), nil
}
