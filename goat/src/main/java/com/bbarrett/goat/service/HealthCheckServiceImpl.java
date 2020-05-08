package com.bbarrett.goat.service;

import com.bbarrett.goat.healthcheck.HealthCheck;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

@Service
public class HealthCheckServiceImpl implements HealthCheck
{
    @Autowired
    private NameServiceImpl nameService;

    public String healthCheck()
    {
        return nameService.getName();
    }
}
