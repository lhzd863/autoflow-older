ps -ef|grep conf.yaml|grep -v grep|awk -F ' ' '{print $2}'|xargs kill -9

cd /home/k8s/Go/ws/src/github.com/lhzd863/autoflow/apiserver
bash start.sh

cd /home/k8s/Go/ws/src/github.com/lhzd863/autoflow/slv/slv-0
bash start.sh

cd /home/k8s/Go/ws/src/github.com/lhzd863/autoflow/slv/slv-1
bash start.sh

cd /home/k8s/Go/ws/src/github.com/lhzd863/autoflow/slv/slv-2
bash start.sh


cd /home/k8s/Go/ws/src/github.com/lhzd863/autoflow/mst/mst-0
bash start.sh


cd /home/k8s/Go/ws/src/github.com/lhzd863/autoflow/mst/mst-1
bash start.sh


cd /home/k8s/Go/ws/src/github.com/lhzd863/autoflow/mst/mst-2
bash start.sh




