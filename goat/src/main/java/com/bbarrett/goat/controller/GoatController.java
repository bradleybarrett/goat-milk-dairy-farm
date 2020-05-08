package com.bbarrett.goat.controller;

import com.bbarrett.goat.healthcheck.HealthCheck;
import com.bbarrett.goat.service.HealthCheckServiceImpl;
import com.bbarrett.goat.service.NameServiceImpl;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController
public class GoatController
{
    @Autowired
    private HealthCheck goatHealthCheck;

    @Autowired
    private NameServiceImpl nameService;

    @GetMapping("/health")
    public String healthCheck()
    {
        return goatHealthCheck.healthCheck();
    }

    @GetMapping("/milk")
    public ResponseEntity<String> getMilk()
    {
        String name = nameService.getName();
        return new ResponseEntity<>(name, HttpStatus.OK);
    }
}
