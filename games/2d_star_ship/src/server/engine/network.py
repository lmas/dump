
import logging

from tornado.tcpserver import TCPServer


class NetworkServer(TCPServer):
    def __init__(self, connect_cb, disconnect_cb, data_cb, port=12345, host=''):
        TCPServer.__init__(self)
        self.connect_callback = connect_cb
        self.disconnect_callback = disconnect_cb
        self.data_callback = data_cb
        self.clients = []
        self.listen(port, host)
        logging.info('Listening on port: %i', port)

    def broadcast(self, data):
        logging.debug('Broadcast: %s', data.strip())
        for client in self.clients:
            client.write(data)

    def broadcast_all_but(self, data, excepted_client):
        logging.debug('Broadcast_all_but: %s %s', excepted_client.ip, data.strip())
        for client in self.clients:
            if client != excepted_client:
                client.write(data)

    def handle_stream(self, stream, address):
        logging.info('Client connected: {} ({})'.format(address[0], address[1]))
        new_client = Client(self, address, stream)
        self.clients.append(new_client)
        self.connect_callback(new_client)

    def client_disconnected(self, client):
        logging.info('Client disconnected: {}'.format(client.ip))
        self.clients.remove(client)
        self.disconnect_callback(client)


class Client(object):
    def __init__(self, server, address, stream):
        self.server = server
        self.ip = address[0]
        self.stream = stream
        self.stream.set_close_callback(self.handle_close)
        self.read()

    def close(self):
        self.stream.close()

    def read(self):
        self.stream.read_until('\n', self.handle_data)

    def write(self, data):
        self.stream.write('{}\n'.format(data))

    def handle_close(self):
        self.server.client_disconnected(self)

    def handle_data(self, data):
        data = data.strip()
        logging.debug('%s sent: %s', self.ip, data)
        self.server.data_callback(self, data)
        self.read()

