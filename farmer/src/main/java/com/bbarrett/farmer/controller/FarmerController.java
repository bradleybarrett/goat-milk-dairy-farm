package com.bbarrett.farmer.controller;

import com.bbarrett.farmer.healthcheck.HealthCheck;
import com.bbarrett.farmer.healthcheck.HealthCheckWithDependencies;
import com.bbarrett.farmer.service.HealthCheckServiceImpl;
import com.bbarrett.farmer.service.MilkBottle;
import com.bbarrett.farmer.service.MilkServiceImpl;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Qualifier;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController
public class FarmerController
{
    @Autowired
    private HealthCheckWithDependencies farmerHealthCheck;

    @Autowired
    private MilkServiceImpl milkService;

    @GetMapping("/health")
    public String healthCheck()
    {
        return farmerHealthCheck.healthCheck();
    }

    @GetMapping("/health/dependencies")
    public ResponseEntity<String> healthCheckDependencies()
    {
        boolean allHealthy = farmerHealthCheck.healthCheckDependencies();
        HttpStatus responseStatus = allHealthy ? HttpStatus.OK : HttpStatus.INTERNAL_SERVER_ERROR;
        return new ResponseEntity<>(responseStatus);
    }

    @GetMapping("/milk")
    public ResponseEntity<MilkBottle> getMilk()
    {
        MilkBottle milkBottle = milkService.getMilk();
        return new ResponseEntity<>(milkBottle, HttpStatus.OK);
    }
}
