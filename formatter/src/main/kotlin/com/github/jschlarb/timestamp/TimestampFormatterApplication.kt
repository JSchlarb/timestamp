package com.github.jschlarb.timestamp

import org.apache.kafka.clients.consumer.ConsumerRecord
import org.springframework.boot.autoconfigure.SpringBootApplication
import org.springframework.boot.context.properties.ConfigurationProperties
import org.springframework.boot.context.properties.EnableConfigurationProperties
import org.springframework.boot.runApplication
import org.springframework.kafka.annotation.KafkaListener
import org.springframework.kafka.core.KafkaTemplate
import org.springframework.stereotype.Component
import java.time.Instant
import java.time.ZoneId
import java.time.format.DateTimeFormatter

@SpringBootApplication
@EnableConfigurationProperties(TimestampProperties::class)
class TimestampFormatterApplication

fun main(args: Array<String>) {
    runApplication<TimestampFormatterApplication>(*args)
}

@Component
class TimestampKafkaConsumer(
    private val kafkaTemplate: KafkaTemplate<String, String>,
    private val props: TimestampProperties,
) {
    companion object {
        private val isoFormatter = DateTimeFormatter.ISO_INSTANT.withZone(ZoneId.systemDefault())
        private val basicIsoDateFormatter = DateTimeFormatter.BASIC_ISO_DATE.withZone(ZoneId.systemDefault())
    }

    @KafkaListener(topics = ["\${timestamp.consumer.topic}"])
    fun consume(record: ConsumerRecord<String, Long>) {
        val timestamp = record.value()
        val isoTimestamp = isoFormatter.format(Instant.ofEpochMilli(timestamp))
        val basicIsoDateTimestamp = basicIsoDateFormatter.format(Instant.ofEpochMilli(timestamp))

        println("Converting $timestamp to [iso: $isoTimestamp] and [basic iso date: $basicIsoDateTimestamp]")

        kafkaTemplate.send(props.producer.topic.iso, record.key(), isoTimestamp)
        kafkaTemplate.send(props.producer.topic.basicIsoDate, record.key(), basicIsoDateTimestamp)
    }
}

@ConfigurationProperties(prefix = "timestamp")
data class TimestampProperties(
    val producer: Producer,
    val consumer: Consumer,
) {
    data class Producer(
        val topic: Topic,
    ) {
        data class Topic(
            val iso: String,
            val basicIsoDate: String,
        )
    }

    data class Consumer(
        val topic: String,
    )
}
