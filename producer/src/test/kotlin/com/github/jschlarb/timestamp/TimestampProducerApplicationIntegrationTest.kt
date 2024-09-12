package com.github.jschlarb.timestamp

import org.junit.jupiter.api.Test
import org.springframework.beans.factory.annotation.Autowired
import org.springframework.boot.test.context.SpringBootTest

@SpringBootTest
class TimestampProducerApplicationIntegrationTest {
    @Autowired
    private lateinit var producer: TimestampKafkaProducer

    @Test
    fun contextLoads() {
        producer.sendTimestamp()
    }
}
