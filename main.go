package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"

	httpSwagger "github.com/swaggo/http-swagger"
	"golang.org/x/sync/errgroup"

	"github.com/MichalPolinkiewicz/roche/cmd/grpc"
	"github.com/MichalPolinkiewicz/roche/cmd/rest"
	"github.com/MichalPolinkiewicz/roche/config"
	"github.com/MichalPolinkiewicz/roche/docs"
	"github.com/MichalPolinkiewicz/roche/pkg/mapper"
	"github.com/MichalPolinkiewicz/roche/pkg/service"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	appConfig, err := config.NewAppConfig("default")
	if err != nil {
		log.Fatal(err)
	}

	if err := start(ctx, appConfig); err != nil {
		log.Println(err)
	}
	log.Println("app done")
}

func start(ctx context.Context, appConfig *config.AppConfig) error {
	// client && ping service init
	postmanClient, err := service.NewPostmanClient(appConfig.PingServiceClientEndpoint, http.DefaultClient)
	if err != nil {
		return err
	}
	pingService, err := service.NewPingService(postmanClient, appConfig.PingServiceTimeout, service.NewPingDecorator(appConfig.Version, appConfig.Env))
	if err != nil {
		return err
	}

	// create REST server
	pingRestHandler, err := rest.NewPingHandler(pingService, &mapper.PingRequestOneToOneMapper{}, &mapper.PingResponseOneToOneMapper{})
	if err != nil {
		return err
	}
	pingRestEndpoint, err := rest.NewEndpoint("/ping", http.MethodPost, pingRestHandler.Handle)
	if err != nil {
		return err
	}

	swaggerEndpoint, err := rest.NewEndpoint("/swagger/", http.MethodGet, httpSwagger.Handler(httpSwagger.URL("http://localhost:"+appConfig.RestPort+"/swagger/swagger.yaml")))
	swaggerServeFileEndpoint, err := rest.NewEndpoint("/swagger/swagger.yaml", http.MethodGet, docs.SwaggerServefile)

	restServer, err := rest.NewRestServer(appConfig.RestPort, []*rest.Endpoint{swaggerEndpoint, pingRestEndpoint, swaggerServeFileEndpoint})
	if err != nil {
		return err
	}

	// create GRPC server
	grpcServer, err := grpc.NewGrpcServer(appConfig.GrpcPort, pingService, &mapper.PingRequestOneToOneMapper{}, &mapper.PingResponseOneToOneMapper{})
	if err != nil {
		return err
	}

	return runServers(ctx, restServer, grpcServer)
}

func runServers(ctx context.Context, restSrv *rest.RestServer, grpcSrv *grpc.GrpcServer) error {
	group, ctx := errgroup.WithContext(ctx)
	group.Go(restSrv.Run(ctx))
	group.Go(grpcSrv.Run(ctx))
	return group.Wait()
}
