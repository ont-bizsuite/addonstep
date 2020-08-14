package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/ont-bizsuite/addonstep/pkg/meta"
	"github.com/ont-bizsuite/addonstep/pkg/service"
)

func main() {
	r := gin.Default()

	// FIXME: service config should init with detail option, this is just the demo, we ignore all the config detail info
	sc := service.NewConfig()

	meta.AppendStep(meta.StepPay)
	meta.AppendStep(meta.StepONS)
	if err := meta.RegistPath(r, sc); err != nil {
		log.Fatal(err)
	}

	r.Run(":8080")
}
