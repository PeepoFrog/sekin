package main

import (
	"github.com/kiracore/sekin/src/shidai/internal/api"
	"github.com/kiracore/sekin/src/shidai/internal/logger"
)

func main() {
	logger.InitLogger()     // Init logger
	defer logger.Log.Sync() // Flush any buffer log entries

	api.Serve()

}
