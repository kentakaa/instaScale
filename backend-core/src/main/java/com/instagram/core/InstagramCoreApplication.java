package com.instagram.core;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.data.mongodb.config.EnableMongoAuditing;

@SpringBootApplication
@EnableMongoAuditing
public class InstagramCoreApplication {
    public static void main(String[] args) {
        SpringApplication.run(InstagramCoreApplication.class, args);
    }
}
