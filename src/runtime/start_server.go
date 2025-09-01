package runtime

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

func Start(lc fx.Lifecycle, router *gin.Engine) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				log.Println("Starting server on :8080")
				if err := router.Run(":8080"); err != nil {
					log.Fatal("Failed to start server:", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Println("Stopping server")
			return nil
		},
	})
}