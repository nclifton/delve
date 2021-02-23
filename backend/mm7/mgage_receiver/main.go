package main

import (
	"github.com/burstsms/mtmo-tp/backend/lib/restbuilder"
	"github.com/burstsms/mtmo-tp/backend/mm7/mgage_receiver/builder"
)

func main() {
	restbuilder.NewFromEnv(builder.NewFromEnv()).Start()
}
