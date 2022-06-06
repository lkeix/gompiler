docker build -t assemble-linux .
docker run --privileged  --rm -v $PWD:/work -it assemble-linux /bin/bash
