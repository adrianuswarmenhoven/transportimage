sudo apt install libpcap-dev -y

The authority for this tutorial is RFC 793, Transmission Control Protocol, found here. I would suggest this is required reading. The RFC's are the "bible" if you will. If someone wants to write an operating system, (OS), or a service that interacts with a network over TCP/IP then they need to make sure that they comply with the RFC or it simply won't work with everything else. This RFC works hand in hand with the IP Protocol, RFC 791, found here, hence the common nomenclature TCP/IP. The IP Protocol does contain flags but they are not to be confused with those in the TCP protocol insofar as they only refer to the fragmentation state of the packets being sent/received.

The Flags

The flags themselves indicate "connection status". Since TCP is a stateful and reliable protocol it is imperative that we know the purpose of each packet. This is what the flags are for, to designate what kind of packet, and therefore the connection state, we are looking at. The flags themselves are contained in a single byte in the packet with each bit having a specific meaning and task. They are as follows:-

Bit:

1 . CWR: Congestion Window Reduced
2 . ECN: Echo
3 . URG: Urgent
4 . ACK: Acknowledge
5 . PSH: Push
6 . RST: Reset
7 . SYN: Synchronize
8 . FIN: Finish/End

For the purpose of this tutorial I will forget about flags 1-3. They are available but not commonly seen in "normal" transmissions and tend not to be commonly used in abuse situations. That leaves us with 4-8 that are commonly seen in transmissions and are regularly "abused".
