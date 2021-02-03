package fixtures

import (
	"fmt"
	"log"
	"net"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/kelseyhightower/envconfig"
	"github.com/phayes/freeport"

	"github.com/golang-migrate/migrate/v4"
	"github.com/orlangure/gnomock"
	"github.com/orlangure/gnomock/preset/postgres"
	"github.com/orlangure/gnomock/preset/rabbitmq"
	"github.com/orlangure/gnomock/preset/redis"
)

type FixturesEnv struct {
	PostgresUser           string   `envconfig:"POSTGRES_USER" default:"gnomock"`
	PostgresUserPassword   string   `envconfig:"POSTGRES_USER_PASSWORD" default:"gnomick"`
	RabbitmqUser           string   `envconfig:"RABBITMQ_USER" default:"gnomock"`
	RabbitmqUserPassword   string   `envconfig:"RABBITMQ_USER_PASSWORD" default:"gnomick"`
	MigrationRoot          string   `envconfig:"MIGRATION_ROOT" default:"file://../migration/sql"`
	HealthCheckHost        string   `envconfig:"RPC_HEALTH_CHECK_HOST" default:"127.0.0.1"`    // all health check services under test listen and serve on the same host but separate ports
	RPCHealthCheckPort     string   `envconfig:"RPC_HEALTH_CHECK_PORT" default:"FREEPORT"`     // port number named "FREEPORT" will allow the fixture to allocate any unused port number
	WorkerHealthCheckPorts []string `envconfig:"WORKER_HEALTH_CHECK_PORTS" default:"FREEPORT"` // comma separated list of PORTS to use or use "FREEPORT"
}
type TestFixtures struct {
	name     string
	env      *FixturesEnv
	Postgres struct {
		ConnStr string
		Stop    func()
	}
	Rabbit struct {
		ConnStr string
	}
	Redis struct {
		Address string
	}
	teardowns             []func()
	GRPCListener          net.Listener
	RPCHealthCheckURI     string
	WorkerHealthCheckURIs []string
	workerPortIndex       int
}

func New(name string) *TestFixtures {

	log.Println("setup fixtures")
	var env FixturesEnv
	if err := envconfig.Process("TEST_FIXTURE", &env); err != nil {
		log.Fatal("failed to read env vars:", err)
	}

	return &TestFixtures{env: &env, name: name}

}

// return any available port if provided port string equals "FREEPORT"
func port(port string) string {
	if port == "FREEPORT" {
		port, err := freeport.GetFreePort()
		if err != nil {
			log.Fatal(err)
		}
		return fmt.Sprintf("%d", port)
	}
	return port
}

func (tf *TestFixtures) SetupPostgres(dbName string) {
	tf.setupPostgresContainer(dbName)
	tf.migrate()
}
func (tf *TestFixtures) SetupRabbit() {
	tf.setupRabbitContainer()
}
func (tf *TestFixtures) SetupRedis() {
	tf.setupRedisContainer()
}

func (tf *TestFixtures) Teardown() {
	log.Println("teardown fixtures")
	for i := len(tf.teardowns) - 1; i >= 0; i-- {
		tf.teardowns[i]()
	}
}
func (tf *TestFixtures) setupPostgresContainer(dbName string) {

	pg := postgres.Preset(
		postgres.WithUser(tf.env.PostgresUser, tf.env.PostgresUserPassword),
		postgres.WithDatabase(dbName),
	)

	container, err := gnomock.Start(pg, gnomock.WithContainerName(fmt.Sprintf("%s-postgres-fixture", tf.name)))
	if err != nil {
		log.Fatal(err.Error())
	}

	stop := func() {
		err = gnomock.Stop(container)
	}

	tf.teardowns = append(tf.teardowns, stop)
	tf.Postgres.Stop = stop

	tf.Postgres.ConnStr = fmt.Sprintf(
		"postgresql://%s:%s@%s/%s?sslmode=%s",
		tf.env.PostgresUser, tf.env.PostgresUserPassword,
		container.DefaultAddress(), dbName, "disable")

}

func (tf *TestFixtures) migrate() {
	m, err := migrate.New(
		tf.env.MigrationRoot,
		tf.Postgres.ConnStr,
	)
	if err != nil {
		log.Fatalf("Failed to initialise golang-migrate connection: %s\n", err)
	}

	// apply migrations
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Failed to run migrations: %s\n", err)
	}

}

func (tf *TestFixtures) setupRabbitContainer() {
	p := rabbitmq.Preset(
		rabbitmq.WithUser(tf.env.RabbitmqUser, tf.env.RabbitmqUserPassword),
	)
	container, err := gnomock.Start(p, gnomock.WithContainerName(fmt.Sprintf("%s-rabbit-fixture", tf.name)))
	if err != nil {
		log.Fatal(err.Error())
	}

	tf.teardowns = append(tf.teardowns, func() {
		err = gnomock.Stop(container)
		if err != nil {
			fmt.Printf("failed to shutdown RabbitMQ container: %v\n", err)
		}
	})

	tf.Rabbit.ConnStr = fmt.Sprintf(
		"amqp://%s:%s@%s",
		tf.env.RabbitmqUser, tf.env.RabbitmqUserPassword,
		container.DefaultAddress(),
	)

}

func (tf *TestFixtures) setupRedisContainer() {
	vs := make(map[string]interface{})

	// Setup Redis
	p := redis.Preset(redis.WithValues(vs))

	container, err := gnomock.Start(p, gnomock.WithContainerName(fmt.Sprintf("%s-redis-fixture", tf.name)))
	if err != nil {
		log.Fatal(err.Error())
	}

	tf.teardowns = append(tf.teardowns, func() {
		err = gnomock.Stop(container)
		if err != nil {
			fmt.Printf("failed to shutdown Redis container: %v\n", err)
		}
	})

	tf.Redis.Address = container.DefaultAddress()
}
