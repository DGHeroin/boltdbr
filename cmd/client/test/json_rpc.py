# -*- coding: utf-8 -*-
import json
import socket
import itertools
import time
import base64

class RPCClient(object):

    def __init__(self, addr, codec=json):
        self._socket = socket.create_connection(addr)
        self._id_iter = itertools.count()
        self._codec = codec

    def _message(self, name, *params):
        return dict(id=self._id_iter.next(),
                    params=list(params),
                    method=name)

    def call(self, name, *params):
        req = self._message(name, *params)
        id = req.get('id')

        mesg = self._codec.dumps(req)
        print(mesg)
        self._socket.sendall(mesg)

        # This will actually have to loop if resp is bigger
        resp = self._socket.recv(4096)
        resp = self._codec.loads(resp)

        if resp.get('id') != id:
            raise Exception("expected id=%s, received id=%s: %s"
                            %(id, resp.get('id'), resp.get('error')))

        if resp.get('error') is not None:
            raise Exception(resp.get('error'))

        return resp.get('result')

    def close(self):
        self._socket.close()


if __name__ == '__main__':
    rpc = RPCClient(("127.0.0.1", 8192))
    mv = dict(Bucket=base64.b64encode("default"), Key=base64.b64encode('my-key'))
    rs = rpc.call("BoltDBR.Get", mv)

    if rs['Error'] == '':
        value = rs['Value']
        print(base64.b64decode(value))