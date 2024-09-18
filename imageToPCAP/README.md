# Image to PCAP

This is a somewhat more complicated method where we hide data in a packet capture (much like a packet capture that is being used in a SOC or how network scanners see things).
In operation though, it is just as easy to use as the other examples.

## Quick start

Whatever you do, you need to install libpcap first.

### Installing libPCAP on Debian style OS's (yes, Ubuntu as well)

```bash
sudo apt install libpcap-dev -y
```

### Installing libPCAP on MacOS
Install MacPorts first.
Then do:
```bash
sudo port install libpcap
```
### Installing libPCAP on Windows
Follow the instructions at [winpcap.org](https://www.winpcap.org)

### If you need to make an executable
```bash
go mod verify && go mod tidy && go build
```

### When you have the executable
```
imageToPCAP <inputfile> <medium> <outputfile>
```

The 'medium' file here is the carrier for your data. To this end you need to find a PCAP file or capture one yourself. 
If you want to capture data yourself or if you want to see and verify the resulting file, go to [Wireshark](https://wireshark.org) and install it.

## About this method

With this method I cheat a bit because I do not deconstruct the image. Rather, I just dump in the whole of the original image file as written on disk. Yes, this means you can stego files other than images as well. But that is not the point.

The raw data is then divided into [nibbles](https://en.wikipedia.org/wiki/Nibble). These nibbles are then encoded in some TCP-flags of the packets from the medium.

To decide which flags we can modify we look at [RFC 9293 - Transmission Control Protocol (TCP)](https://www.rfc-editor.org/info/rfc9293) and we find the following flags:

```
      +========+===================+===========+====================+
      | Bit    | Name              | Reference | Assignment Notes   |
      | Offset |                   |           |                    |
      +========+===================+===========+====================+
      | 4      | Reserved for      | RFC 9293  |                    |
      |        | future use        |           |                    |
      +--------+-------------------+-----------+--------------------+
      | 5      | Reserved for      | RFC 9293  |                    |
      |        | future use        |           |                    |
      +--------+-------------------+-----------+--------------------+
      | 6      | Reserved for      | RFC 9293  |                    |
      |        | future use        |           |                    |
      +--------+-------------------+-----------+--------------------+
      | 7      | Reserved for      | RFC 8311  | Previously used by |
      |        | future use        |           | Historic RFC 3540  |
      |        |                   |           | as NS (Nonce Sum). |
      +--------+-------------------+-----------+--------------------+
      | 8      | CWR (Congestion   | RFC 3168  |                    |
      |        | Window Reduced)   |           |                    |
      +--------+-------------------+-----------+--------------------+
      | 9      | ECE (ECN-Echo)    | RFC 3168  |                    |
      +--------+-------------------+-----------+--------------------+
      | 10     | Urgent pointer    | RFC 9293  |                    |
      |        | field is          |           |                    |
      |        | significant (URG) |           |                    |
      +--------+-------------------+-----------+--------------------+
      | 11     | Acknowledgment    | RFC 9293  |                    |
      |        | field is          |           |                    |
      |        | significant (ACK) |           |                    |
      +--------+-------------------+-----------+--------------------+
      | 12     | Push function     | RFC 9293  |                    |
      |        | (PSH)             |           |                    |
      +--------+-------------------+-----------+--------------------+
      | 13     | Reset the         | RFC 9293  |                    |
      |        | connection (RST)  |           |                    |
      +--------+-------------------+-----------+--------------------+
      | 14     | Synchronize       | RFC 9293  |                    |
      |        | sequence numbers  |           |                    |
      |        | (SYN)             |           |                    |
      +--------+-------------------+-----------+--------------------+
      | 15     | No more data from | RFC 9293  |                    |
      |        | sender (FIN)      |           |                    |
      +--------+-------------------+-----------+--------------------+
```

We can, after reading the RFC, determine that we can use ECE, CWR, URG and PSH flags as 4 bits per packet. A nibble per packet.

We 'rewrite' the 'medium' but modify those TCP flags as per our nibbles and the resulting output has all the checksums corrected for our modifications.

To make sure that the viewer starts 'de-nibbling' at the correct point, we add a 'magic marker' first.

### Advantages

Can be done on a live wire. If the recipient network stack 'knows' about this, the data can be sent piggybacked on a regular transmission.

Checking PCAP's is not something any system will do for fun. There is a lot of 'surplus' data here and just checking on the off-chance that something is in there will not happen with a lot of systems.

It can hide anything.

### Disadvantages

Recipient needs to do significant work to get the data out of the medium. The supplied 'viewer' can do this, but it is hard to debug if some bit has flipped on the way.

PCAPs are huuuuge. We only use 4 bits per TCP packet, which can be multiple Kilobytes in payload. It is very, very inefficient.

TCP flags can be parsed easily (not saying quickly because you still need to process the whole PCAP) so writing a detector is not hard.