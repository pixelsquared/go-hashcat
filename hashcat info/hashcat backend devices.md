# Hashcat get backend devices

## Command

```shell
hashcat --backend-info --quiet
```

## Output

```
OpenCL Info:
============

OpenCL Platform ID #1
  Vendor..: The pocl project
  Name....: Portable Computing Language
  Version.: OpenCL 3.0 PoCL 6.0  Linux, Release, RELOC, SPIR-V, LLVM 18.1.8, SLEEF, DISTRO, POCL_DEBUG

  Backend Device ID #1
    Type...........: CPU
    Vendor.ID......: 1
    Vendor.........: AuthenticAMD
    Name...........: cpu-haswell-AMD Ryzen 9 3900X 12-Core Processor
    Version........: OpenCL 3.0 PoCL HSTR: cpu-x86_64-unknown-linux-gnu-haswell
    Processor(s)...: 24
    Clock..........: 4673
    Memory.Total...: 29934 MB (limited to 4096 MB allocatable in one block)
    Memory.Free....: 14935 MB
    Local.Memory...: 512 KB
    OpenCL.Version.: OpenCL C 1.2 PoCL
    Driver.Version.: 6.0
```

