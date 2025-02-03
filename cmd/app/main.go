package main

import (
	"context"

	"github.com/spf13/viper"
	"karlota.aasumitro.id/config"
	"karlota.aasumitro.id/internal"

	// embed swagger docs files
	_ "karlota.aasumitro.id/docs"
)

//	@version                     0.0.1-dev
//	@title                       KARLOTA MESSENGER
//	@description                 Karlota REST API Documentation
//
//	@contact.name                @aasumitro
//	@contact.url                 https://aasumitro.id/
//	@contact.email               hello@aasumitro.id
//
// @securityDefinitions.apikey  ApiKeyAuth
// @in                          header
// @name                        Authorization
// @description                 Add you JWT here; e.g: Bearer {JWT}

func main() {
	ctx := context.Background()
	// init config
	viper.SetConfigFile(".env")
	cfg := config.Load().With(ctx,
		config.SQLiteConnection(),
		config.RabbitMQConnection(),
		config.GinEngine())
	// init app
	internal.RunApp(ctx, cfg.Infra)
}
