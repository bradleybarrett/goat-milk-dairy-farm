package com.bbarrett.farmer.service;

import com.bbarrett.farmer.FarmerProperties;
import com.bbarrett.farmer.integration.goat.GoatExternalServiceImpl;
import org.apache.commons.lang.StringUtils;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.util.Arrays;

@Service
public class MilkServiceImpl
{
    private FarmerProperties farmerProperties;

    private GoatExternalServiceImpl goatExternalService;

    @Autowired
    public MilkServiceImpl(FarmerProperties farmerProperties, GoatExternalServiceImpl goatExternalService)
    {
        this.farmerProperties = farmerProperties;
        this.goatExternalService = goatExternalService;
    }

    public MilkBottle getMilk()
    {
        String goat = goatExternalService.getMilk();
        String farmer = StringUtils.join(Arrays.asList(
                farmerProperties.getAppName(),
                farmerProperties.getAppPort(),
                farmerProperties.getAppVersion()),
                ":");

        return new MilkBottle(farmer, goat);
    }
}
