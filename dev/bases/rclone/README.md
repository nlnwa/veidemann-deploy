Security in rclone is based on public-key cryptography. 

The rclone server is configured with a list of public keys to authenticate clients.
The client must have a private key matching one of the public keys to authenticate.

### How to generate keypair

Generate keypair (id_rsa, id_rsa.pub):

    # -C comment (usually an email addres)
    # -f output keyfile
    # -N passphrase (for key decryption)
    
    ssh-keygen -C admin@example.com -f id_rsa -N ""

Move `id_rsa` to `client` folder and `id_rsa.pub` to `server` folder.
