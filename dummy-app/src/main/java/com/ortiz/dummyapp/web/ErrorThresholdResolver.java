package com.ortiz.dummyapp.web;

import com.amazonaws.services.simplesystemsmanagement.AWSSimpleSystemsManagement;
import com.amazonaws.services.simplesystemsmanagement.model.GetParametersRequest;
import com.amazonaws.services.simplesystemsmanagement.model.GetParametersResult;
import com.amazonaws.services.simplesystemsmanagement.model.Parameter;

import java.util.Optional;
import java.util.Random;

import static com.ortiz.dummyapp.Constants.PARAMETER_STORE_NAME;

public class ErrorThresholdResolver {

  private String parameterStoreName;
  private AWSSimpleSystemsManagement client;
  private static ThreadLocal<Boolean> flagThrowError = new ThreadLocal<>();


  public ErrorThresholdResolver(AWSSimpleSystemsManagement client) {
    this.client = client;
    this.parameterStoreName = System.getenv(PARAMETER_STORE_NAME);
  }

  public boolean checkIfThrowError() {
    Boolean flagError = flagThrowError.get();
    if (flagError == null) {
      String errorThresholdStr = getParameters(this.parameterStoreName);
      Random rand = new Random(); //instance of random class
      int randomNumber = rand.nextInt(100);
      int errorThreshold = Integer.parseInt(errorThresholdStr);
      flagError = randomNumber < errorThreshold;
      flagThrowError.set(flagError);
    }
    return flagError;
  }

  public void removeLocalData() {
    this.flagThrowError.remove();
  }

  private String getParameters(String parameterName) {
    GetParametersRequest request = new GetParametersRequest();
    request.withNames(parameterName).setWithDecryption(true);
    GetParametersResult parameters = this.client.getParameters(request);
    Optional<Parameter> parameter = parameters.getParameters().stream().findFirst();
    return parameter.get().getValue();
  }
}
