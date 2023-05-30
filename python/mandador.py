import socket
import pickle
import typing

# Ip do servidor
UDP_IP = "localhost"
UDP_PORT = 34754
UDP_PORT_TO_SEND = 34755

# Pacote, contendo o número de sequência e o payload
class Packet:
    def __init__(self, seq, payload):
        self.seq: int = seq
        self.payload: int = payload

def send_file(filename, chunk_size):
    # Abrir o arquivo em modo binário
    with open(filename, 'rb') as file:
        # Criar o socket
        sock = socket.socket(socket.AF_INET, # Internet
                            socket.SOCK_DGRAM) # UDP

        # Ler o arquivo em partes de tamanho chunk_size
        while True:
            data = file.read(chunk_size)
            if not data:
                # Fim do arquivo
                break

            # Enviar os dados para o socket
            packet_to_send = Packet(0, data)

            sock.sendto(pickle.dumps(packet_to_send), (UDP_IP, UDP_PORT_TO_SEND))

        # Fechar o socket
        sock.close()

send_file("payload_sent.txt", 128)










