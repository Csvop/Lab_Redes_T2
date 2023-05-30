import socket
import pickle
import time
import typing

UDP_IP = "localhost"
UDP_PORT = 34755

# Pacote, contendo o número de sequência e o payload
class Packet:
    def __init__(self, seq, payload):
        self.seq = seq
        self.payload = payload

writer = open("payload_received.txt", "wb")

sock = socket.socket(socket.AF_INET, # Internet
                     socket.SOCK_DGRAM) # UDP
sock.bind((UDP_IP, UDP_PORT))

print("UDP target IP: %s" % UDP_IP)

current_data: bytes = b''

while True:
    data, addr = sock.recvfrom(128) # buffer size is 32 bytes

    current_data += data

    data: Packet = pickle.loads(data)
    writer.write(data.payload)
    print("received message: %s" % str(data))
    writer.close()

    # sleep
    time.sleep(3)