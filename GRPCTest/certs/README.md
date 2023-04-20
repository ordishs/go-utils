# Creating a certificate authority

### Generate a private key for the certificate authority:
```
openssl genrsa -out ca.key 2048

   # to create a key with a password, add the -des3 option:
   # openssl genrsa -des3 -out ca.key 2048
```

### Generate a self-signed certificate for the CA:
```
openssl req -new -x509 -key ca.key -out ca.crt -subj "/CN=my-ca"
```

### Generate a private key for the server:
```
openssl genrsa -out server.key 2048
```

### Generate a certificate signing request (CSR) for the server:
```
openssl req -new -key server.key -out server.csr -subj "/CN=localhost"
```

### Sign the server's CSR with the CA's private key to create a signed certificate:
```
openssl x509 -req -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out server.crt
```

Please note that the certificate in this example is bound to the common name of "localhost".  This means that any client will have to disable hostname validation by setting ```InsecureSkipVerify: true``` in the TLS config.


### Create certificate request for client 1...
```
openssl genrsa -out client1.key 2048
openssl req -new -key client1.key -out client1.csr -subj "/CN=localhost"
openssl x509 -req -in client1.csr -CA ca.crt -CAkey ca.key -days 3560 -CAcreateserial -out client1.crt
```

### Create certificate request for client 2...
```
openssl genrsa -out client2.key 2048
openssl req -new -key client2.key -out client2.csr -subj "/CN=localhost"
openssl x509 -req -in client2.csr -CA ca.crt -CAkey ca.key -days 3560 -CAcreateserial -out client2.crt
```


### To view a certificate request:
```
openssl req -in client1.csr -noout -text
```

### To view a certificate
```
openssl x509 -in client1.crt -noout -text
```

### To verify a certificate
```
openssl verify -CAfile ca.crt client1.crt
```
