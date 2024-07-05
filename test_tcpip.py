import socket
import time
import concurrent.futures
from functools import partial

HOST="localhost"
PORT=8080

def make_tcp_connection(
    host: str,
    port: int,
    message: str,
    sleep_time: int,
    id: int,
):
  try:
    with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as sock:
      print("trying connection", id)
      sock.connect((host, port))
      print("connected", id)
      time.sleep(sleep_time)
      sock.sendall(message.encode())
  except ConnectionRefusedError:
    print(f"Connection refused to {host}:{port}")
    return None
  except Exception as e:
    print(f"Error connecting to {host}:{port}: {e}")
    return None


def test_long_connection():
    make_tcp_connection(HOST, PORT, "test message", 10, 0)


def make_connection():

    preset_connection = partial(make_tcp_connection, HOST, PORT, "hello world", 4)
    with concurrent.futures.ThreadPoolExecutor() as executor:
        futures = []
        for idx in range(11):
            futures.append(executor.submit(preset_connection, idx))

        for future in concurrent.futures.as_completed(futures):
            future.result()

if __name__=="__main__":
    make_connection()
