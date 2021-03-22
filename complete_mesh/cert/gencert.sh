#rm *.pem
#
#openssl genrsa -out server.key 2048 #-subj "/C=AT/ST=Vienna/L=Vienna/O=TGM/OU=Education/CN=*.TGM.guru/emailAddress=jlaback@student.tgm.ac.at"
#
#openssl req -new -x509 -sha256 -key server.key \
#              -out server.crt -days 3650 -subj "/C=AT/ST=Vienna/L=Vienna/O=TGM/OU=Education/CN=*.TGM.guru/emailAddress=jlaback@student.tgm.ac.at"
#
#openssl req -new -sha256 -key server.key -out server.csr -subj "/C=AT/ST=Vienna/L=Vienna/O=TGM/OU=Education/CN=*.TGM.guru/emailAddress=jlaback@student.tgm.ac.at"
#
#openssl x509 -req -sha256 -in server.csr -signkey server.key \
#               -out server.crt -days 3650
#
#certstrap init --common-name "CertAuth"
#
#certstrap request-cert --common-name ServiceMesh.ac.at
#
#certstrap sign ServiceMesh.ac.at --CA CertAuth

rm *.pem
openssl genrsa -out server.key
openssl req -new -sha256 -key server.key -out server.csr -subj "/C=AT/ST=Vienna/L=Vienna/O=TGM/OU=Education/CN=*.TGM.guru/emailAddress=jlaback@student.tgm.ac.at"
openssl x509 -req -days 3650 -in server.csr -out server.crt -signkey server.key

openssl genrsa -out server.key
openssl req -new -sha256 -key server.key -out server.csr -subj "/C=AT/ST=Vienna/L=Vienna/O=TGM/OU=Education/CN=*.TGM.guru/emailAddress=jlaback@student.tgm.ac.at"
openssl x509 -req -days 3650 -in server.csr -out server.crt -signkey server.key # -subj "/C=AT/ST=Vienna/L=Vienna/O=TGM/OU=Education/CN=*.TGM.guru/emailAddress=jlaback@student.tgm.ac.at"