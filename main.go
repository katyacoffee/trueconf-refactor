package main

import "trueconf-refactor/service"

func main() {
	s := service.NewService()
	err := s.Init()
	if err != nil {
		return //TODO
	}
}
