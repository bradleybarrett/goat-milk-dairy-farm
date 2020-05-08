package com.bbarrett.farmer.service;

import com.bbarrett.farmer.FarmerProperties;
import org.apache.commons.lang.StringUtils;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.util.Arrays;

@Service
public class NameServiceImpl
{
    private FarmerProperties properties;

    @Autowired
    public NameServiceImpl(FarmerProperties properties)
    {
        this.properties = properties;
    }

    public String getName()
    {
        return StringUtils.join(Arrays.asList(
                properties.getAppName(), properties.getAppPort(), properties.getAppVersion()), ":");
    }
}
