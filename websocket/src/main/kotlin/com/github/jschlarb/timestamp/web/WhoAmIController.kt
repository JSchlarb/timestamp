package com.github.jschlarb.timestamp.web

import org.springframework.web.bind.annotation.GetMapping
import org.springframework.web.bind.annotation.RestController

@RestController
class WhoAmIController {
    @GetMapping("/who-am-i/admin")
    fun amIAdmin(): Boolean = true

    @GetMapping("/who-am-i/reader")
    fun amIReader(): Boolean = true

    @GetMapping("/who-am-i/writer")
    fun amIWriter(): Boolean = true
}
