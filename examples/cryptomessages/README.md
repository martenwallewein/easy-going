# Cryptomessages - A Simple Example of Sending Secure Messages over the Network

In this example we show you how to implement encryption and signing (RSA and HMAC) in Go to send messages over the network in a secure way,
through adding confidentiality and authentication. The following image depicts the idea:



Sender and receiver both have a keypair of private and public key that are known and trusted to both sides. At first, the sender encrypts the message using the receivers public key. Afterwards, the sender signs the message using its own private key and appends the signature to the encrypted message. Next, the sender sends the message including the signature over the network to the receiver. The receiver now splits the received data into the encrypted message and the signature. With its own private key, it can now decrypt the encrypted message, and using the public key of the sender, it can verify the signature. Please note that the same goal could be done by using symmetric cryptography. By combining these two approaches, we achieve the following:
- Confidentiality: The sender can ensure that only the receiver can decrypt the message.
- Authentication: The receiver can ensure that only the sender could have sent this message.

## Usage
Build the application using `go build`. 

Commandline:
```
Usage: ./cryptomessages init|send|receive [args] 
       ./cryptomessages init [rsa|hmac] - Initialize key(pair)s depending on the mode
       ./cryptomessages send [rsa|hmac] address message - Send the given message via udp to the address
       ./cryptomessages receive [rsa|hmac] address - Receive a message via udp listening on the address
```

At first, create RSA keypairs, HMAC keys, or both using the init command 
`./cryptomessages init rsa` or `./cryptomessages init hmac`

Then you need to start the tool twice, at first in sending mode, second in receiving mode, given a valid UDP address:

These settings send the message hello23 using the localhost address and port 4000. `./cryptomessages receive hmac "127.0.0.1:4000"` `./cryptomessages send hmac "127.0.0.1:4000" hello23`. Please ensure that you pass the same address for sender and receiver. Replace hmac for rsa if you want to use RSA for crypto.

**Note:** This is a toy example designed to run on localhost. If you want to deploy it to multiple nodes and sending the messages over real networks, you need to distribute the public keys to all nodes and ensure that always the correct public keys are used, so you need to put a bit of effort into it. And please note that the example assumes that you trust the host's public keys (you are sure that the public keys really belong to the host).