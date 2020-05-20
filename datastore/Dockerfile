FROM elasticsearch:7.7.0

RUN mkdir /data && chown -R elasticsearch:elasticsearch /data \ 
    && echo 'path.data: /data' >> config/elasticsearch.yml

COPY wait-for-it.sh /usr/local/bin/wait-for-it.sh

# Copy the files you may need and your insert script
COPY index.json /opt/elasticsearch/index.json
COPY init.sh /opt/elasticsearch/init.sh

ENV discovery.type single-node
RUN /usr/local/bin/docker-entrypoint.sh elasticsearch -p /tmp/epid & /bin/bash /usr/local/bin/wait-for-it.sh -t 0 localhost:9200 -- /opt/elasticsearch/init.sh; kill $(cat /tmp/epid) && wait $(cat /tmp/epid); exit 0;