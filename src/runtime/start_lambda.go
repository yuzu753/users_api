package runtime

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	ginengine "github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

func StartLambda(lc fx.Lifecycle, router *ginengine.Engine) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			ginLambda := ginadapter.New(router)
			
			go func() {
				log.Println("Starting Lambda handler")
				lambda.Start(ginLambda.ProxyWithContext)
			}()
			
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Println("Stopping Lambda handler")
			return nil
		},
	})
}