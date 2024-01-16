
#run traefik
sudo ./traefik  --configFile=./traefik.yml

#run consul
./consul agent -server -ui -node=server-1 -bootstrap-expect=1 -client="0.0.0.0"   -bind 192.168.1.237 -data-dir ./consuldir

#run loadbalance server instance
GOPROXY=https://goproxy.cn go run main.go --port=10102
GOPROXY=https://goproxy.cn go run main.go --port=10103

#check loadbalance
curl  http://127.0.0.1:80/svcwhoami/greeting/v1/hello/

#dashboard
http://192.168.1.237:8080/dashboard/#/
http://192.168.1.237:8500/ui/dc1/services


- consul do health check, send all available svcs to traefik on each config update, which means consul fileter failed svc before notify traefik

- same routers name can only set to one rule, no duplicated rule,you can set new rule on new router(with new router name)

- treafik log like that if register ok
```
DEBU[2025-01-11T16:48:20+08:00] Configuration received: {"http":{"routers":{"svcwhoami1":{"service":"svcwhoami","rule":"PathPrefix(`/sdm/monkey/`)"}},"services":{"svcwhoami":{"loadBalancer":{"servers":[{"url":"http://100.109.57.39:10103"},{"url":"http://172.17.0.1:10103"},{"url":"http://192.168.1.237:10103"},{"url":"http://192.168.108.1:10103"},{"url":"http://192.168.122.1:10103"},{"url":"http://192.168.49.1:10103"},{"url":"http://192.168.84.1:10103"},{"url":"http://[fd05:aac9:27c7:0:44ba:490a:2c0b:fb92]:10103"},{"url":"http://[fd63:c525:e829:0:4dc2:b780:5f64:a059]:10103"},{"url":"http://[fd7a:115c:a1e0:ab12:4843:cd96:626d:3927]:10103"}],"passHostHeader":true}}}},"tcp":{},"udp":{}}  providerName=consulcatalog
```
