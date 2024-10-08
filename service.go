package main

type IService[T any] interface {
	GetService() T
}

type Service struct {
}

func (s *Service) GetService() *Service {
	return s
}
