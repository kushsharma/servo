# Servo

server management tool
 * Can take backups and upload it to s3
 * Clean logs older than x days
 * For DB backup, only *local* source is supported currently



## notes
### To generate the supported ssh private key, run the following:
```
$ openssl genrsa -des3 -out private.pem 4096
```
Now you have your private key. Now you need to generate the public key.

### To generate the *public* key from your private key, run the following:
```
$ openssl rsa -in private.pem -outform PEM -pubout -out public.pem
```

Now you have a PEM format for your public key. Nice! This can’t be used with SSH’s authorized_keys file though, so we’ll have to do one more conversion:

### To generate the ssh-rsa public key format, run the following:
```
$ ssh-keygen -f public.pem -i -mPKCS8 > id_rsa.pub
```