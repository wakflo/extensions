# Cryptography Integration

## Description

Cryptography Integration: Securely encrypt and decrypt sensitive data within your workflows using industry-standard algorithms such as AES, RSA, and SHA-256. Integrate with popular cryptography libraries like OpenSSL or NaCl to ensure seamless encryption and decryption of data in transit and at rest. Configure custom encryption settings for specific workflow steps, ensuring that sensitive information remains protected throughout the automation process.

**Cryptography Integration Documentation**

**Overview**
The cryptography integration in [Workflow Automation Software] enables secure data transmission and storage by encrypting sensitive information. This feature ensures that confidential data remains protected throughout its lifecycle within the workflow.

**Supported Algorithms**
The following cryptographic algorithms are supported:

* AES-256 (Advanced Encryption Standard with 256-bit key size)
* RSA-4096 (Rivest-Sha256 Algorithm with 4096-bit key size)
* SHA-256 (Secure Hash Algorithm 256)

**Encryption and Decryption**

1. **Encryption**: When a workflow task requires encryption, the software will automatically encrypt the data using the selected algorithm.
2. **Decryption**: Upon receiving encrypted data, the software will decrypt it using the corresponding decryption key.

**Key Management**
The software provides a built-in key management system to generate, store, and manage cryptographic keys. This ensures that sensitive information remains secure and accessible only to authorized users.

**Integration with Workflow Tasks**

1. **Encryption Task**: The encryption task can be added to any workflow process to encrypt data before transmission or storage.
2. **Decryption Task**: The decryption task is used to decrypt encrypted data, making it readable for further processing or analysis.

**Best Practices**
To ensure the security and integrity of sensitive information:

* Use strong passwords and keep them confidential.
* Limit access to cryptographic keys and decryption processes.
* Regularly update software and firmware to prevent vulnerabilities.
* Monitor workflow logs for any suspicious activity.

**Troubleshooting**
In case of issues with encryption or decryption, refer to the troubleshooting guide provided in the software documentation.

## Categories

- tools
- core


## Authors

- Wakflo <integrations@wakflo.com>


## Actions

| Name          | Description                                                                                                                                  | Link                             |
|---------------|----------------------------------------------------------------------------------------------------------------------------------------------|----------------------------------|
| Generate Text | Generates text based on user-defined templates and variables, allowing you to create dynamic and personalized content for various use cases. | [docs](actions/generate_text.md) |## Actions
| Hash Text     | Hashes the input text and returns a unique digital fingerprint (hash value) that can be used to verify the integrity of the original text.   | [docs](actions/hash_text.md)     |