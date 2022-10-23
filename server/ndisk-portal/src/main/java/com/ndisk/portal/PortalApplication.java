package com.ndisk.portal;


import org.apache.catalina.servlets.WebdavServlet;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.context.annotation.ComponentScan;

@SpringBootApplication
public class PortalApplication {
    public static void main(String[] args) throws Exception {
        SpringApplication app = new SpringApplication(PortalApplication.class);
        app.run(args);
    }
}
