
#/bin/bash

#login
aws --profile wpuser ecr-public get-login-password --region us-east-1 | sudo docker login --username AWS --password-stdin public.ecr.aws/y7c9l5j8


go build -tags netgo -a -v


sudo docker build -f ../../Dockerfile . --tag wg:latest

sudo docker tag wg:latest public.ecr.aws/y7c9l5j8/wg:latest

sudo docker push public.ecr.aws/y7c9l5j8/wg:latest


# sudo docker run -it -p 80:8081 wg:latest

# worker-1-large
# ssh w@192.168.122.74