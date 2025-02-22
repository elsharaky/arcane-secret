# Arcane Secret

Arcane Secret is a Kubernetes operator designed to generate various types of cryptographic key pairs and inject them into Kubernetes Secrets for secure use within your cluster. Currently, it supports RSA key pairs, with an option to generate them in SSH format.

## Features

- **Automated Key Pair Generation**: Automatically generates cryptographic key pairs inside Kubernetes clusters.
- **Supports Multiple Algorithms**: Currently supports RSA, with future plans to add ECDSA and ED25519.
- **SSH Format Support**: Allows RSA keys to be generated in SSH-compatible format.
- **Secure Storage**: Injects generated key pairs into Kubernetes Secrets for seamless integration.
- **Easy Deployment**: Install using Helm for quick setup.

## Installation

Arcane Secret can be easily installed using Helm:

```sh
helm repo add arcane-secret https://elsharaky.github.io/arcane-secret/
helm repo update
helm install arcane-secret arcane-secret/arcane-secret --namespace arcane-secret --create-namespace
```

## Usage

To generate an RSA key pair and store it in a Kubernetes Secret, create a `KeyPair` custom resource:

```yaml
apiVersion: api.arcanesecret.io/v1alpha1
kind: KeyPair
metadata:
  name: sample-keypair
  namespace: default
spec:
  algorithm: RSA
  size: 4096
  sshFormat: true
```

Apply the manifest using:

```sh
kubectl apply -f keypair.yaml
```

### Accessing the Generated Secret

Once the `KeyPair` resource is created, Arcane Secret generates a corresponding Kubernetes Secret. To retrieve it:

```sh
kubectl get secret sample-keypair -o yaml
```

## License

Arcane Secret is licensed under the [MIT License](LICENSE).

## Roadmap

- Support for **ECDSA** and **ED25519** key pairs
- Additional secret generation features
- More customization options for key pair storage

## Contributing

Contributions are welcome! Feel free to submit issues or pull requests to improve Arcane Secret.

## Contact

For questions or support, open an issue on the [GitHub repository](https://github.com/elsharaky/arcane-secret).
