# Hashing and Encryption (Data at Rest)

In 2012, LinkedIn was hacked, and 117 million user passwords were stolen and posted online. The passwords were stolen because LinkedIn stored them in plain text (or using an incredibly weak algorithm like SHA-1).

You must protect Data at Rest using Cryptography. 

## 1. Hashing vs Encryption

The most common mistake junior developers make is confusing Hashing with Encryption.

* **Encryption is Two-Way**: You encrypt a Social Security Number using a Secret Key. Later, you use the Secret Key to decrypt it back to the original text. (Used for Credit Cards, PII).
* **Hashing is One-Way**: You hash a Password. It creates a mathematically irreversible string of gibberish. It is theoretically impossible to get the original password back. (Used for Passwords).

## 2. Password Hashing (Bcrypt)

If you use a standard hash function like `SHA-256`, hackers can pre-compute a massive dictionary of every possible password hash (a Rainbow Table) and crack your database in seconds. Even worse, `SHA-256` is designed to be blazingly fast. A modern GPU can guess 100 billion `SHA-256` hashes per second!

You must use a **Cryptographic Key Derivation Function (KDF)** like **Bcrypt**, **Argon2**, or **scrypt**.
These algorithms are intentionally designed to be extremely slow and memory-intensive!

```go
import "golang.org/x/crypto/bcrypt"

// 1. Register: Hash the password
func HashPassword(password string) (string, error) {
    // We use a Cost of 12. 
    // This forces the algorithm to run thousands of iterations, taking ~100ms.
    // If a hacker steals the DB, guessing 1 password takes 100ms.
    // Guessing 1 billion passwords will take them thousands of years!
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
    return string(bytes), err
}

// 2. Login: Compare the incoming text with the stored Hash
func CheckPasswordHash(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}
```
*Note: Bcrypt automatically generates a random "Salt" and embeds it in the final string, making Rainbow Tables mathematically impossible.*

## 3. Symmetric Encryption (AES-GCM)

If you need to store an API Key or a Credit Card number, you must encrypt it so it can be decrypted later. 
The absolute industry standard is **AES-256-GCM** (Advanced Encryption Standard with Galois/Counter Mode).

```go
import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "io"
)

func EncryptAES(plaintext []byte, secretKey []byte) ([]byte, error) {
    // 1. Create the cipher block
    block, _ := aes.NewCipher(secretKey)
    aesgcm, _ := cipher.NewGCM(block)

    // 2. Create a Nonce (Number Used Once). This MUST be unique for every encryption!
    nonce := make([]byte, aesgcm.NonceSize())
    io.ReadFull(rand.Reader, nonce)

    // 3. Encrypt! The Nonce is prepended to the ciphertext so we can read it later.
    ciphertext := aesgcm.Seal(nonce, nonce, plaintext, nil)
    return ciphertext, nil
}
```

## 4. Key Management Systems (KMS)

If you encrypt the entire database using a Secret Key, where do you store the Secret Key?
If you hardcode it in Go, a hacker who gains access to GitHub can decrypt the database. If you put it in a Kubernetes Secret, a hacker who breaches the cluster can read it.

Enterprise companies use a **KMS** (AWS KMS or HashiCorp Vault).
The Secret Key (the Master Key) is generated inside a specialized hardware chip (HSM) at Amazon. The Master Key physically cannot leave the chip. 
If your Go application wants to encrypt a Credit Card, it makes an API call to AWS KMS, sending the plain text over the network. AWS encrypts it inside the secure hardware chip and returns the ciphertext. Your Go application never actually touches the Master Key!
