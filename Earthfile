VERSION 0.8
IMPORT github.com/DenktMit-eG/earthly-lib/java AS java
IMPORT github.com/DenktMit-eG/earthly-lib/sbom AS sbom
IMPORT github.com/DenktMit-eG/earthly-lib/kind AS kind
IMPORT github.com/DenktMit-eG/local-dev-cluster:earthly-experiments AS local-dev-cluster
# Dockerfile for a Maven + Docker environment

maven-build:
    FROM maven:3.9-eclipse-temurin-21

    # SCRIPT
    DO java+JAVA_MVN --goals="clean install" --sources="formatter producer websocket" --project_root="."

    # Artifacts
    SAVE ARTIFACT formatter/target ./formatter/target
    SAVE ARTIFACT producer/target ./producer/target
    SAVE ARTIFACT websocket/target ./websocket/target

# Define a build stage using Paketo buildpacks
build:
    FROM buildpacksio/pack:0.35.0
    RUN ls -la

# Define an image stage
docker-image:
    FROM build
    SAVE IMAGE my-app:latest

sbom:
    FROM sbom+base
    DO sbom+SBOM_FS

vulnerability-scan:
    FROM sbom+base
    DO sbom+VULNERABILITY_SCAN_FS
