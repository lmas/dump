#!/usr/bin/env python
'''Tool for investigating a web server's SSL certificate.

Requires the pyOpenSSL package.'''

import datetime
import ssl
import OpenSSL as ossl

def ssl_cert_details(server, port=443):
    '''Investigate a server's SSL certificate.

    Params:
    server - server address
    port - server's HTTPS port (default 443)

    Returns:
    algorithm - The SSL cert's signature algorithm
    digest - The cert's digest value
    expires - datetime.datetime object of when the cert expires
    '''
    cert = ssl.get_server_certificate((server, port), ssl.PROTOCOL_TLSv1)
    tmp = ossl.crypto.load_certificate(ossl.crypto.FILETYPE_PEM, cert)

    algo = tmp.get_signature_algorithm()
    digest = tmp.digest(tmp.get_signature_algorithm()[:3]).replace(':', '').lower()
    expires = datetime.datetime.strptime(tmp.get_notAfter(), '%Y%m%d%H%M%SZ')

    return algo, digest, expires

