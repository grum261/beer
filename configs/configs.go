package configs

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/cristalhq/aconfig"
	"github.com/cristalhq/aconfig/aconfigdotenv"
)

type DBConfig struct {
	User     string `env:"USER"`
	Password string `env:"PASS"`
	Host     string `env:"HOST"`
	Port     uint   `env:"PORT" default:"5432"`
	Name     string `env:"NAME"`
}

type Argon2Config struct {
	Memory     uint32 `env:"MEMORY" default:"65536"`
	Iterations uint32 `env:"ITERATIONS" default:"3"`
	Threads    uint8  `env:"THREADS" default:"2"`
	SaltLength uint32 `env:"SALT_LENGTH" default:"16"`
	KeyLength  uint32 `env:"KEY_LENGTH" default:"32"`
}

type JWTConfig struct {
	Secret   string        `env:"SECRET"`
	TokenTTL time.Duration `env:"TOKEN_TTL" default:"10m"`
}

type GRPCConfig struct {
	ServerPort  string `env:"SERVER_PORT" default:":8000"`
	GatewayPort string `env:"GATEWAY_PORT" default:":8001"`
}

type FileConfig struct {
	MaxSize int64 `env:"MAX_SIZE" default:"1048576"`
}

type Config struct {
	Argon2 Argon2Config `env:"ARGON2"`
	DB     DBConfig     `env:"DB"`
	JWT    JWTConfig    `env:"JWT"`
	GRPC   GRPCConfig   `env:"GRPC"`
	File   FileConfig   `env:"FILE"`
}

type flags struct {
	IsDebug bool `flag:"is_debug"`
	IsLocal bool `flag:"is_local"`
}

func NewConfig() (Config, error) {
	acfg := aconfig.Config{
		SkipFiles: true,
		SkipEnv:   true,
	}

	var flag flags

	if err := aconfig.LoaderFor(&flag, acfg).Load(); err != nil {
		return Config{}, err
	}

	acfg.SkipFlags = true

	if flag.IsDebug || flag.IsLocal {
		acfg.SkipFiles = false
		acfg.FileDecoders = map[string]aconfig.FileDecoder{
			".env": aconfigdotenv.New(),
		}
		if flag.IsDebug {
			acfg.Files = []string{filepath.Join("..", "..", ".env")}
		} else {
			acfg.Files = []string{filepath.Join(".", ".env")}
		}
	} else {
		acfg.SkipEnv = false
	}

	var cfg Config

	if err := aconfig.LoaderFor(&cfg, acfg).Load(); err != nil {
		return cfg, err
	}

	return cfg, nil
}

func (db *DBConfig) String() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s",
		db.User, db.Password, db.Host, db.Port, db.Name,
	)
}
