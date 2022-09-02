package grpcserver

import (
	"context"
	"net"
	"os/signal"
	"syscall"
	"time"

	middleware "github.com/grpc-ecosystem/go-grpc-middleware/v2"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/auth"
	"github.com/grum261/beer/configs"
	rpcdelivery "github.com/grum261/beer/internal/delivery/grpc"
	"github.com/grum261/beer/internal/repository"
	"github.com/grum261/beer/internal/service"
	"github.com/grum261/beer/proto/userpb"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

func Run(ctx context.Context) error {
	ctxSig, stop := signal.NotifyContext(ctx, syscall.SIGINT)

	cfg, err := configs.NewConfig()
	if err != nil {
		return errors.Wrap(err, "grpcserver.Run: failed to fetch config")
	}

	l, err := zap.NewProduction()
	if err != nil {
		return errors.Wrap(err, "grpcserver.Run: failed to create logger")
	}

	pool, err := pgxpool.Connect(ctxSig, cfg.DB.String())
	if err != nil {
		return errors.Wrap(err, "grpcserver.Run: failed to connect to db")
	}

	l.Info("successfully connected to db", zap.String("dsn", cfg.DB.String()))

	userRepo := repository.NewUserRepository(pool)
	userSvc := service.NewUserService(userRepo, cfg.Argon2)
	userDelivery := rpcdelivery.NewUserDelivery(userSvc, cfg.JWT)

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(
			middleware.ChainUnaryServer(
				auth.UnaryServerInterceptor(rpcdelivery.JWTAuth(cfg.JWT.Secret)),
			),
		),
	)

	userpb.RegisterUserDeliveryServiceServer(grpcServer, userDelivery)

	lis, err := net.Listen("tcp", cfg.GRPC.ServerPort)
	if err != nil {
		return errors.Wrap(err, "grpcserver.Run: failed to create net.Listener")
	}

	g, ctxGroup := errgroup.WithContext(ctxSig)

	l.Info("starting gRPC server", zap.String("port", cfg.GRPC.ServerPort))
	g.Go(func() error {
		return errors.Wrap(grpcServer.Serve(lis), "failed to start grpc server")
	})

	g.Go(func() error {
		<-ctxGroup.Done()
		l.Info("received signal to stop server")
		ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Second)
		defer func() {
			_ = l.Sync()
			pool.Close()
			grpcServer.GracefulStop()
			lis.Close()
			stop()
			cancel()
		}()

		<-ctxTimeout.Done()

		l.Info("gracefully shutted down")

		return nil
	})

	return errors.Wrap(g.Wait(), "grpcserver.Run")
}
