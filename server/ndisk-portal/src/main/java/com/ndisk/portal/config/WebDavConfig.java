package com.ndisk.portal.config;

import com.ndisk.webdav.servlet.WebdavServlet;
import org.springframework.boot.web.servlet.ServletRegistrationBean;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

@Configuration
public class WebDavConfig {

    @Bean
    ServletRegistrationBean<WebdavServlet> webDavServletServletRegistrationBean() {
        ServletRegistrationBean<WebdavServlet> servletRegistrationBean = new ServletRegistrationBean<>(new WebdavServlet(), "/dav/*");
        return servletRegistrationBean;
    }
}
