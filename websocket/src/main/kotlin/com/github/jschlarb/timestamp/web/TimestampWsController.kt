package com.github.jschlarb.timestamp.web

import com.github.jschlarb.timestamp.KafkaTimestampMessage
import org.springframework.messaging.simp.SimpMessagingTemplate
import org.springframework.scheduling.annotation.Async
import org.springframework.stereotype.Controller
import org.springframework.transaction.event.TransactionalEventListener

@Controller
class TimestampWsController(val test: SimpMessagingTemplate) {
    @Async
    @TransactionalEventListener
    fun greeting(msg: KafkaTimestampMessage) {
        val format = msg.topicName.substring(msg.topicName.lastIndexOf(".") + 1)

        test.convertAndSend("/topic/timestamps", Timestamp(msg.topicName, msg.formattedTimestamp, format))
    }
}

data class Timestamp(val topic: String, val timestamp: String, val format: String)
