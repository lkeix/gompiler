docker build -t assemble-linux .
docker run --rm -v $PWD:/work -it assemble-linux  /bin/bash
