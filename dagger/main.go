package main

import (
	"context"
	"dagger/timestamp/internal/dagger"
	"strconv"
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
	// create a cache volume for Maven downloads
	mavenCache := dag.CacheVolume("maven-cache")

	build := dag.Container().
		From("maven:3.9-eclipse-temurin-17").
		WithMountedCache("/root/.m2", mavenCache).
		WithMountedDirectory("/app", source).
		WithWorkdir("/app").
		WithServiceBinding("dind", m.Dind()).
		WithEnvVariable("DOCKER_HOST", "tcp://dind:2375")

	if publish {
		build = build.
			WithSecretVariable("REGISTRY_USER", registryUser).
			WithSecretVariable("REGISTRY_PASSWORD", registryPassword).
			WithExec([]string{
				"mvn",
				"-Dspring-boot.build-image.publish=" + strconv.FormatBool(publish),
				"-Ddocker.publishRegistry.username=$REGISTRY_USER",
				"-Ddocker.publishRegistry.password=$REGISTRY_PASSWORD",
				"clean",
				"spring-boot:build-image",
			})
	} else {
		build = build.
			WithExec([]string{
				"mvn",
				"clean",
				"spring-boot:build-image",
			})
	}

	return build
}

func (m *Timestamp) IntegrationTest(ctx context.Context, source *dagger.Directory) *dagger.Container {
	mavenCache := dag.CacheVolume("maven-cache")
	svc := dag.
		Compose().
		WithFile("./compose.yaml")

	it := svc.
		From("maven:3.9-eclipse-temurin-17").
		WithMountedCache("/root/.m2", mavenCache).
		WithMountedDirectory("/app", source).
		WithWorkdir("/app").
		WithExec([]string{"curl", "http://localhost:8081"})

	return it
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
