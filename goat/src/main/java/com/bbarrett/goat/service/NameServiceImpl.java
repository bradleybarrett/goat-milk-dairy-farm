package com.bbarrett.goat.service;

import com.bbarrett.goat.GoatProperties;
import org.apache.commons.lang.StringUtils;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.util.Arrays;

@Service
public class NameServiceImpl
{
    private GoatProperties properties;

    @Autowired
    public NameServiceImpl(GoatProperties properties)
    {
        this.properties = properties;
    }

    public String getName()
    {
        return StringUtils.join(Arrays.asList(
                properties.getAppName(), properties.getAppPort(), properties.getAppVersion()), ":");
    }
}
