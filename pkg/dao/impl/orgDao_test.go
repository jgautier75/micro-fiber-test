package impl

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"log"
	"micro-fiber-test/pkg/model"
	"os"
	"testing"
)

type TestLogConsumer struct {
	Msgs []string
}

func (g *TestLogConsumer) Accept(l testcontainers.Log) {
	g.Msgs = append(g.Msgs, string(l.Content))
}

func TestDao(t *testing.T) {
	ctx := context.Background()

	// Force docker port on WSL (daemon.json > hosts)
	err := os.Setenv("DOCKER_HOST", "tcp://localhost:2375")
	if err != nil {
		log.Fatal(err)
	}

	// Get current directory
	mydir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Current directory: %s\n", mydir)

	logConsumer := TestLogConsumer{
		Msgs: []string{},
	}

	// Init postgreSQL container with testcontainers
	containerReq := testcontainers.ContainerRequest{
		Image:        "postgres:latest",
		ExposedPorts: []string{"5432/tcp"},
		WaitingFor:   wait.ForListeningPort("5432/tcp"),
		Env: map[string]string{
			"POSTGRES_DB":       "testdb",
			"POSTGRES_PASSWORD": "postgres",
			"POSTGRES_USER":     "postgres",
		},
	}
	dbContainer, _ := testcontainers.GenericContainer(
		context.Background(),
		testcontainers.GenericContainerRequest{
			ContainerRequest: containerReq,
			Started:          true,
		})
	defer func(dbContainer testcontainers.Container, ctx context.Context) {
		err := dbContainer.Terminate(ctx)
		if err != nil {

		}
	}(dbContainer, ctx)

	errLogProd := dbContainer.StartLogProducer(ctx)
	if errLogProd != nil {
		fmt.Printf("Error on log producer: [%v]", errLogProd)
	}
	dbContainer.FollowOutput(&logConsumer)

	// Retrieve postgreSQL container host and port
	host, _ := dbContainer.Host(context.Background())
	port, _ := dbContainer.MappedPort(context.Background(), "5432")
	fmt.Printf("postgreSQL started on [%s]:[%s] \n", host, port)
	pgUrl := fmt.Sprintf("postgres://postgres:postgres@%s:%d/testdb?sslmode=disable", host, port.Int())

	// Run migrations: create tables, sequences, ...
	db, err := sql.Open("postgres", pgUrl)
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	m, err := migrate.NewWithDatabaseInstance(
		"file://../../migrations",
		"postgres", driver)
	errMig := m.Up()
	if errMig != nil {
		panic(errMig)
	}

	// Test create organization
	orgRepo := NewOrgDao()
	org := model.Organization{}
	org.SetStatus(model.OrgStatusActive)
	org.SetLabel("Test Org")
	org.SetCode("test")
	org.SetTenantId(1)
	org.SetType(model.OrgTypeCommunity)
	orgId, err := orgRepo.Create(pgUrl, &org)
	if err != nil {
		fmt.Printf("pgError [%v]", err)
	} else {
		assert.NotNil(t, orgId)
		fmt.Printf("orgId [%d]", orgId)
	}

	fmt.Printf("Container logs: [%v]", logConsumer.Msgs)

}
