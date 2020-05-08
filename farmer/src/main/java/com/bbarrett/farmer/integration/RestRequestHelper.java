package com.bbarrett.farmer.integration;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.ResponseEntity;
import org.springframework.stereotype.Service;
import org.springframework.web.client.RestTemplate;

@Service
public class RestRequestHelper
{
    @Autowired
    private RestTemplate restTemplate;

    public <T> ResponseEntity<T> getForEntity(String serviceName, String path, Class<T> responseType)
    {
        return restTemplate.getForEntity(getUrl("http", serviceName, path), responseType);
    }

    private String getUrl(String protocol, String serviceName, String path)
    {
        return protocol + "://" + serviceName + path;
    }
}
