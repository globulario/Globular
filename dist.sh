rm -r dist
mkdir dist
mkdir dist/globular

#Globular general file
mkdir dist/globular/bin
cp -r bin dist/globular
cp Globular dist/globular/
cp Globular.exe dist/globular/
mkdir dist/globular/WebRoot
mkdir dist/globular/WebRoot/js
mkdir dist/globular/WebRoot/css
cp WebRoot/js/services.js dist/globular/WebRoot/js
cp WebRoot/js/test.js dist/globular/WebRoot/js
cp WebRoot/css/styles.css dist/globular/WebRoot/css
cp WebRoot/index.html dist/globular/WebRoot
cp WebRoot/config.json dist/globular/WebRoot
#echo service
mkdir dist/globular/echo
cp echo/echo_server/echo_server dist/globular/echo
cp echo/echo_server/echo_server.exe dist/globular/echo
cp echo/echo_server/config.json dist/globular/echo
#file service
mkdir dist/globular/file
cp file/file_server/file_server dist/globular/file
cp file/file_server/file_server.exe dist/globular/file
cp file/file_server/config.json dist/globular/file
#ldap service
mkdir dist/globular/ldap
cp ldap/ldap_server/ldap_server dist/globular/ldap
cp ldap/ldap_server/ldap_server.exe dist/globular/ldap
cp ldap/ldap_server/config.json dist/globular/ldap
#sql service
mkdir dist/globular/sql
cp sql/sql_server/sql_server dist/globular/sql
cp sql/sql_server/sql_server.exe dist/globular/sql
cp sql/sql_server/config.json dist/globular/sql
#persistence service
mkdir dist/globular/persistence
cp persistence/persistence_server/persistence_server dist/globular/persistence
cp persistence/persistence_server/persistence_server.exe dist/globular/persistence
cp persistence/persistence_server/config.json dist/globular/persistence
#storage service
mkdir dist/globular/storage
cp storage/storage_server/storage_server dist/globular/storage
cp storage/storage_server/storage_server.exe dist/globular/storage
cp storage/storage_server/config.json dist/globular/storage
#smtp service
mkdir dist/globular/smtp
cp smtp/smtp_server/smtp_server dist/globular/smtp
cp smtp/smtp_server/smtp_server.exe dist/globular/smtp
cp smtp/smtp_server/config.json dist/globular/smtp
#now I will zip the dist/globular file
cd dist
tar -zcvf globular.1.0.tar.gz globular
cp globular.1.0.tar.gz /tmp
cd ../
rm -r dist