package com.bbarrett.farmer.service;

import com.bbarrett.farmer.healthcheck.HealthCheck;
import com.bbarrett.farmer.healthcheck.HealthCheckWithDependencies;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Qualifier;
import org.springframework.stereotype.Service;

import java.util.Collections;
import java.util.List;

@Service
public class HealthCheckServiceImpl implements HealthCheckWithDependencies
{
    @Autowired
    private NameServiceImpl nameService;

    @Qualifier("goat")
    @Autowired
    private HealthCheck goatHealthCheck;

    @Override
    public String healthCheck()
    {
        return nameService.getName();
    }

    @Override
    public boolean healthCheckDependencies()
    {
        List<HealthCheck> healthChecks = Collections.singletonList(goatHealthCheck);
        return healthChecks.stream().allMatch(healthCheck -> healthCheck.healthCheck() != null);
    }
}
