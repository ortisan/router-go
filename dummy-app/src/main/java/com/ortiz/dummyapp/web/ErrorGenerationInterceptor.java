package com.ortiz.dummyapp.web;

import com.ortiz.dummyapp.CustomException;
import org.springframework.web.servlet.HandlerInterceptor;

import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;

/**
 * This class is just to read parameter with number that indicates % of errors.
 * Try to simulate errors with api.
 */
public class ErrorGenerationInterceptor implements HandlerInterceptor {

  private ErrorThresholdResolver errorThresholdResolver;

  public ErrorGenerationInterceptor(ErrorThresholdResolver errorThresholdResolver) {
    this.errorThresholdResolver = errorThresholdResolver;
  }

  @Override
  public boolean preHandle(HttpServletRequest request, HttpServletResponse response, Object handler) throws Exception {
    if (this.errorThresholdResolver.checkIfThrowError()) {
      throw new CustomException("Forcing error to test router redirect.");
    }
    return true;
  }

  @Override
  public void afterCompletion(HttpServletRequest request, HttpServletResponse response, Object handler, Exception ex) throws Exception {
    this.errorThresholdResolver.removeLocalData();
  }
}
