package main

import (
	"github.com/RBAC/internal/repository"
	"github.com/RBAC/internal/router"
)

func main() {

	repository.NewMysql()

	router.Start()

}
