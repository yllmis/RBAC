package main

import (
	"github.com/RBAC/internal/repository"
	"github.com/RBAC/internal/router"
	"github.com/RBAC/pkg/log"
)

func main() {
	log.Init()
	repository.NewMysql()
	repository.NewRedis()

	router.Start()

}
