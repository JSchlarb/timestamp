package main

import (
	"context"
	"dagger/timestamp/internal/dagger"
	"fmt"
)

type Timestamp struct {
	Sock *dagger.Socket
}

func New(sock *dagger.Socket) *Timestamp {
	return &Timestamp{
		Sock: sock,
	}
}

func (m *Timestamp) GradleBuildImage(
	ctx context.Context,
	source *dagger.Directory,
) *dagger.Container {
	// create a cache volume for Gradle downloads
	gradleCache := dag.CacheVolume("gradle-cache")

	build := dag.Container().
		From("gradle:8.10.0-jdk21").
		WithMountedCache("/home/gradle/.gradle", gradleCache).
		WithMountedDirectory("/app", source).
		WithWorkdir("/app").
		WithUnixSocket("/var/run/docker.sock", m.Sock).
		WithExec([]string{
			"./gradlew",
			"clean",
			"bootBuildImage",
		})

	return build
}

func (m *Timestamp) SchemaRegistryService(kafkaService *dagger.Service) *dagger.Service {
	return dag.Container().
		From("docker.io/bitnami/schema-registry:7.6").
		WithEnvVariable("SCHEMA_REGISTRY_LISTENERS", "http://0.0.0.0:8081").
		WithEnvVariable("SCHEMA_REGISTRY_KAFKA_BROKERS", "PLAINTEXT://kafka:9055").
		WithServiceBinding("kafka", kafkaService).
		WithExposedPort(8081).
		AsService()
}

func (m *Timestamp) KafkaService(serviceName string) *dagger.Service {
	return dag.Container().
		From("docker.io/bitnami/kafka:3.7").
		WithEnvVariable("KAFKA_ENABLE_KRAFT", "yes").
		WithEnvVariable("KAFKA_CFG_NODE_ID", "0").
		WithEnvVariable("KAFKA_CFG_PROCESS_ROLES", "controller,broker").
		WithEnvVariable("KAFKA_CFG_CONTROLLER_QUORUM_VOTERS", "0@localhost:9093").
		WithEnvVariable("KAFKA_CFG_LISTENERS", "PLAINTEXT://:9092,CONTROLLER://:9093,EXTERNAL://:9055").
		WithEnvVariable("KAFKA_CFG_ADVERTISED_LISTENERS", fmt.Sprintf("PLAINTEXT://localhost:9092,EXTERNAL://%s:9055", serviceName)).
		WithEnvVariable("KAFKA_CFG_LISTENER_SECURITY_PROTOCOL_MAP", "CONTROLLER:PLAINTEXT,EXTERNAL:PLAINTEXT,PLAINTEXT:PLAINTEXT").
		WithEnvVariable("KAFKA_CFG_CONTROLLER_LISTENER_NAMES", "CONTROLLER").
		WithEnvVariable("KAFKA_CFG_INTER_BROKER_LISTENER_NAME", "PLAINTEXT").
		WithEnvVariable("KAFKA_CFG_AUTO_CREATE_TOPICS_ENABLE", "true").
		WithEnvVariable("ALLOW_PLAINTEXT_LISTENER", "yes").
		WithExposedPort(9055).
		AsService()
}

func (m *Timestamp) GradleTest(
	ctx context.Context,
	source *dagger.Directory,
) *dagger.Container {
	// create a cache volume for Gradle downloads

	build := dag.Container().
		From("gradle:8.10.0-jdk21").
		WithMountedDirectory("/app", source).
		WithWorkdir("/app").
		WithExec([]string{
			"./gradlew",
			"test",
		})

	return build
}

func (m *Timestamp) GradleUnitTest(ctx context.Context, source *dagger.Directory) *dagger.Container {
	return dag.Container().
		From("gradle:8.10.0-jdk21").
		WithMountedDirectory("/app", source).
		WithWorkdir("/app").
		WithExec([]string{"./gradlew", "test"})
}

func (m *Timestamp) GradleIntegrationTest(ctx context.Context, source *dagger.Directory) *dagger.Container {
	kafkaServicename := "kafka"
	kafkaService := m.KafkaService(kafkaServicename)

	return dag.Container().
		From("gradle:8.10.0-jdk21").
		WithMountedDirectory("/app", source).
		WithWorkdir("/app").
		WithServiceBinding(kafkaServicename, kafkaService).
		WithServiceBinding("schema-registry", m.SchemaRegistryService(kafkaService)).
		WithExec([]string{"./gradlew", "integrationTest"})
}

func (m *Timestamp) VulnerabilityCheckFS(
	source *dagger.Directory,
	// +optional
	severity string,
	// +optional
	exitCode string,
) *dagger.Container {
	if severity == "" {
		severity = "HIGH,CRITICAL"
	}
	if exitCode == "" {
		exitCode = "1"
	}
	return dag.Container().
		From("aquasec/trivy:0.53.0").
		WithMountedDirectory("/app", source).
		WithWorkdir("/app").
		WithExec([]string{"trivy", "fs", "--severity", severity, "--exit-code", exitCode, "."})
}

func (m *Timestamp) VulnerabilityCheckImage(
	image string,
	// +optional
	severity string,
	// +optional
	exitCode string,
) *dagger.Container {
	if severity == "" {
		severity = "HIGH,CRITICAL"
	}
	if exitCode == "" {
		exitCode = "1"
	}
	return dag.Container().
		From("aquasec/trivy:0.53.0").
		WithUnixSocket("/var/run/docker.sock", m.Sock).
		WithExec([]string{"trivy", "image", "--severity", severity, "--exit-code", exitCode, image})
}

func (m *Timestamp) Dind() *dagger.Service {
	return dag.
		Docker().
		WithCacheVolume("dind").
		Dind()
}
