package gateway

import (
	"context"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/grum261/beer/configs"
	"github.com/grum261/beer/proto/userpb"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Run(ctx context.Context) error {
	ctxSig, stop := signal.NotifyContext(ctx, syscall.SIGINT)

	cfg, err := configs.NewConfig()
	if err != nil {
		return errors.Wrap(err, "gateway.Run: failed to load config")
	}

	l, err := zap.NewProduction()
	if err != nil {
		return errors.Wrap(err, "gateway.Run: failed to create logger")
	}

	cc, err := grpc.DialContext(ctxSig, cfg.GRPC.ServerPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return errors.Wrap(err, "gateway.Run: failed to create grpc client connection")
	}

	gwMux := runtime.NewServeMux()

	err = userpb.RegisterUserDeliveryServiceHandler(ctxSig, gwMux, cc)
	if err != nil {
		return errors.Wrap(err, "gateway.Run: failed to register grpc-gateway user handler")
	}

	gwServer := &http.Server{
		Addr:    cfg.GRPC.GatewayPort,
		Handler: gwMux,
	}

	g, ctxGroup := errgroup.WithContext(ctxSig)

	l.Info(
		"starting gRPC gateway",
		zap.String("gatewayPort", cfg.GRPC.GatewayPort),
		zap.String("serverPort", cfg.GRPC.ServerPort),
	)
	g.Go(func() error {
		return errors.Wrap(gwServer.ListenAndServe(), "failed to start grpc gateway")
	})

	g.Go(func() error {
		<-ctxGroup.Done()
		ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Second)
		defer func() {
			_ = l.Sync()
			cc.Close()
			stop()
			cancel()
		}()

		return errors.Wrap(gwServer.Shutdown(ctxTimeout), "failed to gracefully stop grpc gateway")
	})

	return errors.Wrap(g.Wait(), "gateway.Run")
}
