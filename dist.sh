
#package the files
mkdir dist
mkdir dist/globular

#complile client services
cd client
npx webpack
mv dist/services.js ../WebRoot/js
cd ../

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
#event service
mkdir dist/globular/event
cp event/event_server/event_server dist/globular/event
cp event/event_server/event_server.exe dist/globular/event
cp event/event_server/config.json dist/globular/event

#now I will zip the dist/globular file
cd dist
tar -zcvf globular.1.0.tar.gz globular
sudo cp globular.1.0.tar.gz /tmp

#remove the globular dir
rm -r globular
# recreate it empty
mkdir globular
mkdir globular/WebRoot
cd ../

# I will compile the website with polymer.
cd WebRoot/website
#remove previous build
rm -r build
polymer build
cd ../../
cp -r WebRoot/website/build/default/* dist/globular/WebRoot
cp -r sslforfree dist/globular
mkdir dist/globular/WebRoot/image
cp -r WebRoot/website/image/* dist/globular/WebRoot/image

# move configuration from document. 
cp /home/dave/Documents/config/config.json dist/globular/WebRoot/config.json
cp -r WebRoot/website/build/default/node_modules dist/globular/WebRoot
mkdir dist/globular/event
mkdir dist/globular/echo
mkdir dist/globular/file
mkdir dist/globular/ldap
mkdir dist/globular/sql
mkdir dist/globular/persistence
mkdir dist/globular/storage
mkdir dist/globular/smtp
cp /home/dave/Documents/config/event/config.json dist/globular/event
cp /home/dave/Documents/config/echo/config.json dist/globular/echo
cp /home/dave/Documents/config/file/config.json dist/globular/file
cp /home/dave/Documents/config/ldap/config.json dist/globular/ldap
cp /home/dave/Documents/config/sql/config.json dist/globular/sql
cp /home/dave/Documents/config/persistence/config.json dist/globular/persistence
cp /home/dave/Documents/config/storage/config.json dist/globular/storage
cp /home/dave/Documents/config/smtp/config.json dist/globular/smtp
# set the dist folder to give acces to binary distribution of globular.
#sudo cp  /tmp/globular.1.0.tar.gz dist/globular/WebRoot/dist

cd dist
#the website will be merge with base globular on ec2
tar -zcvf website.tar.gz globular
sudo cp website.tar.gz /tmp

#remove the dist folder
cd ../
rm -r dist

#copy it to ec2
scp -i ~/Globular.app.pem /tmp/globular.1.0.tar.gz  ubuntu@ec2-34-214-248-201.us-west-2.compute.amazonaws.com:~/
scp -i ~/Globular.app.pem /tmp/website.tar.gz  ubuntu@ec2-34-214-248-201.us-west-2.compute.amazonaws.com:~/
