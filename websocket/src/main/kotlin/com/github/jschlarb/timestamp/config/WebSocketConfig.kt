package com.github.jschlarb.timestamp.config

import org.springframework.context.annotation.Configuration
import org.springframework.core.convert.converter.Converter
import org.springframework.security.authentication.AbstractAuthenticationToken
import org.springframework.security.config.annotation.web.messaging.MessageSecurityMetadataSourceRegistry
import org.springframework.security.config.annotation.web.socket.AbstractSecurityWebSocketMessageBrokerConfigurer
import org.springframework.security.core.GrantedAuthority
import org.springframework.security.core.authority.SimpleGrantedAuthority
import org.springframework.security.oauth2.jwt.Jwt
import org.springframework.security.oauth2.server.resource.authentication.JwtAuthenticationToken
import org.springframework.security.oauth2.server.resource.authentication.JwtGrantedAuthoritiesConverter
import org.springframework.stereotype.Component
import org.springframework.web.socket.config.annotation.EnableWebSocketMessageBroker
import java.util.stream.Collectors
import java.util.stream.Stream

@Configuration
@EnableWebSocketMessageBroker
class WebSocketSecurityConfig : AbstractSecurityWebSocketMessageBrokerConfigurer() {
    override fun configureInbound(messages: MessageSecurityMetadataSourceRegistry) {
        messages
            .nullDestMatcher().permitAll()
            .simpSubscribeDestMatchers("/topic/timestamps").hasRole("read-action")
            .simpSubscribeDestMatchers("/topic/whoami").authenticated()
            .simpMessageDestMatchers("/app/whoami").authenticated()
            .anyMessage().denyAll()
    }

    // Disable CSRF for WebSockets
    override fun sameOriginDisabled(): Boolean = true
}

// pls don't hate me
@Component
class CustomJwtAuthenticationConverter : Converter<Jwt, AbstractAuthenticationToken> {
    private val defaultGrantedAuthoritiesConverter = JwtGrantedAuthoritiesConverter()

    override fun convert(source: Jwt): AbstractAuthenticationToken {
        val authorities =
            Stream.concat(
                extractResourceAdditionals(source),
                defaultGrantedAuthoritiesConverter.convert(source)!!.stream(),
            ).collect(Collectors.toSet())

        return JwtAuthenticationToken(source, authorities, source.getClaimAsString("preferred_username"))
    }

    private fun extractResourceAdditionals(source: Jwt): Stream<out GrantedAuthority> =
        Stream.concat(
            extractResourceScopes(source),
            extractResourceRoles(source),
        ).map { SimpleGrantedAuthority(it) }

    private fun extractResourceScopes(source: Jwt): Stream<String> {
        val resourceAccess = source.getClaim<Map<String, Any>?>("resource_access")

        val resource = resourceAccess["account"] as Map<String, Any>?
        if (resource != null) {
            val resourceRoles = resource["roles"] as Collection<String>?
            if (resourceRoles != null) {
                return resourceRoles.stream()
            }
        }

        return Stream.empty()
    }

    private fun extractResourceRoles(source: Jwt): Stream<String> {
        val realmAccess = source.getClaim<Map<String, Any>?>("realm_access")

        val resourceRoles = realmAccess["roles"] as Collection<String>?
        if (resourceRoles != null) {
            return resourceRoles.stream().map { "ROLE_$it" }
        }

        return Stream.empty()
    }
}
