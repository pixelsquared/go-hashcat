# Hashcat Cracking

## Command

```shell
hashcat --optimized-kernel-enable --quiet --status --status-json --status-timer 1 --hash-type 0 --attack-mode 3 '25f9e794323b453885f5181f1b624d0b' '?a?a?a?a?a?a'
```

## Output

```
{ "session": "hashcat", "guess": { "guess_base": "?a?a?a?a?a?a", "guess_base_count": 1, "guess_base_offset": 1, "guess_base_percent": 100.00, "guess_mask_length": 6, "guess_mod": null, "guess_mod_count": 1, "guess_mod_offset": 1, "guess_mod_percent": 100.00, "guess_mode": 9 }, "status": 3, "target": "25f9e794323b453885f5181f1b624d0b", "progress": [1646960640, 735091890625], "restore_point": 172032, "recovered_hashes": [0, 1], "recovered_salts": [0, 1], "rejected": 0, "devices": [ { "device_id": 1, "device_name": "cpu-haswell-AMD Ryzen 9 3900X 12-Core Processor", "device_type": "CPU", "speed": 1742256954, "temp": 65, "util": 33 } ], "time_start": 1742933386, "estimated_stop": 1742933807 }
{ "session": "hashcat", "guess": { "guess_base": "?a?a?a?a?a?a", "guess_base_count": 1, "guess_base_offset": 1, "guess_base_percent": 100.00, "guess_mask_length": 6, "guess_mod": null, "guess_mod_count": 1, "guess_mod_offset": 1, "guess_mod_percent": 100.00, "guess_mode": 9 }, "status": 3, "target": "25f9e794323b453885f5181f1b624d0b", "progress": [3371016192, 735091890625], "restore_point": 368640, "recovered_hashes": [0, 1], "recovered_salts": [0, 1], "rejected": 0, "devices": [ { "device_id": 1, "device_name": "cpu-haswell-AMD Ryzen 9 3900X 12-Core Processor", "device_type": "CPU", "speed": 1715424049, "temp": 69, "util": 66 } ], "time_start": 1742933386, "estimated_stop": 1742933814 }
{ "session": "hashcat", "guess": { "guess_base": "?a?a?a?a?a?a", "guess_base_count": 1, "guess_base_offset": 1, "guess_base_percent": 100.00, "guess_mask_length": 6, "guess_mod": null, "guess_mod_count": 1, "guess_mod_offset": 1, "guess_mod_percent": 100.00, "guess_mode": 9 }, "status": 3, "target": "25f9e794323b453885f5181f1b624d0b", "progress": [5139111936, 735091890625], "restore_point": 565248, "recovered_hashes": [0, 1], "recovered_salts": [0, 1], "rejected": 0, "devices": [ { "device_id": 1, "device_name": "cpu-haswell-AMD Ryzen 9 3900X 12-Core Processor", "device_type": "CPU", "speed": 1718843715, "temp": 72, "util": 0 } ], "time_start": 1742933386, "estimated_stop": 1742933813 }
{ "session": "hashcat", "guess": { "guess_base": "?a?a?a?a?a?a", "guess_base_count": 1, "guess_base_offset": 1, "guess_base_percent": 100.00, "guess_mask_length": 6, "guess_mod": null, "guess_mod_count": 1, "guess_mod_offset": 1, "guess_mod_percent": 100.00, "guess_mode": 9 }, "status": 3, "target": "25f9e794323b453885f5181f1b624d0b", "progress": [6900916224, 735091890625], "restore_point": 761856, "recovered_hashes": [0, 1], "recovered_salts": [0, 1], "rejected": 0, "devices": [ { "device_id": 1, "device_name": "cpu-haswell-AMD Ryzen 9 3900X 12-Core Processor", "device_type": "CPU", "speed": 1718929550, "temp": 74, "util": 0 } ], "time_start": 1742933386, "estimated_stop": 1742933813 }
{ "session": "hashcat", "guess": { "guess_base": "?a?a?a?a?a?a", "guess_base_count": 1, "guess_base_offset": 1, "guess_base_percent": 100.00, "guess_mask_length": 6, "guess_mod": null, "guess_mod_count": 1, "guess_mod_offset": 1, "guess_mod_percent": 100.00, "guess_mode": 9 }, "status": 3, "target": "25f9e794323b453885f5181f1b624d0b", "progress": [8700469248, 735091890625], "restore_point": 958464, "recovered_hashes": [0, 1], "recovered_salts": [0, 1], "rejected": 0, "devices": [ { "device_id": 1, "device_name": "cpu-haswell-AMD Ryzen 9 3900X 12-Core Processor", "device_type": "CPU", "speed": 1727200381, "temp": 75, "util": 0 } ], "time_start": 1742933386, "estimated_stop": 1742933811 }
{ "session": "hashcat", "guess": { "guess_base": "?a?a?a?a?a?a", "guess_base_count": 1, "guess_base_offset": 1, "guess_base_percent": 100.00, "guess_mask_length": 6, "guess_mod": null, "guess_mod_count": 1, "guess_mod_offset": 1, "guess_mod_percent": 100.00, "guess_mode": 9 }, "status": 3, "target": "25f9e794323b453885f5181f1b624d0b", "progress": [10462273536, 735091890625], "restore_point": 1155072, "recovered_hashes": [0, 1], "recovered_salts": [0, 1], "rejected": 0, "devices": [ { "device_id": 1, "device_name": "cpu-haswell-AMD Ryzen 9 3900X 12-Core Processor", "device_type": "CPU", "speed": 1728915073, "temp": 76, "util": 0 } ], "time_start": 1742933386, "estimated_stop": 1742933811 }
{ "session": "hashcat", "guess": { "guess_base": "?a?a?a?a?a?a", "guess_base_count": 1, "guess_base_offset": 1, "guess_base_percent": 100.00, "guess_mask_length": 6, "guess_mod": null, "guess_mod_count": 1, "guess_mod_offset": 1, "guess_mod_percent": 100.00, "guess_mode": 9 }, "status": 3, "target": "25f9e794323b453885f5181f1b624d0b", "progress": [12261826560, 735091890625], "restore_point": 1351680, "recovered_hashes": [0, 1], "recovered_salts": [0, 1], "rejected": 0, "devices": [ { "device_id": 1, "device_name": "cpu-haswell-AMD Ryzen 9 3900X 12-Core Processor", "device_type": "CPU", "speed": 1733559892, "temp": 77, "util": 0 } ], "time_start": 1742933386, "estimated_stop": 1742933809 }
{ "session": "hashcat", "guess": { "guess_base": "?a?a?a?a?a?a", "guess_base_count": 1, "guess_base_offset": 1, "guess_base_percent": 100.00, "guess_mask_length": 6, "guess_mod": null, "guess_mod_count": 1, "guess_mod_offset": 1, "guess_mod_percent": 100.00, "guess_mode": 9 }, "status": 3, "target": "25f9e794323b453885f5181f1b624d0b", "progress": [14048796672, 735091890625], "restore_point": 1548288, "recovered_hashes": [0, 1], "recovered_salts": [0, 1], "rejected": 0, "devices": [ { "device_id": 1, "device_name": "cpu-haswell-AMD Ryzen 9 3900X 12-Core Processor", "device_type": "CPU", "speed": 1735607423, "temp": 77, "util": 0 } ], "time_start": 1742933386, "estimated_stop": 1742933809 }
{ "session": "hashcat", "guess": { "guess_base": "?a?a?a?a?a?a", "guess_base_count": 1, "guess_base_offset": 1, "guess_base_percent": 100.00, "guess_mask_length": 6, "guess_mod": null, "guess_mod_count": 1, "guess_mod_offset": 1, "guess_mod_percent": 100.00, "guess_mode": 9 }, "status": 3, "target": "25f9e794323b453885f5181f1b624d0b", "progress": [15835766784, 735091890625], "restore_point": 1744896, "recovered_hashes": [0, 1], "recovered_salts": [0, 1], "rejected": 0, "devices": [ { "device_id": 1, "device_name": "cpu-haswell-AMD Ryzen 9 3900X 12-Core Processor", "device_type": "CPU", "speed": 1736669529, "temp": 78, "util": 100 } ], "time_start": 1742933386, "estimated_stop": 1742933809 }
{ "session": "hashcat", "guess": { "guess_base": "?a?a?a?a?a?a", "guess_base_count": 1, "guess_base_offset": 1, "guess_base_percent": 100.00, "guess_mask_length": 6, "guess_mod": null, "guess_mod_count": 1, "guess_mod_offset": 1, "guess_mod_percent": 100.00, "guess_mode": 9 }, "status": 3, "target": "25f9e794323b453885f5181f1b624d0b", "progress": [17610153984, 735091890625], "restore_point": 1941504, "recovered_hashes": [0, 1], "recovered_salts": [0, 1], "rejected": 0, "devices": [ { "device_id": 1, "device_name": "cpu-haswell-AMD Ryzen 9 3900X 12-Core Processor", "device_type": "CPU", "speed": 1737255286, "temp": 78, "util": 0 } ], "time_start": 1742933386, "estimated_stop": 1742933808 }
{ "session": "hashcat", "guess": { "guess_base": "?a?a?a?a?a?a", "guess_base_count": 1, "guess_base_offset": 1, "guess_base_percent": 100.00, "guess_mask_length": 6, "guess_mod": null, "guess_mod_count": 1, "guess_mod_offset": 1, "guess_mod_percent": 100.00, "guess_mode": 9 }, "status": 3, "target": "25f9e794323b453885f5181f1b624d0b", "progress": [19384541184, 735091890625], "restore_point": 2138112, "recovered_hashes": [0, 1], "recovered_salts": [0, 1], "rejected": 0, "devices": [ { "device_id": 1, "device_name": "cpu-haswell-AMD Ryzen 9 3900X 12-Core Processor", "device_type": "CPU", "speed": 1737573366, "temp": 78, "util": 0 } ], "time_start": 1742933386, "estimated_stop": 1742933808 }
{ "session": "hashcat", "guess": { "guess_base": "?a?a?a?a?a?a", "guess_base_count": 1, "guess_base_offset": 1, "guess_base_percent": 100.00, "guess_mask_length": 6, "guess_mod": null, "guess_mod_count": 1, "guess_mod_offset": 1, "guess_mod_percent": 100.00, "guess_mode": 9 }, "status": 3, "target": "25f9e794323b453885f5181f1b624d0b", "progress": [21171511296, 735091890625], "restore_point": 2334720, "recovered_hashes": [0, 1], "recovered_salts": [0, 1], "rejected": 0, "devices": [ { "device_id": 1, "device_name": "cpu-haswell-AMD Ryzen 9 3900X 12-Core Processor", "device_type": "CPU", "speed": 1738184930, "temp": 78, "util": 0 } ], "time_start": 1742933386, "estimated_stop": 1742933808 }
{ "session": "hashcat", "guess": { "guess_base": "?a?a?a?a?a?a", "guess_base_count": 1, "guess_base_offset": 1, "guess_base_percent": 100.00, "guess_mask_length": 6, "guess_mod": null, "guess_mod_count": 1, "guess_mod_offset": 1, "guess_mod_percent": 100.00, "guess_mode": 9 }, "status": 3, "target": "25f9e794323b453885f5181f1b624d0b", "progress": [22920732672, 735091890625], "restore_point": 2531328, "recovered_hashes": [0, 1], "recovered_salts": [0, 1], "rejected": 0, "devices": [ { "device_id": 1, "device_name": "cpu-haswell-AMD Ryzen 9 3900X 12-Core Processor", "device_type": "CPU", "speed": 1735660714, "temp": 78, "util": 66 } ], "time_start": 1742933386, "estimated_stop": 1742933809 }
{ "session": "hashcat", "guess": { "guess_base": "?a?a?a?a?a?a", "guess_base_count": 1, "guess_base_offset": 1, "guess_base_percent": 100.00, "guess_mask_length": 6, "guess_mod": null, "guess_mod_count": 1, "guess_mod_offset": 1, "guess_mod_percent": 100.00, "guess_mode": 9 }, "status": 3, "target": "25f9e794323b453885f5181f1b624d0b", "progress": [24695119872, 735091890625], "restore_point": 2727936, "recovered_hashes": [0, 1], "recovered_salts": [0, 1], "rejected": 0, "devices": [ { "device_id": 1, "device_name": "cpu-haswell-AMD Ryzen 9 3900X 12-Core Processor", "device_type": "CPU", "speed": 1735683916, "temp": 78, "util": 0 } ], "time_start": 1742933386, "estimated_stop": 1742933809 }
{ "session": "hashcat", "guess": { "guess_base": "?a?a?a?a?a?a", "guess_base_count": 1, "guess_base_offset": 1, "guess_base_percent": 100.00, "guess_mask_length": 6, "guess_mod": null, "guess_mod_count": 1, "guess_mod_offset": 1, "guess_mod_percent": 100.00, "guess_mode": 9 }, "status": 3, "target": "25f9e794323b453885f5181f1b624d0b", "progress": [26475798528, 735091890625], "restore_point": 2924544, "recovered_hashes": [0, 1], "recovered_salts": [0, 1], "rejected": 0, "devices": [ { "device_id": 1, "device_name": "cpu-haswell-AMD Ryzen 9 3900X 12-Core Processor", "device_type": "CPU", "speed": 1757145618, "temp": 78, "util": 33 } ], "time_start": 1742933386, "estimated_stop": 1742933804 }
```

