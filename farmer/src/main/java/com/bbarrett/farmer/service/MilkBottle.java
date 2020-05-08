package com.bbarrett.farmer.service;

public class MilkBottle
{
    private String farmer;
    private String goat;

    public MilkBottle(String farmer, String goat)
    {
        this.farmer = farmer;
        this.goat = goat;
    }

    public String getFarmer()
    {
        return farmer;
    }

    public void setFarmer(String farmer)
    {
        this.farmer = farmer;
    }

    public String getGoat()
    {
        return goat;
    }

    public void setGoat(String goat)
    {
        this.goat = goat;
    }
}
