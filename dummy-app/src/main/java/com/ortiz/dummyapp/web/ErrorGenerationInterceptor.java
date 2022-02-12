package com.ortiz.dummyapp.web;

import com.amazonaws.services.simplesystemsmanagement.AWSSimpleSystemsManagement;
import com.amazonaws.services.simplesystemsmanagement.model.GetParametersRequest;
import com.amazonaws.services.simplesystemsmanagement.model.GetParametersResult;
import com.amazonaws.services.simplesystemsmanagement.model.Parameter;
import org.springframework.http.HttpStatus;
import org.springframework.web.server.ResponseStatusException;
import org.springframework.web.servlet.HandlerInterceptor;

import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;
import java.util.Optional;
import java.util.Random;

import static com.ortiz.dummyapp.Constants.PARAMETER_STORE_NAME;

public class ErrorGenerationInterceptor implements HandlerInterceptor {

  private String parameterStoreName;
  private AWSSimpleSystemsManagement client;

  public ErrorGenerationInterceptor(AWSSimpleSystemsManagement client) {
    this.client = client;
    this.parameterStoreName = System.getenv(PARAMETER_STORE_NAME);
  }

  @Override
  public boolean preHandle(HttpServletRequest request, HttpServletResponse response, Object handler) throws Exception {
    String errorThresholdStr = getParameters(this.parameterStoreName);
    Random rand = new Random(); //instance of random class
    int randomNumber = rand.nextInt(100);
    int errorThreshold = Integer.parseInt(errorThresholdStr);
    if (randomNumber <= errorThreshold) {
      throw new ResponseStatusException(HttpStatus.INTERNAL_SERVER_ERROR, "Forcing error to test router redirect.");
    }
    return true;
  }

  private String getParameters(String parameterName) {
    GetParametersRequest request = new GetParametersRequest();
    request.withNames(parameterName).setWithDecryption(true);
    GetParametersResult parameters = this.client.getParameters(request);
    Optional<Parameter> parameter = parameters.getParameters().stream().findFirst();
    return parameter.get().getValue();
  }
}
