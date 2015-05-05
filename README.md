# agent_check

Inspired by [this post](http://blog.loadbalancer.org/open-source-windows-service-for-reporting-server-load-back-to-haproxy-load-balancer-feedback-agent/)

Responds to telnet request with CPU Idle %. 
Also opens up a command channel on localhost that will accept 
management commands (i.e. Drain, Stop, Ready) Can be used with [HAProxy’s](http://www.haproxy.org/) 
agent-check to weight traffic based on CPU load.  
Note: The agent check functionality was added as of version 1.5

[user@aServer ~]$ nc localhost 5309  
UP 99%  
[user@aServer ~]$ echo "DRAIN" | nc localhost 8675  
DRAIN OK  
[user@aServer ~]$ nc localhost 5309  
DRAIN 99%  

user@aServer ~]$ echo "READY" | nc localhost 8675  
READY OK  
[user@aServer ~]$ nc localhost 5309  
READY 100%  


