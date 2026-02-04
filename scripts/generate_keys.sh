#!/bin/bash

# Script to generate RSA keys for JWT signing
# These keys are used for RS256 algorithm

echo "Generating RSA keys for JWT..."

# Create keys directory if it doesn't exist
mkdir -p keys

# Generate private key (4096 bits for strong security)
echo "Generating private key..."
openssl genrsa -out keys/private_key.pem 4096

if [ $? -ne 0 ]; then
    echo "Error: Failed to generate private key"
    exit 1
fi

# Generate public key from private key
echo "Generating public key..."
openssl rsa -in keys/private_key.pem -pubout -out keys/public_key.pem

if [ $? -ne 0 ]; then
    echo "Error: Failed to generate public key"
    exit 1
fi

# Set appropriate permissions
chmod 600 keys/private_key.pem
chmod 644 keys/public_key.pem

echo ""
echo "✅ RSA keys generated successfully!"
echo ""
echo "Private key: keys/private_key.pem"
echo "Public key:  keys/public_key.pem"
echo ""
echo "⚠️  Important: Keep your private key secure and never commit it to version control!"
