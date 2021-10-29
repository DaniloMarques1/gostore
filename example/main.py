from lib.pylib.client import Client 

HOST = 'localhost'
PORT = 5000

def main():
    client = Client(HOST, PORT)
    client.connect()

    response = client.read_operation('name')
    print(response)

    response = client.store_operation('numbers', [1, 2, 3])
    print(response)

    response = client.list_operation()
    print(response)

    response = client.keys_operation()
    print(response)

    response = client.delete_operation('numbers')
    print(response)

    client.disconnect()

main()
