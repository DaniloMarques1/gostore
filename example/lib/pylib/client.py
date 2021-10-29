import socket
import json

class Response:
    def __init__(self, s):
        self.code = 0
        self.message = ''
        self.parse_response(s)

    def parse_response(self, s):
        s = s.replace('\n', '')
        splited = s.split(';')

        codeSplit = splited[0].split('=');
        code = codeSplit[1]
        self.code = code

        msgSplit = splited[1].split('=');
        msg = msgSplit[1]
        self.message = msg

    def __str__(self):
        return f'code={self.code};message={self.message}'

class ReadOperation:
    def __init__(self, key):
        self.key = key
    
    def parse_operation(self):
        return bytes(f'op=read;key={self.key};\n', 'utf-8')

class ListOperation:
    def parse_operation(self):
        return bytes(f'op=list;\n', 'utf-8')

class KeysOperation:
    def parse_operation(self):
        return bytes(f'op=keys;\n', 'utf-8')

class StoreOperation:
    def __init__(self, key, value):
        self.key = key
        self.value = value

    def parse_operation(self):
        if type(self.value) == list or type(self.value) == dict or type(self.value) == tuple:
            self.value = json.dumps(self.value)
        return bytes(f'op=store;key={self.key};value={self.value};\n', 'utf-8')

class DeleteOperation:
    def __init__(self, key):
        self.key = key

    def parse_operation(self):
        return bytes(f'op=delete;key={self.key};\n', 'utf-8')

class Client:
    def __init__(self, host, port):
        self.host = host
        self.port = port
        self.conn = socket.socket(socket.AF_INET, socket.SOCK_STREAM)

    def connect(self):
        self.conn.connect((self.host, self.port))

    def disconnect(self):
        self.conn.close()

    def read_operation(self, key):
        read_op = ReadOperation(key)
        response = self.__send_message(read_op)
        return response
    
    def list_operation(self):
        list_op = ListOperation()
        response = self.__send_message(list_op)
        return response

    def keys_operation(self):
        keys_op = KeysOperation()
        response = self.__send_message(keys_op)
        return response

    def store_operation(self, key, value):
        store_op = StoreOperation(key, value)
        response = self.__send_message(store_op)
        return response
    
    def delete_operation(self, key):
        delete_op = DeleteOperation(key)
        response = self.__send_message(delete_op)
        return response

    def __send_message(self, operation):
        b = operation.parse_operation()
        self.conn.send(b)
        response = Response(self.conn.recv(8192).decode('utf-8')) # TODO how to read everything from conn
        return response
