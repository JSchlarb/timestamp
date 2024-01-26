package com.github.jschlarb.timestamp

import org.springframework.boot.autoconfigure.SpringBootApplication
import org.springframework.boot.context.properties.ConfigurationProperties
import org.springframework.boot.context.properties.EnableConfigurationProperties
import org.springframework.boot.runApplication
import org.springframework.kafka.core.KafkaTemplate
import org.springframework.scheduling.annotation.EnableScheduling
import org.springframework.scheduling.annotation.Scheduled
import org.springframework.stereotype.Component

@SpringBootApplication
@EnableScheduling
@EnableConfigurationProperties(TimestampProperties::class)
class TimestampProducerApplication

fun main(args: Array<String>) {
    runApplication<TimestampProducerApplication>(*args)
}


@Component
class TimestampKafkaProducer(
    private val kafkaTemplate: KafkaTemplate<String, Long>,
    private val props: TimestampProperties
) {
    @Scheduled(fixedDelay = 1_000)
    fun sendTimestamp() {
        kafkaTemplate.send(props.producer.topic, 0, "CURRENT_TIME", System.currentTimeMillis())
    }
}


@ConfigurationProperties(prefix = "timestamp")
data class TimestampProperties(
    val producer: Producer,
) {
    data class Producer(
        val topic: String
    )
}
