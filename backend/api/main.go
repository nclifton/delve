package main

import (
	"github.com/burstsms/mtmo-tp/backend/api/builder"
	"github.com/burstsms/mtmo-tp/backend/lib/restbuilder"
)

func main() {
	restbuilder.NewFromEnv(builder.NewFromEnv()).Start()
}
