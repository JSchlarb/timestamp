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

JAVA_MVN_CONTAINER_DEPS:
    FUNCTION
    RUN mkdir target/dependency
    RUN (cd target/dependency; jar -xf ../*.jar)
    SAVE ARTIFACT target target

JAVA_MVN_BUILD_CONTAINER:
    FUNCTION
    ARG snapshot_repository="snapshot-registry.denktmit.tech"
    ARG project_name="example"
    ARG project_version="0.0.1"
    ARG main_class="com.github.denktmit.example.web.Application.Kt"
    RUN cp -r target/dependency/BOOT-INF/lib /app/lib
    RUN cp -r target/dependency/META-INF /app/META-INF
    RUN cp -r target/dependency/BOOT-INF/classes/* /app
    ENTRYPOINT ["java","-cp","app:app/lib/*",$main_class]
    SAVE IMAGE $snapshot_repository/$project_name:$project_version

docker-formatter-deps:
    FROM maven:3.9-eclipse-temurin-21
    COPY +maven-build/formatter/target target
    DO +JAVA_MVN_CONTAINER_DEPS

docker-formatter:
    FROM eclipse-temurin:21-jre
    WORKDIR app
    COPY +docker-formatter-deps/target target
    DO +JAVA_MVN_BUILD_CONTAINER \
        --main_class "com.github.jschlarb.timestamp.TimestampFormatterApplicationKt" \
        --project_name "formatter" \
        --project_version="0.1.7"

sbom:
    FROM sbom+base
    DO sbom+SBOM_FS

vulnerability-scan:
    FROM sbom+base
    DO sbom+VULNERABILITY_SCAN_FS
