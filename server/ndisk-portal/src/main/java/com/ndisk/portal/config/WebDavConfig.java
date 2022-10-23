package com.ndisk.portal.config;

import com.ndisk.webdav.servlet.WebDavServlet;
import org.springframework.beans.factory.annotation.Configurable;
import org.springframework.boot.web.servlet.ServletRegistrationBean;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

@Configuration
public class WebDavConfig {

    @Bean
    ServletRegistrationBean<WebDavServlet> webDavServletServletRegistrationBean() {
        ServletRegistrationBean<WebDavServlet> servletRegistrationBean = new ServletRegistrationBean<>(new WebDavServlet(), "/dav/*");
        return servletRegistrationBean;
    }
}
