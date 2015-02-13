# agent_check

[user@web03 ~]$ nc localhost 5309  
UP 99%  
[user@web03 ~]$ echo "DRAIN" | nc localhost 8675  
DRAIN OK  
[user@web03 ~]$ nc localhost 5309  
DRAIN 99%  
![and it's drained]  
(https://transfer.sh/n9i1c/shot-2015-02-11t18-22-18.png)

euser@web03 ~]$ echo "READY" | nc localhost 8675  
READY OK  
[user@web03 ~]$ nc localhost 5309  
READY 100%  

![and it's back]
(https://transfer.sh/1fQLgP/shot-2015-02-11t18-25-33.png)

the % reflects the %idle on the machine so this can be used for anything that's being loadbalanced by haproxy
