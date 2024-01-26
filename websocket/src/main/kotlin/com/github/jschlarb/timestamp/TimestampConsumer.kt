package com.github.jschlarb.timestamp

import org.apache.kafka.clients.consumer.ConsumerRecord
import org.springframework.context.ApplicationEvent
import org.springframework.context.ApplicationEventPublisher
import org.springframework.kafka.annotation.KafkaListener
import org.springframework.stereotype.Component

@Component
class TimestampConsumer(private val publisher: ApplicationEventPublisher) {

    @KafkaListener(topics = ["\${timestamp.consumer.topic.iso}", "\${timestamp.consumer.topic.basicIsoDate}"])
    fun publishCustomEventIso(record: ConsumerRecord<String, String>) {
        publisher.publishEvent(KafkaTimestampMessage(this, record.value(), record.topic()))
    }
}

class KafkaTimestampMessage(source: Any, val formattedTimestamp: String, val topicName: String) :
    ApplicationEvent(source)
