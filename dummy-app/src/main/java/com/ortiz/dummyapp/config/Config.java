package com.ortiz.dummyapp.config;

import com.amazonaws.client.builder.AwsClientBuilder;
import com.amazonaws.services.simplesystemsmanagement.AWSSimpleSystemsManagement;
import com.amazonaws.services.simplesystemsmanagement.AWSSimpleSystemsManagementClientBuilder;
import com.ortiz.dummyapp.web.ErrorGenerationInterceptor;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.web.servlet.config.annotation.InterceptorRegistry;
import org.springframework.web.servlet.config.annotation.WebMvcConfigurer;

@Configuration
public class Config implements WebMvcConfigurer {

  @Value("${cloud.aws.region.static}")
  private String awsRegion;

  @Value("${ssm.service-endpoint}")
  private String ssmServiceEndpoint;

  @Bean
  public AWSSimpleSystemsManagement ssmClient() {
    AWSSimpleSystemsManagement client = AWSSimpleSystemsManagementClientBuilder.standard().withEndpointConfiguration(
      new AwsClientBuilder.EndpointConfiguration(ssmServiceEndpoint, awsRegion)).build();
    return client;
  }

  @Override
  public void addInterceptors(InterceptorRegistry registry) {
    registry.addInterceptor(new ErrorGenerationInterceptor(ssmClient())).addPathPatterns("/**");
    // .excludePathPatterns("/actuator/**");
  }
}
