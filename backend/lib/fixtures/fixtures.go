package fixtures

import (
	"fmt"
	"log"
	"net"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/golang-migrate/migrate/v4"
	"github.com/orlangure/gnomock"
	"github.com/orlangure/gnomock/preset/postgres"
	"github.com/orlangure/gnomock/preset/rabbitmq"
	"github.com/orlangure/gnomock/preset/redis"
)

type Config struct {
	Name                   string
	PostgresUser           string
	PostgresUserPassword   string
	RabbitmqUser           string
	RabbitmqUserPassword   string
	MigrationRoot          string
	HealthCheckHost        string   // all health check services under test listen and serve on the same host but separate ports
	RPCHealthCheckPort     string   // port number named "FREEPORT" will allow the fixture to allocate any unused port number
	WorkerHealthCheckPorts []string // comma separated list of PORTS to use or use "FREEPORT"
}

var defaults = Config{
	Name:                   "",
	PostgresUser:           "gnomock",
	PostgresUserPassword:   "gnomick",
	RabbitmqUser:           "gnomock",
	RabbitmqUserPassword:   "gnomick",
	MigrationRoot:          "file://../migration/sql", // the domain's migration directory is expected to be a sibling directory of the work directory
	HealthCheckHost:        "127.0.0.1",
	RPCHealthCheckPort:     "18086",
	WorkerHealthCheckPorts: []string{"18087", "18088", "18089"},
}

type TestFixtures struct {
	config   Config
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

func New(config Config) *TestFixtures {

	log.Println("setup fixtures")

	return &TestFixtures{
		config: Config{
			Name:                   config.Name,
			PostgresUser:           ifEmptyStringThen(config.PostgresUser, defaults.PostgresUser),
			PostgresUserPassword:   ifEmptyStringThen(config.PostgresUser, defaults.PostgresUser),
			RabbitmqUser:           ifEmptyStringThen(config.RabbitmqUser, defaults.RabbitmqUser),
			RabbitmqUserPassword:   ifEmptyStringThen(config.RabbitmqUser, defaults.RabbitmqUser),
			MigrationRoot:          ifEmptyStringThen(config.MigrationRoot, defaults.MigrationRoot),
			HealthCheckHost:        ifEmptyStringThen(config.HealthCheckHost, defaults.HealthCheckHost),
			RPCHealthCheckPort:     ifEmptyStringThen(config.RPCHealthCheckPort, defaults.RPCHealthCheckPort),
			WorkerHealthCheckPorts: ifEmptyStringArrayThen(config.WorkerHealthCheckPorts, defaults.WorkerHealthCheckPorts),
		}}

}
func ifEmptyStringThen(value string, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}
func ifEmptyStringArrayThen(value []string, defaultValue []string) []string {
	if len(value) == 0 {
		return defaultValue
	}
	return value
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
		postgres.WithUser(tf.config.PostgresUser, tf.config.PostgresUserPassword),
		postgres.WithDatabase(dbName),
	)

	container, err := gnomock.Start(pg, gnomock.WithContainerName(fmt.Sprintf("%s-postgres-fixture", tf.config.Name)))
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
		tf.config.PostgresUser, tf.config.PostgresUserPassword,
		container.DefaultAddress(), dbName, "disable")

}

func (tf *TestFixtures) migrate() {
	m, err := migrate.New(
		tf.config.MigrationRoot,
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
		rabbitmq.WithUser(tf.config.RabbitmqUser, tf.config.RabbitmqUserPassword),
	)
	container, err := gnomock.Start(p, gnomock.WithContainerName(fmt.Sprintf("%s-rabbit-fixture", tf.config.Name)))
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
		tf.config.RabbitmqUser, tf.config.RabbitmqUserPassword,
		container.DefaultAddress(),
	)

}

func (tf *TestFixtures) setupRedisContainer() {
	vs := make(map[string]interface{})

	// Setup Redis
	p := redis.Preset(redis.WithValues(vs))

	container, err := gnomock.Start(p, gnomock.WithContainerName(fmt.Sprintf("%s-redis-fixture", tf.config.Name)))
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
