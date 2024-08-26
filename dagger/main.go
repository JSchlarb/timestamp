package main

import (
	"context"
	"dagger/timestamp/internal/dagger"
)

type Timestamp struct{}

func (m *Timestamp) Build(
	ctx context.Context,
	source *dagger.Directory,
	publish bool,
	// +optional
	registryUser *dagger.Secret,
	// +optional
	registryPassword *dagger.Secret,
) *dagger.Container {
	// create a cache volume for Gradle downloads
	gradleCache := dag.CacheVolume("gradle-cache")

	build := dag.Container().
		From("gradle:8.10.0-jdk21").
		WithMountedCache("/home/gradle/.gradle", gradleCache).
		WithMountedDirectory("/app", source).
		WithWorkdir("/app").
		WithServiceBinding("dind", m.Dind()).
		WithEnvVariable("DOCKER_HOST", "tcp://dind:2375")

	if publish {
		build = build.
			WithEnvVariable("DOCKER_REGISTRY_URL", "my.repo.url").
			WithSecretVariable("DOCKER_USERNAME", registryUser).
			WithSecretVariable("DOCKER_PASSWORD", registryPassword).
			WithExec([]string{
				"./gradlew",
				"clean",
				"bootBuildImage",
				"--publishImage",
			})
	} else {
		build = build.
			WithExec([]string{
				"./gradlew",
				"clean",
				"bootBuildImage",
			})
	}

	return build
}

func (m *Timestamp) SchemaRegistryService() *dagger.Service {
	return dag.Container().
		From("docker.io/bitnami/schema-registry:7.6").
		WithEnvVariable("SCHEMA_REGISTRY_LISTENERS", "http://0.0.0.0:8081").
		WithEnvVariable("SCHEMA_REGISTRY_KAFKA_BROKERS", "PLAINTEXT://kafka:9055").
		WithServiceBinding("kafka", m.KafkaService()).
		WithExposedPort(8081).
		AsService()
}

func (m *Timestamp) KafkaService() *dagger.Service {
	return dag.Container().
		From("docker.io/bitnami/kafka:3.7").
		WithEnvVariable("KAFKA_ENABLE_KRAFT", "yes").
		WithEnvVariable("KAFKA_CFG_NODE_ID", "0").
		WithEnvVariable("KAFKA_CFG_PROCESS_ROLES", "controller,broker").
		WithEnvVariable("KAFKA_CFG_CONTROLLER_QUORUM_VOTERS", "0@localhost:9093").
		WithEnvVariable("KAFKA_CFG_LISTENERS", "PLAINTEXT://:9092,CONTROLLER://:9093,EXTERNAL://:9055").
		WithEnvVariable("KAFKA_CFG_ADVERTISED_LISTENERS", "PLAINTEXT://localhost:9092,EXTERNAL://kafka:9055").
		WithEnvVariable("KAFKA_CFG_LISTENER_SECURITY_PROTOCOL_MAP", "CONTROLLER:PLAINTEXT,EXTERNAL:PLAINTEXT,PLAINTEXT:PLAINTEXT").
		WithEnvVariable("KAFKA_CFG_CONTROLLER_LISTENER_NAMES", "CONTROLLER").
		WithEnvVariable("KAFKA_CFG_INTER_BROKER_LISTENER_NAME", "PLAINTEXT").
		WithEnvVariable("KAFKA_CFG_AUTO_CREATE_TOPICS_ENABLE", "true").
		WithEnvVariable("ALLOW_PLAINTEXT_LISTENER", "yes").
		WithExposedPort(9055).
		AsService()
}

func (m *Timestamp) IntegrationTest(ctx context.Context, source *dagger.Directory) *dagger.Container {
	mavenCache := dag.CacheVolume("maven-cache")

	return dag.Container().
		From("maven:3.9-eclipse-temurin-17").
		WithMountedCache("/root/.m2", mavenCache).
		WithMountedDirectory("/app", source).
		WithWorkdir("/app").
		WithServiceBinding("schema-registry", m.SchemaRegistryService()).
		WithExec([]string{"curl", "schema-registry:8081/v1/metadata/id"})
}

func (m *Timestamp) VulnerabilityCheck(
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

func (m *Timestamp) Dind() *dagger.Service {
	return dag.
		Docker().
		WithCacheVolume("dind").
		Dind()
}
