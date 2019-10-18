FROM ubuntu
ADD libplctag.so /usr/local/lib
ADD libgrpc++.so /usr/local/lib
ADD libgrpc++.so.1 /usr/local/lib
ADD libgrpc++.so.1.20.0 /usr/local/lib
ADD libprotobuf.so /usr/local/lib
ADD libprotobuf.so.20 /usr/local/lib
ADD libprotobuf.so.20.0.1 /usr/local/lib
ADD libgrpc.so /usr/local/lib
ADD libgrpc.so.7 /usr/local/lib
ADD libgrpc.so.7.0.0 /usr/local/lib
ADD libgpr.so /usr/local/lib
ADD libgpr.so.7 /usr/local/lib
ADD libgpr.so.7.0.0 /usr/local/lib
RUN apt-get update && apt-get install -y gnupg2 \
    wget \
  && rm -rf /var/lib/apt/lists/*
RUN wget https://s3-eu-west-1.amazonaws.com/deb.robustperception.io/41EFC99D.gpg && apt-key add 41EFC99D.gpg
RUN apt-get update && apt-get install -y \
  build-essential \
  curl \
  mongodb \
  prometheus 
RUN curl http://www.unixodbc.org/unixODBC-2.3.7.tar.gz --output unixODBC-2.3.7.tar.gz
RUN tar -xvf unixODBC-2.3.7.tar.gz
RUN rm unixODBC-2.3.7.tar.gz
WORKDIR unixODBC-2.3.7
RUN ./configure && make all install clean && ldconfig && mkdir /globular && cd /globular
WORKDIR /globular
ADD Globular /globular
ADD bin/grpcwebproxy bin/grpcwebproxy
ADD proto/ressource.proto proto/ressource.proto
ADD proto/admin.proto proto/admin.proto
ADD proto/echo.proto proto/echo.proto
ADD proto/event.proto proto/event.proto
ADD proto/file.proto proto/file.proto
ADD proto/ldap.proto proto/ldap.proto
ADD proto/monitoring.proto proto/monitoring.proto
ADD proto/persistence.proto proto/persistence.proto
ADD proto/plc.proto proto/plc.proto
ADD proto/plc_link.proto proto/plc_link.proto
ADD proto/smtp.proto proto/smtp.proto
ADD proto/sql.proto proto/sql.proto
ADD proto/storage.proto proto/storage.proto
ADD echo_server/echo_server echo/echo_server
ADD echo_server/config.json echo/config.json
ADD file_server/file_server file/file_server
ADD file_server/config.json file/config.json
ADD event_server/event_server event/event_server
ADD event_server/config.json event/config.json
ADD ldap_server/ldap_server ldap/ldap_server
ADD ldap_server/config.json ldap/config.json
ADD monitoring_server/monitoring_server monitoring/monitoring_server
ADD monitoring_server/config.json monitoring/config.json
ADD persistence_server/persistence_server persistence/persistence_server
ADD persistence_server/config.json persistence/config.json
ADD plc_exporter/plc_exporter plc_exporter/plc_exporter
ADD plc_exporter/config.json plc_exporter/config.json
ADD plc_link_server/plc_link_server plc_link/plc_link_server
ADD plc_link_server/config.json plc_link/config.json
ADD plc_server_ab/plc_server_ab plc_server_ab/plc_server_ab
ADD plc_server_ab/config.json plc_server_ab/config.json
ADD plc_server_siemens/plc_server_siemens plc_server_siemens/plc_server_siemens
ADD plc_server_siemens/config.json plc_server_siemens/config.json
ADD smtp_server/smtp_server smtp/smtp_server
ADD smtp_server/config.json smtp/config.json
ADD sql_server/sql_server sql/sql_server
ADD sql_server/config.json sql/config.json
ADD storage_server/storage_server storage/storage_server
ADD storage_server/config.json storage/config.json
#CMD /globular/Globular

