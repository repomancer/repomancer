# Repomancer
Software repository change management

## Design Notes

- Use `gh` as the interface to GitHub because:
  - It's already good at handling edge cases like different repository formats
  - Will return JSON
  - Handles Authentication and credential storage well

- Where to put code:
  - Custom widgets hold layout
  - Windows hold logic that wires things together (sets callbacks, etc)
  - Project struct holds project state, can be loaded/saved