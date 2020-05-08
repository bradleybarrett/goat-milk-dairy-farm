package com.bbarrett.farmer.integration.goat;

import com.bbarrett.farmer.healthcheck.HealthCheck;
import com.bbarrett.farmer.integration.RestRequestHelper;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.ResponseEntity;
import org.springframework.stereotype.Service;

@Service("goat")
public class GoatExternalServiceImpl implements GoatExternalService, HealthCheck
{
	private final String serviceName = "lb-goat";
	private RestRequestHelper restRequestHelper;

	@Autowired
	public GoatExternalServiceImpl(RestRequestHelper restRequestHelper)
	{
		this.restRequestHelper = restRequestHelper;
	}

	@Override
	public String getMilk()
	{
		ResponseEntity<String> response = restRequestHelper.getForEntity(serviceName, "/milk", String.class);
		return response.getBody();
	}

	@Override
	public String healthCheck()
	{
		ResponseEntity<String> response = restRequestHelper.getForEntity(serviceName, "/health", String.class);
		return response.getStatusCode().is2xxSuccessful() ? response.getBody() : null;
	}
}
