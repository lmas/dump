#!/bin/bash

# path to iptables
IPT="/sbin/iptables"
IPT6="/sbin/ip6tables"

# Localhost Interface
LO_IFACE="lo"

###############################################################################

# Disable and block all IPv6
echo 1 > /proc/sys/net/ipv6/conf/all/disable_ipv6

# flush all rules
$IPT6 -F
#$IPT6 -t nat -F
$IPT6 -t mangle -F

# erase all non-default chains
$IPT6 -X
#$IPT6 -t nat -X
$IPT6 -t mangle -X

# set default policies
$IPT6 -P INPUT DROP
$IPT6 -P FORWARD DROP
$IPT6 -P OUTPUT DROP

###############################################################################

# flush all rules
$IPT -F
$IPT -t nat -F
$IPT -t mangle -F

# erase all non-default chains
$IPT -X
$IPT -t nat -X
$IPT -t mangle -X

# set default policies
$IPT -P INPUT DROP
$IPT -P FORWARD DROP
$IPT -P OUTPUT ACCEPT

# allow localhost
$IPT -A INPUT -i lo -j ACCEPT
$IPT -A OUTPUT -o lo -j ACCEPT

# enable traffick accounting
$IPT -A INPUT -i eth0 -m comment --comment 'TRAFFICK_IN'
$IPT -A OUTPUT -o eth0 -m comment --comment 'TRAFFICK_OUT'

# block invalid packets
#$IPT -A INPUT -m conntrack --ctstate INVALID -j LOG --log-prefix "iptables: INVALID "
$IPT -A INPUT -m conntrack --ctstate INVALID -j DROP
$IPT -A OUTPUT -m conntrack --ctstate INVALID -j DROP

# accept established connections
$IPT -A INPUT -m conntrack --ctstate ESTABLISHED,RELATED -j ACCEPT
$IPT -A OUTPUT -m conntrack --ctstate ESTABLISHED,RELATED -j ACCEPT

# accept dhcp
#$IPT -A INPUT -d 255.255.255.255 -p udp -m multiport --ports 67,68 -j ACCEPT

# block spammy broadcasts
#$IPT -A INPUT -d 192.168.1.255 -j LOG --log-prefix "iptables: BROADCAST "
#$IPT -A INPUT -d 192.168.1.255 -j DROP
#$IPT -A INPUT -d 255.255.255.255 -j LOG --log-prefix "iptables: BROADCAST "
$IPT -A INPUT -d 255.255.255.255 -j DROP

# block "new not syn" packets
#$IPT -A INPUT -p tcp ! --syn -m conntrack --ctstate NEW -j LOG --log-prefix "iptables: NEW_NOT_SYN "
$IPT -A INPUT -p tcp ! --syn -m conntrack --ctstate NEW -j DROP

# block fragmented ICMP packets
#$IPT -A INPUT -p ICMP --fragment -j LOG --log-prefix "iptables: ICMP_FRAG "
$IPT -A INPUT -p ICMP --fragment -j DROP

###############################################################################
# CUSTOM IN/OUT PORTS

# web
#$IPT -t nat -A PREROUTING  -i eth0 -p tcp --dport 80 -j REDIRECT --to-port 8000
#$IPT -A INPUT -p tcp --dport 8000 -j ACCEPT

# ssh
$IPT -A INPUT -p tcp --dport 8822 -j ACCEPT

###############################################################################

# log unwated incomming packets
#$IPT -A INPUT -j LOG --log-prefix "iptables: BAD_IN "
#$IPT -A OUTPUT -j LOG --log-prefix "iptables: BAD_OUT "
