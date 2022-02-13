package com.ortiz.dummyapp.web;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.actuate.health.Health;
import org.springframework.boot.actuate.health.HealthIndicator;
import org.springframework.stereotype.Component;

@Component
public class CustomHealthCheck implements HealthIndicator {

  @Autowired
  private ErrorThresholdResolver errorThresholdResolver;

  @Override
  public Health health() {
    Health.Builder status = Health.up();
    if (errorThresholdResolver.checkIfThrowError()) {
      status.down();
    }
    return status.build();
  }
}
