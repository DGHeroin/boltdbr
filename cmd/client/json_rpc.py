# -*- coding: utf-8 -*-
import requests

"""
JSON RPC client using Python
"""
def rpc_call():
    url = 'http://localhost:8192/'
    payload = {
        'id': 1,
        'method': 'BoltDBR.Get',
        'params': [{'Bucket': 'default', 'Key': 'my-key'}]
    }
    r = requests.post(url, json=payload)
    print r.text


if __name__ == '__main__':
    rpc_call()