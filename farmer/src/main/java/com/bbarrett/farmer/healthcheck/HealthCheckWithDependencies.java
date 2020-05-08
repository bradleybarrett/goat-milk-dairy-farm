package com.bbarrett.farmer.healthcheck;

public interface HealthCheckWithDependencies extends HealthCheck
{
	boolean healthCheckDependencies();
}
