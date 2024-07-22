// A generated module for Timestamp functions
//
// This module has been generated via dagger init and serves as a reference to
// basic module structure as you get started with Dagger.
//
// Two functions have been pre-created. You can modify, delete, or add to them,
// as needed. They demonstrate usage of arguments and return types using simple
// echo and grep commands. The functions can be called from the dagger CLI or
// from one of the SDKs.
//
// The first line in this comment block is a short description line and the
// rest is a long description with more detail on the module's purpose or usage,
// if appropriate. All modules should have a short description.

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
		WithExec([]string{
			"trivy",
			"fs",
			"--severity",
			severity,
			"--exit-code",
			exitCode,
			".",
		})
}

func (m *Timestamp) Dind() *dagger.Service {
	return dag.
		Docker().
		WithCacheVolume("dind").
		Dind()
}
