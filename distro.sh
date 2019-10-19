
#run the script from /media/dave/60B6E593B6E569CC/Project/src/github.com/davecourtois/Globular

cd /media/dave/60B6E593B6E569CC/Project/src/github.com/davecourtois/Globular
rm -r ~/distro
mkdir ~/distro
mkdir ~/distro/linux_64
cp Globular ~/distro/linux_64
mkdir ~/distro/linux_64/bin
cp bin/grpcwebproxy ~/distro/linux_64/bin
mkdir ~/distro/linux_64/certs
mkdir ~/distro/linux_64/creds
mkdir ~/distro/linux_64/WebRoot
cp -r WebRoot/css ~/distro/linux_64/WebRoot
cp -r WebRoot/js ~/distro/linux_64/WebRoot
cp WebRoot/index.html ~/distro/linux_64/WebRoot
mkdir ~/distro/linux_64/echo
cp echo/echo_server/echo_server ~/distro/linux_64/echo/echo_server
cp echo/echo_server/config.json ~/distro/linux_64/echo/config.json
mkdir ~/distro/linux_64/event
cp event/event_server/event_server ~/distro/linux_64/event/event_server
cp event/event_server/config.json ~/distro/linux_64/event/config.json
mkdir ~/distro/linux_64/file
cp file/file_server/file_server ~/distro/linux_64/file/file_server
cp file/file_server/config.json ~/distro/linux_64/file/config.json
mkdir ~/distro/linux_64/ldap
cp ldap/ldap_server/ldap_server ~/distro/linux_64/ldap/ldap_server
cp ldap/ldap_server/config.json ~/distro/linux_64/ldap/config.json
mkdir ~/distro/linux_64/monitoring
cp monitoring/monitoring_server/monitoring_server ~/distro/linux_64/monitoring/monitoring_server
cp monitoring/monitoring_server/config.json ~/distro/linux_64/monitoring/config.json
mkdir ~/distro/linux_64/persistence
cp persistence/persistence_server/persistence_server ~/distro/linux_64/persistence/persistence_server
cp persistence/persistence_server/config.json ~/distro/linux_64/persistence/config.json
mkdir ~/distro/linux_64/plc_link
cp plc_link/plc_link_server/plc_link_server ~/distro/linux_64/plc_link/plc_link_server
cp plc_link/plc_link_server/config.json ~/distro/linux_64/plc_link/config.json
mkdir ~/distro/linux_64/smtp
cp smtp/smtp_server/smtp_server ~/distro/linux_64/smtp/smtp_server
cp smtp/smtp_server/config.json ~/distro/linux_64/smtp/config.json
mkdir ~/distro/linux_64/sql
cp sql/sql_server/sql_server ~/distro/linux_64/sql/sql_server
cp sql/sql_server/config.json ~/distro/linux_64/sql/config.json
mkdir ~/distro/linux_64/storage
cp storage/storage_server/storage_server ~/distro/linux_64/storage/storage_server
cp storage/storage_server/config.json ~/distro/linux_64/storage/config.json
mkdir ~/distro/linux_64/plc_exporter
cp plc/plc_exporter/plc_exporter ~/distro/linux_64/plc_exporter/plc_exporter
cp plc/plc_exporter/config.json ~/distro/linux_64/plc_exporter/config.json
mkdir ~/distro/linux_64/plc_server_ab
cp plc/plc_server/plc_server_ab/build-plc_server_ab-Desktop_Qt_5_11_1_GCC_64bit-Release/plc_server_ab ~/distro/linux_64/plc_server_ab/plc_server_ab
cp plc/plc_server/plc_server_ab/build-plc_server_ab-Desktop_Qt_5_11_1_GCC_64bit-Release/config.json ~/distro/linux_64/plc_server_ab/config.json
mkdir ~/distro/linux_64/plc_server_siemens
cp plc/plc_server/plc_server_siemens/plc_server_siemens ~/distro/linux_64/plc_server_siemens/plc_server_siemens
cp plc/plc_server/plc_server_siemens/config.json ~/distro/linux_64/plc_server_siemens/config.json