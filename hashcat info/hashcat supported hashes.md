# Hashcat get supported hashes

## Command

```shell
hashcat --machine-readable --hash-info --quiet
```

## Output (Truncated)

```json
{
    "0": {
        "name": "MD5",
        "category": "Raw Hash",
        "slow_hash": false,
        "password_len_min": 0,
        "password_len_max": 256,
        "is_salted": false,
        "kernel_type": [
            "pure",
            "optimized"
        ],
        "example_hash_format": "plain",
        "example_hash": "8743b52063cd84097a65d1633f5c74f5",
        "example_pass": "hashcat",
        "benchmark_mask": "?b?b?b?b?b?b?b",
        "benchmark_charset1": "N/A",
        "autodetect_enabled": true,
        "self_test_enabled": true,
        "potfile_enabled": true,
        "custom_plugin": false,
        "plaintext_encoding": [
            "ASCII",
            "HEX"
        ]
    },
    "10": {
        "name": "md5($pass.$salt)",
        "category": "Raw Hash salted and/or iterated",
        "slow_hash": false,
        "password_len_min": 0,
        "password_len_max": 256,
        "is_salted": true,
        "salt_type": "generic",
        "salt_len_min": 0,
        "salt_len_max": 256,
        "kernel_type": [
            "pure",
            "optimized"
        ],
        "example_hash_format": "plain",
        "example_hash": "3d83c8e717ff0e7ecfe187f088d69954:343141",
        "example_pass": "hashcat",
        "benchmark_mask": "?b?b?b?b?b?b?b",
        "benchmark_charset1": "N/A",
        "autodetect_enabled": true,
        "self_test_enabled": true,
        "potfile_enabled": true,
        "custom_plugin": false,
        "plaintext_encoding": [
            "ASCII",
            "HEX"
        ]
    },
    "11": {
        "name": "Joomla < 2.5.18",
        "category": "Forums, CMS, E-Commerce",
        "slow_hash": false,
        "password_len_min": 0,
        "password_len_max": 256,
        "is_salted": true,
        "salt_type": "generic",
        "salt_len_min": 0,
        "salt_len_max": 256,
        "kernel_type": [
            "pure",
            "optimized"
        ],
        "example_hash_format": "plain",
        "example_hash": "b78f863f2c67410c41e617f724e22f34:89384528665349271307465505333378",
        "example_pass": "hashcat",
        "benchmark_mask": "?b?b?b?b?b?b?b",
        "benchmark_charset1": "N/A",
        "autodetect_enabled": true,
        "self_test_enabled": true,
        "potfile_enabled": true,
        "custom_plugin": false,
        "plaintext_encoding": [
            "ASCII",
            "HEX"
        ]
    },
    ...
}
```

