package com.github.jschlarb.timestamp

import org.springframework.boot.autoconfigure.SpringBootApplication
import org.springframework.boot.context.properties.ConfigurationProperties
import org.springframework.boot.context.properties.EnableConfigurationProperties
import org.springframework.boot.runApplication

@SpringBootApplication
@EnableConfigurationProperties(TimestampProperties::class)
class TimestampWebSocketApplication

fun main(args: Array<String>) {
    runApplication<TimestampWebSocketApplication>(*args)
}


@ConfigurationProperties(prefix = "timestamp")
data class TimestampProperties(
    val consumer: Consumer
) {
    data class Consumer(
        val topic: Topic
    ) {
        data class Topic(
            val iso: String,
            val basicIsoDate: String
        )
    }
}
