package com.github.jschlarb.timestamp.config

import org.springframework.context.annotation.Bean
import org.springframework.context.event.ApplicationEventMulticaster
import org.springframework.context.event.SimpleApplicationEventMulticaster
import org.springframework.core.task.SimpleAsyncTaskExecutor
import org.springframework.stereotype.Component

@Component
class AsyncConfig {
    @Bean
    fun applicationEventMulticaster(): ApplicationEventMulticaster =
        SimpleApplicationEventMulticaster().apply {
            setTaskExecutor(SimpleAsyncTaskExecutor())
        }
}
