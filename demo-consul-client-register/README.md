
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
