#!/bin/bash

sudo ifconfig wlan0 up
sudo iwconfig wlan0 essid ardrone2_280143
sudo ifconfig wlan0 192.168.1.3
