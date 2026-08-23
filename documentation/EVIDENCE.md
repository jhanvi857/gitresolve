# Evidence and Validation Materials

This directory contains validation artifacts, screenshots, and evidence used to verify the project's functionality and safety characteristics.

## Contents

- **_repos/**: Real-world repository references used for testing
  - `chi/`: Evidence from the chi HTTP router
  - `etcd/`: Evidence from etcd distributed configuration
  
- **raw/**: Raw validation data
  - `chi.json`: Conflict resolution results from chi
  - `cobra.json`: Conflict resolution results from cobra
  - `etcd.json`: Conflict resolution results from etcd
  
- **tools/**: Helper utilities for evidence collection and validation

## Purpose

These materials demonstrate:
- Resolution accuracy across real-world Go projects
- Conflict classification correctness
- Safety invariant validation
- Edge case handling and recovery

See the main README for more information about validation and testing strategy.
