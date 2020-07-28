package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/ont-bizsuite/addonstep/pkg/meta"
	"github.com/ont-bizsuite/addonstep/pkg/service"
)

func main() {
	r := gin.Default()

	// FIXME: pc should init with detail option, this is just the demo, we ignore all the config detail info
	sc := service.NewConfig()

	if err := meta.RegistPath(r, sc); err != nil {
		log.Fatal(err)
	}

	r.Run(":8080")
}
