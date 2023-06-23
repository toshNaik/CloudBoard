#!/bin/bash
source sample.env
secrets=(access-token-private-key access-token-public-key refresh-token-private-key refresh-token-public-key redis-url db-dsn)
envs=(ACCESS_TOKEN_PRIVATE_KEY ACCESS_TOKEN_PUBLIC_KEY REFRESH_TOKEN_PRIVATE_KEY REFRESH_TOKEN_PUBLIC_KEY REDIS_URL DB_DSN)
len=${#envs[@]}

for ((idx=0; idx<$len; idx++))
do
    SECRET=${secrets[idx]}
    ENV=${envs[idx]}
    
    eval "echo -n \$${ENV}" > temp.txt
    # create secret on gcp
    gcloud secrets create $SECRET --data-file=temp.txt
done

rm temp.txt