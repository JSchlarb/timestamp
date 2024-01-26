package com.github.jschlarb.timestamp.web

import org.springframework.messaging.handler.annotation.MessageMapping
import org.springframework.messaging.handler.annotation.SendTo
import org.springframework.messaging.simp.SimpMessagingTemplate
import org.springframework.security.core.context.SecurityContextHolder
import org.springframework.stereotype.Controller


@Controller
class WhoAmIWsController(val test: SimpMessagingTemplate) {
    @MessageMapping("/whoami")
    @SendTo("/topic/whoami")
    fun tester(): String {
        val authentication = SecurityContextHolder.getContext().authentication

        return authentication.name
    }
}