## Output JSON Fields

```json
{
  "session": "hashcat", // Session Name
  "guess": {
    "guess_base": "?a?a?a?a?a?a",
    "guess_base_count": 1,
    "guess_base_offset": 1,
    "guess_base_percent": 100.00,
    "guess_mask_length": 6,
    "guess_mod": null,
    "guess_mod_count": 1,
    "guess_mod_offset": 1,
    "guess_mod_percent": 100.00,
    "guess_mode": 9
  },
  "status": 3, // 3 (running), 5 (exhausted), 6 (cracked), 7 (aborted), 8 (quit)
  "target": "25f9e794323b453885f5181f1b624d0b", // The target hash or hash file being attacked
  "progress": [
      26475798528, // count of candidates processed
      735091890625 // total count of candidates
  ],
  "restore_point": 2924544, // How long until the next restore point is reached
  "recovered_hashes": [
    0, // target hashes that have been recovered
    1 // total target hashes
  ],
  "recovered_salts": [
    0, // target salts that have been recovered
    1 // total target salts
  ],
  "rejected": 0, // How many candidates were rejected as inapplicable (too long, etc)
  "devices": [
    {
      "device_id": 1,
      "device_name": "cpu-haswell-AMD Ryzen 9 3900X 12-Core Processor",
      "device_type": "CPU",
      "speed": 1757145618,
      "temp": 78, // Temperature in Celsius
      "util": 33 // Percent of utilization
    }
  ],
  "time_start": 1742933386,
  "estimated_stop": 1742933804
}
```

