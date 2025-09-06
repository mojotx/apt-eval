# TLS Certificates Directory

Place your TLS certificates in this directory:

- `wildcard.crt`: Your wildcard SSL/TLS certificate file
- `wildcard.key`: Your private key file (keep this secure!)

## Important Security Notes

1. **Never commit your private key files to version control**
2. Make sure the private key has appropriate file permissions (e.g., `chmod 600 wildcard.key`)
3. In production environments, consider using a certificate manager or secure storage solution

## Using Custom Certificate Paths

You can specify custom paths for your certificates using environment variables:

```bash
CERT_FILE=/path/to/your/certificate.crt KEY_FILE=/path/to/your/private.key ./apt-eval
```
