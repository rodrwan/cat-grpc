#/bin/bash
#
eval $(docker-machine env)
docker run -i -t -p 10000:10000 --name catty-server catty
