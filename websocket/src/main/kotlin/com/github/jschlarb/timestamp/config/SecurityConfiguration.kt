package com.github.jschlarb.timestamp.config

import org.springframework.context.annotation.Bean
import org.springframework.context.annotation.Configuration
import org.springframework.http.HttpMethod
import org.springframework.security.config.Customizer
import org.springframework.security.config.annotation.web.builders.HttpSecurity
import org.springframework.security.config.annotation.web.configuration.EnableWebSecurity
import org.springframework.security.config.http.SessionCreationPolicy
import org.springframework.security.web.SecurityFilterChain
import org.springframework.web.cors.CorsConfiguration
import org.springframework.web.cors.CorsConfigurationSource
import org.springframework.web.cors.UrlBasedCorsConfigurationSource


@Configuration
@EnableWebSecurity
class SecurityConfiguration(val jwtConverter: CustomJwtAuthenticationConverter) {

    @Bean
    // completely disable cors ... it's just for dev and illustration purpose
    fun corsConfigurationSource(): CorsConfigurationSource {
        val configuration = CorsConfiguration()
        configuration.allowedMethods = listOf(
            HttpMethod.GET.name(),
            HttpMethod.HEAD.name(),
            HttpMethod.POST.name(),
            HttpMethod.OPTIONS.name()
        );
        configuration.allowCredentials = true
        configuration.setAllowedOriginPatterns(listOf("*"))
        configuration.allowedHeaders = listOf(
            "X-Requested-With",
            "Origin",
            "Content-Type",
            "Accept",
            "Authorization"
        );
        val source = UrlBasedCorsConfigurationSource()
        source.registerCorsConfiguration("/**", configuration)
        return source
    }

    @Bean
    fun configure(http: HttpSecurity): SecurityFilterChain {
        http
            .httpBasic { it.disable() }
            .csrf { it.disable() }
            .cors { }
            .sessionManagement { it.sessionCreationPolicy(SessionCreationPolicy.STATELESS) }
            .authorizeHttpRequests { authorizeRequests ->
                authorizeRequests
                    .requestMatchers("/timestamps/**").permitAll()
                    .requestMatchers(HttpMethod.OPTIONS).permitAll()
                    .requestMatchers("/actuator/**").permitAll()
                    .requestMatchers("/who-am-i/reader").hasRole("read-action")
                    .requestMatchers("/who-am-i/writer").hasRole("write-action")
                    .requestMatchers("/who-am-i/admin").hasRole("admin-action")
                    .anyRequest().authenticated()
            }
            .oauth2ResourceServer {
                // todo
                it.jwt { it.jwtAuthenticationConverter(jwtConverter) }
            }
            .oauth2Client(Customizer.withDefaults())
        return http.build()
    }
}
