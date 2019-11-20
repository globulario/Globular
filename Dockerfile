 FROM ubuntu
ADD /usr/local/lib/libplctag.so /usr/local/lib
ADD /usr/local/lib/libgrpc++.so /usr/local/lib
ADD /usr/local/lib/libgrpc++.so.1 /usr/local/lib
ADD /usr/local/lib/libgrpc++.so.1.20.0 /usr/local/lib
ADD /usr/local/lib/libprotobuf.so /usr/local/lib
ADD /usr/local/lib/libprotobuf.so.20 /usr/local/lib
ADD /usr/local/lib/libprotobuf.so.20.0.1 /usr/local/lib
ADD /usr/local/lib/libgrpc.so /usr/local/lib
ADD /usr/local/lib/libgrpc.so.7 /usr/local/lib
ADD /usr/local/lib/libgrpc.so.7.0.0 /usr/local/lib
ADD /usr/local/lib/libgpr.so /usr/local/lib
ADD /usr/local/lib/libgpr.so.7 /usr/local/lib
ADD /usr/local/lib/libgpr.so.7.0.0 /usr/local/lib
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
ADD ressource/ressource.proto proto/ressource.proto
ADD admin/admin.proto proto/admin.proto
ADD echo/echo_server/echo_server echo/echo_server
ADD echo/echo_server/config.json echo/config.json
ADD echo/echopb/echo.proto proto/echo.proto
ADD file/file_server/file_server file/file_server
ADD file/file_server/config.json file/config.json
ADD file/filepb/file.proto proto/file.proto
ADD event/event_server/event_server event/event_server
ADD event/event_server/config.json event/config.json
ADD event/eventpb/event.proto proto/event.proto
ADD ldap/ldap_server/ldap_server ldap/ldap_server
ADD ldap/ldap_server/config.json ldap/config.json
ADD ldap/ldappb/ldap.proto proto/ldap.proto
ADD monitoring/monitoring_server/monitoring_server monitoring/monitoring_server
ADD monitoring/monitoring_server/config.json monitoring/config.json
ADD monitoring/monitoringpb/monitoring.proto proto/monitoring.proto
ADD persistence/persistence_server/persistence_server persistence/persistence_server
ADD persistence/persistence_server/config.json persistence/config.json
ADD persistence/persistencepb/persistence.proto proto/persistence.proto
ADD plc/plc_exporter/plc_exporter plc_exporter/plc_exporter
ADD plc/plc_exporter/config.json plc_exporter/config.json
ADD plc_link/plc_link_server/plc_link_server plc_link/plc_link_server
ADD plc_link/plc_link_server/config.json plc_link/config.json
ADD plc_link/plc_linkpb/plc_link.proto proto/plc_link.proto
ADD plc/plc_server_ab/build-plc_server_ab-Desktop_Qt_5_11_1_GCC_64bit-Release/plc_server_ab plc_server_ab/plc_server_ab
ADD plc/plc_server_ab/build-plc_server_ab-Desktop_Qt_5_11_1_GCC_64bit-Release/config.json plc_server_ab/config.json
ADD plc/plc_server_siemens/plc_server_siemens plc_server_siemens/plc_server_siemens
ADD plc/plc_server_siemens/config.json plc_server_siemens/config.json
ADD plc/plcpb/plc.proto proto/plc.proto
ADD smtp/smtp_server/smtp_server smtp/smtp_server
ADD smtp/smtp_server/config.json smtp/config.json
ADD smtp/smtppb/smtp.proto proto/smtp.proto
ADD sql/sql_server/sql_server sql/sql_server
ADD sql/sql_server/config.json sql/config.json
ADD sql/sqlpb/sql.proto proto/sql.proto
ADD storage/storage_server/storage_server storage/storage_server
ADD storage/storage_server/config.json storage/config.json
ADD storage/storagepb/storage.proto proto/storage.proto
#CMD /globular/Globular

