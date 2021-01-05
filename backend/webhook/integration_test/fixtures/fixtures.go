package fixtures

import (
	"fmt"
	"log"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/kelseyhightower/envconfig"

	"github.com/golang-migrate/migrate/v4"
	"github.com/orlangure/gnomock"
	"github.com/orlangure/gnomock/preset/postgres"
	"github.com/orlangure/gnomock/preset/rabbitmq"
	"github.com/orlangure/gnomock/preset/redis"
)

type FixturesEnv struct {
	PostgresUser         string `envconfig:"POSTGRES_USER"`
	PostgresUserPassword string `envconfig:"POSTGRES_USER_PASSWORD"`
	RabbitmqUser         string `envconfig:"RABBITMQ_USER"`
	RabbitmqUserPassword string `envconfig:"RABBITMQ_USER_PASSWORD"`
}
type TestFixtures struct {
	env      *FixturesEnv
	Postgres struct {
		ConnStr string
	}
	Rabbit struct {
		ConnStr string
	}
	Redis struct {
		Address string
	}
	teardowns []func()
}

func New() *TestFixtures {

	var env FixturesEnv
	if err := envconfig.Process("INTEGRATION_TEST_FIXTURE", &env); err != nil {
		log.Fatal("failed to read env vars:", err)
	}

	return &TestFixtures{env: &env}
}

func (tf *TestFixtures) SetupPostgres(dbName string, migrationRoot string) {
	tf.setupPostgresContainer(dbName)
	tf.migrate(migrationRoot)
}
func (tf *TestFixtures) SetupRabbit() {
	tf.setupRabbitContainer()
}
func (tf *TestFixtures) SetupRedis() {
	tf.setupRedisContainer()
}

func (tf *TestFixtures) Teardown() {
	for i := len(tf.teardowns) - 1; i >= 0; i-- {
		tf.teardowns[i]()
	}
}
func (tf *TestFixtures) setupPostgresContainer(dbName string) {

	pg := postgres.Preset(
		postgres.WithUser(tf.env.PostgresUser, tf.env.PostgresUserPassword),
		postgres.WithDatabase(dbName),
	)

	container, err := gnomock.Start(pg)
	if err != nil {
		log.Fatal(err.Error())
	}

	tf.teardowns = append(tf.teardowns, func() {
		err = gnomock.Stop(container)
		if err != nil {
			log.Printf("failed to shutdown Postgres container: %v\n", err)
		}
	})

	tf.Postgres.ConnStr = fmt.Sprintf(
		"postgresql://%s:%s@%s:%d/%s?sslmode=%s",
		tf.env.PostgresUser, tf.env.PostgresUserPassword,
		container.Host, container.DefaultPort(), "webhook", "disable")

}

func (tf *TestFixtures) migrate(sourceUrl string) {
	m, err := migrate.New(
		sourceUrl,
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
	container, err := gnomock.Start(p)
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

	container, err := gnomock.Start(p)
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
