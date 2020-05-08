package com.bbarrett.goat;

import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Service;

@Service
public class GoatProperties
{
    @Value("${host.port}")
    private int appPort;

    @Value("${spring.application.name}")
    private String appName;

    @Value("${version}")
    private String appVersion;

    public int getAppPort()
    {
        return appPort;
    }

    public String getAppName()
    {
        return appName;
    }

    public String getAppVersion()
    {
        return appVersion;
    }
}
