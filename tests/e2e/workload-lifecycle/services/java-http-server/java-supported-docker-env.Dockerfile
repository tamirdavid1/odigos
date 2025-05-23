FROM maven:3.8.5-openjdk-17 AS build
COPY src /home/app/src
COPY pom.xml /home/app
RUN mvn -f /home/app/pom.xml clean package

FROM eclipse-temurin:17.0.12_7-jre-jammy
ENV JAVA_OPTS="-Dthis.does.not.exist=true"
ENV LD_PRELOAD=""
COPY --from=build /home/app/target/*.jar /app/java-supported-docker-env.jar
USER 15000
CMD ["sh", "-c", "java $JAVA_OPTS -jar /app/java-supported-docker-env.jar"]
